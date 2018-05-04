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
			Name:    "bootimg",
			Aliases: []string{"b"},
			Usage:   "parse and extract android bootimg",
			Flags: []cli.Flag{
				cli.StringFlag{Name: "image", Value: "", Usage: "boot image file"},
				cli.BoolFlag{Name: "extract", Usage: "extract bootimage"},
			},
			Action: func(c *cli.Context) error {
				b := images.Bootimg{ImagePath: c.String("image")}

				// dump the header
				v, err := b.Hdr()
				if err != nil {
					log.Fatalln(err)
				}

				fmt.Println(v)

				if c.Bool("extract") {
					b.Unpack()
				}
				return nil
			},
		},

		{
			Name:    "ramdisk",
			Aliases: []string{"r"},
			Usage:   "extract the ramdisk",
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

		{
			Name:    "kernel",
			Aliases: []string{"k"},
			Usage:   "dump kernel info and configs",
			Flags: []cli.Flag{
				cli.StringFlag{Name: "image", Value: "", Usage: "kernel image file"},
				cli.BoolFlag{Name: "configs", Usage: "enable kernel config dump"},
			},
			Action: func(c *cli.Context) error {
				r := images.Kernel{ImagePath: c.String("image")}
				v, err := r.Version()
				if err != nil {
					log.Fatalln(err)
				}

				fmt.Println(v)
				// check if it is arm64 image
				arm64 := images.Arm64Image{ImagePath: c.String("image")}
				hdr, err := arm64.Hdr()

				if err == nil {
					fmt.Println("Arm 64 Linux Kernel Image, info")
					fmt.Println(hdr)
				}

				// check if it has something appended
				if arm64.IsSomethingAppended() {
					fmt.Println("something is appended after the kernel")
					ks, _ := arm64.ActualKernelSize()
					fmt.Println("Actualy Kernel Size", ks)

					arm64.Split()

				} else {
					fmt.Println("Pure and simple kernel image")
				}

				// dump kernel build configs
				if c.Bool("configs") {
					configs, err := r.Configs()
					if err != nil {
						fmt.Println(configs)
					} else {
						fmt.Println(configs)
					}
				}
				return nil
			},
		},

		{
			Name:    "dtb",
			Aliases: []string{"d"},
			Usage:   "dump dtb header and decompile dtb",
			Flags: []cli.Flag{
				cli.StringFlag{Name: "image", Value: "", Usage: "dtb image file"},
				cli.BoolFlag{Name: "hdr", Usage: "dump hdr info"},
				cli.BoolFlag{Name: "dump", Usage: "dump dtb"},
			},
			Action: func(c *cli.Context) error {
				dtb := images.Dtb{ImagePath: c.String("image")}

				// dump the header
				if !dtb.IsDtb() {
					fmt.Println("Not a DTB file")
					return nil
				}

				fmt.Println("Is a DTB file")

				if c.Bool("hdr") {
					fmt.Printf("Header Info:\n %#v\n", dtb.Hdr())
				}

				if c.Bool("dump") {
					dtb.ToDts()
				}
				return nil
			},
		},
	}

	app.Run(os.Args)
}
