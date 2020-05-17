package ip

import (
	"bytes"
	"fmt"
	"github.com/google/uuid"
	"github.com/malc0mn/ptp-ip/internal"
	ipInternal "github.com/malc0mn/ptp-ip/ip/internal"
	"github.com/malc0mn/ptp-ip/ptp"
	"io"
	"log"
	"os"
	"testing"
)

var (
	address         = "127.0.0.1"
	okPort          = DefaultPort
	failPort uint16 = 25740
)

func TestMain(m *testing.M) {
	go newLocalOkResponder(address, okPort)
	go newLocalFailResponder(address, failPort)
	code := m.Run()
	os.Exit(code)
}

func (c *Client) sendAnyPacket(w io.Writer, p Packet) error {
	log.Printf("[ip_test] sendAnyPacket() sending %T", p)

	pl := ipInternal.MarshalLittleEndian(p)
	pll := len(pl)

	// The packet length MUST include the header, so we add 8 bytes for that!
	h := ipInternal.MarshalLittleEndian(Header{uint32(pll + HeaderSize), p.PacketType()})

	// Send header.
	n, err := w.Write(h)
	internal.FailOnError(err)
	if n != HeaderSize {
		return fmt.Errorf(BytesWrittenMismatch.Error(), n, HeaderSize)
	}
	log.Printf("[ip_test] sendAnyPacket() header bytes written %d", n)

	// Send payload.
	n, err = w.Write(pl)
	if n != pll {
		return fmt.Errorf(BytesWrittenMismatch.Error(), n, pll)
	}
	internal.FailOnError(err)
	log.Printf("[ip_test] sendAnyPacket() payload bytes written %d", n)

	return nil
}

func TestNewDefaultInitiator(t *testing.T) {
	got, err := NewDefaultInitiator()
	if err != nil {
		t.Errorf("NewDefaultInitiator() err = %s; want <nil>", err)
	}
	if got.GUID == uuid.Nil {
		t.Errorf("NewDefaultInitiator() GUID = %s; want valid non-empty UUID", got.GUID)
	}
	if got.FriendlyName != InitiatorFriendlyName {
		t.Errorf("NewDefaultInitiator() Friendlyname = %s; want %s", got.FriendlyName, InitiatorFriendlyName)
	}
}

func TestNewInitiatorWithFriendlyName(t *testing.T) {
	got, err := NewInitiator("Friendly test", "")
	if err != nil {
		t.Errorf("NewInitiator() err = %s; want <nil>", err)
	}
	if got.GUID == uuid.Nil {
		t.Errorf("NewInitiator() GUID = %s; want valid non-empty UUID", got.GUID)
	}
	wantName := "Friendly test"
	if got.FriendlyName != wantName {
		t.Errorf("NewInitiator() Friendlyname = %s; want %s", got.FriendlyName, wantName)
	}
}

func TestNewResponder(t *testing.T) {
	got := NewResponder(DefaultIpAddress, DefaultPort)
	if got.IpAddress != DefaultIpAddress {
		t.Errorf("NewResponder() IpAddress = %s; want %s", got.IpAddress, DefaultIpAddress)
	}
	if got.Port != DefaultPort {
		t.Errorf("NewResponder() IpAddress = %d; want %d", got.Port, DefaultPort)
	}
	if got.GUID != uuid.Nil {
		t.Errorf("NewResponder() FriendlyName = %s; want <nil>", got.GUID)
	}
	if got.FriendlyName != "" {
		t.Errorf("NewResponder() FriendlyName = %s; want <nil>", got.FriendlyName)
	}
}

func TestNewClient(t *testing.T) {
	got, err := NewClient(DefaultIpAddress, DefaultPort, "", "")
	if err != nil {
		t.Errorf("NewClient() err = %s; want <nil>", err)
	}
	if got.commandDataConn != nil {
		t.Errorf("NewClient() commandDataConn = %v; want <nil>", got.commandDataConn)
	}
	if got.eventConn != nil {
		t.Errorf("NewClient() eventConn = %v; want <nil>", got.eventConn)
	}
	if got.initiator == nil {
		t.Errorf("NewClient() initiator = %v; want Initiator", got.initiator)
	}
	if got.responder == nil {
		t.Errorf("NewClient() responder = %v; want Responder", got.responder)
	}
}

