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
)

func main() {
	exe = filepath.Base(os.Args[0])

	initFlags()

	if noArgs := len(os.Args) < 2; noArgs || help == true {
		usage()
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

	/*sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)*/

	cl, err := ip.NewClient(conf.host, uint16(conf.port), conf.fname, conf.guid)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating PTP/IP client - %s\n", err)
		os.Exit(errCreateClient)
	}
	defer cl.Close()

	fmt.Printf("Created new client with name '%s' and GUID '%s'.\n", cl.InitiatorFriendlyName(), cl.InitiatorGUIDAsString())
	fmt.Printf("Attempting to connect to %s\n", cl.String())
	err = cl.Dial()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error connecting to responder - %s\n", err)
		os.Exit(errResponderConnect)
	}

	if server == true {
		launchServer()
	}

	os.Exit(ok)
}
