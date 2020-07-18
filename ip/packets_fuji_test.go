package ip

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/malc0mn/ptp-ip/ptp"
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
	c, err := NewClient("fuji", address, fujiPort, "testèr", "67bace55-e7a4-4fbc-8e31-5122ee73a17c", LevelDebug)
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
