package main

import (
	"github.com/malc0mn/ptp-ip/ip"
	"sort"
)

func init() {
	registerCommand(&help{})
}

type help struct{}

func (help) name() string {
	return "help"
}

func (help) alias() []string {
	return []string{}
}

func (help) execute(_ *ip.Client, f []string, _ chan<- string) string {
	if len(f) == 0 {
		names := make([]string, 0, len(commands))
		for name := range commands {
			names = append(names, name)
		}
		sort.Strings(names)

		txt := "\nSupported commands:\n\n"
		for _, name := range names {
			txt += commands[name].help() + "\n\n"
		}
		return txt
	}

	if cmd, exists := commands[f[0]]; exists {
		return "\n" + cmd.help()
	}

	return "\nUnknown command " + f[0] + "!\n"
}

func (h help) help() string {
	help := `"` + h.name() + `" displays help for all commands or for a single one.` + "\n"

	if args := h.arguments(); len(args) > 0 {
		help += helpAddArgumentsTitle()
		for i, arg := range args {
			switch i {
			case 0:
				help += "\t- " + arg + " to get help for\n"
			}
		}
	}

	return help
}

func (help) arguments() []string {
	return []string{"command"}
}
