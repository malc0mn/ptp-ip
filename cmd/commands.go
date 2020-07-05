package main

import (
	"encoding/hex"
	"fmt"
	"github.com/malc0mn/ptp-ip/ip"
	"github.com/malc0mn/ptp-ip/ptp"
	"log"
)

type command func(*ip.Client, []string) string

func commandByName(n string) command {
	switch n {
	case "info":
		return info
	case "getval":
		return getval
	case "opreq":
		return opreq
	case "state":
		return state
	default:
		return unknown
	}
}

func unknown(_ *ip.Client, _ []string) string {
	return "unknown command"
}

func info(c *ip.Client, f []string) string {
	res, err := c.GetDeviceInfo()
log.Printf("%v - %T", res, res)

	if err != nil {
		res = err.Error()
	}

	return formatDeviceInfo(c.ResponderVendor(), res, f)
}

func getval(c *ip.Client, f []string) string {
	errorFmt := "getval error: %s\n"

	cod, err := ip.HexStringToUint64(f[0], 16)
	if err != nil {
		return fmt.Sprintf(errorFmt, err)
	}
	c.Debugf("Converted uint16: %#x", cod)

	v, err := c.GetDevicePropertyValue(ptp.DevicePropCode(cod))
	if err != nil {
		return fmt.Sprintf(errorFmt, err)
	}

	// TODO: add vendor check here or create a generic *AsString() function somehow...?
	return formatDevicePropVal(c.ResponderVendor(), ptp.DevicePropCode(cod), int64(v)) + fmt.Sprintf(" (%#x)", v)
}

func opreq(c *ip.Client, f []string) string {
	var res string
	errorFmt := "opreq error: %s\n"

	cod, err := ip.HexStringToUint64(f[0], 16)
	if err != nil {
		return fmt.Sprintf(errorFmt, err)
	}
	c.Debugf("Converted uint16: %#x", cod)

	params := f[1:]
	p := make([]uint32, len(params))
	for i, param := range params {
		conv, err := ip.HexStringToUint64(param, 64)
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

log.Printf("opreq received %d packets.", len(d))
	for _, raw := range d {
		res += fmt.Sprintf("\nReceived %d bytes. HEX dump:\n%s", len(raw), hex.Dump(raw))
	}

	return res
}

func state(c *ip.Client, f []string) string {
	res, err := c.GetDeviceState()
log.Printf("%v - %T, %s", res, err, err)

	if err != nil {
		res = err.Error()
	}

	return formatDeviceInfo(c.ResponderVendor(), res, f)
}
