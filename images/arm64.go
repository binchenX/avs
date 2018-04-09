package images

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"os"
)

const (
	linuxARM64ImageMagic         = "ARM\x64"
	KERNEL_IMAGE_STEXT_OFFSET    = 0x12C
	KERNEL_IMAGE_RAW_SIZE_OFFSET = 0x130
)

// Arm64ImageHeader is Image format for arm64 kernel Image, see [1]
// [1] https://www.kernel.org/doc/Documentation/arm64/booting.txt
type Arm64ImageHeader struct {
	Code0      uint32  /* Executable code */
	Code1      uint32  /* Executable code */
	TextOffset uint64  /* Image load offset, LE */
	ImageSize  uint64  /* Effective Image size, LE */
	Flags      uint64  /* Kernel flags, LE */
	Res2       uint64  /* reserved */
	Res3       uint64  /* reserved */
	Res4       uint64  /* reserved */
	Magic      [4]byte /* Magic number, LE, "ARM\x64" */
	Res5       uint32  /* reserved (used for PE COFF offset) */
}

func (hdr *Arm64ImageHeader) String() string {
	var s = ""
	s += fmt.Sprintf("TextOffset 0x%x\n", hdr.TextOffset)
	s += fmt.Sprintf("ImageSize 0x%x(%d)\n", hdr.ImageSize, hdr.ImageSize)
	s += fmt.Sprintf("Flag 0x%x\n", hdr.Flags)
	return s
}

// Arm64Image point to a path
type Arm64Image struct {
	ImagePath string
}

// Hdr return the head info of the arm64 Image, or err
func (i *Arm64Image) Hdr() (*Arm64ImageHeader, error) {
	var hdr Arm64ImageHeader
	f, err := os.Open(i.ImagePath)

	defer f.Close()

	if err != nil {
		return nil, err
	}

	r := bufio.NewReader(f)
	binary.Read(r, binary.LittleEndian, &hdr)

	if string(hdr.Magic[:]) != linuxARM64ImageMagic {
		fmt.Printf("Not valid arm64 image %s\n", hdr.Magic)
		return nil, fmt.Errorf("Not ARM64 Linux Image")
	}

	return &hdr, nil
}

// ActualKernelSize return the actual kernel size
func (i *Arm64Image) ActualKernelSize() (int64, error) {
	f, err := os.Open(i.ImagePath)
	defer f.Close()

	if err != nil {
		return 0, err
	}

	// kernel real-size, reset and seek to the off
	f.Seek(KERNEL_IMAGE_STEXT_OFFSET, os.SEEK_SET)
	t := make([]byte, 4)
	f.Read(t)
	stext := binary.LittleEndian.Uint32(t)

	f.Seek(KERNEL_IMAGE_RAW_SIZE_OFFSET, os.SEEK_SET)
	f.Read(t)
	rawSize := binary.LittleEndian.Uint32(t)

	return int64(stext + rawSize), nil
}

// IsSomethingAppended return true if file size > ActualKernelSize
func (i *Arm64Image) IsSomethingAppended() bool {

	info, _ := os.Stat(i.ImagePath)
	s, _ := i.ActualKernelSize()

	if info.Size() > s {
		return true
	}

	return false
}

// Split will split the kernel and the stuff appended, usually dtb
func (i *Arm64Image) Split() {

	if !i.IsSomethingAppended() {
		return
	}

	imageFile, _ := os.Open(i.ImagePath)
	info, _ := imageFile.Stat()
	s, _ := i.ActualKernelSize()
	kernel := fmt.Sprintf("%s.kernel", i.ImagePath)
	fkernel, _ := os.Create(kernel)
	extract(imageFile, 0, uint32(s), fkernel)

	dtb := fmt.Sprintf("%s.dtb", i.ImagePath)
	fdtb, _ := os.Create(dtb)
	extract(imageFile, s, uint32(info.Size()-s), fdtb)
}
