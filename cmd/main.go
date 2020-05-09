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

	if help == true || len(os.Args) < 2 {
		usage()
		os.Exit(0)
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
