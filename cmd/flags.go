package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/malc0mn/ptp-ip/ip"
	"strconv"
)

const (
	binary    = "ptpip"
	defaultIp = "127.0.0.1"
)

var (
	host  string
	port  = uint16Value(ip.DefaultPort)
	fname string
	guid  string

	cmd    string
	config string

	server bool
	saddr  string
	sport  = uint16Value(ip.DefaultPort)

	help    bool
	version bool
)

// Custom flag type that will only accept uint16 values, ideal for ports!
type uint16Value uint64

func (i *uint16Value) Set(s string) error {
	v, err := strconv.ParseUint(s, 0, 16)
	if err != nil {
		err = errors.New("value out of range")
	}
	*i = uint16Value(v)

	return err
}

func (i *uint16Value) String() string {
	return strconv.FormatInt(int64(*i), 10)
}

func initFlags() {
	flag.StringVar(&host, "h", ip.DefaultIpAddress, "The responder host to connect to.")
	flag.Var(&port, "p", "The responder port to connect to.")
	flag.StringVar(&fname, "n", "", "A custom friendly name to use for the initiator.")
	flag.StringVar(&guid, "g", "", "A custom GUID to use for the initiator. (default random)")

	flag.StringVar(&cmd, "c", "", "The command to send to the responder.")
	flag.StringVar(&config, "f", "", "Read all settings from a config file.")

	flag.BoolVar(&server, "s", false, fmt.Sprintf("This will run the %s command as a server", exe))
	flag.StringVar(&saddr, "sa", defaultIp, "To be used in combination with '-s': this defines the server address to listen on.")
	flag.Var(&sport, "sp", "To be used in combination with '-s': this defines the server port to listen on.")

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
