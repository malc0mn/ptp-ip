package main

import (
	"fmt"
	"github.com/malc0mn/ptp-ip/ip"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
)

const (
	ok                  = 0
	errGeneral          = 1
	errInvalidArgs      = 2
	errOpenConfig       = 102
	errCreateClient     = 104
	errResponderConnect = 105
)

var (
	version   = "0.0.0"
	buildTime = "unknown"
	exe       string
)

func main() {
	exe = filepath.Base(os.Args[0])

	initFlags()

	if noArgs := len(os.Args) < 2; noArgs || help {
		printUsage()
		exit := ok
		if noArgs {
			exit = errGeneral
		}
		os.Exit(exit)
	}

	if ver {
		fmt.Printf("%s version %s built on %s\n", exe, version, buildTime)
		os.Exit(ok)
	}

	if file != "" {
		loadConfig()
	}

	checkPorts()

	if cmd != "" && (interactive || server) || (interactive && server) {
		fmt.Fprintln(os.Stderr, "Too many arguments: either run in server mode OR interactive mode OR execute a single command; not all at once!")
		os.Exit(errInvalidArgs)
	}

	// TODO: finish this implementation so CTRL+C will also abort client.Dial() etc. properly.
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	done := make(chan bool, 1)
	go func() {
		sig := <-sigs
		fmt.Printf("Received signal %s, shutting down...\n", sig)
		done <- true
	}()

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

	if cmd != "" {
		f := strings.Fields(cmd)
		fmt.Print(commandByName(f[0])(client, f[1:]))
	}

	if server || interactive {
		if interactive {
			go iShell(client)
		}

		if server {
			go launchServer(client)
		}

		<-done
		fmt.Println("Bye bye!")
	}

	os.Exit(ok)
}
