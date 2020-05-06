package ip

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"github.com/google/uuid"
	"github.com/malc0mn/ptp-ip/ip/internal"
	"io"
	"log"
	"net"
)

var lmp = "[Mocked responder]"

func newLocalResponder(address string, port int) {
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
		go handleMessage(conn)
	}
}

func handleMessage(conn net.Conn) {
	var err error

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
	case PKT_InitEventRequest:
		log.Printf("%s responding to InitEventRequest", lmp)
		res = &InitEventAckPacket{}
	default:
		log.Printf("%s unknown packet type %#x", lmp, h.PacketType)
		return
	}
	if res != nil {
		err = sendPacket(rw, res)
		if err != nil {
			log.Printf("%s error responding: %s", lmp, err)
		}
		rw.Flush()
	}
}

func sendPacket(w io.Writer, p Packet) error {
	pl := internal.MarshalLittleEndian(p)
	pll := len(pl)

	// The packet length MUST include the header, so we add 8 bytes for that!
	h := internal.MarshalLittleEndian(Header{uint32(pll + HeaderSize), p.PacketType()})

	// Send header.
	n, err := w.Write(h)
	if err != nil {
		return err
	}

	if n != HeaderSize {
		return fmt.Errorf(BytesWrittenMismatch.Error(), n, HeaderSize)
	}
	fmt.Errorf("%s sendPacket bytes written %d", lmp, n)

	// Send payload.
	n, err = w.Write(pl)
	if n != pll {
		return fmt.Errorf(BytesWrittenMismatch.Error(), n, pll)
	}
	if err != nil {
		return err
	}

	fmt.Errorf("%s sendPacket bytes written %d", lmp, n)

	return nil
}
