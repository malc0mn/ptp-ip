package ip

import (
	"bytes"
	"fmt"
	"github.com/google/uuid"
	"testing"
)

func TestNewInitiator(t *testing.T) {
	got := NewInitiator("")
	if got.GUID == uuid.Nil {
		t.Errorf("NewInitiator() GUID = %s; want valid non-empty UUID", got.GUID)
	}
	if got.FriendlyName != InitiatorFriendlyName {
		t.Errorf("NewInitiator() Friendlyname = %s; want %s", got.FriendlyName, InitiatorFriendlyName)
	}
}

func TestNewInitiatorWithFriendlyName(t *testing.T) {
	got := NewInitiator("Friendly test")
	if got.GUID == uuid.Nil {
		t.Errorf("NewInitiator() GUID = %s; want valid non-empty UUID", got.GUID)
	}
	if got.FriendlyName != "Friendly test" {
		t.Errorf("NewInitiator() Friendlyname = %s; want %s", got.FriendlyName, "Friendly test")
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
	got := NewClient(DefaultIpAddress, DefaultPort, "")
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
	c := NewClient(DefaultIpAddress, DefaultPort, "writèr")
	p := NewInitCommandRequestPacketForClient(c)
	want := "240000000100000077000000720000006900000074000000e80000007200000000000000"

	var buf bytes.Buffer
	c.sendPacket(&buf, p)
	got := fmt.Sprintf("%x", buf.Bytes())

	if got != want {
		t.Errorf("sendPacket() buffer = %s; want %s", got, want)
	}
}

func TestClient_ReadResponse(t *testing.T) {
	c := NewClient(DefaultIpAddress, DefaultPort, "writèr")
	guid, _ := uuid.NewRandom()
	p := &InitCommandAckPacket{1, guid, "remote", uint32(PV_VersionOnePointZero)}

	var buf bytes.Buffer
	c.sendPacket(&buf, p)

	rp, err := c.readResponse(&buf)
	if err != nil {
		t.Errorf("readResponse() error = %s; want <nil>", err)
	}
	fmt.Printf("%T", rp)
}