func TestClient_incrementTransactionId(t *testing.T) {
	c := Client{}

	got := c.TransactionId()
	want := ptp.TransactionID(0)
	if got != want {
		t.Errorf("TransactionId() = %#x; want %#x", got, want)
	}

	c.incrementTransactionId()
	got = c.TransactionId()
	want = ptp.TransactionID(1)
	if got != want {
		t.Errorf("TransactionId() = %#x; want %#x", got, want)
	}

	c.transactionId = 0xFFFFFFFE
	c.incrementTransactionId()
	got = c.TransactionId()
	want = ptp.TransactionID(1)
	if got != want {
		t.Errorf("TransactionId() = %#x; want %#x", got, want)
	}
}

func TestClient_sendPacket(t *testing.T) {
	c, err := NewClient(DefaultIpAddress, DefaultPort, "writèr", "e462b590-b516-474a-9db8-a465b370fabd")
	if err != nil {
		t.Errorf("sendPacket() err = %s; want <nil>", err)
	}

	p := NewInitCommandRequestPacketForClient(c)

	want := "[00100011 00000000 00000000 00000000 00000001 00000000 00000000 00000000 11100100 01100010 10110101 10010000 10110101 00010110 01000111 01001010 10011101 10111000 10100100 01100101 10110011 01110000 11111010 10111101 01110111 01110010 01101001 01110100 11000011 10101000 01110010 00000000 00000000 00000001 00000000]"

	var buf bytes.Buffer
	c.sendPacket(&buf, p)
	got := fmt.Sprintf("%.8b", buf.Bytes())

	if got != want {
		t.Errorf("sendPacket() buffer = %s; want %s", got, want)
	}
}

func TestClient_readResponse(t *testing.T) {
	c, err := NewClient(DefaultIpAddress, DefaultPort, "writèr", "d6555687-a599-44b8-a4af-279d599a92f6")
	if err != nil {
		t.Errorf("readResponse() err = %s; want <nil>", err)
	}

	guidR, _ := uuid.Parse("7c946ae4-6d6a-4589-90ed-d059f8cc426b")
	p := &InitCommandAckPacket{uint32(1), guidR, "remôte", uint32(0x00020005)}

	var b bytes.Buffer
	c.sendAnyPacket(&b, p)

	rp, err := c.readResponse(&b)
	if err != nil {
		t.Errorf("readResponse() error = %s; want <nil>", err)
	}

	want := "*ip.InitCommandAckPacket"
	if fmt.Sprintf("%T", rp) != want {
		t.Errorf("readResponse() PaketType = %T; want %s", rp, want)
	}

	gotType := rp.PacketType()
	wantType := PKT_InitCommandAck
	if gotType != wantType {
		t.Errorf("readResponse() PaketType = %x; want %x", gotType, wantType)
	}

	gotNum := rp.(*InitCommandAckPacket).ConnectionNumber
	wantNum := uint32(1)
	if gotNum != wantNum {
		t.Errorf("readResponse() ConnectionNumber = %d; want %d", gotNum, wantNum)
	}

	gotGuid := rp.(*InitCommandAckPacket).ResponderGUID
	wantGuid, _ := uuid.Parse("7c946ae4-6d6a-4589-90ed-d059f8cc426b")
	if gotGuid != wantGuid {
		t.Errorf("readResponse() ResponderGUID = %s; want %s", gotGuid, wantGuid)
	}

	gotName := rp.(*InitCommandAckPacket).ResponderFriendlyName
	wantName := "remôte"
	if gotName != wantName {
		t.Errorf("readResponse() ResponderFriendlyName = %s (%#x); want %s (%#x)", gotName, gotName, wantName, wantName)
	}

	gotVer := rp.(*InitCommandAckPacket).ResponderProtocolVersion
	wantVer := uint32(0x00020005)
	if gotVer != wantVer {
		t.Errorf("readResponse() ResponderProtocolVersion = %#x; want %#x", gotVer, wantVer)
	}
}

