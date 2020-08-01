package ip

import (
	"bytes"
	"fmt"
	"github.com/google/uuid"
	"github.com/malc0mn/ptp-ip/ptp"
	"io/ioutil"
	"reflect"
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
	c, err := NewClient("fuji", DefaultIpAddress, DefaultPort, "test", "", logLevel)
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
		Parameter1:    uint32(DPC_Fuji_FilmSimulation),
	}

	pl := oreq.Payload()
	got := fmt.Sprintf("%.8b", pl)
	want := "[00000001 00000000 00010101 00010000 00000001 00000000 00000000 00000000 00000001 11010000 00000000 00000000 00000000 00000000 00000000 00000000 00000000 00000000 00000000 00000000 00000000 00000000 00000000 00000000 00000000 00000000 00000000 00000000]"
	if got != want {
		t.Errorf("payload() buffer = %s; want %s", got, want)
	}
}

func TestFujiInitCommandDataConn(t *testing.T) {
	c, err := NewClient("fuji", address, fujiCmdPort, "testèr", "67bace55-e7a4-4fbc-8e31-5122ee73a17c", logLevel)
	defer c.Close()
	if err != nil {
		t.Fatal(err)
	}
	err = FujiInitCommandDataConn(c)
	if err != nil {
		t.Errorf("FujiInitCommandDataConn() error = %s; want <nil>", err)
	}

	got := c.TransactionId()
	want := ptp.TransactionID(5)
	if got != want {
		t.Errorf("TransactionId() got = %#x; want %#x", got, want)
	}
}

func TestFujiSetDeviceProperty(t *testing.T) {
	c, err := NewClient("fuji", address, fujiCmdPort, "testèr", "67bace55-e7a4-4fbc-8e31-5122ee73a17c", logLevel)
	defer c.Close()
	if err != nil {
		t.Fatal(err)
	}

	err = c.Dial()
	if err != nil {
		t.Fatal(err)
	}

	err = FujiSetDeviceProperty(c, DPC_Fuji_FilmSimulation, uint32(FS_Fuji_Astia))
	if err != nil {
		t.Errorf("FujiSetDeviceProperty() error = %s; want <nil>", err)
	}

	got := c.TransactionId()
	want := ptp.TransactionID(6)
	if got != want {
		t.Errorf("TransactionId() got = %#x; want %#x", got, want)
	}
}

func TestFujiSetDevicePropertyFail(t *testing.T) {
	c, err := NewClient("fuji", address, fujiCmdPort, "testèr", "67bace55-e7a4-4fbc-8e31-5122ee73a17c", logLevel)
	defer c.Close()
	if err != nil {
		t.Fatal(err)
	}
	err = FujiSetDeviceProperty(c, DPC_Fuji_FilmSimulation, uint32(FS_Fuji_Astia))
	want := "not connected"
	if err.Error() != want {
		t.Errorf("FujiSetDeviceProperty() error = %s; want %s", err, want)
	}
}

func TestFujiGetEndOfDataPacket(t *testing.T) {
	c, err := NewClient("fuji", address, fujiCmdPort, "testèr", "67bace55-e7a4-4fbc-8e31-5122ee73a17c", logLevel)
	defer c.Close()
	if err != nil {
		t.Fatal(err)
	}

	got, err := FujiGetEndOfDataPacket(c, &FujiOperationResponsePacket{
		DataPhase:             uint16(DP_NoDataOrDataIn),
		OperationResponseCode: RC_Fuji_GetDeviceInfo,
		TransactionID:         10,
	})
	if err != nil {
		t.Errorf("FujiGetEndOfDataPacket() error = %s; want <nil>", err)
	}

	var want *FujiOperationResponsePacket
	if got != want {
		t.Errorf("FujiGetEndOfDataPacket() got = %#v; want %#v", got, want)
	}

	got, err = FujiGetEndOfDataPacket(c, &FujiOperationResponsePacket{
		DataPhase:             uint16(DP_Unknown),
		OperationResponseCode: RC_Fuji_GetDeviceInfo,
		TransactionID:         10,
	})
	if err != nil {
		t.Errorf("FujiGetEndOfDataPacket() error = %s; want <nil>", err)
	}

	if got != want {
		t.Errorf("FujiGetEndOfDataPacket() got = %#v; want %#v", got, want)
	}

	// TODO: how to actually test DP_DataOut here properly so we can drop the two tests above?
}

