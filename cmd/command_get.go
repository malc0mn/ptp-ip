package main

import (
	"fmt"
	ptpfmt "github.com/malc0mn/ptp-ip/fmt"
	"github.com/malc0mn/ptp-ip/ip"
)

func init() {
	registerCommand(&get{})
}

type get struct{}

func (get) name() string {
	return "get"
}

func (get) alias() []string {
	return []string{}
}

func (get) execute(c *ip.Client, f []string) string {
	errorFmt := "get error: %s\n"

	cod, err := formatDeviceProperty(c, f[0])
	if err != nil {
		return fmt.Sprintf(errorFmt, err)
	}

	v, err := c.GetDevicePropertyValue(cod)
	if err != nil {
		return fmt.Sprintf(errorFmt, err)
	}

	return ptpfmt.DevicePropValAsString(c.ResponderVendor(), cod, int64(v)) + fmt.Sprintf(" (%#x)", v)
}

func (g get) help() string {
	help := `"` + g.name() + `" gets the current value for the given property.` + "\n"

	if args := g.arguments(); len(args) > 0 {
		help += helpAddArgumentsTitle()
		for i, arg := range args {
			switch i {
			case 0:
				help += "\t- " + arg + ": a hexadecimal field code in the form of '0x5001' or one of the supported unified field names:\n" + helpAddUnifiedFieldNames()
			}
		}
	}

	return help
}

func (get) arguments() []string {
	return []string{"property"}
}
