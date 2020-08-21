package main

import (
	"bufio"
	"fmt"
	"github.com/malc0mn/ptp-ip/ip"
	"os"
	"time"
)

// TODO: add a channel to receive output from async processes, like the multi capture command
func iShell(c *ip.Client) {
	rw := bufio.NewReadWriter(bufio.NewReader(os.Stdin), bufio.NewWriter(os.Stdout))
	fmt.Print("Interactive shell ready to receive commands.\n")
	for {
		// TODO: find a good way (not sleep) to "separate" the outputs so that the '> ' below does not get 'mixed' with
		//  the Dial() debug output from the client...
		time.Sleep(1 * time.Second)

		fmt.Print("> ")
		readAndExecuteCommand(rw, c, "[iShell]")
		fmt.Print("\n\n")
	}
}
