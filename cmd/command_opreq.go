package main

import (
	"encoding/hex"
	"fmt"
	ptpfmt "github.com/malc0mn/ptp-ip/fmt"
	"github.com/malc0mn/ptp-ip/ip"
	"github.com/malc0mn/ptp-ip/ptp"
)

func init() {
	registerCommand(&opreq{})
}

type opreq struct{}

func (opreq) name() string {
	return "opreq"
}

func (opreq) alias() []string {
	return []string{}
}

func (opreq) execute(c *ip.Client, f []string) string {
	var res string
	errorFmt := "opreq error: %s\n"

	cod, err := ptpfmt.HexStringToUint64(f[0], 16)
	if err != nil {
		return fmt.Sprintf(errorFmt, err)
	}
	c.Debugf("Converted uint16: %#x", cod)

	params := f[1:]
	p := make([]uint32, len(params))
	for i, param := range params {
		conv, err := ptpfmt.HexStringToUint64(param, 64)
		if err != nil {
			return fmt.Sprintf(errorFmt, err)
		}
		p[i] = uint32(conv)
	}

	c.Debugf("Converted params: %#x", p)

	d, err := c.OperationRequestRaw(ptp.OperationCode(cod), p)
	if err != nil {
		return fmt.Sprintf(errorFmt, err)
	}

	for _, raw := range d {
		res += fmt.Sprintf("\nReceived %d bytes. HEX dump:\n%s", len(raw), hex.Dump(raw))
	}

	return res
}

func (o opreq) help() string {
	help := `"` + o.name() + `" This command is intended for reverse engineering and/or debugging purposes. The output will always be a hexadecimal dump of the packets received from the responder.` + "\n"

	if args := o.arguments(); len(args) > 0 {
		help += helpAddArgumentsTitle()
		for i, arg := range args {
			switch i {
			case 0:
				help += "\t- " + arg + ": a hexadecimal operation code in the form of '0x1014'. The supported operation codes will vary from vendor to vendor.\n"
			case 1:
				help += "\t- " + arg + ": depending on the operation code, an additional parameter might be required. It is expected to be in hexadecimal form, e.g. '0x5003'\n"
			}
		}
	}

	return help
}

func (opreq) arguments() []string {
	return []string{"opcode", "param"}
}
