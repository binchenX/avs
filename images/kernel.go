package images

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/pierrchen/avs/utils"
)

// Kernel image
type Kernel struct {
	// absolution path or the relative to current dir where the command is calling
	ImagePath string
}

func commandInstalled(cmd string) bool {
	cmdLine := "which " + cmd
	c := exec.Command("bash", "-c", cmdLine)
	o, _ := c.Output()
	if string(o) == "" {
		return false
	}
	return true
}

// Configs return the configs file of the Kernel
// https://github.com/torvalds/linux/blob/master/scripts/extract-ikconfig
// It assume extract-ikconfig is inside of the $PATH
func (k *Kernel) Configs() (string, error) {
	img, err := filepath.Abs(k.ImagePath)
	if err != nil {
		return "", fmt.Errorf("can't find image %s", img)
	}

	if r, _ := utils.FileExists(img); r == false {
		return "", fmt.Errorf("can't find image %s", img)
	}

	extractTool := "extract-ikconfig"

	if !commandInstalled(extractTool) {
		return "", fmt.Errorf("command %s ins't installed, get it from %s and install it to $PATH",
			extractTool, "https://github.com/torvalds/linux/blob/master/scripts/extract-ikconfig")
	}

	cmdLine := "extract-ikconfig " + img
	cmd := exec.Command("bash", "-c", cmdLine)

	config, err := cmd.Output()
	return string(config), nil
}

// Version returns the version of the Kernel
func (k *Kernel) Version() (string, error) {

	img, err := filepath.Abs(k.ImagePath)
	if err != nil {
		return "", fmt.Errorf("can't find image %s", img)
	}

	if r, _ := utils.FileExists(img); r == false {
		return "", fmt.Errorf("can't find image %s", img)
	}

	cmdLine := "strings " + img + ` | grep "Linux version" `
	cmd := exec.Command("bash", "-c", cmdLine)

	version, err := cmd.Output()

	if err != nil {
		return "", err
	}
	// the 2nd line is garbage, contains Linux version %s (%s)
	v := strings.Split(string(version), "\n")[0]
	return v, nil
}
