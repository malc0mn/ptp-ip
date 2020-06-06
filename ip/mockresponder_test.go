package ip

import (
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
	runResponder(address, port, handleMessages, "[Mocked OK responder]")
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
		log.Printf("%s new connection %v...", mr.lmp, conn)
		go mr.handler(conn, mr.lmp)
	}
}

func readMessage(r io.Reader, lmp string) (Header, PacketOut, error) {
	var err error

	var h Header
	log.Printf("%s awaiting packet header...", lmp)
	err = binary.Read(r, binary.LittleEndian, &h)
	if err != nil {
		if err == io.EOF {
			log.Printf("%s client disconnected", lmp)
		} else {
			log.Printf("%s error reading header: %s", lmp, err)
		}
		return h, nil, err
	}
	pkt, err := NewPacketOutFromPacketType(h.PacketType)
	if err != nil {
		log.Printf("%s error creating packet: %s", lmp, err)
		return h, nil, err
	}

	vs := int(h.Length) - HeaderSize - internal.TotalSizeOfFixedFields(pkt)
	err = internal.UnmarshalLittleEndian(r, pkt, int(h.Length), vs)
	if err != nil {
		log.Printf("%s error reading packet %T data %s", lmp, pkt, err)
		return h, nil, err
	}
	log.Printf("%v %T", pkt, pkt)
	return h, pkt, nil
}

func writeMessage(w io.Writer, pkt Packet, lmp string) {
	err := sendPacket(w, pkt, lmp)
	if err != nil {
		log.Printf("%s error responding: %s", lmp, err)
	}
}

func handleMessages(conn net.Conn, lmp string) {
	// NO defer conn.Close() here since we need to mock a real responder and thus need to keep the connections open when
	// established and continuously listen for messages in a loop.
	for {
		h, pkt, err := readMessage(conn, lmp)
		if err == io.EOF {
			conn.Close()
			break
		}
		if pkt == nil {
			continue
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
		case PKT_OperationRequest:
			log.Printf("%s responding to OperationRequest", lmp)
			res = &OperationResponsePacket{}
		default:
			log.Printf("%s unknown packet type %#x", lmp, h.PacketType)
			continue
		}
		if res != nil {
			writeMessage(conn, res, lmp)
		}
	}
}

func alwaysFailMessage(conn net.Conn, lmp string) {
	// TCP connections are closed by the Responder on failure!
	defer conn.Close()
	_, pkt, _ := readMessage(conn, lmp)
	if pkt == nil {
		return
	}

	writeMessage(conn, &InitFailPacket{
		Reason: FR_FailRejectedInitiator,
	}, lmp)
}

func sendPacket(w io.Writer, p Packet, lmp string) error {
	log.Printf("%s sendPacket %T", lmp, p)

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
	log.Printf("%s sendPacket header bytes written %d", lmp, n)

	// Send payload.
	if pll == 0 {
		log.Printf("%s packet has no payload", lmp)
		return nil
	}

	n, err = w.Write(pl)
	if err != nil {
		return err
	}
	if n != pll {
		return fmt.Errorf(BytesWrittenMismatch.Error(), n, pll)
	}

	log.Printf("%s sendPacket payload bytes written %d", lmp, n)

	return nil
}
