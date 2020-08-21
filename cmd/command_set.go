package main

import (
	"fmt"
	ptpfmt "github.com/malc0mn/ptp-ip/fmt"
	"github.com/malc0mn/ptp-ip/ip"
)

func init() {
	registerCommand(&set{})
}

type set struct{}

func (set) name() string {
	return "set"
}

func (set) alias() []string {
	return []string{}
}

func (set) execute(c *ip.Client, f []string) string {
	errorFmt := "set error: %s\n"

	cod, err := formatDeviceProperty(c, f[0])
	if err != nil {
		return fmt.Sprintf(errorFmt, err)
	}

	// TODO: add support for "string" values such as "astia" for film simulation.
	val, err := ptpfmt.HexStringToUint64(f[1], 32)
	if err != nil {
		return fmt.Sprintf(errorFmt, err)
	}
	c.Debugf("Converted value to: %#x", val)

	err = c.SetDeviceProperty(cod, uint32(val))
	if err != nil {
		return fmt.Sprintf(errorFmt, err)
	}

	return fmt.Sprintf("property %s successfully set to %#x\n", f[0], val)
}

func (s set) help() string {
	help := `"` + s.name() + `" sets the given value for the given property. Depending on the camera operation mode (aperture priority, shutter priority, manual or auto), not all properties might be settable!` + "\n"

	if args := s.arguments(); len(args) > 0 {
		help += helpAddArgumentsTitle()
		for i, arg := range args {
			switch i {
			case 0:
				help += "\t- " + arg + " is a hexadecimal field code in the form of '0x5001' or one of the supported unified field names:\n" + helpAddUnifiedFieldNames()
			case 1:
				help += "\t- " + arg + " is a hexadecimal value to set the field to. E.g. '0x6'\n"
			}
		}
	}

	return help
}

func (set) arguments() []string {
	return []string{"property", "value"}
}
