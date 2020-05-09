package main

import (
	"fmt"
	"os"
	"path/filepath"
)

var (
	Version   = "0.0.0"
	BuildTime = "unknown"
	exe string
)

func main() {
	initFlags()

	exe = filepath.Base(os.Args[0])

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

	if server == true {
		launchServer()
	}

	os.Exit(0)
}
