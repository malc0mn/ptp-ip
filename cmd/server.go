package main

import (
	"bufio"
	"log"
	"net"
)

func validateAddress() {
	if ip := net.ParseIP(conf.saddr); ip == nil {
		log.Fatalf("Invalid IP address '%s'", conf.saddr)
	}
}

func launchServer() {
	validateAddress()

	lmp := "[Local server]"
	sock, err := net.Listen("tcp", net.JoinHostPort(conf.saddr, conf.sport.String()))
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
		go handleMessages(conn, lmp)
	}
}

func handleMessages(conn net.Conn, lmp string) {
	defer conn.Close()
	rw := bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn))

	msg, err := rw.ReadString('\n')
	if err != nil {
		log.Printf("%s Error reading message '%s'", lmp, err)
		return
	}
	log.Printf("%s Message received: '%s'", lmp, msg)

	// TODO: actual message handling goes here.
}
