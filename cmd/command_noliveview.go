// +build !with_lv

package main

import "github.com/malc0mn/ptp-ip/ip"

var nolv = "Binary not compiled with live view support!"

func init() {
	registerCommand(&liveview{})
}

type liveview struct{}

func (liveview) name() string {
	return "liveview"
}

func (liveview) alias() []string {
	return []string{}
}

func (liveview) execute(_ *ip.Client, _ []string, _ chan<- string) string {
	return nolv + "\n"
}

func (l liveview) help() string {
	return `"` + l.name() + `" is not supported in this build!`
}

func (liveview) arguments() []string {
	return []string{}
}

func mainThread() {
	return
}

func preview(_ []byte) string {
	return nolv
}