func TestClient_initCommandDataConn(t *testing.T) {
	c, err := NewClient(address, okPort, "testèr", "67bace55-e7a4-4fbc-8e31-5122ee73a17c")
	defer c.Close()
	if err != nil {
		t.Fatal(err)
	}
	err = c.initCommandDataConn()
	if err != nil {
		t.Errorf("initCommandDataConn() error = %s; want <nil>", err)
	}

	got := c.TransactionId()
	want := ptp.TransactionID(0)
	if got != want {
		t.Errorf("TransactionId() got = %#x; want %#x", got, want)
	}
}
func TestClient_initCommandDataConnFail(t *testing.T) {
	c, err := NewClient(address, failPort, "testér", "b3ca53e9-bb61-4c85-9fcd-3b446a9e81e6")
	defer c.Close()
	if err != nil {
		t.Fatal(err)
	}
	err = c.initCommandDataConn()
	if err == nil {
		t.Errorf("initCommandDataConn() error = %s; want rejected: device not allowed", err)
	}

	got := c.TransactionId()
	want := ptp.TransactionID(0)
	if got != want {
		t.Errorf("TransactionId() got = %#x; want %#x", got, want)
	}
}

func TestClient_initEventConn(t *testing.T) {
	c, err := NewClient(address, okPort, "testèr", "67bace55-e7a4-4fbc-8e31-5122ee73a17c")
	defer c.Close()
	if err != nil {
		t.Fatal(err)
	}
	err = c.initEventConn()
	if err != nil {
		t.Errorf("initEventConn() error = %s; want <nil>", err)
	}

	got := c.TransactionId()
	want := ptp.TransactionID(1)
	if got != want {
		t.Errorf("TransactionId() got = %#x; want %#x", got, want)
	}
}

func TestClient_initEventConnFail(t *testing.T) {
	c, err := NewClient(address, failPort, "testér", "733e8d71-0f05-4aba-9745-ea9294dd2278")
	defer c.Close()
	if err != nil {
		t.Fatal(err)
	}
	err = c.initEventConn()
	if err == nil {
		t.Errorf("initEventConn() error = %s; want rejected: device not allowed", err)
	}

	got := c.TransactionId()
	want := ptp.TransactionID(0)
	if got != want {
		t.Errorf("TransactionId() got = %#x; want %#x", got, want)
	}
}

func TestClient_Dial(t *testing.T) {
	c, err := NewClient(address, okPort, "testèr", "7e5ac7d3-46b7-4c50-b0d9-ba56c0e599f0")
	defer c.Close()
	if err != nil {
		t.Fatal(err)
	}

	err = c.Dial()
	if err != nil {
		t.Errorf("Dial() err = %s; want <nil>", err)
	}

	c, err = NewClient(address, failPort, "testér", "f62b41f8-a094-4dab-b537-99afd04c6024")
	defer c.Close()
	if err != nil {
		t.Fatal(err)
	}

	err = c.Dial()
	if err == nil {
		t.Errorf("Dial() err = %s; want rejected: device not allowed", err)
	}
}

func TestClient_GetDeviceInfo(t *testing.T) {
	c, err := NewClient(address, okPort, "tèster", "558acd44-f794-4b26-9129-d460b2a29e8d")
	defer c.Close()
	if err != nil {
		t.Fatal(err)
	}

	err = c.Dial()
	if err != nil {
		t.Fatal(err)
	}

	got, err := c.GetDeviceInfo()
	if err != nil {
		t.Errorf("GetDeviceInfo() err = %s; want <nil>", err)
	}
	if got == nil {
		t.Errorf("GetDeviceInfo() got = %v; want *ip.OperationResponsePacket", got)
	}
}
