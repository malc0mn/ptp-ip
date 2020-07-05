package main

import (
	"encoding/hex"
	"fmt"
	"github.com/malc0mn/ptp-ip/ip"
	"log"
)

type command func(*ip.Client, []string) string

func commandByName(n string) command {
	switch n {
	case "info":
		return info
	case "state":
		return state
	case "opreq":
		return opreq
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

func state(c *ip.Client, f []string) string {
	res, err := c.GetDeviceState()
log.Printf("%v - %T, %s", res, err, err)

	if err != nil {
		res = err.Error()
	}

	return formatDeviceInfo(c.ResponderVendor(), res, f)
}

func opreq(c *ip.Client, f []string) string {
	var res string

	d, err := c.OperationRequestRaw(f[0], f[1:])
	if err != nil {
		res = fmt.Sprintf("opreq error: %s", err)
	} else {
		log.Printf("opreq received %d packets.", len(d))
		for _, raw := range d {
			res += fmt.Sprintf("\nReceived %d bytes. HEX dump:\n%s", len(raw), hex.Dump(raw))
		}
	}

	return res
}
