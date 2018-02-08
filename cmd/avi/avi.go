package main

import (
	"fmt"
	"log"
	"os"

	"github.com/pierrchen/avs/images"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "avi"
	app.Version = "0.0.1"
	app.Usage = "Andoid Image Tool"
	app.Action = func(c *cli.Context) error {
		fmt.Println("Please use avi subcommands, see avi -h")
		return nil
	}

	app.Commands = []cli.Command{
		{
			Name:    "ramdisk",
			Aliases: []string{"r"},
			Usage:   "regerate the device config",
			Flags: []cli.Flag{
				cli.StringFlag{Name: "image", Value: "", Usage: "ramdisk.img image for unpack"},
				cli.StringFlag{Name: "extract", Value: "avs_ramdisk", Usage: "extracted out dir"},
			},
			Action: func(c *cli.Context) error {
				r := images.Ramdisk{ImagePath: c.String("image")}
				err := r.Unpack(c.String("extract"))
				if err != nil {
					log.Fatalln(err)
				}
				return nil
			},
		},
	}

	app.Run(os.Args)
}
