package main

import (
	"fmt"
	"github.com/go-ini/ini"
	"github.com/malc0mn/ptp-ip/ip"
	"log"
	"os"
)

type config struct {
	host  string
	port  uint16Value
	fname string
	guid  string

	saddr string
	sport uint16Value
}

var (
	conf = &config{
		host:  ip.DefaultIpAddress,
		port:  uint16Value(ip.DefaultPort),
		saddr: defaultIp,
		sport: uint16Value(ip.DefaultPort),
	}
)

func loadConfig() {
	f, err := ini.Load(file)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening config file - %s\n", err)
		os.Exit(errOpenConfig)
	}

	// Initiator
	if i, err := f.GetSection("initiator"); err == nil {
		if k, err := i.GetKey("friendly_name"); err == nil {
			conf.fname = k.String()
		}
		if k, err := i.GetKey("guid"); err == nil {
			conf.guid = k.String()
		}
	}

	// Responder
	if i, err := f.GetSection("responder"); err == nil {
		if k, err := i.GetKey("host"); err == nil {
			conf.host = k.String()
		}
		if k, err := i.GetKey("port"); err == nil {
			if err := conf.port.Set(k.String()); err != nil {
				log.Fatal(valueOutOfRange)
			}
		}
	}

	// Server
	if i, err := f.GetSection("server"); err == nil {
		if k, err := i.GetKey("enabled"); err == nil {
			if v, err := k.Bool(); err == nil {
				server = v
			}
		}
		if k, err := i.GetKey("address"); err == nil {
			conf.saddr = k.String()
		}
		if k, err := i.GetKey("port"); err == nil {
			if err := conf.sport.Set(k.String()); err != nil {
				log.Fatal(valueOutOfRange)
			}
		}
	}
}