func TestFujiGetDevicePropertyValue(t *testing.T) {
	c, err := NewClient("fuji", address, fujiCmdPort, "testèr", "67bace55-e7a4-4fbc-8e31-5122ee73a17c", logLevel)
	defer c.Close()
	if err != nil {
		t.Fatal(err)
	}

	err = c.Dial()
	if err != nil {
		t.Fatal(err)
	}

	got, err := FujiGetDevicePropertyValue(c, DPC_Fuji_AppVersion)
	if err != nil {
		t.Errorf("FujiGetDevicePropertyValue() error = %s; want <nil>", err)
	}

	want := uint32(PM_Fuji_AppVersion)
	if got != want {
		t.Errorf("FujiGetDevicePropertyValue() got = %#x; want %#x", got, want)
	}
}

func TestFujiSendOperationRequest(t *testing.T) {
	c, err := NewClient("fuji", address, fujiCmdPort, "testèr", "67bace55-e7a4-4fbc-8e31-5122ee73a17c", logLevel)
	defer c.Close()
	if err != nil {
		t.Fatal(err)
	}

	err = c.Dial()
	if err != nil {
		t.Fatal(err)
	}

	// We use close session here because our fuji mock will not respond to it.
	err = FujiSendOperationRequest(c, ptp.OC_CloseSession, PM_Fuji_NoParam)
	if err != nil {
		t.Errorf("FujiSendOperationRequest() error = %s; want <nil>", err)
	}
}

func TestFujiSendOperationRequestAndGetResponse(t *testing.T) {
	c, err := NewClient("fuji", address, fujiCmdPort, "testèr", "67bace55-e7a4-4fbc-8e31-5122ee73a17c", logLevel)
	defer c.Close()
	if err != nil {
		t.Fatal(err)
	}

	err = c.Dial()
	if err != nil {
		t.Fatal(err)
	}

	gotPar, gotPkt, xs, err := FujiSendOperationRequestAndGetResponse(c, ptp.OC_GetDevicePropValue, uint32(DPC_Fuji_AppVersion), 4)
	if len(xs) > 0 {
		t.Errorf("FujiSendOperationRequestAndGetResponse() excess bytes = %d; want <nil>", len(xs))
	}
	if err != nil {
		t.Errorf("FujiSendOperationRequestAndGetResponse() error = %s; want <nil>", err)
	}

	wantPkt := "*ip.FujiOperationResponsePacket"
	if fmt.Sprintf("%T", gotPkt) != wantPkt {
		t.Errorf("FujiSendOperationRequestAndGetResponse() got = %T; want %s", gotPkt, wantPkt)
	}

	wantPar := uint32(PM_Fuji_AppVersion)
	if gotPar != wantPar {
		t.Errorf("FujiSendOperationRequestAndGetResponse() got = %#x; want %#x", gotPar, wantPar)
	}
}

