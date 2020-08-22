package main

import (
	"github.com/malc0mn/ptp-ip/ip"
)

// No init function here!!!

type unknown struct{}

func (unknown) name() string {
	return "unknown"
}

func (unknown) alias() []string {
	return []string{}
}

func (unknown) execute(_ *ip.Client, _ []string, _ chan<- string) string {
	return "unknown command\n"
}

func (c unknown) help() string {
	return ""
}

func (unknown) arguments() []string {
	return []string{}
}
