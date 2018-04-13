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
	// A/BUpdate Config
	// https://source.android.com/devices/tech/ota/ab_updates#build-variables
	// VerifiedBoot Config
	// https://android.googlesource.com/platform/external/avb/+show/master/README.md
	Target     *Target     `json:"target"`
	Bootloader *Bootloader `json:"bootloader,omitempty"`
	SELinux    *SELinux    `json:"selinux,omitempty"`
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

// Kernel includes the kernel command line, and where to look for the kernel image and DTB.
type Kernel struct {
	// Kernel command line need only include product specific stuff
	// Following configs will be automatically be added
	// androidboot.hardware=${spec.Product.Device}
	// androidboot.selinux=${spec.BoardConfig.SELinux.Mode}
	// see template.go/FullCmdLine
	CmdLine     string `json:"cmd_line"`
	LocalKernel string `json:"local_kernel"`
	// TODO:TARGET_COMPRESSED_KERNEL=true
	Compressed string `json:"compressed,omitempty"`
	LocalDTB   string `json:"local_dtb,omitempty"`
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

// VendorRaw are the raw instructions that will be copied directly to the device.mk.
// It is a fallback for things that can't be expressed nicely in current specification.
// Use it *rarely*.
type VendorRaw struct {
	Instructions []string `json:"instructions"`
}
