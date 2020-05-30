package main

import (
	"fmt"
	"github.com/malc0mn/ptp-ip/ip"
	"os"
	"path/filepath"
)

const (
	ok                  = 0
	errGeneral          = 1
	errOpenConfig       = 102
	errCreateClient     = 104
	errResponderConnect = 105
)

var (
	Version   = "0.0.0"
	BuildTime = "unknown"
	exe       string
	client    ip.Client
)

func main() {
	exe = filepath.Base(os.Args[0])

	initFlags()

	if noArgs := len(os.Args) < 2; noArgs || help == true {
		printUsage()
		exit := ok
		if noArgs {
			exit = errGeneral
		}
		os.Exit(exit)
	}

	if version == true {
		fmt.Printf("%s version %s built on %s\n", exe, Version, BuildTime)
		os.Exit(ok)
	}

	if file != "" {
		loadConfig()
	}

	checkPorts()

	/*sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)*/

	client, err := ip.NewClient(conf.vendor, conf.host, uint16(conf.port), conf.fname, conf.guid)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating PTP/IP client - %s\n", err)
		os.Exit(errCreateClient)
	}
	defer client.Close()

	if conf.cport != 0 {
		client.SetCommandDataPort(uint16(conf.cport))
	}
	if conf.eport != 0 {
		client.SetEventPort(uint16(conf.eport))
	}
	if conf.sport != 0 {
		client.SetStreamerPort(uint16(conf.sport))
	}

	fmt.Printf("Created new client with name '%s' and GUID '%s'.\n", client.InitiatorFriendlyName(), client.InitiatorGUIDAsString())
	fmt.Printf("Attempting to connect to %s\n", client.CommandDataAddres())
	err = client.Dial()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error connecting to responder - %s\n", err)
		os.Exit(errResponderConnect)
	}

	if server == true {
		launchServer()
	}

	os.Exit(ok)
}
