// Package vdts validates the device config files.
package vdts

import (
	"errors"
	"fmt"

	"github.com/pierrchen/avs/spec"
)

// IVal is the validator interface,
type IVal func(spec *spec.Spec, genDir string) error

// ValdiateSpec validates the spec, and it is the entry validator.
// It validator not only the spec but also the artifact geneated from the spec
// TODO: add validate stage1 and stage2, stage1 won't validate anything depend on
// the generation
func ValdiateSpec(spec *spec.Spec, absDeviceDir string) (err error) {
	if spec == nil {
		return errors.New("nil spec")
	}

	validateAll(spec, absDeviceDir, []IVal{
		validateBootImage,
		ValidateSystemImage})
	return nil
}

// validateAll is a helper function that will call all the validators pass in
func validateAll(spec *spec.Spec, genDir string, validators []IVal) error {
	var hasError bool
	for _, v := range validators {
		if err := v(spec, genDir); err != nil {
			fmt.Println("[avs v]", err)
			hasError = true
		}
	}
	if hasError {
		return fmt.Errorf("validation failed")
	}

	return nil
}
