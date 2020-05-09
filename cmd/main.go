package main

import (
	"fmt"
	"github.com/malc0mn/ptp-ip/ip"
	"os"
	"path/filepath"
)

var (
	Version   = "0.0.0"
	BuildTime = "unknown"
	exe string
)

func main() {
	exe = filepath.Base(os.Args[0])

	initFlags()

	if noArgs := len(os.Args) < 2; noArgs || help == true {
		usage()
		exit := 0
		if noArgs {
			exit = 1
		}
		os.Exit(exit)
	}

	if version == true {
		fmt.Printf("%s version %s built on %s\n", exe, Version, BuildTime)
		os.Exit(0)
	}

	/*sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)*/

	c, err := ip.NewClient(host, int(port), fname, guid)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating PTP/IP client: %s\n", err)
		os.Exit(2)
	}

	fmt.Printf("Created new client with name '%s' and GUID '%s'.\n", c.InitiatorFriendlyName(), c.InitiatorGUIDAsString())
	fmt.Printf("Attempting to connect to: %s\n", c.String())
	err = c.Dial()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating PTP/IP client: %s\n", err)
		os.Exit(3)
	}

	if server == true {
		launchServer()
	}

	os.Exit(0)
}
