package specconv

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/pierrchen/avs/spec"
	"github.com/pierrchen/avs/tmpl"
	"github.com/pierrchen/avs/utils"
)

var tmlMap = map[string]string{
	"vendorsetup.sh":     tplVendorSetup,
	"AndroidProducts.mk": tplAndriodProduct,
	"BoardConfig.mk":     tplBoard,
	"device.mk":          tplDevice,
	"manifest.xml":       tplManifest,
}

// all the templates
const (
	tplVendorSetup    string = "vendorsetup.tpl"
	tplAndriodProduct string = "androidproducts.tpl"
	tplBoard          string = "boardconfig.tpl"
	tplDevice         string = "device.tpl"
	tplManifest       string = "manifests.tpl"
	tplProduct        string = "product.tpl"
	tplUevent         string = "uevent.tpl"
	tplFstab          string = "fstab.tpl"
	tplUsbRc          string = "usb.tpl"
	tplInitRc         string = "initrc.tpl"
)

const (
	// src path
	featureFileSrc string = "frameworks/native/data/etc"
	productDir     string = "$(SRC_TARGET_DIR)/product"
	// dest path
	defaultFeatureFileDst   string = "system/etc/permissions"
	defaultFirmwareDst      string = outVendorDir + "/firmware"
	defaultKernelModuleDst  string = outVendorDir + "/lib/modules"
	defaultRuntimeConfigDst string = "system/etc"
)

// some variables used and expected by android build system
const (
	// set by avs and used by Android Build system
	outVendorDir string = "$(TARGET_COPY_OUT_VENDOR)"
	// used by Android build system to find a file in host system
	// as copy source path
	copyLocal string = "$(LOCAL_PATH)"
)

// getGenFileName return the path for geneated file
func getGenFileName(name string) string {
	return name + ".gen"
}

// called by GenerateScaffold to setup product specfic file and template mapping
func addProductSpecificFileMapping(spec *spec.Spec) {
	productMK := fmt.Sprintf("%s.mk", spec.Product.Name)
	tmlMap[productMK] = tplProduct

	// uevent.rc
	rc := spec.BootImage.Rootfs.UeventRc
	if rc.File == "" {
		ueventRc := rc.Name
		if ueventRc == "" {
			ueventRc = getGenFileName("ueventd.rc")
		}
		tmlMap[ueventRc] = tplUevent
	}

	// fstab.hw
	fs := spec.BootImage.Rootfs.Fstab
	fileName := "fstab." + spec.Product.Name
	if fs.Name != "" {
		fileName = fs.Name
	}
	tmlMap[fileName] = tplFstab

	// init.hw.usb.rc
	if spec.BoardConfig.USBGadget != nil {
		usbRcFile := fmt.Sprintf("rootfs/init.%s.usb.rc", spec.Product.Name)
		tmlMap[usbRcFile] = tplUsbRc
	}
}

// hardcoded by Android framework and used the Android Device configure files
func getFeatureFileSrcDir() string {
	return featureFileSrc
}

func getFeatureFileDestDir() string {
	return defaultFeatureFileDst
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
			dst = outVendorDir + "/lib"
		} else {
			dst = outVendorDir + "/bin"
		}
	} else {
		dst = join(outVendorDir, cp.DestDir)
	}

	return cp.Src + ":" + join(dst, filepath.Base(cp.Src))
}

func getInheritProductMkDir(product string) string {
	return join(productDir, product+".mk")
}

// removePackageTag remove the possible tag
func removePackageTag(packageName string) string {
	ps := strings.Split(packageName, ":")
	return ps[0]
}

func rcInstallDest(rc *spec.RcScripts) string {
	if rc.ServicRc == "true" {
		// rc to start a service for a particular hal copy to here
		return join(outVendorDir, "etc/init")
	}
	return "root"
}

func getInitRcCopyStatement(rcs []spec.RcScripts) string {
	var s []string
	for _, rc := range rcs {
		name := rc.File
		if name == "" {
			name = rc.Name
		}

		filepath.Base(name)
		t := join(copyLocal, name) + ":" + join(rcInstallDest(&rc), filepath.Base(name)) + ` \`
		s = append(s, t)
	}

	return strings.Join(s[:], "\n")
}

func join(dir string, file string) string {
	return filepath.Join(dir, file)
}

func getUeventdCopySrc(uc spec.UeventRc) string {
	if uc.File != "" {
		return join(copyLocal, uc.File)
	}

	src := "ueventd.rc.gen"
	if uc.Name != "" {
		src = uc.Name
	}
	return join(copyLocal, src)
}

func getFstabCopySrc(fs spec.Fstab) string {
	src := "fstab.rc.gen"
	if fs.Name != "" {
		src = fs.Name
	}
	return copyLocal + "/" + src
}

// RuntimeConfigInstructions turns the RuntimeConfig to a Android statement
func RuntimeConfigInstructions(config spec.RuntimeConfig) string {
	from := config.Src
	dstDir := defaultRuntimeConfigDst
	if config.DestDir != "" {
		dstDir = config.DestDir
	}
	dst := join(dstDir, filepath.Base(from))
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

// for kernel command, absolute path in target system
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
	return src + ":" + defaultFirmwareDst + "/" + filepath.Base(src)
}

// InstsallDriver is the instruction to install drivers on target
func InstsallDriver(src string) string {
	// TODO: generate the insmod instruction in the service.rc
	return src + ":" + defaultKernelModuleDst + "/" + filepath.Base(src)
}

func executeTemplate2(f *os.File, tmpName string, tmpContent string, spec *spec.Spec) (err error) {
	return nil
}

func executeTemplate(out *os.File, tmpName string, tmpContent string, spec *spec.Spec) (err error) {
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
		"InstsallDriver":            InstsallDriver,
	}

	tmpl, err := template.New(tmpName).Funcs(funcMap).Parse(string(tmpContent))
	if err != nil {
		fmt.Println("create template failed", tmpName, err)
		return err
	}

	return generate(tmpl, out, spec)
}

// executeTemplateForRc genereate Rcscript only
func executeTemplateForRc(f *os.File, rc *spec.RcScripts) (err error) {
	funcMap := template.FuncMap{}

	tmpl, err := template.New(tplInitRc).Funcs(funcMap).Parse(tmpl.Initrc)
	if err != nil {
		fmt.Println("rcScript template failed", tplInitRc, err)
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
