package vdts

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/pierrchen/avs/spec"
	"github.com/pierrchen/avs/utils"
)

// ValidateSystemImage valdiate if the HAL OK
func ValidateSystemImage(spec *spec.Spec, genDir string) error {
	return validateAll(spec, genDir, []IVal{
		validateFeatureFiles,
		validateHalRuntimeConfigs,
		//validateHalPackagesBuild,
		validateHalPackagesCopy,
		validateHalInitRc,
	})
}

// validateHalInitRc validate the hal initrc
func validateHalInitRc(s *spec.Spec, genDir string) error {
	// pre-generation validation

	var rcs []spec.RcScripts

	for _, hal := range s.Hals {
		if hal.InitRc != nil {
			rcs = append(rcs, hal.InitRc...)
		}
	}

	for _, rc := range rcs {
		if (rc.File != "" && rc.Name != "") ||
			(rc.File == "" && rc.Name == "") {
			fmt.Printf("rc.File %v and rc.Name %s can't both empty or non-empty", rc.File, rc.Name)
		}
	}
	// post-generattion validation

	// making sure all the init.rc are generated
	for _, rc := range rcs {
		srcFile := rc.Name
		if rc.File != "" {
			srcFile = rc.File
		}
		srcFile = "$(LOCAL_PATH)/" + srcFile
		err := validateCopySrc(srcFile, genDir)
		if err != nil {
			fmt.Printf("[avs v] can't find copy sources: %s\n", srcFile)
		}
	}

	return nil
}

// validateHalPackagesCopy vaildiate the binary blob copy, including
// the share libary, the firmware and the kernel drivers
func validateHalPackagesCopy(spec *spec.Spec, genDir string) error {
	var allCopyPkgs []string

	for _, h := range spec.Hals {
		if h.Packages != nil && h.Packages.Copy != nil {
			for _, cp := range h.Packages.Copy {
				allCopyPkgs = append(allCopyPkgs, cp.Src)
			}
		}

		if h.Firmwares != nil {
			for _, f := range []string(*(h.Firmwares)) {
				allCopyPkgs = append(allCopyPkgs, f)
			}
		}

		if h.Drivers != nil {
			for _, d := range []string(*(h.Drivers)) {
				allCopyPkgs = append(allCopyPkgs, d)
			}
		}
	}

	var invalid []string
	for _, p := range allCopyPkgs {
		if err := validateCopySrc(p, genDir); err != nil {
			invalid = append(invalid, p)
		}
	}

	if len(invalid) != 0 {
		return fmt.Errorf("can't find copy sources: %v", invalid)
	}
	return nil
}

func validateHalRuntimeConfigs(spec *spec.Spec, genDir string) error {
	var invalid []string
	for _, h := range spec.Hals {
		for _, c := range h.RuntimeConfigs {
			if err := validateCopySrc(c.Src, genDir); err != nil {
				invalid = append(invalid, c.Src)
			}
		}
	}

	if len(invalid) != 0 {
		return fmt.Errorf("can't find copy sources: %v", invalid)
	}
	return nil
}

// validateFeatureFiles validate all the feautres files, it can be found in
// either HAL declaration or BoardConfigs
func validateFeatureFiles(spec *spec.Spec, genDir string) error {
	var invalid []string

	var allFeatures []string

	// fixed Board Features
	for _, f := range spec.BoardConfig.BoardFeatures {
		allFeatures = append(allFeatures, f)
	}

	// HAL features
	for _, h := range spec.Hals {
		for _, f := range h.Features {
			allFeatures = append(allFeatures, f)
		}
	}

	// validate all the feature files
	featureFileDir := "frameworks/native/data/etc/"
	for _, f := range allFeatures {
		if err := validateCopySrc(featureFileDir+f, genDir); err != nil {
			invalid = append(invalid, featureFileDir+f)
		}
	}

	if len(invalid) != 0 {
		return fmt.Errorf("can't find features files: %v", invalid)
	}
	return nil
}

// path start with "$(LOCAL_PATH)" must in $genDir
// all the other paths are relative the to ${ANDROID_BUILD_TOP}
func validateCopySrc(src string, genDir string) error {
	//	fmt.Printf("validate src copy %s\n", src)
	L := "$(LOCAL_PATH)"
	if strings.HasPrefix(src, L) {
		p := filepath.Join(genDir, src[len(L):])
		if r, _ := utils.FileExists(p); r == false {
			return fmt.Errorf("%s(:%s)dont't exsits", src, p)
		}
	} else {
		androiTop := os.Getenv("ANDROID_BUILD_TOP")
		if androiTop == "" {
			fmt.Println("warning: ignore feature file validating, since ${ANDROID_BUILD_TOP} was not set")
			return nil
		}

		p := filepath.Join(androiTop, src)
		if r, _ := utils.FileExists(p); r == false {
			return fmt.Errorf("%s dont't exsits", p)
		}

	}
	return nil
}

func validateHalPackagesBuild(spec *spec.Spec, genDir string) error {
	for feature, vendorPackges := range GetVendorPackges(spec) {
		// TODO: What if there are multiply Android.mk for this feature HAL
		// Or it is not in the location we assume
		mk := filepath.Join(genDir, feature, "Android.mk")
		if r, _ := utils.FileExists(mk); r == false {
			return fmt.Errorf("can't find %s", mk)
		}

		for _, module := range vendorPackges {
			// make sure the Android.mk has a Model called with the name in check
			if r, _ := HasModulesInFile(mk, module); r == true {
				return nil
			}
			return fmt.Errorf("no moduel %s in %s", module, mk)
		}

	}
	return nil
}

// HasModulesInFile test if there is a module in file
func HasModulesInFile(file string, module string) (bool, error) {
	f, err := os.Open(file)

	if err != nil {
		fmt.Printf("can't open %s", file)
		return false, err
	}

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		t := scanner.Text()
		patten := fmt.Sprintf("LOCAL_MODULE := %s", module)
		matched, _ := regexp.MatchString(patten, t)
		if matched {
			return true, nil
		}
	}

	return false, nil
}

// GetVendorPackges return the packages that is build from vendor HAL implemenation
// that is the package with a tag of "v" at the end of the package name
func GetVendorPackges(spec *spec.Spec) map[string][]string {
	// feature-> vendor packages
	var ppp map[string][]string
	ppp = make(map[string][]string)
	for _, h := range spec.Hals {
		var packages []string
		if h.Packages != nil && h.Packages.Build != nil {
			for _, p := range h.Packages.Build {
				ps := strings.Split(p, ":")
				if len(ps) == 2 {
					if ps[1] == "v" {
						packages = append(packages, ps[0])
					}
				}
			}
		}

		if len(packages) != 0 {
			ppp[h.Name] = packages
		}
	}

	return ppp
}
