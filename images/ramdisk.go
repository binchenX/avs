// Package images manipulates the Android images.
package images

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/pierrchen/avs/utils"
)

// Ramdisk is Android ramdisk image
type Ramdisk struct {
	// absolution path or the relative to current dir where the command is calling
	ImagePath string
}

// Unpack a Andriod ramdisk to dir
// dir is either absolute dir or relative to current dir where the command is calling
func (r *Ramdisk) Unpack(dir string) error {
	img, err := filepath.Abs(r.ImagePath)
	if err != nil {
		return fmt.Errorf("can't find ramdisk image %s", img)
	}

	if r, _ := utils.FileExists(img); r == false {
		return fmt.Errorf("can't find ramdisk image %s", img)
	}

	outDir := dir

	// 1. copy ramdisk.img to ramdisk.img.gz in tmp dir
	tmpDir, err := ioutil.TempDir("", "avs")
	defer os.RemoveAll(tmpDir)
	tmpImage := filepath.Join(tmpDir, "ramdisk.img.gz")
	copyCmd := "cp " + img + " " + tmpImage
	cmd := exec.Command("bash", "-c", copyCmd)
	err = cmd.Run()
	if err != nil {
		return err
	}

	// 2. unzip the tmp .gz image and extract it from cpio into the dir specified
	cmdLine := "gzip -dc " + tmpImage + " | cpio -id"
	// fmt.Println(cmdLine)
	cur, _ := os.Getwd()

	if r, _ := utils.FileExists(outDir); r == true {
		log.Fatalf("%s exsit, remove that first\n", outDir)
	}

	ramdiskDir := filepath.Join(cur, outDir)
	utils.CreateGenDirIfNotExsit(ramdiskDir)

	fmt.Println("Press Enter to finish.")
	cmd = exec.Command("bash", "-c", cmdLine)
	cmd.Dir = ramdiskDir
	err = cmd.Run()
	if err != nil {
		return err
	}

	fmt.Printf("ramdisk %s is extracted to %s\n", img, ramdiskDir)
	return nil
}

// Pack pack dir to ramdisk image
func Pack(dir string, image Ramdisk) {
	// TODO:
}
