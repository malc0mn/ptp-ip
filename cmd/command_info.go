package main

import (
	"github.com/malc0mn/ptp-ip/ip"
)

func init() {
	registerCommand(&info{})
}

type info struct{}

func (info) name() string {
	return "info"
}

func (info) alias() []string {
	return []string{}
}

func (info) execute(c *ip.Client, f []string) string {
	res, err := c.GetDeviceInfo()

	if err != nil {
		res = err.Error()
	}

	return formatDeviceInfo(c.ResponderVendor(), res, f)
}

func (i info) help() string {
	help := `"` + i.name() + `" displays the device info. The data returned can vary from vendor to vendor.` + "\n"

	if args := i.arguments(); len(args) > 0 {
		help += helpAddArgumentsTitle()
		for i, arg := range args {
			switch i {
			case 0:
				help += "\t- " + `"` + arg + `" to output the data in parsable json format` + "\n"
			case 1:
				help += "\t- " + `"` + arg + `" to be used together with "` + args[0] + `": format the output in a human readable way` + "\n"
			}
		}
	}

	return help
}

func (info) arguments() []string {
	return []string{"json", "pretty"}
}
