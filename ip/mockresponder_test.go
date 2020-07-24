package ip

import (
	"encoding/binary"
	"flag"
	"fmt"
	"github.com/malc0mn/ptp-ip/ip/internal"
	"github.com/malc0mn/ptp-ip/ptp"
	"io"
	"log"
	"net"
	"os"
	"testing"
)

const MockResponderGUID string = "3e8626cc-5059-4225-bdd6-d160b2e6a60f"

var (
	address            = "127.0.0.1"
	okPort             = DefaultPort
	fujiCmdPort uint16 = 55740
	fujiEvtPort uint16 = 55741
	failPort    uint16 = 25740
	logLevel           = LevelSilent
	lgr         Logger
	// TODO: the eventChan needs to be moved, see the comment in TestMain()
	evtChan = make(chan uint32, 10)
)

func TestMain(m *testing.M) {
	flag.Parse()

	if testing.Verbose() {
		logLevel = LevelDebug
	}

	lgr = NewLogger(logLevel, os.Stderr, "", log.LstdFlags)

	go newLocalOkResponder(DefaultVendor, address, okPort)
	// TODO: this is not good, we need to integrate the event handler in a single mock responder run...
	go newLocalOkResponder("fuji", address, fujiCmdPort)
	go newLocalOkResponder("fuji-event", address, fujiEvtPort)
	go newLocalFailResponder(address, failPort)
	os.Exit(m.Run())
}

type msgHandler func(net.Conn, chan uint32, string)

type MockResponder struct {
	vendor  ptp.VendorExtension
	address string
	port    uint16
	handler msgHandler
	lmp     string
}

func runResponder(vendor ptp.VendorExtension, address string, port uint16, handler msgHandler, lmp string) {
	mr := &MockResponder{
		vendor:  vendor,
		address: address,
		port:    port,
		handler: handler,
		lmp:     lmp,
	}

	mr.run()
}

func newLocalOkResponder(vendor string, address string, port uint16) {
	var handler msgHandler
	switch vendor {
	case "fuji":
		handler = handleFujiMessages
	case "fuji-event":
		handler = handleFujiEvents
	default:
		handler = handleGenericMessages
	}

	runResponder(ptp.VendorStringToType(vendor), address, port, handler, fmt.Sprintf("[Mocked %s OK responder]", vendor))
}

func newLocalFailResponder(address string, port uint16) {
	runResponder(ptp.VendorExtension(0), address, port, alwaysFailMessage, "[Mocked FAIL responder]")
}

func (mr *MockResponder) run() {
	ln, err := net.Listen("tcp", fmt.Sprintf("%s:%d", mr.address, mr.port))
	defer ln.Close()
	if err != nil {
		lgr.Errorf("%s error %s...", mr.lmp, err)
		return
	}
	lgr.Infof("%s listening on %s...", mr.lmp, ln.Addr().String())

	for {
		conn, err := ln.Accept()
		if err != nil {
			lgr.Errorf("%s accept error %s...", mr.lmp, err)
			continue
		}
		lgr.Infof("%s new connection %v...", mr.lmp, conn)
		go mr.handler(conn, evtChan, mr.lmp)
	}
}

func readMessage(r io.Reader, lmp string) (Header, PacketOut, error) {
	var err error

	var h Header
	lgr.Infof("%s awaiting packet header...", lmp)
	err = binary.Read(r, binary.LittleEndian, &h)
	if err != nil {
		if err == io.EOF {
			lgr.Infof("%s client disconnected", lmp)
		} else {
			lgr.Errorf("%s error reading header: %s", lmp, err)
		}
		return h, nil, err
	}
	pkt, err := NewPacketOutFromPacketType(h.PacketType)
	if err != nil {
		lgr.Errorf("%s error creating packet: %s", lmp, err)
		return h, nil, err
	}

	vs := int(h.Length) - HeaderSize - internal.TotalSizeOfFixedFields(pkt)
	err = internal.UnmarshalLittleEndian(r, pkt, int(h.Length), vs)
	if err != nil {
		lgr.Errorf("%s error reading packet %T data %s", lmp, pkt, err)
		return h, nil, err
	}

	return h, pkt, nil
}

func readMessageRaw(r io.Reader, lmp string) (uint32, []byte, error) {
	var err error

	var l uint32
	lgr.Infof("%s awaiting packet length...", lmp)
	err = binary.Read(r, binary.LittleEndian, &l)
	if err != nil {
		if err == io.EOF {
			lgr.Infof("%s client disconnected", lmp)
		} else {
			lgr.Errorf("%s error reading packet length: %s", lmp, err)
		}
		return l, nil, err
	}

	b := make([]byte, int(l)-4)
	if err := binary.Read(r, binary.LittleEndian, &b); err != nil {
		lgr.Errorf("%s error reading payload: %s", lmp, err)
		return l, nil, err
	}

	return l, b, nil
}

func sendMessage(w io.Writer, pkt Packet, extra []byte, lmp string) {
	err := sendAnyPacket(w, pkt, extra, lmp)
	if err != nil {
		lgr.Errorf("%s error responding: %s", lmp, err)
	}
}

func alwaysFailMessage(conn net.Conn, _ chan uint32, lmp string) {
	// TCP connections are closed by the Responder on failure!
	defer conn.Close()
	if _, pkt, _ := readMessage(conn, lmp); pkt == nil {
		return
	}

	sendMessage(conn, &InitFailPacket{
		Reason: FR_FailRejectedInitiator,
	}, nil, lmp)
}

func sendAnyPacket(w io.Writer, p Packet, extra []byte, lmp string) error {
	lgr.Infof("%s sendAnyPacket() %T", lmp, p)

	pl := internal.MarshalLittleEndian(p)
	pll := len(pl)
	if extra != nil {
		pll += len(extra)
	}

	// An invalid packet type means it does not adhere to the PTP/IP standard, so we only send the length field here.
	if p.PacketType() == PKT_Invalid {
		// Send length only. The length must include the size of the length field, so we add 4 bytes for that!
		if _, err := w.Write(internal.MarshalLittleEndian(uint32(pll + 4))); err != nil {
			return err
		}
	} else {
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
		lgr.Infof("%s sendAnyPacket() header bytes written %d", lmp, n)
	}

	// Send payload.
	if pll == 0 {
		lgr.Infof("%s sendAnyPacket() packet has no payload", lmp)
		return nil
	}

	n, err := w.Write(pl)
	if err != nil {
		return err
	}

	if extra != nil {
		nn, err := w.Write(extra)
		if err != nil {
			return err
		}
		n += nn
	}

	if n != pll {
		return fmt.Errorf(BytesWrittenMismatch.Error(), n, pll)
	}

	lgr.Infof("%s sendAnyPacket() payload bytes written %d", lmp, n)

	return nil
}
