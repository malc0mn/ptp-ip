// +build !with_lv

package main

import "github.com/malc0mn/ptp-ip/ip"

var nolv = "Binary not compiled with live view support!"

func mainThread() {
	return
}

func liveview(_ *ip.Client, _ []string) string {
	return nolv + "\n"
}

func preview(_ []byte) string {
	return nolv
}
