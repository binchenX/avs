package images

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"os"
)

// check [1] for the default value
// [1] https://android.googlesource.com/platform/system/core/+/master/mkbootimg/mkbootimg
const (
	// These are the configurations we can config in the device config
	DefaultLoadBaseAddr  = 0x10000000
	DefaultKernelOffset  = 0x8000
	DefaultRamdiskOffset = 0x1000000

	// These are the value end up in the final boot image
	DefaultKernelLoadAddr  = DefaultLoadBaseAddr + DefaultKernelOffset
	DefaultRamdiskLoadAddr = DefaultLoadBaseAddr + DefaultRamdiskOffset
)

// Bootimg is to handle android bootimg
type Bootimg struct {
	// absolution path or the relative to current dir where the command is calling
	ImagePath string
}

// The format of the bootimg follows [1], and will be created using [2]
// [1] https://android.googlesource.com/platform/system/core/+/master/mkbootimg/include/bootimg/bootimg.h
// [2] https://android.googlesource.com/platform/system/core/+/master/mkbootimg/mkbootimg

// BootImgHdr is the Android Boot Image header
type BootImgHdr struct {
	Magic      [8]byte
	KernelSize uint32
	// KernelAddr is where the bootloader should load the kernel into (the memory)
	// Bootloader may either honor, ignore or rectify it.
	KernelAddr        uint32
	RamdiskSize       uint32
	RamdiskAddr       uint32
	SecondSize        uint32
	SecondAddr        uint32
	TagAddr           uint32
	PageSize          uint32
	BootImgHdrVersion uint32
	OsVersion         uint32
	ProductName       [16]byte
	CmdLine           [512]byte
	ID                [8]byte
	ExtraCmdline      [1024]byte
}

// ToString print hdr info
func (h *BootImgHdr) String() string {
	var s = ""
	/* os_version = ver << 11 | lvl */
	osVer := h.OsVersion >> 11
	osLvl := h.OsVersion & ((1 << 11) - 1)

	/* ver = A << 14 | B << 7 | C         (7 bits for each of A, B, C)
	 * lvl = ((Y - 2000) & 127) << 4 | M  (7 bits for Y, 4 bits for M) */
	s += fmt.Sprintf("OsVersion   :0x%x (Android Version: %d.%d.%d, Patch Level: %d.%d)\n",
		h.OsVersion,
		(osVer>>7)&0x7F, (osVer>>14)&0x7F, osVer&0x7F,
		(osLvl>>4)+2000, osLvl&0x0F)

	s += fmt.Sprintf("Product:%s\n", h.ProductName)
	s += fmt.Sprintf("CmdLine:%s\n", h.CmdLine)

	s += fmt.Sprintf("KernelAddr  :0x%x", h.KernelAddr)
	if h.KernelAddr == DefaultKernelLoadAddr {
		s += fmt.Sprintf("(default)")
	}
	s += "\n"

	s += fmt.Sprintf("Ramdisk     :0x%x", h.RamdiskAddr)
	if h.RamdiskAddr == DefaultRamdiskLoadAddr {
		s += fmt.Sprintf("(default)")
	}
	s += "\n"

	s += fmt.Sprintf("PageSize    :0x%x(%d)\n", h.PageSize, h.PageSize)
	s += fmt.Sprintf("KernelSize  :0x%x(%d)\n", h.KernelSize, h.KernelSize)
	s += fmt.Sprintf("RamdiskSize :0x%x(%d)\n", h.RamdiskSize, h.RamdiskSize)

	if h.SecondSize != 0 {
		s += fmt.Sprintf("SecondSize  :0x%x(%d)\n", h.SecondSize, h.SecondSize)
		s += fmt.Sprintf("SecondAddr  :0x%x\n", h.SecondAddr)
	}

	return s
}

// kernel, ramdisk, and 2nd are page size aligned
func align(size uint32, pageSize uint32) uint32 {
	return (size + pageSize - 1) / pageSize * pageSize
}

// Hdr return the Hdr of the boot image
func (b *Bootimg) Hdr() (*BootImgHdr, error) {
	f, err := os.Open(b.ImagePath)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	r := bufio.NewReader(f)
	var hdr BootImgHdr
	if err := binary.Read(r, binary.LittleEndian, &hdr); err != nil {
		fmt.Printf("Fail to read the bootimg %s, err %s", b.ImagePath, err)
		return nil, err
	}

	if string(hdr.Magic[:]) != "ANDROID!" {
		err := fmt.Errorf("Not An ANdroid Boot Image")
		return nil, err
	}

	return &hdr, nil
}

func extract(file *os.File, start int64, size uint32, outFile *os.File) error {
	_, err := file.Seek(int64(start), os.SEEK_SET)

	if err != nil {
		return fmt.Errorf("Err seek", err)
	}

	buf := make([]byte, size)
	n, err := file.Read(buf)
	if err != nil || n < int(size) {
		return fmt.Errorf("Err read", err)
	}

	n, err = outFile.Write(buf)
	if err != nil || n < int(size) {
		return fmt.Errorf("Err write", err)
	}

	return nil
}

// Unpack dump the bootimage header infor and extact all the stuff within (kernel, ramdisk, dtb)
func (b *Bootimg) Unpack() error {
	hdr, err := b.Hdr()
	f, err := os.Open(b.ImagePath)

	if err != nil {
		return err
	}

	type dump struct {
		In    *os.File
		Start uint32
		Size  uint32
		Out   *os.File
	}

	var dumps []dump

	if hdr.KernelSize != 0 {
		// kernel
		var d dump
		d.In = f
		d.Start = 1 * hdr.PageSize
		d.Size = hdr.KernelSize
		d.Out, _ = os.Create("kernel.out")
		dumps = append(dumps, d)
	}

	if hdr.RamdiskSize != 0 {
		// ramdisk
		var d dump
		d.In = f
		d.Start = 1*hdr.PageSize + align(hdr.KernelSize, hdr.PageSize)
		d.Size = hdr.RamdiskSize
		d.Out, _ = os.Create("ramdisk.out")
		dumps = append(dumps, d)
	}

	if hdr.SecondSize != 0 {
		// second
		var d dump
		d.In = f
		d.Start = 1*hdr.PageSize + align(hdr.KernelSize, hdr.PageSize) + align(hdr.RamdiskSize, hdr.PageSize)
		d.Size = hdr.SecondSize
		d.Out, _ = os.Create("second.out")
		dumps = append(dumps, d)
	}

	for _, d := range dumps {
		err = extract(d.In, int64(d.Start), d.Size, d.Out)

		if err != nil {
			fmt.Print("unpack %s failed, %s\n", err)
		} else {
			fmt.Printf("unpack %s OK\n", d.Out.Name())
		}
	}

	return err
}
