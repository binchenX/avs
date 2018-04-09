package images

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"os"
	"os/exec"
)

const (
	dtbHeaderMagic = 0xd00dfeed
)

// Dtb is the dtb file
type Dtb struct {
	ImagePath string
}

// DtbHeader dtb header info
type DtbHeader struct {
	Magic           uint32
	TotalSize       uint32
	OffDtStruct     uint32
	OffDtStrings    uint32
	OffMemRsvmap    uint32
	Version         uint32
	LastCompVersion uint32
	BootCpuidPhys   uint32
	SizeDtStrings   uint32
	SizeDtStruct    uint32
}

// IsDtb is dtb or not
func (d *Dtb) IsDtb() bool {
	f, err := os.Open(d.ImagePath)
	defer f.Close()

	if err != nil {
		return false
	}

	t := make([]byte, 4)
	f.Read(t)

	magic := binary.BigEndian.Uint32(t)

	if magic == dtbHeaderMagic {
		return true
	}

	return false
}

// Hdr return the dtb header
func (d *Dtb) Hdr() *DtbHeader {
	f, err := os.Open(d.ImagePath)
	defer f.Close()
	if err != nil {
		return nil
	}

	var hdr DtbHeader
	r := bufio.NewReader(f)
	err = binary.Read(r, binary.BigEndian, &hdr)
	if err != nil {
		return nil
	}
	return &hdr
}

// ToDts dump to dts, two options
// 1. dtc -I dtb -O dts -o d.ImagePath.dts d.ImagePath
// 2. fdtdump d.ImagePath
// choose 1
func (d *Dtb) ToDts() error {
	// dtc -I dtb -O dts -o d.ImagePath.dts d.ImagePath
	cmdline := "dtc " + " -I dtb -O dts -o " + d.ImagePath + ".dts " + d.ImagePath
	cmd := exec.Command("bash", "-c", cmdline)
	err := cmd.Run()
	if err != nil {
		fmt.Println("dts dumped to ", d.ImagePath+".dts")
	}
	return err
}
