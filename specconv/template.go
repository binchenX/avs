package specconv

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/pierrchen/avs/spec"
	"github.com/pierrchen/avs/utils"
)

var tmlMap = map[string]string{
	"vendorsetup.sh":     "vendorsetup.tpl",
	"AndroidProducts.mk": "androidproducts.tpl",
	"BoardConfig.mk":     "boardconfig.tpl",
	"device.mk":          "device.tpl",
	"manifest.xml":       "manifests.tpl",
}

// getGenFileName return the path for geneated file
func getGenFileName(name string) string {
	return name + ".gen"
}

// called by GenerateScaffold to setup product specfic file and template mapping
func addProductSpecificFileMapping(spec *spec.Spec) {
	productMK := fmt.Sprintf("%s.mk", spec.Product.Name)
	tmlMap[productMK] = "product.tpl"

	// uevent.rc
	rc := spec.BootImage.Rootfs.UeventRc
	if rc.File == "" {
		ueventRc := rc.Name
		if ueventRc == "" {
			ueventRc = getGenFileName("ueventd.rc")
		}
		tmlMap[ueventRc] = "uevent.tpl"
	}

	// fstab
	fs := spec.BootImage.Rootfs.Fstab
	fileName := "fstab." + spec.Product.Name
	if fs.Name != "" {
		fileName = fs.Name
	}
	tmlMap[fileName] = "fstab.tpl"
}

// hardcoded by Android framework and used the Android Device configure files
func getFeatureFileSrcDir() string {
	return "frameworks/native/data/etc"
}

func getFeatureFileDestDir() string {
	return "system/etc/permissions"
}

func hasVendorPartition(pt *spec.PartitionTable) bool {
	for _, p := range pt.Partitions {
		if p.Name == "vendor" {
			return true
		}
	}
	return false
}

func getVendorOut(pt *spec.PartitionTable) string {
	if hasVendorPartition(pt) {
		return "vendor"
	} else {
		return "system/vendor"
	}
}

// CopyPackage.Src
// Always return a relative path that is relative to $(ANDROID_BUILD_TOP)
// Default use "$(LOCAL_PATH)" means it is relative to the device config dir
// For framework/**, use it directly
// we can do some smart checking here, e.g:
// 1. for feature files, it will always in framework
// 2. for binary copy it must be in vendor/xx
// 3. for rc file it should be in device config dir, i.e LOCAL_PATH
func getCopyInstruction(cp spec.CopyPackage) string {
	var dst string
	if cp.DestDir == "" {
		if strings.HasSuffix(cp.Src, ".so") {
			dst = "$(TARGET_COPY_OUT_VENDOR)" + "/lib"
		} else {
			dst = "$(TARGET_COPY_OUT_VENDOR)" + "/bin"
		}
	} else {
		dst = "$(TARGET_COPY_OUT_VENDOR)" + "/" + cp.DestDir
	}

	return cp.Src + ":" + dst + "/" + filepath.Base(cp.Src)
}

func getInheritProductMkDir(product string) string {
	return "$(SRC_TARGET_DIR)/product/" + product + ".mk"
}

// removePackageTag remove the possible tag
func removePackageTag(packageName string) string {
	ps := strings.Split(packageName, ":")
	return ps[0]
}

func rcInstallDest(rc *spec.RcScripts) string {
	if rc.ServicRc == "true" {
		// rc to start a service for a particular hal copy to here
		return "$(TARGET_COPY_OUT_VENDOR)/etc/init/"
	}
	return "root/"
}

