package main

import (
	"fmt"
	"github.com/malc0mn/ptp-ip/ip"
)

func init() {
	registerCommand(&describe{})
}

type describe struct{}

func (describe) name() string {
	return "describe"
}

func (describe) alias() []string {
	return []string{}
}

func (describe) execute(c *ip.Client, f []string) string {
	errorFmt := "describe error: %s\n"

	cod, err := formatDeviceProperty(c, f[0])
	if err != nil {
		return fmt.Sprintf(errorFmt, err)
	}

	res, err := c.GetDevicePropertyDescription(cod)
	if err != nil {
		return fmt.Sprintf(errorFmt, err)
	}

	if res == nil {
		return fmt.Sprintf(errorFmt, fmt.Sprintf("cannot describe property %#x", cod))
	}

	return fujiFormatDeviceProperty(res, f[1:])
}

func (d describe) help() string {
	help := `"` + d.name() + `" describes the given property.` + "\n"

	if args := d.arguments(); len(args) > 0 {
		help += helpAddArgumentsTitle()
		for i, arg := range args {
			switch i {
			case 0:
				help += "\t- " + arg + ": a hexadecimal field code in the form of '0x5005' or one of the supported unified field names:\n" + helpAddUnifiedFieldNames()
			case 1:
				help += "\t- " + `"` + arg + `" to output the data in parsable json format` + "\n"
			case 2:
				help += "\t- " + `"` + arg + `" to be used together with "` + args[1] + `": format the output in a human readable way` + "\n"
			}
		}
	}

	return help
}

func (describe) arguments() []string {
	return []string{"property", "json", "pretty"}
}
