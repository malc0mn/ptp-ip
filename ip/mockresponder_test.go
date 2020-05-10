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

type MockedResponder struct {
	address string
	port    uint16
	handler func(net.Conn, string)
	lmp     string
}

func runResponder(address string, port uint16, handler func(net.Conn, string), lmp string) {
	mr := &MockedResponder{
		address: address,
		port:    port,
		handler: handler,
		lmp:     lmp,
	}

	mr.run()
}

func newLocalOkResponder(address string, port uint16) {
	runResponder(address, port, handleMessage, "[Mocked OK responder]")
}

func newLocalFailResponder(address string, port uint16) {
	runResponder(address, port, alwaysFailMessage, "[Mocked FAIL responder]")
}

func (mr *MockedResponder) run() {
	ln, err := net.Listen("tcp", fmt.Sprintf("%s:%d", mr.address, mr.port))
	defer ln.Close()
	if err != nil {
		log.Printf("%s error %s...", mr.lmp, err)
		return
	}
	log.Printf("%s listening on %s...", mr.lmp, ln.Addr().String())

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Printf("%s accept error %s...", mr.lmp, err)
			continue
		}
		go mr.handler(conn, mr.lmp)
	}
}

func readMessage(rw *bufio.ReadWriter, lmp string) (Header, PacketOut) {
	var err error

	log.Printf("%s received message", lmp)

	var h Header
	log.Printf("%s reading packet header...", lmp)
	err = binary.Read(rw, binary.LittleEndian, &h)
	if err != nil {
		log.Printf("%s error reading header: %s", lmp, err)
		return h, nil
	}
	pkt, err := NewPacketOutFromPacketType(h.PacketType)
	if err != nil {
		log.Printf("%s error creating packet: %s", lmp, err)
		return h, nil
	}

	vs := int(h.Length) - HeaderSize - internal.TotalSizeOfFixedFields(pkt)
	err = internal.UnmarshalLittleEndian(rw, pkt, vs)
	if err != nil {
		log.Printf("%s error reading packet %T data %s", lmp, pkt, err)
		return h, nil
	}

	return h, pkt
}

func writeMessage(rw *bufio.ReadWriter, pkt Packet, lmp string) {
	err := sendPacket(rw, pkt, lmp)
	if err != nil {
		log.Printf("%s error responding: %s", lmp, err)
	}
	rw.Flush()
}

func handleMessage(conn net.Conn, lmp string) {
	defer conn.Close()
	rw := bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn))
	h, pkt := readMessage(rw, lmp)
	if pkt == nil {
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
		writeMessage(rw, res, lmp)
	}
}

func alwaysFailMessage(conn net.Conn, lmp string) {
	defer conn.Close()
	rw := bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn))
	_, pkt := readMessage(rw, lmp)
	if pkt == nil {
		return
	}

	writeMessage(rw, &InitFailPacket{
		Reason: FR_FailRejectedInitiator,
	}, lmp)
}

func sendPacket(w io.Writer, p Packet, lmp string) error {
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
