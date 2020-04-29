package ip

import (
	"fmt"
	"testing"
)

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
