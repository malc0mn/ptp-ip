package main

import (
	"github.com/malc0mn/ptp-ip/ip"
)

func init() {
	registerCommand(&state{})
}

type state struct{}

func (state) name() string {
	return "state"
}

func (state) alias() []string {
	return []string{}
}

func (state) execute(c *ip.Client, f []string, _ chan<- string) string {
	res, err := c.GetDeviceState()

	if err != nil {
		res = err.Error()
	}

	return formatDeviceInfo(c.ResponderVendor(), res, f)
}

func (i state) help() string {
	help := `"` + i.name() + `" displays the current device state. This currently is a Fuji specific command!` + "\n"

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

func (state) arguments() []string {
	return []string{"json", "pretty"}
}
