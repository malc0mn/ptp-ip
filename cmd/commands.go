package main

import (
	"encoding/hex"
	"fmt"
	"github.com/malc0mn/ptp-ip/ip"
	"log"
)

type command interface {
	execute() string
}

func newCommandByName(n string, c *ip.Client, f []string) command {
	switch n {
	case "info":
		return &info{
			c: c,
			f: f,
		}
	case "state":
		return &state{
			c: c,
			f: f,
		}
	case "opreq":
		return &opreq{
			c: c,
			f: f,
		}
	default:
		return &unknown{}
	}
}

type unknown struct{}

func (u unknown) execute() string {
	return "unknown command"
}

type info struct {
	c *ip.Client
	f []string
}

func (i info) execute() string {
	res, err := i.c.GetDeviceInfo()
log.Printf("%v - %T", res, res)

	if err != nil {
		res = err.Error()
	}

	return formatDeviceInfo(i.c.ResponderVendor(), res, i.f)
}

type state struct {
	c *ip.Client
	f []string
}

func (s state) execute() string {
	res, err := s.c.GetDeviceState()
log.Printf("%v - %T, %s", res, err, err)

	if err != nil {
		res = err.Error()
	}

	return formatDeviceInfo(s.c.ResponderVendor(), res, s.f)
}

type opreq struct {
	c *ip.Client
	f []string
}

func (o opreq) execute() string {
	var res string

	d, err := o.c.OperationRequestRaw(o.f[0], o.f[1:])
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
