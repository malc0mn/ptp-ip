package main

import (
	"bufio"
	ptpfmt "github.com/malc0mn/ptp-ip/fmt"
	"github.com/malc0mn/ptp-ip/ip"
	"log"
	"strings"
	"sync"
)

var (
	commandsMu sync.RWMutex
	commands   = make(map[string]command)
	aliases    = make(map[string]string)
)

type command interface {
	name() string
	alias() []string
	// TODO: is there a more elegant solution to drop at least the async output channel argument here...?
	execute(*ip.Client, []string, chan<- string) string
	help() string
	arguments() []string
}

func registerCommand(cmd command) {
	commandsMu.Lock()
	defer commandsMu.Unlock()
	if cmd == nil {
		panic("cmd: registerCommand command is nil")
	}

	name := cmd.name()
	if _, dup := commands[name]; dup {
		panic("cmd: registerCommand called twice for command " + name)
	}
	commands[name] = cmd

	for _, alias := range cmd.alias() {
		if _, dup := aliases[alias]; dup {
			panic("cmd: registerCommand double alias " + alias)
		}
		aliases[alias] = name
	}
}

func helpAddAliases(aliases []string) string {
	var help string

	if len(aliases) > 0 {
		help += "\n\t" + `Possible aliases: "` + strings.Join(aliases, `", "`) + `"` + "\n"
	}

	return help
}

func helpAddArgumentsTitle() string {
	return "\tAllowed arguments:\n"
}

func helpAddUnifiedFieldNames() string {
	return "\t" + `  "` + strings.Join(ptpfmt.UnifiedFieldNames, `", "`) + `"` + "\n"
}

func readAndExecuteCommand(rw *bufio.ReadWriter, c *ip.Client, lmp string) {
	msg, err := rw.ReadString('\n')
	if err != nil {
		log.Printf("%s error reading message '%s'", lmp, err)
		return
	}
	msg = strings.TrimSuffix(msg, "\n")
	if msg == "" {
		log.Printf("%s ignoring empty message!", lmp)
		return
	}
	log.Printf("%s message received: '%s'", lmp, msg)

	executeCommand(msg, rw.Writer, c, lmp)
}

func executeCommand(msg string, w *bufio.Writer, c *ip.Client, lmp string) {
	var wg sync.WaitGroup
	f := strings.Fields(msg)
	asyncOut := make(chan string)

	// Launch async output routine.
	wg.Add(1)
	go func() {
		for msg := range asyncOut {
			if _, err := w.Write([]byte(msg + "\n")); err != nil {
				log.Printf("%s error writing response: '%s'", lmp, err)
				continue
			}
			if err := w.Flush(); err != nil {
				log.Printf("%s error flushing buffer: '%s'", lmp, err)
			}
		}
		wg.Done()
	}()

	_, err := w.Write([]byte(commandByName(f[0]).execute(c, f[1:], asyncOut)))
	close(asyncOut)
	wg.Wait()
	if err != nil {
		log.Printf("%s error writing response: '%s'", lmp, err)
		return
	}
	err = w.Flush()
	if err != nil {
		log.Printf("%s error flushing buffer: '%s'", lmp, err)
	}
}

func commandByName(n string) command {
	if name, exists := aliases[n]; exists {
		n = name
	}

	if cmd, exists := commands[n]; exists {
		return cmd
	}

	return &unknown{}
}
