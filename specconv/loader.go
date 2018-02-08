package specconv

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/pierrchen/avs/spec"
)

// SaveSpecToJSON save spec to path as json format
// path - absolution path to spec file
func SaveSpecToJSON(jsonData interface{}, path string) error {
	data, err := json.MarshalIndent(jsonData, "", "    ")
	if err != nil {
		fmt.Println("wrong spec")
		return err
	}

	if filepath.IsAbs(path) == false {
		fmt.Printf("%s isn't absolute path\n", path)
		return errors.New("spec path should avs")
	}

	if err := ioutil.WriteFile(path, data, 0777); err != nil {
		fmt.Println("fail to write to output file", err)
		return err
	}

	return nil
}

// LoadSpecFromString return a spec from jsonSting, return error on error
func LoadSpecFromString(jsonString string) (spec *spec.Spec, err error) {
	if err = json.NewDecoder(strings.NewReader(jsonString)).Decode(&spec); err != nil {
		fmt.Printf("%#v", err)
		return nil, err
	}
	return spec, nil
}

// LoadSpec load config.json file and return an Spec object
func LoadSpec(configFile string) (spec *spec.Spec, err error) {
	cf, err := os.Open(configFile)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("JSON specification file %s not found", configFile)
		}
		return nil, err
	}
	defer cf.Close()

	if err = json.NewDecoder(cf).Decode(&spec); err != nil {
		fmt.Printf("%#v", err)
		return nil, err
	}
	return spec, nil
}

// LoadHalSpec load a HAL spec from json file
func LoadHalSpec(configFile string) (spec *spec.HAL, err error) {
	cf, err := os.Open(configFile)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("JSON specification file %s not found", configFile)
		}
		return nil, err
	}
	defer cf.Close()

	if err = json.NewDecoder(cf).Decode(&spec); err != nil {
		fmt.Printf("%#v", err)
		return nil, err
	}
	return spec, nil
}

// GetAvsInstallDir return the avs instsall dir
func GetAvsInstallDir() string {
	return filepath.Join(os.Getenv("GOPATH"), "src", avsPackage)
}

func getDefaultConfigFile() string {
	return "configs/default.json"
}

// TODO:be able to select diffent type of scaffolding, say arm64, qemu, full
func loadTemplateSpec() (*spec.Spec, error) {
	return LoadSpec(filepath.Join(GetAvsInstallDir(), getDefaultConfigFile()))
}
