package ip

import (
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"strconv"
)

func newLocalResponder(address string, port int) {
	lmp := "[Local responder]"

	ln, err := net.Listen("tcp", fmt.Sprintf("%s:%d", address, port))
	defer ln.Close()
	if err != nil {
		log.Printf("%s error %s...", lmp, err)
		return
	}
	log.Printf("%s listening on %s...", lmp, ln.Addr().String())

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Printf("%s accept error %s...", lmp, err)
			continue
		}
		go func(conn net.Conn) {
			log.Printf("%s received message", lmp)
			defer conn.Close()
			var h Header
			log.Printf("%s reading packet header...", lmp)
			err = binary.Read(conn, binary.LittleEndian, &h)
			if err != nil {
				log.Printf("%s error reading header: %s", lmp, err)
				return
			}
			switch h.PacketType {
			case PKT_InitCommandRequest:
				log.Printf("%s need to respond to %#x", lmp, h.PacketType)
			}
		}(conn)
	}
}

func getListenerHostAndPort(ln net.Listener) (string, int, error) {
	var err error

	if host, port, err := net.SplitHostPort(ln.Addr().String()); err == nil {
		p, _ := strconv.Atoi(port)
		return host, p, nil
	}

	return "", 0, err
}
