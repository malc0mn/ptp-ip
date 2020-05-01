package ip

import (
	"bytes"
	"fmt"
	"github.com/google/uuid"
	"testing"
)

func TestNewInitiator(t *testing.T) {
	got := NewDefaultInitiator()
	if got.GUID == uuid.Nil {
		t.Errorf("NewInitiator() GUID = %s; want valid non-empty UUID", got.GUID)
	}
	if got.FriendlyName != InitiatorFriendlyName {
		t.Errorf("NewInitiator() Friendlyname = %s; want %s", got.FriendlyName, InitiatorFriendlyName)
	}
}

func TestNewInitiatorWithFriendlyName(t *testing.T) {
	got := NewInitiator("Friendly test", uuid.Nil)
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
	got := NewClient(DefaultIpAddress, DefaultPort, "", uuid.Nil)
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

func TestClient_SendPacket(t *testing.T) {
	guid, _ := uuid.Parse("e462b590-b516-474a-9db8-a465b370fabd")
	c := NewClient(DefaultIpAddress, DefaultPort, "writèr", guid)
	p := NewInitCommandRequestPacketForClient(c)

	want := "[00100011 00000000 00000000 00000000 00000001 00000000 00000000 00000000 11100100 01100010 10110101 10010000 10110101 00010110 01000111 01001010 10011101 10111000 10100100 01100101 10110011 01110000 11111010 10111101 01110111 01110010 01101001 01110100 11000011 10101000 01110010 00000000 00000000 00000001 00000000]"

	var buf bytes.Buffer
	c.sendPacket(&buf, p)
	got := fmt.Sprintf("%.8b", buf.Bytes())

	if got != want {
		t.Errorf("sendPacket() buffer = %s; want %s", got, want)
	}
}

func TestClient_ReadResponse(t *testing.T) {
	guidC, _ := uuid.Parse("d6555687-a599-44b8-a4af-279d599a92f6")
	c := NewClient(DefaultIpAddress, DefaultPort, "writèr", guidC)
	guidR, _ := uuid.Parse("7c946ae4-6d6a-4589-90ed-d059f8cc426b")
	p := &InitCommandAckPacket{uint32(1), guidR, "remôte", uint32(PV_VersionOnePointZero)}

	var b bytes.Buffer
	c.sendPacket(&b, p)

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
		t.Errorf("readResponse() ResponderFriendlyName = %s; want %s", gotName, wantName)
	}

	gotVer := rp.(*InitCommandAckPacket).ResponderProtocolVersion
	wantVer := uint32(PV_VersionOnePointZero)
	if gotVer != wantVer {
		t.Errorf("readResponse() ResponderProtocolVersion = %x; want %x", gotVer, wantVer)
	}
}
