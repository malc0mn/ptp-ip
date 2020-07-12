package ip

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/malc0mn/ptp-ip/ptp"
	"testing"
)

func TestNewFujiInitCommandRequestPacket(t *testing.T) {
	uuid, _ := uuid.NewRandom()
	got := NewFujiInitCommandRequestPacket(uuid, "têst")
	want := "têst"

	if got.GetFriendlyName() != want {
		t.Errorf("NewFujiInitCommandRequestPacket() FriendlyName = %s; want %s", got.GetFriendlyName(), want)
	}
	if got.GetProtocolVersion() != PV_Fuji {
		t.Errorf("NewFujiInitCommandRequestPacket() ProtocolVersion = %#x; want %#x", got.GetProtocolVersion(), PV_Fuji)
	}
}

func TestNewFujiInitCommandRequestPacketForClient(t *testing.T) {
	c, err := NewClient("fuji", DefaultIpAddress, DefaultPort, "test", "", LevelDebug)
	if err != nil {
		t.Errorf("NewClient() err = %s; want <nil>", err)
	}

	got := NewFujiInitCommandRequestPacketForClient(c)
	want := "test"

	if got.GetFriendlyName() != want {
		t.Errorf("NewFujiInitCommandRequestPacketForClient() FriendlyName = %s; want %s", got.GetFriendlyName(), want)
	}
	if got.GetProtocolVersion() != PV_Fuji {
		t.Errorf("NewFujiInitCommandRequestPacketForClient() ProtocolVersion = %#x; want %#x", got.GetProtocolVersion(), PV_Fuji)
	}
}

func TestNewFujiInitCommandRequestPacketWithVersion(t *testing.T) {
	uuid, _ := uuid.NewRandom()
	got := NewFujiInitCommandRequestPacketWithVersion(uuid, "versíon", 0x00020005)
	wantName := "versíon"
	wantVersion := ProtocolVersion(0x00020005)

	if got.GetFriendlyName() != wantName {
		t.Errorf("NewFujiInitCommandRequestPacketWithVersion() FriendlyName = %s; want %s", got.GetFriendlyName(), wantName)
	}
	if got.GetProtocolVersion() != wantVersion {
		t.Errorf("NewFujiInitCommandRequestPacketWithVersion() ProtocolVersion = %#x; want %#x", got.GetProtocolVersion(), wantVersion)
	}
}

func TestFujiOperationRequestPacket_Payload(t *testing.T) {
	oreq := &FujiOperationRequestPacket{
		DataPhaseInfo: uint16(DP_NoDataOrDataIn),
		OperationCode: ptp.OC_GetDevicePropValue,
		TransactionID: 1,
		Parameter1: uint32(DPC_Fuji_FilmSimulation),
	}

	pl := oreq.Payload()
	got := fmt.Sprintf("%.8b", pl)
	want := "[00000001 00000000 00010101 00010000 00000001 00000000 00000000 00000000 00000001 11010000 00000000 00000000 00000000 00000000 00000000 00000000 00000000 00000000 00000000 00000000 00000000 00000000 00000000 00000000 00000000 00000000 00000000 00000000]"
	if got != want {
		t.Errorf("payload() buffer = %s; want %s", got, want)
	}
}
