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

	interactive bool
	server bool

	help bool
	ver  bool

	verbosity ip.LogLevel
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
	flag.StringVar(&conf.vendor, "t", ip.DefaultVendor, "The vendor of the responder that will be connected to.")
	flag.StringVar(&conf.host, "h", ip.DefaultIpAddress, "The responder host to connect to.")
	flag.Var(&conf.port, "p", "The responder port to connect to. Use this flag when the responder has only ONE port for all channels!")
	flag.Var(&conf.cport, "pc", "The responder port used for the Command/Data connection.")
	flag.Var(&conf.eport, "pe", "The responder port used for the Event connection.")
	flag.Var(&conf.sport, "ps", "The responder port used for the streamer or 'live view' connection.")
	flag.StringVar(&conf.fname, "n", "", "A custom friendly name to use for the initiator.")
	flag.StringVar(&conf.guid, "g", "", "A custom GUID to use for the initiator. (default random)")

	flag.BoolVar(&interactive, "i", false, fmt.Sprintf("This will run the %s command with an interactive shell.", exe))

	flag.StringVar(&cmd, "c", "", "The command to send to the responder.")
	flag.StringVar(&file, "f", "", "Read all settings from a config file. The config file will override any command line flags present.")

	flag.BoolVar(&server, "s", false, fmt.Sprintf("This will run the %s command as a server", exe))
	flag.StringVar(&conf.srvAddr, "sa", defaultIp, "To be used in combination with '-s': this defines the server address to listen on.")
	flag.Var(&conf.srvPort, "sp", "To be used in combination with '-s': this defines the server port to listen on.")

	flag.BoolVar(&help, "?", false, "Display usage information.")
	flag.BoolVar(&ver, "version", false, "Display version info.")

	flag.Var(&verbosity, "v", "PTP/IP log level verbosity: ranges from v to vvv.")

	// Set a custom usage function.
	flag.Usage = printUsage

	flag.Parse()
}

// TODO: customise.
func printUsage() {
	fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s:\n", exe)
	flag.PrintDefaults()
}
