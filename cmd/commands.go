package main

import (
	"bufio"
	"encoding/hex"
	"fmt"
	ptpfmt "github.com/malc0mn/ptp-ip/fmt"
	"github.com/malc0mn/ptp-ip/ip"
	"github.com/malc0mn/ptp-ip/ptp"
	"io/ioutil"
	"log"
	"strings"
)

var (
	liveview  command
	lvEnabled bool
)

func init() {
	if !lvEnabled {
		liveview = func(_ *ip.Client, _ []string) string {
			return "Binary not compiled with live view support!"
		}
	}
}

type command func(*ip.Client, []string) string

func readAndExecuteCommand(rw *bufio.ReadWriter, c *ip.Client, lmp string) {
	msg, err := rw.ReadString('\n')
	if err != nil {
		log.Printf("%s error reading message '%s'", lmp, err)
		return
	}
	msg = strings.TrimSuffix(msg, "\n")
	if msg == "" {
		log.Printf("%s ignoring empty message!", lmp)
		return
	}
	log.Printf("%s message received: '%s'", lmp, msg)

	f := strings.Fields(msg)
	_, err = rw.Write([]byte(commandByName(f[0])(c, f[1:])))
	if err != nil {
		log.Printf("%s error writing response: '%s'", lmp, err)
		return
	}
	err = rw.Flush()
	if err != nil {
		log.Printf("%s error flushing buffer: '%s'", lmp, err)
	}
}

func commandByName(n string) command {
	switch n {
	case "capture", "shoot", "shutter", "snap":
		return capture
	case "describe":
		return describe
	case "get":
		return get
	// TODO: add "help" command that can output usage for all supported commands
	//case "help":
	//	return help
	case "info":
		return info
	case "liveview":
		return liveview
	case "opreq":
		return opreq
	case "set":
		return set
	case "state":
		return state
	default:
		return unknown
	}
}

func unknown(_ *ip.Client, _ []string) string {
	return "unknown command\n"
}

func capture(c *ip.Client, f []string) string {
	res, err := c.InitiateCapture()
	if err != nil {
		return err.Error()
	}
	if len(f) == 1 {
		if err := ioutil.WriteFile(f[0], res, 0644); err != nil {
			return err.Error() + "\n"
		}

		return fmt.Sprintf("Image preview saved to %s\n", f[0])
	}

	return "Image captured, check the camera\n"
}

func describe(c *ip.Client, f []string) string {
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

func info(c *ip.Client, f []string) string {
	res, err := c.GetDeviceInfo()

	if err != nil {
		res = err.Error()
	}

	return formatDeviceInfo(c.ResponderVendor(), res, f)
}

func get(c *ip.Client, f []string) string {
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

func set(c *ip.Client, f []string) string {
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

func opreq(c *ip.Client, f []string) string {
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

func state(c *ip.Client, f []string) string {
	res, err := c.GetDeviceState()

	if err != nil {
		res = err.Error()
	}

	return formatDeviceInfo(c.ResponderVendor(), res, f)
}
