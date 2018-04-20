package vdts

import (
	"errors"
	"fmt"
	"reflect"
	"sort"

	"github.com/pierrchen/avs/spec"
	"github.com/pierrchen/avs/utils"
)

func validateBootImage(spec *spec.Spec, absDeviceDir string) (err error) {
	return validateAll(spec, absDeviceDir, []IVal{
		//validatKernelDTB,
		validateRootfs,
		validateParititions,
		validateMkBootImgArgs,
	})
}

func validatKernelDTB(spec *spec.Spec, absDeviceDir string) error {
	if spec.BootImage.Kernel == nil {
		return errors.New("must have KernelConfig")
	}

	if spec.BootImage.Kernel.CmdLine == "" {
		return errors.New("must have kernel command line")
	}

	// check if kernel and dtb exsit, just print warning
	if err := validateCopySrc(spec.BootImage.Kernel.LocalKernel, absDeviceDir); err != nil {
		fmt.Printf("[avs v] can't find kernel Image %s\n", spec.BootImage.Kernel.LocalKernel)
	}

	if err := validateCopySrc(spec.BootImage.Kernel.LocalDTB, absDeviceDir); err != nil {
		fmt.Printf("[avs v] can't find dtb Image %s\n", spec.BootImage.Kernel.LocalDTB)
	}

	return nil
}

func validateRootfs(spec *spec.Spec, absDeviceDir string) error {
	for _, rc := range spec.BootImage.Rootfs.InitRc {
		if rc.File != "" &&
			(rc.Name != "" || rc.Actions != nil || rc.Imports != nil || rc.Services != nil) {
			fmt.Printf("rc.File (%s) is not empty, all other other attribute will be ignored\n", rc.File)
		}

		rcName := rc.File
		if rcName == "" {
			rcName = rc.Name
		}
		if rcName == "" {
			fmt.Printf("[avs v] rc file don't have a name")
			break
		}
		if err := validateCopySrc("$(LOCAL_PATH)/"+rcName, absDeviceDir); err != nil {
			fmt.Printf("[avs v] can't find rc file %s\n", rcName)
		}
	}

	return nil
}

// - BoardConfig.PartitionTable.Partitions
// - BootImage.Rootfs.Fstab
// Must contain at least 3 partitions (system, userdata, cache) and they must be in sync (e.g type)
func validateParititions(spec *spec.Spec, absDeviceDir string) error {

	P := []string{"system", "userdata", "cache"}

	var parts []string
	for _, p := range spec.BoardConfig.PartitionTable.Partitions {
		parts = append(parts, p.Name)
	}

	if utils.IncludedIn(P, parts) != true {
		fmt.Printf("missing partitions table declaration, has only %v, need at least %v\n", parts, P)
	}

	var mounts []string
	for _, m := range spec.BootImage.Rootfs.Fstab.Mounts {

		// ignore the "auto", which are managed by volume managers
		// usually for usb and sdcard
		if m.Dst != "auto" {
			// remove the leading '/' in Dst
			p := m.Dst[1:]
			// same partition
			if p == "data" {
				p = "userdata"
			}
			mounts = append(mounts, p)
		}
	}

	if utils.IncludedIn(P, mounts) != true {
		fmt.Printf("missing partitions in fstab, has only %v, need at least %v\n", mounts, P)
	}

	// very partition in the partition table must have correspoinding entry in fstab
	sort.Strings(parts)
	sort.Strings(mounts)

	if !reflect.DeepEqual(parts, mounts) {
		fmt.Printf("partitoins %v and mounts %s doesn't match.\n", parts, mounts)
	}

	return nil
}

func validateMkBootImgArgs(spec *spec.Spec, absDeviceDir string) error {
	args := spec.BootImage.Args
	if args != nil && args.Lda != nil {
		lda := args.Lda
		if !((lda.LoadBase != "" && lda.KernelOffset != "" && lda.RamdiskOffset != "") ||
			(lda.LoadBase == "" && lda.KernelOffset == "" && lda.RamdiskOffset == " ")) {
			return fmt.Errorf("MkBootImageLoadArgsLoadAddress should either all default or has value")
		}
		return nil
	}

	return nil
}
