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
		t.Errorf("NewInitCommandRequestPacket() protocolVersion = %#x; want %#x", got.ProtocolVersion, PV_VersionOnePointZero)
	}
}

func TestNewInitCommandRequestPacketForClient(t *testing.T) {
	c, err := NewClient(DefaultIpAddress, DefaultPort, "test", uuid.Nil)
	if err != nil {
		t.Errorf("NewInitCommandRequestPacketForClient() err = %s; want <nil>", err)
	}

	got := NewInitCommandRequestPacketForClient(c)
	want := "test"

	if got.FriendlyName != want {
		t.Errorf("NewInitCommandRequestPacketForClient() friendlyName = %s; want %s", got.FriendlyName, want)
	}
	if got.ProtocolVersion != PV_VersionOnePointZero {
		t.Errorf("NewInitCommandRequestPacketForClient() protocolVersion = %#x; want %#x", got.ProtocolVersion, PV_VersionOnePointZero)
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
		t.Errorf("NewInitCommandRequestPacket() protocolVersion = %#x; want %#x", got.ProtocolVersion, wantVersion)
	}
}

func TestNewPacketOutFromPacketType(t *testing.T) {
	types := map[PacketType]string{
		PKT_InitCommandRequest: "InitCommandRequest",
		PKT_InitEventRequest:   "InitEventRequest",
		PKT_OperationRequest:   "OperationRequest",
		PKT_StartData:          "StartData",
		PKT_Data:               "Data",
		PKT_Cancel:             "Cancel",
		PKT_EndData:            "EndData",
		PKT_ProbeRequest:       "ProbeRequest",
		PKT_ProbeResponse:      "ProbeResponse",
	}

	for typ, want := range types {
		got, err := NewPacketOutFromPacketType(typ)
		want := fmt.Sprintf("*ip.%sPacket", want)
		if err != nil {
			t.Errorf("NewPacketOutFromPacketType() err = %s; want <nil>", err)
		}

		name := fmt.Sprintf("%T", got)
		if name != want {
			t.Errorf("NewPacketOutFromPacketType() returned %s; want %s", name, want)
		}

		if got.PacketType() != typ {
			t.Errorf("NewPacketOutFromPacketType() type = %#x; want %#x", got.PacketType(), typ)
		}
	}
}

func TestNewPacketOutFromPacketTypeFail(t *testing.T) {
	types := []PacketType{
		PKT_InitCommandAck,
		PKT_InitEventAck,
		PKT_InitFail,
		PKT_OperationResponse,
		PKT_Event,
	}
	for _, typ := range types {
		got, err := NewPacketOutFromPacketType(typ)
		if err == nil {
			t.Errorf("NewPacketOutFromPacketType() err = %s; want unknown packet type %#x", err, typ)
		}
		if got != nil {
			t.Errorf("NewPacketOutFromPacketType() got = %T; want <nil>", got)
		}
	}
}

func TestNewPacketInFromPacketType(t *testing.T) {
	types := map[PacketType]string{
		PKT_InitCommandAck:    "InitCommandAck",
		PKT_InitEventAck:      "InitEventAck",
		PKT_InitFail:          "InitFail",
		PKT_OperationResponse: "OperationResponse",
		PKT_Event:             "Event",
		PKT_StartData:         "StartData",
		PKT_Data:              "Data",
		PKT_Cancel:            "Cancel",
		PKT_EndData:           "EndData",
		PKT_ProbeRequest:      "ProbeRequest",
		PKT_ProbeResponse:     "ProbeResponse",
	}

	for typ, want := range types {
		got, err := NewPacketInFromPacketType(typ)
		want := fmt.Sprintf("*ip.%sPacket", want)
		if err != nil {
			t.Errorf("NewPacketInFromPacketType() err = %s; want <nil>", err)
		}

		name := fmt.Sprintf("%T", got)
		if name != want {
			t.Errorf("NewPacketInFromPacketType() returned %s; want %s", name, want)
		}

		if got.PacketType() != typ {
			t.Errorf("NewPacketInFromPacketType() type = %#x; want %#x", got.PacketType(), typ)
		}
	}
}

func TestNewPacketInFromPacketTypeFail(t *testing.T) {
	types := []PacketType{
		PKT_InitCommandRequest,
		PKT_InitEventRequest,
		PKT_OperationRequest,
	}
	for _, typ := range types {
		got, err := NewPacketInFromPacketType(typ)
		if err == nil {
			t.Errorf("NewPacketInFromPacketType() err = %s; want unknown packet type %#x", err, typ)
		}
		if got != nil {
			t.Errorf("NewPacketInFromPacketType() got = %T; want <nil>", got)
		}
	}
}
