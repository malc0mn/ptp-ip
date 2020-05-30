package main

import (
	"errors"
	"fmt"
	"github.com/go-ini/ini"
	"github.com/malc0mn/ptp-ip/ip"
	"log"
	"os"
)

type config struct {
	vendor string
	host   string
	port   uint16Value
	cport  uint16Value
	eport  uint16Value
	sport  uint16Value
	fname  string
	guid   string

	srvAddr string
	srvPort uint16Value
}

var (
	portSpecAmbiguous = errors.New("ambiguous port specification: use a single port OR define multiple ports")

	conf = &config{
		vendor:  ip.DefaultVendor,
		host:    ip.DefaultIpAddress,
		port:    uint16Value(ip.DefaultPort),
		srvAddr: defaultIp,
		srvPort: uint16Value(ip.DefaultPort),
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
		if k, err := i.GetKey("vendor"); err == nil {
			conf.vendor = k.String()
		}
		if k, err := i.GetKey("host"); err == nil {
			conf.host = k.String()
		}
		if k, err := i.GetKey("port"); err == nil {
			if err := conf.port.Set(k.String()); err != nil {
				log.Fatal(valueOutOfRange)
			}
		}
		if k, err := i.GetKey("cmd_data_port"); err == nil {
			if err := conf.cport.Set(k.String()); err != nil {
				log.Fatal(valueOutOfRange)
			}
		}
		if k, err := i.GetKey("event_port"); err == nil {
			if err := conf.eport.Set(k.String()); err != nil {
				log.Fatal(valueOutOfRange)
			}
		}
		if k, err := i.GetKey("stream_port"); err == nil {
			if err := conf.sport.Set(k.String()); err != nil {
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
			conf.srvAddr = k.String()
		}
		if k, err := i.GetKey("port"); err == nil {
			if err := conf.srvPort.Set(k.String()); err != nil {
				log.Fatal(valueOutOfRange)
			}
		}
	}
}

func checkPorts() {
	if conf.cport != 0 && conf.eport != 0 {
		conf.port = 0
	}

	if conf.port != 0 && (conf.cport != 0 || conf.eport != 0) {
		log.Fatal(portSpecAmbiguous)
	}
}
