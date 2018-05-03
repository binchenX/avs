package spec

// this file includes all things related to partition, mount, fstab, file system

// valid partition names
const (
	SYSTEM string = "system"
	DATA   string = "userdata"
	VENDOR string = "vendor"
	CACHE  string = "cache"
)

// FsType is the file system type
type FsType string

// valid file systems
const (
	EXT4   FsType = "ext4"
	SQUASH FsType = "squashfs"
)

// valid partition schemes
const (
	MBR string = "mbr"
	GPT string = "gpt"
)

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
	Type      FsType `json:"type"`
	MntFlag   string `json:"mnt_flag"`
	FsMgrFlag string `json:"fs_mgr_flag"`
}
