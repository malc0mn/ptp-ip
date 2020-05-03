package ip

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"github.com/google/uuid"
	"github.com/malc0mn/ptp-ip/ip/internal"
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
			rw := bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn))

			var h Header
			log.Printf("%s reading packet header...", lmp)
			err = binary.Read(rw, binary.LittleEndian, &h)
			if err != nil {
				log.Printf("%s error reading header: %s", lmp, err)
				return
			}
			pkt, err := NewPacketOutFromPacketType(h.PacketType)
			if err != nil {
				log.Printf("%s error creating packet: %s", lmp, err)
				return
			}

			vs := int(h.Length) - HeaderSize - internal.TotalSizeOfFixedFields(pkt)
			err = internal.UnmarshalLittleEndian(rw, pkt, vs)
			if err != nil {
				log.Printf("%s error reading packet %T data %s", lmp, pkt, err)
				return
			}

			var res PacketIn
			switch h.PacketType {
			case PKT_InitCommandRequest:
				log.Printf("%s responding to InitCommandRequest", lmp)
				uuid, _ := uuid.Parse("3e8626cc-5059-4225-bdd6-d160b2e6a60f")
				res = &InitCommandAckPacket{
					ConnectionNumber:         1,
					ResponderGUID:            uuid,
					ResponderFriendlyName:    lmp,
					ResponderProtocolVersion: uint32(PV_VersionOnePointZero),
				}
			}
			if res != nil {
				binary.Write(rw, binary.LittleEndian, res)
				rw.Flush()
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