func TestFujiSendOperationRequestAndGetRawResponse(t *testing.T) {
	c, err := NewClient("fuji", address, fujiCmdPort, "testèr", "67bace55-e7a4-4fbc-8e31-5122ee73a17c", logLevel)
	defer c.Close()
	if err != nil {
		t.Fatal(err)
	}

	err = c.Dial()
	if err != nil {
		t.Fatal(err)
	}

	got, err := FujiSendOperationRequestAndGetRawResponse(c, ptp.OC_GetDevicePropDesc, []uint32{uint32(DPC_Fuji_FilmSimulation)})
	if err != nil {
		t.Errorf("FujiSendOperationRequestAndGetRawResponse() error = %s; want <nil>", err)
	}

	want := [][]byte{
		// Raw response.
		{0x2e, 0x00, 0x00, 0x00, 0x02, 0x00, 0x14, 0x10, 0x06, 0x00, 0x00, 0x00, 0x01, 0xd0, 0x04, 0x00, 0x01, 0x01, 0x00, 0x01, 0x00, 0x02, 0x0b, 0x00, 0x01, 0x00, 0x02, 0x00, 0x03, 0x00, 0x04, 0x00, 0x05, 0x00, 0x06, 0x00, 0x07, 0x00, 0x08, 0x00, 0x09, 0x00, 0x0a, 0x00, 0x0b, 0x00},
		// Raw end of data packet.
		{0xc, 0x0, 0x0, 0x0, 0x3, 0x0, 0x1, 0x20, 0x6, 0x0, 0x0, 0x0},
	}
	for i, g := range got {
		if bytes.Compare(g, want[i]) != 0 {
			t.Errorf("FujiSendOperationRequestAndGetRawResponse() got = %#v; want %#v", got, want)
			break
		}
	}
}

func TestFujiGetDevicePropertyDesc(t *testing.T) {
	c, err := NewClient("fuji", address, fujiCmdPort, "testèr", "67bace55-e7a4-4fbc-8e31-5122ee73a17c", logLevel)
	defer c.Close()
	if err != nil {
		t.Fatal(err)
	}

	err = c.Dial()
	if err != nil {
		t.Fatal(err)
	}

	got, err := FujiGetDevicePropertyDesc(c, ptp.DPC_WhiteBalance)
	if err != nil {
		t.Errorf("FujiGetDevicePropertyDesc() error = %s; want <nil>", err)
	}

	want := &ptp.DevicePropDesc{
		DevicePropertyCode:  ptp.DPC_WhiteBalance,
		DataType:            ptp.DTC_UINT16,
		GetSet:              ptp.DPD_GetSet,
		FactoryDefaultValue: []uint8{0x2, 0x0},
		CurrentValue:        []uint8{0x2, 0x0},
		FormFlag:            ptp.DPF_FormFlag_Enum,
		Form: &ptp.EnumerationForm{
			NumberOfValues: 10,
			SupportedValues: [][]uint8{
				{0x02, 0x00}, {0x04, 0x00}, {0x06, 0x80}, {0x01, 0x80}, {0x02, 0x80}, {0x03, 0x80}, {0x06, 0x00},
				{0x0a, 0x80}, {0x0b, 0x80}, {0x0c, 0x80},
			},
		},
	}
	want.Form.SetDevicePropDesc(want)

	if !reflect.DeepEqual(got, want) {
		t.Errorf("FujiGetDevicePropertyDesc() got = %#v; want %#v", got, want)
	}
}

