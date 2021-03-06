package ip

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/malc0mn/ptp-ip/ptp"
	"testing"
)

func TestNewInitCommandRequestPacket(t *testing.T) {
	uuid, _ := uuid.NewRandom()
	got := NewInitCommandRequestPacket(uuid, "têst")
	want := "têst"

	if got.GetFriendlyName() != want {
		t.Errorf("NewInitCommandRequestPacket() FriendlyName = %s; want %s", got.GetFriendlyName(), want)
	}
	if got.GetProtocolVersion() != PV_VersionOnePointZero {
		t.Errorf("NewInitCommandRequestPacket() ProtocolVersion = %#x; want %#x", got.GetProtocolVersion(), PV_VersionOnePointZero)
	}
}

func TestNewInitCommandRequestPacketForClient(t *testing.T) {
	c, err := NewClient(DefaultVendor, DefaultIpAddress, DefaultPort, "test", "", logLevel)
	if err != nil {
		t.Errorf("NewClient() err = %s; want <nil>", err)
	}

	got := NewInitCommandRequestPacketForClient(c)
	want := "test"

	if got.GetFriendlyName() != want {
		t.Errorf("NewInitCommandRequestPacketForClient() FriendlyName = %s; want %s", got.GetFriendlyName(), want)
	}
	if got.GetProtocolVersion() != PV_VersionOnePointZero {
		t.Errorf("NewInitCommandRequestPacketForClient() ProtocolVersion = %#x; want %#x", got.GetProtocolVersion(), PV_VersionOnePointZero)
	}
}

func TestNewInitCommandRequestPacketWithVersion(t *testing.T) {
	uuid, _ := uuid.NewRandom()
	got := NewInitCommandRequestPacketWithVersion(uuid, "versíon", 0x00020005)
	wantName := "versíon"
	wantVersion := ProtocolVersion(0x00020005)

	if got.GetFriendlyName() != wantName {
		t.Errorf("NewInitCommandRequestPacketWithVersion() FriendlyName = %s; want %s", got.GetFriendlyName(), wantName)
	}
	if got.GetProtocolVersion() != wantVersion {
		t.Errorf("NewInitCommandRequestPacketWithVersion() ProtocolVersion = %#x; want %#x", got.GetProtocolVersion(), wantVersion)
	}
}

func TestNewInitEventRequestPacket(t *testing.T) {
	got := NewInitEventRequestPacket(5)
	want := uint32(5)
	if got.GetConnectionNumber() != want {
		t.Errorf("NewInitEventRequestPacket() ConnectionNumber = %#x; want %#x", got.GetConnectionNumber(), want)
	}
}

func TestInitFailPacket_ReasonAsError(t *testing.T) {
	errs := map[FailReason]string{
		FR_FailBusy:              "busy: too many active connections",
		FR_FailRejectedInitiator: "rejected: device not allowed",
		FR_FailUnspecified:       "reason unspecified",
		FR_Fuji_DeviceBusy:       "fuji: invalid friendly name or camera state: allow to 'change' client or 'reset' connection",
		FR_Fuji_InvalidParameter: "fuji: unknown protocol version",
		FailReason(0x5032000):    "unknown failure reason returned 0x5032000",
	}

	for reason, want := range errs {
		ifp := InitFailPacket{
			Reason: reason,
		}
		got := ifp.ReasonAsError().Error()
		if got != want {
			t.Errorf("ReasonAsError() Reason = '%s'; want '%s'", got, want)
		}
	}
}

func TestNewPacketOutFromPacketType(t *testing.T) {
	types := map[PacketType]string{
		PKT_InitCommandRequest: "GenericInitCommandRequest",
		PKT_InitEventRequest:   "GenericInitEventRequest",
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
		PKT_Event:             "GenericEvent",
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

func TestOperationRequestPacket_Payload(t *testing.T) {
	oreq := &OperationRequestPacket{
		DataPhaseInfo:    DP_NoDataOrDataIn,
		OperationRequest: ptp.GetDeviceInfo(2),
	}

	pl := oreq.Payload()
	got := fmt.Sprintf("%.8b", pl)
	want := "[00000001 00000000 00000000 00000000 00000001 00010000 00000010 00000000 00000000 00000000 00000000 00000000 00000000 00000000 00000000 00000000 00000000 00000000 00000000 00000000 00000000 00000000 00000000 00000000 00000000 00000000 00000000 00000000 00000000 00000000]"
	if got != want {
		t.Errorf("payload() buffer = %s; want %s", got, want)
	}
}
