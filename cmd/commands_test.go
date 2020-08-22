package main

import (
	"fmt"
	"github.com/malc0mn/ptp-ip/ip"
	"testing"
)

func TestCommandByName(t *testing.T) {
	cmds := map[string]command{
		"capture":  &capture{},
		"describe": &describe{},
		"get":      &get{},
		"help":     &help{},
		"info":     &info{},
		"liveview": &liveview{},
		"opreq":    &opreq{},
		"shoot":    &capture{},
		"shutter":  &capture{},
		"snap":     &capture{},
		"set":      &set{},
		"state":    &state{},
	}
	for name, want := range cmds {
		got := commandByName(name)
		if fmt.Sprintf("%v", got) != fmt.Sprintf("%v", want) {
			t.Errorf("commandByName(%s) got = %v; want %v", name, got, want)
		}
	}
}

func TestUnknown(t *testing.T) {
	got := unknown{}.execute(&ip.Client{}, []string{}, make(chan string))
	want := "unknown command\n"
	if got != want {
		t.Errorf("got = '%s'; want '%s'", got, want)
	}
}
