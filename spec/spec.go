// Package spec includes the Android build configration specfication.
package spec

// Spec is the specification for Android device configration.
// Attributes without "omitempty" are required, otherwise it is schema error.
type Spec struct {
	Version          *Version          `json:"version"`
	Product          *Product          `json:"product"`
	BoardConfig      *BoardConfig      `json:"boardConfig"`
	BootImage        *BootImage        `json:"boot_image"`
	FrameworkConfigs *FrameworkConfigs `json:"framework_configs,omitempty"`
	Hals             []HAL             `json:"hals"`
	VendorRaw        *VendorRaw        `json:"vendor_raw,omitempty"`
}

// Version describe the avs spec version, as well as the Android version this spec applies for.
type Version struct {
	Schema  string `json:"schema"`
	Android string `json:"android"`
}

// Product describe the product information and which base products it "inherits".
type Product struct {
	Name        string `json:"name"`
	Device      string `json:"device"`
	Brand       string `json:"brand"`
	Model       string `json:"model"`
	Manufacture string `json:"manufacture"`
	// The base product configrations.
	// In theory you can put any configration here you want to inherit, but most products
	// will use the definition from here [1] (e.g "full_base.mk", "core.mk".) and it is the path
	// we currently looking for.
	// [1] https://android.googlesource.com/platform/build/+/master/target/product/
	InheritProducts []string `json:"inherit_products,omitempty"`
}

// BoardConfig is configrations for the Board. All the build time configurations belongs to a
// particular HAL should be put into HAL.BuildConfigs.
type BoardConfig struct {
	PartitionTable PartitionTable `json:"partition_table"`
	Target         *Target        `json:"target"`
	Bootloader     *Bootloader    `json:"bootloader,omitempty"`
	SEPolicy       *SEPolicy      `json:"sepolicy,omitempty"`
	// BoardFeatures are the features that require no HAL support. Such as software features
	// implemented by frameworks (e.g android.software.webview), or features directly supported by
	// kernel (e.g android.hardware.usb.host.xml).
	// All the valid features are list here [1].
	// [1] https://android.googlesource.com/platform/frameworks/native/+/master/data/etc
	BoardFeatures []string `json:"board_features,omitempty"`
}

// PartitionTable is the partion table for the device.
type PartitionTable struct {
	FlashBockSize string `json:"flash_block_size"`
	// mbr/ebr, gpt, others
	Scheme     string      `json:"scheme"`
	Partitions []Partition `json:"partitions"`
}

// Partition is the configration for each partition.
type Partition struct {
	Name string `json:"name"`
	Type string `json:"type"`
	Size string `json:"size"`
}

// Target is the static build configuration that impacts how the host machine
// will built the images.
type Target struct {
	Archs      []Arch `json:"archs"`
	NoRecovery bool   `json:"no_recovery"`
	NoRadio    bool   `json:"no_radioimage"`
	Binder     string `json:"binder"`
	// Default to Product.Name
	BoardPlatform string `json:"boardPlatform,omitempty"`
}

// Arch configration.
type Arch struct {
	Name    string `json:"name"`
	Variant string `json:"variant"`
	CPU     *CPU   `json:"cpu"`
}

// CPU configuration.
type CPU struct {
	Variant string `json:"variant"`
	Abi     string `json:"abi"`
	Abi2    string `json:"abi2,omitempty"`
}

// Bootloader is the configrations for bootloader.
type Bootloader struct {
	Has2ndBootloader bool `json:"has_2nd_bootloader,omitempty"`
	// Default to Product Name
	BoardName string `json:"board_name,omitempty"`
}

// BootImage contains all the configrations that will compromise the final boot image.
type BootImage struct {
	Args   *MkBootImageArgs `json:"args,omitempty"`
	Kernel *Kernel          `json:"kernel"`
	Rootfs *RootfsOverlay   `json:"rootfs_overlay"`
}

// MkBootImageArgs is the args passing to mkbootimage scripts when creating the boot image.
type MkBootImageArgs struct {
	RamdiskOffset string `json:"ramkdis_offset,omitempty"`
	KernelOffset  string `json:"kernel_offset,omitempty"`
}

// SEPolicy is the sepolicy config.
type SEPolicy struct {
	Dir string `json:"dir"`
}

// Kernel includes the kernel command line, and where to look for the kernel image and DTB.
type Kernel struct {
	CmdLine     string `json:"cmd_line"`
	LocalKernel string `json:"local_kernel"`
	LocalDTB    string `json:"local_dtb,omitempty"`
}

