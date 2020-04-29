package ip

import (
	"fmt"
	"github.com/google/uuid"
	"testing"
)

func TestNewInitCommandRequestPacket(t *testing.T) {
	uuid, _ := uuid.NewRandom()
	got := NewInitCommandRequestPacket(uuid, "têst")
	want := "têst"

	if got.FriendlyName != want {
		t.Errorf("NewInitCommandRequestPacket() friendlyName = %s; want %s", got.FriendlyName, want)
	}
	if got.ProtocolVersion != PV_VersionOnePointZero {
		t.Errorf("NewInitCommandRequestPacket() protocolVersion = %x; want %x", got.ProtocolVersion, PV_VersionOnePointZero)
	}
}

func TestNewInitCommandRequestPacketForClient(t *testing.T) {
	c := NewClient(DefaultIpAddress, DefaultPort, "test")
	got := NewInitCommandRequestPacketForClient(c)
	want := "test"

	if got.FriendlyName != want {
		t.Errorf("NewInitCommandRequestPacketForClient() friendlyName = %s; want %s", got.FriendlyName, want)
	}
	if got.ProtocolVersion != PV_VersionOnePointZero {
		t.Errorf("NewInitCommandRequestPacketForClient() protocolVersion = %x; want %x", got.ProtocolVersion, PV_VersionOnePointZero)
	}
}

func TestNewInitCommandRequestPacketWithVersion(t *testing.T) {
	uuid, _ := uuid.NewRandom()
	got := NewInitCommandRequestPacketWithVersion(uuid, "versíon", 0x00020005)
	wantName := "versíon"
	wantVersion := ProtocolVersion(0x00020005)

	if got.FriendlyName != wantName {
		t.Errorf("NewInitCommandRequestPacket() friendlyName = %s; want %s", got.FriendlyName, wantName)
	}
	if got.ProtocolVersion != wantVersion {
		t.Errorf("NewInitCommandRequestPacket() protocolVersion = %x; want %x", got.ProtocolVersion, wantVersion)
	}
}

func TestNewPacketFromPacketType(t *testing.T) {
	types := map[PacketType]string{
		PKT_InitCommandRequest: "InitCommandRequest",
		PKT_InitCommandAck:     "InitCommandAck",
		PKT_InitEventRequest:   "InitEventRequest",
		PKT_InitEventAck:       "InitEventAck",
		PKT_InitFail:           "InitFail",
		PKT_OperationRequest:   "OperationRequest",
		PKT_OperationResponse:  "OperationResponse",
		PKT_Event:              "Event",
		PKT_StartData:          "StartData",
		PKT_Data:               "Data",
		PKT_Cancel:             "Cancel",
		PKT_EndData:            "EndData",
		PKT_ProbeRequest:       "ProbeRequest",
		PKT_ProbeResponse:      "ProbeResponse",
	}

	for typ, want := range types {
		got, err := NewPacketFromPacketType(typ)
		want := fmt.Sprintf("*ip.%sPacket", want)
		if err != nil {
			t.Errorf("NewPacketFromPacketType() err = %s; want <nil>", err)
		}

		name := fmt.Sprintf("%T", got)
		if name != want {
			t.Errorf("NewPacketFromPacketType() returned %s; want %s", name, want)
		}

		if got.PacketType() != typ {
			t.Errorf("NewPacketFromPacketType() type = %x; want %x", got.PacketType(), typ)
		}
	}
}
