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

// Configs return the configs file of the Kernel
func (k *Kernel) Configs() string {
	// https://github.com/torvalds/linux/blob/master/scripts/extract-ikconfig
	return ""
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