// RootfsOverlay includes the files that will be included in the rootfs (part of the boot image).
type RootfsOverlay struct {
	// only 1 fstab is allowed in system. The destination must be that is
	// fstab.$(ro.hardware)
	Fstab *Fstab `json:"fstab"`
	// we can have multipy init.rc files but the entry is always init.$(ro.hardware).rc,
	// which will import other init.rc
	InitRc []RcScripts `json:"init.rc"`
	// only 1 ueventRc is allowed in the system. The destination file must be
	// ueventd.$(ro.hardware).rc
	UeventRc *UeventRc `json:"uevent.rc"`
}

// Fstab is the mount configration.
type Fstab struct {
	// Name will be the file name for the generated script and it will be copied to device
	// with same name. Default name ueventd.rc.gen
	Name   string  `json:"name,omitempty"`
	Mounts []Mount `json:"mounts"`
}

// RcScripts is a script follow Android init syntax and semantics, see[1].
// [1]https://android.googlesource.com/platform/system/core/+/master/init/README.md
type RcScripts struct {
	// If Files is no nil, we will cp it directly $(LOCAL_PATH)/File to destination
	// All the other attribution will be ignored
	File string `json:"file,omitempty"`
	// Name will be the file name for the generated script and it will be copied to device
	// with same name. At least one of the File, or Name should be non empty.
	Name     string      `json:"name,omitempty"`
	Imports  []string    `json:"imports,omitempty"`
	Actions  []RcAction  `json:"actions,omitempty"`
	Services []RcService `json:"services,omitempty"`
}

// RcAction is the Action statement.
type RcAction struct {
	Triggers string   `json:"triggers"`
	Commands []string `json:"commands"`
}

// RcService is the Service statement.
type RcService struct {
	Name    string   `json:"name"`
	Path    string   `json:"path"`
	Args    string   `json:"args,omitempty"`
	Options []string `json:"options,omitempty"`
}

// RcImport is the Import statement.
type RcImport struct {
	ImportPath string `json:"path"`
}

// UeventRc is the rules for eventd.
type UeventRc struct {
	// If Files is no nil, we will cp it directly $(LOCAL_PATH)/File to destination
	// All the other attribution will be ignored
	File string `json:"file,omitempty"`
	// Name will be the file name for the generated script and it will be copied to device
	// with same name. Default name ueventd.rc.gen
	Name  string       `json:"name,omitempty"`
	Rules []UeventRule `json:"rules"`
}

// UeventRule is rule for eventd.
type UeventRule struct {
	Node string `json:"node"`
	Attr string `json:"attr,omitempty"`
	Mode string `json:"mode"`
	UID  string `json:"uid"`
	GUID string `json:"guid"`
}

// Mount is mount intruction.
type Mount struct {
	Src       string `json:"src"`
	Dst       string `json:"dst"`
	Type      string `json:"type"`
	MntFlag   string `json:"mnt_flag"`
	FsMgrFlag string `json:"fs_mgr_flag"`
}

// FrameworkConfigs includes the build time and runtime configurations for frameworks
// components(e.g dalvik).
type FrameworkConfigs struct {
	// goes to BoardConfig.mk
	BuildConfigs []string `json:"build_configs,omitempty"`
	// goes to device.mk
	Properties []string `json:"properties,omitempty"`
}

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

// SEPolicyF is the sepolicy configration.
type SEPolicyF struct {
	// create new file type
	FileTe []string `json:"file.te,omitempty"`
	// create new process domain
	ServiceTe []string `json:"service.te,omitempty"`
	// lable the files
	FileContexts []string `json:"file_contexts,omitempty"`
	// lable the proccess/services
	ServiceContexts []string `json:"service_contexts,omitempty"`
}

// Manifest is the manifest for the interface (Treble only).
type Manifest struct {
	Name      string           `json:"name"`
	Format    string           `json:"format"`
	Transport *Transport       `json:"transport"`
	Impl      *Impl            `json:"impl,omitempty"`
	Version   string           `json:"version"`
	Interface *ServiceInterace `json:"interface"`
}

// Transport is the transport type could be hwbinder, passthrough.
type Transport struct {
	Arch string `json:"arch,,omitempty"`
	Mode string `json:"mode"`
}

// Impl is the implementation.
type Impl struct {
	Level string `json:"level"`
}

// ServiceInterace is the interface this HAL model implemented.
type ServiceInterace struct {
	Name     string `json:"name"`
	Instance string `json:"instance"`
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

// CopyPackage is the package that will be copied directly.
type CopyPackage struct {
	Src     string `json:"src"`
	DestDir string `json:"destDir"`
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

// VendorRaw are the raw instructions that will be copied directly to the device.mk.
// It is a fallback for things that can't be expressed nicely in current specification.
// Use it *rarely*.
type VendorRaw struct {
	Instructions []string `json:"instructions"`
}
