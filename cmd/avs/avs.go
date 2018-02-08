package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/pierrchen/avs/specconv"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "avs"
	app.Version = "0.0.1"
	app.Usage = "Andoid build specification cli"
	app.Action = func(c *cli.Context) error {
		fmt.Println("Please use avs subcommands, see avs -h")
		return nil
	}

	app.Commands = []cli.Command{
		{
			Name:    "init",
			Aliases: []string{"i"},
			Usage:   "init device config: avs s --vendor v --device d",
			Flags: []cli.Flag{
				cli.StringFlag{Name: "vendor", Value: "", Usage: "vendor name"},
				cli.StringFlag{Name: "device", Value: "", Usage: "device name"},
				cli.StringFlag{Name: "config", Value: "", Usage: "config file"},
			},
			Action: func(c *cli.Context) error {
				if c.String("vendor") == "" || c.String("device") == "" {
					log.Fatalln("must specify --vendor and --device for avs s subcmd")
				}

				if err := specconv.InitDeviceConfig(c.String("vendor"), c.String("device"), c.String("config")); err != nil {
					log.Fatalln("[avs s] Error generating the scafffolding", err)
				}

				fmt.Println("[avs s] OK")
				return nil
			},
		},
		{
			Name:    "validate",
			Aliases: []string{"v"},
			Usage:   "validate device config",
			Flags: []cli.Flag{
				cli.StringFlag{Name: "dir", Value: "", Usage: " dir for validation, default is current dir"},
			},
			Action: func(c *cli.Context) error {
				absGenDir := checkDir(c, true)
				err := specconv.ValdiateDeviceConfig(absGenDir)
				if err != nil {
					fmt.Println("[avs v] spec validation failed, please fix the errors!")
				} else {
					fmt.Println("[avs v] OK")
				}
				return err
			},
		},
		{
			Name:    "update",
			Aliases: []string{"u"},
			Usage:   "update the device config",
			Flags: []cli.Flag{
				cli.StringFlag{Name: "dir", Value: "", Usage: "dir for update, default is current dir"},
			},
			Action: func(c *cli.Context) error {
				absGenDir := checkDir(c, true)
				if err := specconv.UpdateDeviceConfigs(absGenDir); err != nil {
					log.Fatalln("[avs s] Error updating the config file", err)
				}
				fmt.Println("[avs u] OK")
				return nil
			},
		},
		{
			Name:    "clean",
			Aliases: []string{"c"},
			Usage:   "clean up all the geneated files",
			Flags: []cli.Flag{
				cli.StringFlag{Name: "dir", Value: "", Usage: "dir for cleanup, default is current dir"},
			},
			Action: func(c *cli.Context) error {
				absGenDir := checkDir(c, true)
				if err := specconv.CleanDeviceConfigs(absGenDir); err != nil {
					log.Fatalln("[avs s] Error updating the config file", err)
				}
				fmt.Println("[avs u] OK")
				return nil
			},
		},
	}

	app.Run(os.Args)

}

// return abs dir, otherwise, exit
func checkDir(c *cli.Context, check bool) string {
	dir := c.String("dir")
	if dir == "" {
		dir = filepath.Join(c.String("vendor"), c.String("device"))
	}
	absDir, err := getGenDir(dir)
	if err != nil {
		log.Fatalln("Error when getting the directory for generating the scaffold", err)
	}

	if check {
		if r, _ := dirExists(absDir); r == false {
			log.Fatalf("%s don't exsit\n", absDir)
		}
	}

	return absDir
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