func getInitRcCopyStatement(rcs []spec.RcScripts) string {
	var s []string
	for _, rc := range rcs {
		name := rc.File
		if name == "" {
			name = rc.Name
		}

		filepath.Base(name)
		t := `    $(LOCAL_PATH)/` + name + ":" + rcInstallDest(&rc) + filepath.Base(name) + ` \`
		s = append(s, t)
	}

	return strings.Join(s[:], "\n")
}

func getUeventdCopySrc(uc spec.UeventRc) string {
	if uc.File != "" {
		return `$(LOCAL_PATH)/` + uc.File
	}

	src := "ueventd.rc.gen"
	if uc.Name != "" {
		src = uc.Name
	}
	return `$(LOCAL_PATH)/` + src
}

func getFstabCopySrc(fs spec.Fstab) string {
	src := "fstab.rc.gen"
	if fs.Name != "" {
		src = fs.Name
	}
	return `$(LOCAL_PATH)/` + src
}

// RuntimeConfigInstructions turns the RuntimeConfig to a Android statement
func RuntimeConfigInstructions(config spec.RuntimeConfig) string {
	from := config.Src
	dstDir := "system/etc"
	if config.DestDir != "" {
		dstDir = config.DestDir
	}
	dst := dstDir + "/" + filepath.Base(from)
	return from + ":" + dst
}

func generate(tmpl *template.Template, f *os.File, data interface{}) error {
	// TODO: remove the genDir from the f.Name()
	fmt.Printf("generate file %s\n", f.Name())
	return tmpl.Execute(f, data)
}

// UserImageExt4 return true if any of the user images (system, cache, userdata) is
// ext4 format. All of those image will usre mkuserimg.sh script to create the correct
// image with correct file system.
// TARGET_USERIMAGES_USE_EXT4 := true need to be set so that
// 2. Build system knows we need ext filesystem, and it is ext4 variant
// 3) Serveral packages (MKEXTUSERIMG) $(MAKE_EXT4FS) $(E2FSCK)) will be built in the host
// ideally, we should also have userImageExt2, userImageExt3 but nobody should
// use that ext2, 3 nowadays.
func UserImageExt4(boardConfig *spec.BoardConfig) bool {
	for _, pt := range boardConfig.PartitionTable.Partitions {
		if strings.Contains(pt.Type, "ext4") {
			return true
		}
	}
	return false
}

func getFirmwareLocation(spec *spec.Spec) string {
	if hasVendorPartition(&spec.BoardConfig.PartitionTable) {
		return "/vendor/firmware"
	}
	return "/system/etc/firmware"
}

// getFullKernelCommand return the full kernel command line
func getFullKernelCommand(spec *spec.Spec) string {
	var s []string

	// 1. androidboot.xxx
	s = append(s, "androidboot.hardware="+spec.Product.Device)
	selinuxMode := spec.BoardConfig.SELinux.Mode
	// default to enforcing
	if selinuxMode == "" {
		selinuxMode = "enforcing"
	}
	s = append(s, "androidboot.selinux="+selinuxMode)
	// 2. stanard kernel cmdline
	s = append(s, "firmware_class.path="+getFirmwareLocation(spec))
	// 3. vendor specific kernel cmdline
	s = append(s, spec.BootImage.Kernel.CmdLine)
	return strings.Join(s, " ")
}

// InstsallFirmware is the instruction to install firmware on target
func InstsallFirmware(src string) string {
	// FIXME: always use /vendor/firmware now
	return src + ":" + "/vendor/firmware/" + filepath.Base(src)
}

// f - generated output file
// templateFile - template
// spec - spec
func executeTemplate(f *os.File, templateFile string, spec *spec.Spec) (err error) {
	funcMap := template.FuncMap{
		"ToUpper":                   strings.ToUpper,
		"FeatureFileSrcDir":         getFeatureFileSrcDir,
		"FeatureFileDestDir":        getFeatureFileDestDir,
		"CopyInstruction":           getCopyInstruction,
		"InheritProduct":            getInheritProductMkDir,
		"getGenFileName":            getGenFileName,
		"removeTag":                 removePackageTag,
		"getInitRcCopyStatement":    getInitRcCopyStatement,
		"getUeventdCopySrc":         getUeventdCopySrc,
		"getFstabCopySrc":           getFstabCopySrc,
		"RuntimeConfigInstructions": RuntimeConfigInstructions,
		"UserImageExt4":             UserImageExt4,
		"getFullKernelCommand":      getFullKernelCommand,
		"getVendorOut":              getVendorOut,
		"InstsallFirmware":          InstsallFirmware,
	}
	tmpFile := filepath.Join(avsInstallDir, "tmpl", templateFile)

	tmpl, err := template.New(templateFile).Funcs(funcMap).ParseFiles(tmpFile)
	if err != nil {
		fmt.Println("create template failed", tmpFile, err)
		return err
	}

	return generate(tmpl, f, spec)
}

// executeTemplateForRc genereate Rcscript only
func executeTemplateForRc(f *os.File, rc *spec.RcScripts) (err error) {
	templateFile := "initrc.tpl"
	funcMap := template.FuncMap{}
	tmpFile := filepath.Join(avsInstallDir, "tmpl", templateFile)

	tmpl, err := template.New(templateFile).Funcs(funcMap).ParseFiles(tmpFile)
	if err != nil {
		fmt.Println("rcScript template failed", tmpFile, err)
		return err
	}

	return generate(tmpl, f, rc)
}

// createOrUpdateHALDirs create the HAL dirs
// TODO: dont create any HAL dir?
func createOrUpdateHALDirs(spec *spec.Spec, absGenDir string) error {
	for _, h := range spec.Hals {
		path := filepath.Join(absGenDir, h.Name)
		if r, _ := utils.FileExists(path); r == false {
			// need x permission so be able to cd into
			if err := os.Mkdir(path, 0777); err != nil {
				fmt.Println("faild to create hal dir", path)
			}
		}
	}

	return nil
}
