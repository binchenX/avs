package spec

// HAL is all the HAL related configrations (other than the HAL code itself) for this device.
type HAL struct {
	// Name is the name of this HAL.
	Name string `json:"name"`
	// Manifests is the manifest required by Treble (Android O).
	Manifests []Manifest `json:"manifests,omitempty"`
	// Features are features supported by this HAL, it will be copied to device.
	// All the features files must exsits in frameworks/native/data/etc/
	Features []string `json:"features,omitempty"`
	// BuildConfigs are the build configrations related with this features. It can either be a build
	// configration for framework (e.g Wifi has generic framework support, and that configrations should
	// be put in Wifi HAL.BuildConfigs), or for the particular HAL implementation (say your gralloc
	// implementation has some build configrations for different chips, you can put it Graphic HAL.BuildConfigs,
	// but maybe better to handle that in gralloc/Android.mk so people don't need to care).
	// Those BuildConfigs configurations will end up in the geneated BoardConfig.mk, logically those
	// configurations are part of the HAL configration, they are put into the BoardConfig.mk just to
	// optimize the compiling process.
	BuildConfigs []string `json:"build_configs,omitempty"`
	// Packages are libaries and bins required.
	// It can be empty, since the required package maybe included in the base product config,
	// as specified in the Spec.Product.
	Packages *Packages `json:"required_packages,omitempty"`
	// RawInstructions are the instructions that aren't modeled at the moment and will be copied
	// directly to the device.mk.
	RawInstructions []string `json:"raw,omitempty"`
	// InitRc is the rc script need by the init to start the HAL service.
	InitRc []RcScripts `json:"init.rc,omitempty"`
	// RuntimeConfigs are config files needed in runtime, it will be copied to device.
	RuntimeConfigs []RuntimeConfig `json:"runtime_configs,omitempty"`
	// Properties are properties for this HAL.
	Properties []string `json:"properties,omitempty"`
	// Device nodes needed for this features. It will be aggreated to the
	// BootImage.Rootfs.UEventrc file.
	// Note that HAL are suppose to use standard device node, other than vendor
	// specific node to a feature, say use /dev/graphic/fb instead of /dev/myFb.
	// In another word, if all HAL use standard device nodes, there would be no
	// need for such a field in HAL.
	UeventRules []UeventRule `json:"uevent_rules,omitempty"`
	// SEPolicy is the SEPoplicy required for this HAL
	// There are several cases you will need this:
	// 1. This HAL implementation needs to access a *non-standard* device node ( a
	//	UeventRules will be needed in this case as well) or a sysfs file.
	// 2. This HAL implementation has a daemon.
	// see spec_example.go for examples
	SEPolicy *SEPolicyF `json:"sepolicy,omitempty"`
}

// Packages to build or copy.
type Packages struct {
	// Build are packages that will be build. Each component will take care of its installation.
	// To indicate it is a vendor package to enable additional verification (during `avs v`),
	// a tag of ":v" can be added at the end of the package name, e.g hwcomposer.poplar:v
	Build []string `json:"build,omitempty"`
	// Copy are packages that will be copied
	Copy []CopyPackage `json:"copy,omitempty"`
}

// CopyPackage is the libary and executable that will be copied directly.
// The sole purpose for its existence is to deal with vendor propriety binary
// that has no source code.
//	 Src: the package to be copied, must be a file (not directory)
// 	 the path is either relative to $(ANDROID_BUILD_TOP) or the device directory.
//   The first is for the cases when vendor want to release the binary as a separate package
//   which bundles a click through license, e.g Mali gles.
//   The second is for the cases when pacakges reside in the device directory but with a EULA,
//   e.g bt firmware.
// DestDir: The copy destination directory. It is optional, and when no dir specified,
// following rules apply:
// 1. all .so will be copied to $(vendor)/lib;
// 2. all executable will be copied to $(vendor)/bin
// If specified, the DestDir is *relative* to the $(vendor_out) eg. hw, /lib/egl
// and with no ending "/".
// The exactly dest dir for $(vendor) depend on the partition configrations, if no dedicated vendor
// partition, it will copy to $OUT/system/vendor; if yes, it will copy to $OUT/vendor. This is taking
// care of by avs automatically.
// 	Tag: reserved.
// The the copy command will be equivlenat to
// cp Src DestDir/basename(Src)
type CopyPackage struct {
	Src     string `json:"src"`
	DestDir string `json:"destDir,omitempty"`
	Tag     string `json:"tag,omitempty"`
}

// RuntimeConfig is the RuntimeConfig File that will be installed on the device.
type RuntimeConfig struct {
	Src string `json:"src"`
	// Default value is "system/etc/Basename(.Src)"
	DestDir string `json:"destDir,omitempty"`
}

// Property is the property setting.
type Property struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}
