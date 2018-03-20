// Package specconv generates all you need from the device config file.
package specconv

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/pierrchen/avs/spec"
	"github.com/pierrchen/avs/utils"
	"github.com/pierrchen/avs/vdts"
)

var (
	avsInstallDir         = GetAvsInstallDir()
	defaultConfigFile     = getDefaultConfigFile()
	avsPackage            = "github.com/pierrchen/avs"
	defaultConfigJSONName = "config.json"
)

func createDeviceAndKernelDir(vendor, device string) error {
	// create device directory
	absGenDir, _ := filepath.Abs(filepath.Join(vendor, device))
	utils.CreateGenDirIfNotExsit(absGenDir)
	// create kernel directory
	absKernelDir, _ := filepath.Abs(filepath.Join(vendor, device+"-kernel"))
	utils.CreateGenDirIfNotExsit(absKernelDir)

	return nil
}

var avsstate = AvsState{}

// InitDeviceConfig generate scaffolding
// genDir - absolute path for putting the generated stuff
func InitDeviceConfig(vendor, device string, config string) error {
	createDeviceAndKernelDir(vendor, device)

	var spec *spec.Spec
	var err error
	if config == "" {
		spec, err = loadTemplateSpec()

	} else {
		f, _ := filepath.Abs(config)
		spec, err = LoadSpec(f)
	}

	if err != nil {
		log.Fatalln("Error when creating scaffold config", err)
	}

	enrichTemplateSpec(spec, vendor, device)

	deviceDir, _ := filepath.Abs(filepath.Join(vendor, device))
	f := filepath.Join(deviceDir, defaultConfigJSONName)
	SaveSpecToJSON(spec, f)

	avsstate.GenDir = deviceDir
	// no need to validate the spec, the default one is always valid
	generateAll(spec, deviceDir)

	avsstate.Update()
	return nil
}

// fix a few things from the template json basing on the vendor, device name
func enrichTemplateSpec(spec *spec.Spec, vendor, device string) {
	kernelDir := device + "-kernel"
	kernelImage := "Image"
	dtbImage := device + ".dtb"

	spec.Product.Name = device
	spec.Product.Device = device
	spec.Product.Brand = device
	spec.Product.Model = device
	spec.Product.Manufacture = vendor
	spec.BootImage.Kernel.LocalKernel = filepath.Join("device", vendor, kernelDir, kernelImage)
	spec.BootImage.Kernel.LocalDTB = filepath.Join("device", vendor, kernelDir, dtbImage)
	spec.BoardConfig.SEPolicy.Dir = filepath.Join("device", vendor, device, "sepolicy")
}

func generateAll(spec *spec.Spec, genDir string) error {
	addProductSpecificFileMapping(spec)

	for file, tmpl := range tmlMap {
		path := filepath.Join(genDir, file)
		outFile, err := os.Create(path)
		if err != nil {
			log.Printf("faild to create %s\n", path)
			return err
		}
		defer outFile.Close()

		executeTemplate(outFile, tmpl, spec)
		avsstate.GenereatedFiles = append(avsstate.GenereatedFiles, outFile.Name())
	}

	generateRcScripts(spec, genDir)
	return nil
}

func generateRcScripts(spec *spec.Spec, genDir string) error {
	for _, rc := range spec.BootImage.Rootfs.InitRc {
		// use rc.File directly, don't generate
		if rc.File != "" {
			break
		}

		path := filepath.Join(genDir, rc.Name)
		outFile, err := os.Create(path)
		if err != nil {
			log.Printf("faild to create %s", path)
			return err
		}
		defer outFile.Close()
		executeTemplateForRc(outFile, &rc)
		avsstate.GenereatedFiles = append(avsstate.GenereatedFiles, outFile.Name())
	}
	return nil
}

// ValdiateDeviceConfig validate the default config file in the path as specified by absGenDir
func ValdiateDeviceConfig(absGenDir string) (err error) {
	specFile := filepath.Join(absGenDir, defaultConfigJSONName)
	spec, err := LoadSpec(specFile)
	if err != nil {
		log.Fatalln(err)
	}

	pass := true

	spec = override(spec, absGenDir)

	if err = vdts.ValdiateSpec(spec, absGenDir); err != nil {
		pass = false
	}

	if pass == false {
		return fmt.Errorf("validation failed")
	}
	return nil
}

func hasHal(spec *spec.Spec, name string) (int, bool) {
	for i, hal := range spec.Hals {
		if hal.Name == name {
			return i, true
		}
	}
	return -1, false
}

// find hal.*.overlay file in the dir, the * is the hal feature name.
// if there is already a corresponding feature spec in the spec passed in,
// the new one will override it; otherwise, the new one will be added.
func override(spec *spec.Spec, dir string) *spec.Spec {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return spec
	}

	for _, f := range files {
		if strings.HasSuffix(filepath.Base(f.Name()), ".overlay") &&
			strings.HasPrefix(filepath.Base(f.Name()), "hal.") {
			halSpec, err := LoadHalSpec(f.Name())
			if err != nil {
				fmt.Printf("Fail to load hal override spec %s\n", f.Name())
				break
			}
			fmt.Printf("Loading overlay spec %s\n", filepath.Base(f.Name()))
			index, has := hasHal(spec, halSpec.Name)
			if has {
				spec.Hals[index] = *halSpec
			} else {
				spec.Hals = append(spec.Hals, *halSpec)
			}

		}
	}

	return spec
}

// UpdateDeviceConfigs updates the device configrations.
// There must be already a config.json in path. Everthing will be regenerated at the moment.
// TODO: rengerate only the things that changed for rebuild performance, especiall the stuff in
// BoardConfig.mk
func UpdateDeviceConfigs(deviceDir string) error {

	specFile := filepath.Join(deviceDir, defaultConfigJSONName)
	spec, err := LoadSpec(specFile)

	if err != nil {
		fmt.Printf("err loading the spec file %s", specFile)
		return fmt.Errorf("err loading the spec file %s", specFile)
	}

	spec = override(spec, deviceDir)

	err = vdts.ValdiateSpec(spec, deviceDir)

	if err != nil {
		fmt.Println("spec validation failed, please fix the errors first")
		return fmt.Errorf("spec validation failed, please fix the errors first")
	}

	generateAll(spec, deviceDir)

	avsstate.Update()
	return nil
}

// CleanDeviceConfigs clean up all the genereated files
func CleanDeviceConfigs(deviceDir string) error {
	state, err := LoadAvsState(deviceDir)
	if err != nil {
		return err
	}

	for _, f := range state.GenereatedFiles {
		os.Remove(f)
	}

	return nil
}
