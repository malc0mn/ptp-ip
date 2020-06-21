package main

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"github.com/malc0mn/ptp-ip/ip"
	"log"
	"net"
	"strings"
)

func validateAddress() {
	if ip := net.ParseIP(conf.srvAddr); ip == nil {
		log.Fatalf("Invalid IP address '%s'", conf.srvAddr)
	}
}

func launchServer(c *ip.Client) {
	validateAddress()

	lmp := "[Local server]"
	sock, err := net.Listen("tcp", net.JoinHostPort(conf.srvAddr, conf.srvPort.String()))
	defer sock.Close()
	if err != nil {
		log.Printf("%s error %s...", lmp, err)
		return
	}
	log.Printf("%s listening on %s...", lmp, sock.Addr().String())
	log.Printf("%s awaiting messages... (CTRL+C to quit)", lmp)

	for {
		conn, err := sock.Accept()
		if err != nil {
			log.Printf("%s accept error %s...", lmp, err)
			continue
		}
		go handleMessages(conn, c, lmp)
	}
}

func handleMessages(conn net.Conn, c *ip.Client, lmp string) {
	defer conn.Close()
	rw := bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn))

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

	f := strings.Fields(msg)
	switch f[0] {
	case "info":
		res, err := c.GetDeviceInfo()
		log.Printf("%v - %T", res, res)

		if err != nil {
			res = err.Error()
		}
		conn.Write([]byte(PrintDeviceInfo(c.ResponderVendor(), res)))
	case "state":
		res, err := c.GetDeviceState()
		log.Printf("%v - %T, %s", res, err, err)

		if err != nil {
			res = err.Error()
		}
		conn.Write([]byte(PrintDeviceInfo(c.ResponderVendor(), res)))
	case "opreq":
		var res string
		d, err := c.OperationRequestRaw(f[1], f[2:])
		if err != nil {
			res = fmt.Sprintf("%s error: %s", lmp, err)
		} else {
			log.Printf("%s received %d packets.", lmp, len(d))
			for _, raw := range d {
				res += fmt.Sprintf("\nReceived %d bytes. HEX dump:\n%s", len(raw), hex.Dump(raw))
			}
		}

		conn.Write([]byte(res))
	}
}
