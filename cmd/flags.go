package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/malc0mn/ptp-ip/ip"
	"strconv"
)

const (
	defaultIp = "127.0.0.1"
)

var (
	valueOutOfRange = errors.New("value out of range")

	cmd  string
	file string

	server bool

	help    bool
	version bool
)

// Custom flag type that will only accept uint16 values, ideal for ports!
type uint16Value uint64

func (i *uint16Value) Set(s string) error {
	v, err := strconv.ParseUint(s, 0, 16)
	if err != nil || v == 0 {
		err = valueOutOfRange
	}
	*i = uint16Value(v)

	return err
}

func (i *uint16Value) String() string {
	return strconv.FormatInt(int64(*i), 10)
}

func initFlags() {
	flag.StringVar(&conf.host, "h", ip.DefaultIpAddress, "The responder host to connect to.")
	flag.Var(&conf.port, "p", "The responder port to connect to.")
	flag.StringVar(&conf.fname, "n", "", "A custom friendly name to use for the initiator.")
	flag.StringVar(&conf.guid, "g", "", "A custom GUID to use for the initiator. (default random)")

	flag.StringVar(&cmd, "c", "", "The command to send to the responder.")
	flag.StringVar(&file, "f", "", "Read all settings from a config file.")

	flag.BoolVar(&server, "s", false, fmt.Sprintf("This will run the %s command as a server", exe))
	flag.StringVar(&conf.saddr, "sa", defaultIp, "To be used in combination with '-s': this defines the server address to listen on.")
	flag.Var(&conf.sport, "sp", "To be used in combination with '-s': this defines the server port to listen on.")

	flag.BoolVar(&help, "?", false, "Display usage information.")
	flag.BoolVar(&version, "v", false, "Display version info.")

	// Set a custom usage function.
	flag.Usage = usage

	flag.Parse()
}

// TODO: customise.
func usage() {
	fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s:\n", exe)
	flag.PrintDefaults()
}
