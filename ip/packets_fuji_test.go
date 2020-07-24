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

func TestFujiDevicePropCodeAsString(t *testing.T) {
	check := map[ptp.DevicePropCode]string{
		DPC_Fuji_FilmSimulation:     "film simulation",
		DPC_Fuji_ImageQuality:       "image quality",
		DPC_Fuji_RecMode:            "rec mode",
		DPC_Fuji_CommandDialMode:    "command dial mode",
		DPC_Fuji_ExposureIndex:      "ISO",
		DPC_Fuji_MovieISO:           "movie ISO",
		DPC_Fuji_FocusMeteringMode:  "focus point",
		DPC_Fuji_FocusLock:          "focus lock",
		DPC_Fuji_DeviceError:        "device error",
		DPC_Fuji_ImageSpaceSD:       "image space SD",
		DPC_Fuji_MovieRemainingTime: "movie remaining time",
		DPC_Fuji_ShutterSpeed:       "shutter speed",
		DPC_Fuji_ImageAspectRatio:   "image size",
		DPC_Fuji_BatteryLevel:       "battery level",
		DPC_Fuji_InitSequence:       "init sequence",
		DPC_Fuji_AppVersion:         "app version",
		ptp.DevicePropCode(0):       "",
	}

	for code, want := range check {
		got := FujiDevicePropCodeAsString(code)
		if got != want {
			t.Errorf("FujiDevicePropCodeAsString() return = '%s', want '%s'", got, want)
		}
	}
}

func TestFujiBatteryLevelAsString(t *testing.T) {
	check := map[FujiBatteryLevel]string{
		BAT_Fuji_3bCritical: "critical",
		BAT_Fuji_3bOne:      "1/3",
		BAT_Fuji_3bTwo:      "2/3",
		BAT_Fuji_3bFull:     "3/3",
		BAT_Fuji_5bCritical: "critical",
		BAT_Fuji_5bOne:      "1/5",
		BAT_Fuji_5bTwo:      "2/5",
		BAT_Fuji_5bThree:    "3/5",
		BAT_Fuji_5bFour:     "4/5",
		BAT_Fuji_5bFull:     "5/5",
		FujiBatteryLevel(0): "",
	}

	for code, want := range check {
		got := FujiBatteryLevelAsString(code)
		if got != want {
			t.Errorf("FujiBatteryLevelAsString() return = '%s', want '%s'", got, want)
		}
	}
}

func TestFujiCommandDialModeAsString(t *testing.T) {
	check := map[FujiCommandDialMode]string{
		CMD_Fuji_Both:          "both",
		CMD_Fuji_Aperture:      "aperture",
		CMD_Fuji_ShutterSpeed:  "shutter speed",
		CMD_Fuji_None:          "none",
		FujiCommandDialMode(4): "",
	}

	for code, want := range check {
		got := FujiCommandDialModeAsString(code)
		if got != want {
			t.Errorf("FujiCommandDialModeAsString() return = '%s', want '%s'", got, want)
		}
	}
}

func TestFujiDeviceErrorAsString(t *testing.T) {
	check := map[FujiDeviceError]string{
		DE_Fuji_None:       "none",
		FujiDeviceError(4): "",
	}

	for code, want := range check {
		got := FujiDeviceErrorAsString(code)
		if got != want {
			t.Errorf("FujiDeviceErrorAsString() return = '%s', want '%s'", got, want)
		}
	}
}

func TestFujiExposureIndexAsString(t *testing.T) {
	check := map[uint32]string{
		uint32(EDX_Fuji_Auto): "auto",
		0x00000140:            "320",
		0x00001900:            "6400",
		0x40000064:            "L 100",
		0x40003200:            "H 12800",
		0x40006400:            "H 25600",
		0x4000C800:            "H 51200",
		0x800000C8:            "S 200",
		0x80000190:            "S 400",
		0x80000320:            "S 800",
		0x80000640:            "S 1600",
		0x80000C80:            "S 3200",
		0x80001900:            "S 6400",
	}

	for code, want := range check {
		got := FujiExposureIndexAsString(FujiExposureIndex(code))
		if got != want {
			t.Errorf("FujiExposureIndexAsString() return = '%s', want '%s'", got, want)
		}
	}
}

func TestFujiFilmSimulationAsString(t *testing.T) {
	check := map[FujiFilmSimulation]string{
		FS_Fuji_Provia:             "PROVIA",
		FS_Fuji_Velvia:             "Velvia",
		FS_Fuji_Astia:              "ASTIA",
		FS_Fuji_Monochrome:         "Monochrome",
		FS_Fuji_Sepia:              "Sepia",
		FS_Fuji_ProNegHigh:         "PRO Neg. Hi",
		FS_Fuji_ProNegStandard:     "PRO Neg. Std",
		FS_Fuji_MonochromeYeFilter: "Monochrome + Ye Filter",
		FS_Fuji_MonochromeRFilter:  "Monochrome + R Filter",
		FS_Fuji_MonochromeGFilter:  "Monochrome + G Filter",
		FS_Fuji_ClassicChrome:      "Classic Chrome",
		FS_Fuji_ACROS:              "ACROS",
		FS_Fuji_ACROSYe:            "ACROS Ye",
		FS_Fuji_ACROSR:             "ACROS R",
		FS_Fuji_ACROSG:             "ACROS G",
		FS_Fuji_ETERNA:             "ETERNA",
		FujiFilmSimulation(0):      "",
	}

	for code, want := range check {
		got := FujiFilmSimulationAsString(code)
		if got != want {
			t.Errorf("FujiFilmSimulationAsString() return = '%s', want '%s'", got, want)
		}
	}
}