func TestFujiGetDeviceInfo(t *testing.T) {
	c, err := NewClient("fuji", address, fujiCmdPort, "testèr", "67bace55-e7a4-4fbc-8e31-5122ee73a17c", logLevel)
	defer c.Close()
	if err != nil {
		t.Fatal(err)
	}

	err = c.Dial()
	if err != nil {
		t.Fatal(err)
	}

	got, err := FujiGetDeviceInfo(c)
	if err != nil {
		t.Errorf("FujiGetDeviceInfo() error = %s; want <nil>", err)
	}

	want := []*ptp.DevicePropDesc{
		{
			DevicePropertyCode:  ptp.DPC_CaptureDelay,
			DataType:            ptp.DTC_UINT16,
			GetSet:              ptp.DPD_GetSet,
			FactoryDefaultValue: []uint8{0x0, 0x0},
			CurrentValue:        []uint8{0x0, 0x0},
			FormFlag:            ptp.DPF_FormFlag_Enum,
			Form: &ptp.EnumerationForm{
				NumberOfValues: 3,
				SupportedValues: [][]uint8{
					{0x00, 0x00}, {0x02, 0x00}, {0x04, 0x00},
				},
			},
		},
		{
			DevicePropertyCode:  ptp.DPC_FlashMode,
			DataType:            ptp.DTC_UINT16,
			GetSet:              ptp.DPD_GetSet,
			FactoryDefaultValue: []uint8{0x2, 0x0},
			CurrentValue:        []uint8{0x9, 0x80},
			FormFlag:            ptp.DPF_FormFlag_Enum,
			Form: &ptp.EnumerationForm{
				NumberOfValues: 2,
				SupportedValues: [][]uint8{
					{0x09, 0x80}, {0x0a, 0x80},
				},
			},
		},
		{
			DevicePropertyCode:  ptp.DPC_WhiteBalance,
			DataType:            ptp.DTC_UINT16,
			GetSet:              ptp.DPD_GetSet,
			FactoryDefaultValue: []uint8{0x2, 0x0},
			CurrentValue:        []uint8{0x2, 0x0},
			FormFlag:            ptp.DPF_FormFlag_Enum,
			Form: &ptp.EnumerationForm{
				NumberOfValues: 10,
				SupportedValues: [][]uint8{
					{0x02, 0x00}, {0x04, 0x00}, {0x06, 0x80}, {0x01, 0x80}, {0x02, 0x80}, {0x03, 0x80}, {0x06, 0x00},
					{0x0a, 0x80}, {0x0b, 0x80}, {0x0c, 0x80},
				},
			},
		},
		{
			DevicePropertyCode:  ptp.DPC_ExposureBiasCompensation,
			DataType:            ptp.DTC_INT16,
			GetSet:              ptp.DPD_GetSet,
			FactoryDefaultValue: []uint8{0x0, 0x0},
			CurrentValue:        []uint8{0x0, 0x0},
			FormFlag:            ptp.DPF_FormFlag_Enum,
			Form: &ptp.EnumerationForm{
				NumberOfValues: 19,
				SupportedValues: [][]uint8{
					{0x48, 0xf4}, {0x95, 0xf5}, {0xe3, 0xf6}, {0x30, 0xf8}, {0x7d, 0xf9}, {0xcb, 0xfa}, {0x18, 0xfc},
					{0x65, 0xfd}, {0xb3, 0xfe}, {0x00, 0x00}, {0x4d, 0x01}, {0x9b, 0x02}, {0xe8, 0x03}, {0x35, 0x05},
					{0x83, 0x06}, {0xd0, 0x07}, {0x1d, 0x09}, {0x6b, 0x0a}, {0xb8, 0x0b},
				},
			},
		},
		{
			DevicePropertyCode:  DPC_Fuji_FilmSimulation,
			DataType:            ptp.DTC_UINT16,
			GetSet:              ptp.DPD_GetSet,
			FactoryDefaultValue: []uint8{0x1, 0x0},
			CurrentValue:        []uint8{0x2, 0x0},
			FormFlag:            ptp.DPF_FormFlag_Enum,
			Form: &ptp.EnumerationForm{
				NumberOfValues: 11,
				SupportedValues: [][]uint8{
					{0x01, 0x00}, {0x02, 0x00}, {0x03, 0x00}, {0x04, 0x00}, {0x05, 0x00}, {0x06, 0x00}, {0x07, 0x00}, {0x08, 0x00},
					{0x09, 0x00}, {0x0a, 0x00}, {0x0b, 0x00},
				},
			},
		},
		{
			DevicePropertyCode:  DPC_Fuji_ExposureIndex,
			DataType:            ptp.DTC_UINT32,
			GetSet:              ptp.DPD_GetSet,
			FactoryDefaultValue: []uint8{0xff, 0xff, 0xff, 0xff},
			CurrentValue:        []uint8{0x0, 0x19, 0x0, 0x80},
			FormFlag:            ptp.DPF_FormFlag_Enum,
			Form: &ptp.EnumerationForm{
				NumberOfValues: 25,
				SupportedValues: [][]uint8{
					{0x90, 0x01, 0x00, 0x80}, {0x20, 0x03, 0x00, 0x80}, {0x40, 0x06, 0x00, 0x80}, {0x80, 0x0c, 0x00, 0x80},
					{0x00, 0x19, 0x00, 0x80}, {0x64, 0x00, 0x00, 0x40}, {0xc8, 0x00, 0x00, 0x00}, {0xfa, 0x00, 0x00, 0x00},
					{0x40, 0x01, 0x00, 0x00}, {0x90, 0x01, 0x00, 0x00}, {0xf4, 0x01, 0x00, 0x00}, {0x80, 0x02, 0x00, 0x00},
					{0x20, 0x03, 0x00, 0x00}, {0xe8, 0x03, 0x00, 0x00}, {0xe2, 0x04, 0x00, 0x00}, {0x40, 0x06, 0x00, 0x00},
					{0xd0, 0x07, 0x00, 0x00}, {0xc4, 0x09, 0x00, 0x00}, {0x80, 0x0c, 0x00, 0x00}, {0xa0, 0x0f, 0x00, 0x00},
					{0x88, 0x13, 0x00, 0x00}, {0x00, 0x19, 0x00, 0x00}, {0x00, 0x32, 0x00, 0x40}, {0x00, 0x64, 0x00, 0x40},
					{0x00, 0xc8, 0x00, 0x40},
				},
			},
		},
		{
			DevicePropertyCode:  DPC_Fuji_RecMode,
			DataType:            ptp.DTC_UINT16,
			GetSet:              ptp.DPD_GetSet,
			FactoryDefaultValue: []uint8{0x1, 0x0},
			CurrentValue:        []uint8{0x1, 0x0},
			FormFlag:            ptp.DPF_FormFlag_Enum,
			Form: &ptp.EnumerationForm{
				NumberOfValues: 2,
				SupportedValues: [][]uint8{
					{0x0, 0x0}, {0x1, 0x0},
				},
			},
		},
		{
			DevicePropertyCode:  DPC_Fuji_FocusMeteringMode,
			DataType:            ptp.DTC_UINT32,
			GetSet:              ptp.DPD_GetSet,
			FactoryDefaultValue: []uint8{0x0, 0x0, 0x0, 0x0},
			CurrentValue:        []uint8{0x2, 0x7, 0x2, 0x3},
			FormFlag:            ptp.DPF_FormFlag_Range,
			Form: &ptp.RangeForm{
				MinimumValue: []uint8{0x00, 0x00, 0x00, 0x00},
				MaximumValue: []uint8{0x07, 0x07, 0x09, 0x10},
				StepSize:     []uint8{0x01, 0x00, 0x00, 0x00},
			},
		},
	}

	for _, f := range want {
		f.Form.SetDevicePropDesc(f)
	}

	for i, g := range got.([]*ptp.DevicePropDesc) {
		if !reflect.DeepEqual(g, want[i]) {
			t.Errorf("FujiGetDeviceInfo() got = %#v; want %#v", got, want)
			break
		}
	}
}

