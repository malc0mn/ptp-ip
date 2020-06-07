package main

import (
	"bufio"
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
		log.Printf("%s Error reading message '%s'", lmp, err)
		return
	}
	msg = strings.TrimSuffix(msg, "\n")
	log.Printf("%s Message received: '%s'", lmp, msg)

	switch msg {
	case "info":
		res, err := c.GetDeviceInfo()
		log.Printf("%v - %T, %s", res, err, err)
	}
}