func TestFujiFlashModeAsString(t *testing.T) {
	check := map[ptp.FlashMode]string{
		FM_Fuji_On:         "on",
		FM_Fuji_RedEye:     "red eye",
		FM_Fuji_RedEyeOn:   "red eye on",
		FM_Fuji_RedEyeSync: "red eye sync",
		FM_Fuji_RedEyeRear: "red eye rear",
		FM_Fuji_SlowSync:   "slow sync",
		FM_Fuji_RearSync:   "rear sync",
		FM_Fuji_Commander:  "commander",
		FM_Fuji_Disabled:   "disabled",
		FM_Fuji_Enabled:    "enabled",
		ptp.FlashMode(0):   "",
	}

	for code, want := range check {
		got := FujiFlashModeAsString(code)
		if got != want {
			t.Errorf("FujiFlashModeAsString() return = '%s', want '%s'", got, want)
		}
	}
}

func TestFujiFocusLockAsString(t *testing.T) {
	check := map[FujiFocusLock]string{
		FL_Fuji_On:       "on",
		FL_Fuji_Off:      "off",
		FujiFocusLock(2): "",
	}

	for code, want := range check {
		got := FujiFocusLockAsString(code)
		if got != want {
			t.Errorf("FujiFocusLockAsString() return = '%s', want '%s'", got, want)
		}
	}
}

func TestFujiFocusMeteringModeAsString(t *testing.T) {
	check := map[uint32]string{
		0x10000101: "1x1",
		0x02000203: "2x3",
		0x00300504: "5x4",
		0x00040707: "7x7",
	}

	for code, want := range check {
		got := FujiFocusMeteringModeAsString(code)
		if got != want {
			t.Errorf("FujiFocusMeteringModeAsString() return = '%s', want '%s'", got, want)
		}
	}
}

func TestFujiFocusModeAsString(t *testing.T) {
	check := map[ptp.FocusMode]string{
		FCM_Fuji_Single_Auto:     "single auto",
		FCM_Fuji_Continuous_Auto: "continuous auto",
		ptp.FocusMode(2):         "",
	}

	for code, want := range check {
		got := FujiFocusModeAsString(code)
		if got != want {
			t.Errorf("FujiFocusModeAsString() return = '%s', want '%s'", got, want)
		}
	}
}

func TestFujiImageAspectRatioAsString(t *testing.T) {
	check := map[FujiImageSize]string{
		IS_Fuji_Small_3x2:   "S 3:2",
		IS_Fuji_Small_16x9:  "S 16:9",
		IS_Fuji_Small_1x1:   "S 1:1",
		IS_Fuji_Medium_3x2:  "M 3:2",
		IS_Fuji_Medium_16x9: "M 16:9",
		IS_Fuji_Medium_1x1:  "M 1:1",
		IS_Fuji_Large_3x2:   "L 3:2",
		IS_Fuji_Large_16x9:  "L 16:9",
		IS_Fuji_Large_1x1:   "L 1:1",
		FujiImageSize(0):    "",
	}

	for code, want := range check {
		got := FujiImageAspectRatioAsString(code)
		if got != want {
			t.Errorf("FujiImageAspectRatioAsString() return = '%s', want '%s'", got, want)
		}
	}
}

func TestFujiImageQualityAsString(t *testing.T) {
	check := map[FujiImageQuality]string{
		IQ_Fuji_Fine:         "fine",
		IQ_Fuji_Normal:       "normal",
		IQ_Fuji_FineAndRAW:   "fine + RAW",
		IQ_Fuji_NormalAndRAW: "normal + RAW",
		IQ_Fuji_RAW:          "RAW",
		FujiImageQuality(0):  "",
	}

	for code, want := range check {
		got := FujiImageQualityAsString(code)
		if got != want {
			t.Errorf("FujiImageQualityAsString() return = '%s', want '%s'", got, want)
		}
	}
}

func TestFujiWhiteBalanceAsString(t *testing.T) {
	check := map[ptp.WhiteBalance]string{
		WB_Fuji_Fluorescent1: "fluorescent 1",
		WB_Fuji_Fluorescent2: "fluorescent 2",
		WB_Fuji_Fluorescent3: "fluorescent 3",
		WB_Fuji_Shade:        "shade",
		WB_Fuji_Underwater:   "underwater",
		WB_Fuji_Temperature:  "temprerature",
		WB_Fuji_Custom:       "custom",
		ptp.WhiteBalance(0):  "",
	}

	for code, want := range check {
		got := FujiWhiteBalanceAsString(code)
		if got != want {
			t.Errorf("FujiWhiteBalanceAsString() return = '%s', want '%s'", got, want)
		}
	}
}

func TestFujiSelfTimerAsString(t *testing.T) {
	check := map[FujiSelfTimer]string{
		ST_Fuji_1Sec:     "1 second",
		ST_Fuji_2Sec:     "2 seconds",
		ST_Fuji_5Sec:     "5 seconds",
		ST_Fuji_10Sec:    "10 seconds",
		ST_Fuji_Off:      "off",
		FujiSelfTimer(5): "",
	}

	for code, want := range check {
		got := FujiSelfTimerAsString(code)
		if got != want {
			t.Errorf("FujiSelfTimerAsString() return = '%s', want '%s'", got, want)
		}
	}
}

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

	gotPar, gotPkt, err := FujiSendOperationRequestAndGetResponse(c, ptp.OC_GetDevicePropValue, uint32(DPC_Fuji_AppVersion), 4)
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

	want, _ := ioutil.ReadFile("testdata/preview.jpg")
	got, err := FujiInitiateCapture(c)
	if err != nil {
		t.Errorf("FujiInitiateCapture() error = %s; want <nil>", err)
	}

	if bytes.Compare(got, want) != 0 {
		t.Errorf("FujiInitiateCapture() imgdata = %#v; want %#v", got, want)
	}
}