func TestFujiGetDeviceState(t *testing.T) {
	c, err := NewClient("fuji", address, fujiCmdPort, "testèr", "67bace55-e7a4-4fbc-8e31-5122ee73a17c", logLevel)
	defer c.Close()
	if err != nil {
		t.Fatal(err)
	}

	err = c.Dial()
	if err != nil {
		t.Fatal(err)
	}

	got, err := FujiGetDeviceState(c)
	if err != nil {
		t.Errorf("FujiGetDeviceState() error = %s; want <nil>", err)
	}

	want := []*ptp.DevicePropDesc{
		{DevicePropertyCode: ptp.DPC_BatteryLevel, DataType: ptp.DTC_UINT32, CurrentValue: []uint8{0x02, 0x00, 0x00, 0x00}},
		{DevicePropertyCode: DPC_Fuji_ImageAspectRatio, DataType: ptp.DTC_UINT32, CurrentValue: []uint8{0x0a, 0x00, 0x00, 0x00}},
		{DevicePropertyCode: ptp.DPC_WhiteBalance, DataType: ptp.DTC_UINT32, CurrentValue: []uint8{0x02, 0x00, 0x00, 0x00}},
		{DevicePropertyCode: ptp.DPC_FocusMode, DataType: ptp.DTC_UINT32, CurrentValue: []uint8{0x01, 0x80, 0x00, 0x00}},
		{DevicePropertyCode: ptp.DPC_FlashMode, DataType: ptp.DTC_UINT32, CurrentValue: []uint8{0x0a, 0x80, 0x00, 0x00}},
		{DevicePropertyCode: ptp.DPC_ExposureProgramMode, DataType: ptp.DTC_UINT32, CurrentValue: []uint8{0x02, 0x00, 0x00, 0x00}},
		{DevicePropertyCode: ptp.DPC_ExposureBiasCompensation, DataType: ptp.DTC_UINT32, CurrentValue: []uint8{0xb3, 0xfe, 0x00, 0x00}},
		{DevicePropertyCode: ptp.DPC_CaptureDelay, DataType: ptp.DTC_UINT32, CurrentValue: []uint8{0x00, 0x00, 0x00, 0x00}},
		{DevicePropertyCode: DPC_Fuji_FilmSimulation, DataType: ptp.DTC_UINT32, CurrentValue: []uint8{0x02, 0x00, 0x00, 0x00}},
		{DevicePropertyCode: DPC_Fuji_ImageQuality, DataType: ptp.DTC_UINT32, CurrentValue: []uint8{0x04, 0x00, 0x00, 0x00}},
		{DevicePropertyCode: DPC_Fuji_CommandDialMode, DataType: ptp.DTC_UINT32, CurrentValue: []uint8{0x00, 0x00, 0x00, 0x00}},
		{DevicePropertyCode: DPC_Fuji_ExposureIndex, DataType: ptp.DTC_UINT32, CurrentValue: []uint8{0x00, 0x19, 0x00, 0x80}},
		{DevicePropertyCode: DPC_Fuji_FocusMeteringMode, DataType: ptp.DTC_UINT32, CurrentValue: []uint8{0x02, 0x07, 0x02, 0x03}},
		{DevicePropertyCode: DPC_Fuji_FocusLock, DataType: ptp.DTC_UINT32, CurrentValue: []uint8{0x00, 0x00, 0x00, 0x00}},
		{DevicePropertyCode: DPC_Fuji_DeviceError, DataType: ptp.DTC_UINT32, CurrentValue: []uint8{0x00, 0x00, 0x00, 0x00}},
		{DevicePropertyCode: DPC_Fuji_ImageSpaceSD, DataType: ptp.DTC_UINT32, CurrentValue: []uint8{0xd6, 0x05, 0x00, 0x00}},
		{DevicePropertyCode: DPC_Fuji_MovieRemainingTime, DataType: ptp.DTC_UINT32, CurrentValue: []uint8{0x8f, 0x06, 0x00, 0x00}},
	}

	for i, g := range got.([]*ptp.DevicePropDesc) {
		if !reflect.DeepEqual(g, want[i]) {
			t.Errorf("FujiGetDeviceState() got = %#v; want %#v", got, want)
			break
		}
	}
}

func TestFujiInitiateCapture(t *testing.T) {
	c, err := NewClient("fuji", address, fujiCmdPort, "testèr", "67bace55-e7a4-4fbc-8e31-5122ee73a17c", logLevel)
	c.SetEventPort(fujiEvtPort)

	defer c.Close()
	if err != nil {
		t.Fatal(err)
	}

	err = c.Dial()
	if err != nil {
		t.Fatal(err)
	}

	want, err := ioutil.ReadFile("testdata/preview.jpg")
	if err != nil {
		t.Fatal(err)
	}
	got, err := FujiInitiateCapture(c)
	if err != nil {
		t.Errorf("FujiInitiateCapture() error = %s; want <nil>", err)
	}

	if bytes.Compare(got, want) != 0 {
		t.Errorf("FujiInitiateCapture() imgdata = %#v; want %#v", got, want)
	}
}
