package main

import (
	"fmt"
	"github.com/malc0mn/ptp-ip/ip"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
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

	// TODO: finish this implementation so CTRL+C will also abort client.Dial() etc. properly.
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	done := make(chan bool, 1)

	client, err := ip.NewClient(conf.vendor, conf.host, uint16(conf.port), conf.fname, conf.guid, verbosity)
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
	fmt.Printf("Attempting to connect to %s\n", client.CommandDataAddress())
	err = client.Dial()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error connecting to responder - %s\n", err)
		os.Exit(errResponderConnect)
	}

	if server == true {
		go launchServer(client)
		go func() {
			sig := <-sigs
			fmt.Printf("Received signal %s, shutting down...\n", sig)
			done <- true
		}()
	}

	<-done
	fmt.Println("Bye bye!")
	os.Exit(ok)
}
