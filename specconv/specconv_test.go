package specconv

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/pierrchen/avs/spec"
	"github.com/pierrchen/avs/vdts"
	"github.com/stretchr/testify/assert"
)

func example() *spec.Spec {
	return &spec.Spec{
		Version: &spec.Version{
			Schema:  "0.1",
			Android: "Android o",
		},

		Product: &spec.Product{
			Name:        "poplar",
			Device:      "Poplar",
			Brand:       "Poplar",
			Model:       "Poplar",
			Manufacture: "linaro",
		},
	}
}

const jsonStream = `{
	"version": {
			"schema": "0.1",
			"android": "Android o"
	},
	"product": {
			"name": "poplar",
			"device": "Poplar",
			"brand": "Poplar",
			"model": "Poplar",
			"manufacture": "linaro"
	},
	"boardConfig": null,
	"boot_image": null,
	"hals": null
}`

func TestSpecToJson(t *testing.T) {
	spec := example()
	data, err := json.MarshalIndent(spec, "", "\t")

	if err != nil {
		fmt.Println("wrong spec")
	}

	f, _ := ioutil.TempFile("", "avs")

	fmt.Printf("create file %s\n", f.Name())
	defer os.Remove(f.Name())

	if err := ioutil.WriteFile(f.Name(), data, 0666); err != nil {
		fmt.Println("fail to write to output file")
	}
}

func TestJsonToSpec(t *testing.T) {

	spec, err := LoadSpecFromString(jsonStream)
	if err != nil {
		fmt.Println("malformed config file")
	}

	assert.Equal(t, spec, example(), "should equal")
}

func TestPoplarJsonToSpec(t *testing.T) {
	origTestFile := "../testFixtures/config.json"
	recoverFile := ".config_recover.json"
	spec, err := LoadSpec(origTestFile)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// convert back to json
	data, err := json.MarshalIndent(spec, "", "    ")
	assert.Nil(t, err, "wrong spec")

	if err := ioutil.WriteFile(recoverFile, data, 0666); err != nil {
		fmt.Println("fail to write to output file")
		os.Exit(1)
	}

	od, err := ioutil.ReadFile(origTestFile)
	assert.Nil(t, err)

	rd, err := ioutil.ReadFile(recoverFile)
	assert.Nil(t, err)

	assert.Equal(t, od, rd, "conversion failed")
}

func TestEmptyKernelCmdLine(t *testing.T) {
	// const jsonStream = `
	// {
	// 	"version": {
	// 		"schema": "0.1",
	// 		"android": "Android O"
	// 	}
	// }
	// `
	// spec, err := specconv.LoadSpecFromString(jsonStream)
	// if err != nil {
	// 	fmt.Println(err)
	// 	os.Exit(1)
	// }

	// err = vdts.ValdiateSpec(spec)
	// assert.NotNil(t, err, "should show error message")
}

func TestCreateFile(t *testing.T) {

	err := ioutil.WriteFile("/home/binchen/go/src/github.com/pierrchen/avs/dont/config.json", []byte("foo"), 0644)
	assert.NotNil(t, err, "create should fail since there is no dir called (dont)")
}

func TestGetGenDir(t *testing.T) {
	curDir, _ := os.Getwd()
	dir, _ := getGenDir("")
	assert.Equal(t, dir, filepath.Join(curDir, "."), "specify no dir, the default dir is current dir")
	dir, _ = getGenDir("test")
	assert.Equal(t, dir, filepath.Join(curDir, "test"), "specify relative path should concat it to $PWD")
	dir, _ = getGenDir(filepath.Join(curDir, "gen"))
	assert.Equal(t, dir, filepath.Join(curDir, "gen"), "use absolute path if that is what given")
}

func TestCheckGenDir(t *testing.T) {
	assert.Equal(t, false, checkGenDir(".xxxx"), "there shouldn't be `.xxxx` dir")
}

func TestDirExsit(t *testing.T) {

	// curDir, _ := os.Getwd()
	// r, _ := dirExists(filepath.Join(curDir, "gen"))
	// assert.Equal(t, true, r, "xxx")

	// r, _ = dirExists(filepath.Join(curDir, "xxxx"))
	// assert.Equal(t, true, r, "xxx")

	// r, _ = dirExists(filepath.Join(curDir, "yyyy"))
	// assert.Equal(t, false, r, "yyyy")
}

func TestValdiateHAL(t *testing.T) {
	const jsonStream = `
	{
	"hals": [
		{
			"name": "wifi",
			"build_configs": [
				"WPA_SUPPLICANT_VERSION := VER_0_8_X",
				"BOARD_WPA_SUPPLICANT_DRIVER := NL80211",
				"BOARD_WPA_SUPPLICANT_PRIVATE_LIB := lib_driver_cmd_bcmdhd",
				"BOARD_HOSTAPD_DRIVER := NL80211",
				"BOARD_HOSTAPD_PRIVATE_LIB := lib_driver_cmd_bcmdhd",
				"BOARD_WLAN_DEVICE := bcmdhd"
			],
			"features": [
				"android.hardware.wifi"
			],
			"runtime_configs": [
				{
					"src": "wpa_supplicant.conf",
					"destDir": "system/etc/wifi/"
				},
				{
					"src": "wpa_supplicant.rc",
					"destDir": "system/etc/init/"
				}
			],
			"packages": {
				"build": [
					"android.hardware.wifi@1.0-service:f",
					"wificond:f",
					"wificond.rc:f",
					"libwpa_client:f",
					"wpa_supplicant:f",
					"hostapd:f",
					"wpa_cli:f",
					"libwifi-hal:f",
					"wifi_hal:v"
				]
			}
		}
	]
	}
	`
	spec, err := LoadSpecFromString(jsonStream)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	vdts.GetVendorPackges(spec)

}

func TestHasModulesInFile(t *testing.T) {
	//dir, _ := filepath.Abs("gen/wifi/Android.mk")
	//r, _ := vdts.HasModulesInFile(dir, "wifi_hal")
	//assert.True(t, r, "should have `wifi_hal` modle in gen/wifi/Andriod.mk ")
}

// return the gen dir absoluted path
func getGenDir(para string) (string, error) {
	genDir := para
	if para == "" {
		genDir = "."
	}

	absGenDir, err := filepath.Abs(genDir)
	if err != nil {
		return "", err
	}
	return absGenDir, nil
}

// check if the Gen dir as specified in the para exsit or not
func checkGenDir(para string) bool {

	absPath, err := getGenDir(para)
	if err != nil {
		return false
	}

	if r, _ := dirExists(absPath); r == false {
		return false
	}
	return true
}

// for update, update all the *.mk in top dir, don't touch any hal dir
func createOrUpdateConfig() error {
	return nil
}

// dirExists check if file exsist
func dirExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}
