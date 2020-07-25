package fmt

import (
	"fmt"
	"github.com/malc0mn/ptp-ip/ip"
	"github.com/malc0mn/ptp-ip/ptp"
	"testing"
)

func TestFujiDevicePropCodeAsString(t *testing.T) {
	check := map[ptp.DevicePropCode]string{
		ip.DPC_Fuji_FilmSimulation:     "film simulation",
		ip.DPC_Fuji_ImageQuality:       "image quality",
		ip.DPC_Fuji_RecMode:            "rec mode",
		ip.DPC_Fuji_CommandDialMode:    "command dial mode",
		ip.DPC_Fuji_ExposureIndex:      "ISO",
		ip.DPC_Fuji_MovieISO:           "movie ISO",
		ip.DPC_Fuji_FocusMeteringMode:  "focus point",
		ip.DPC_Fuji_FocusLock:          "focus lock",
		ip.DPC_Fuji_DeviceError:        "device error",
		ip.DPC_Fuji_ImageSpaceSD:       "image space SD",
		ip.DPC_Fuji_MovieRemainingTime: "movie remaining time",
		ip.DPC_Fuji_ShutterSpeed:       "shutter speed",
		ip.DPC_Fuji_ImageAspectRatio:   "image size",
		ip.DPC_Fuji_BatteryLevel:       "battery level",
		ip.DPC_Fuji_InitSequence:       "init sequence",
		ip.DPC_Fuji_AppVersion:         "app version",
		ptp.DevicePropCode(0):          "",
	}

	for code, want := range check {
		got := FujiDevicePropCodeAsString(code)
		if got != want {
			t.Errorf("FujiDevicePropCodeAsString() return = '%s', want '%s'", got, want)
		}
	}
}

func TestFujiPropToDevicePropCode(t *testing.T) {
	check := map[string]ptp.DevicePropCode{
		PRP_Effect:            ip.DPC_Fuji_FilmSimulation,
		PRP_FocusMeteringMode: ip.DPC_Fuji_FocusMeteringMode,
		PRP_ISO:               ip.DPC_Fuji_ExposureIndex,
		"recmode":             ip.DPC_Fuji_RecMode,
	}

	for prop, want := range check {
		got, err := FujiPropToDevicePropCode(prop)
		if err != nil {
			t.Errorf("FujiPropToDevicePropCode() error = %s, want <nil>", err)
		}
		if got != want {
			t.Errorf("FujiPropToDevicePropCode() return = '%#x', want '%#x'", got, want)
		}
	}

	prop := "test"
	got, err := FujiPropToDevicePropCode(prop)
	wantE := fmt.Sprintf("unknown field name '%s'", prop)
	if err.Error() != wantE {
		t.Errorf("FujiPropToDevicePropCode() error = %s, want %s", err, wantE)
	}
	wantC := ptp.DevicePropCode(0)
	if got != wantC {
		t.Errorf("FujiPropToDevicePropCode() return = %d, want %d", got, wantC)
	}
}

func TestFujiBatteryLevelAsString(t *testing.T) {
	check := map[ip.FujiBatteryLevel]string{
		ip.BAT_Fuji_3bCritical: "critical",
		ip.BAT_Fuji_3bOne:      "1/3",
		ip.BAT_Fuji_3bTwo:      "2/3",
		ip.BAT_Fuji_3bFull:     "3/3",
		ip.BAT_Fuji_5bCritical: "critical",
		ip.BAT_Fuji_5bOne:      "1/5",
		ip.BAT_Fuji_5bTwo:      "2/5",
		ip.BAT_Fuji_5bThree:    "3/5",
		ip.BAT_Fuji_5bFour:     "4/5",
		ip.BAT_Fuji_5bFull:     "5/5",
		ip.FujiBatteryLevel(0): "",
	}

	for code, want := range check {
		got := FujiBatteryLevelAsString(code)
		if got != want {
			t.Errorf("FujiBatteryLevelAsString() return = '%s', want '%s'", got, want)
		}
	}
}

func TestFujiCommandDialModeAsString(t *testing.T) {
	check := map[ip.FujiCommandDialMode]string{
		ip.CMD_Fuji_Both:          "both",
		ip.CMD_Fuji_Aperture:      "aperture",
		ip.CMD_Fuji_ShutterSpeed:  "shutter speed",
		ip.CMD_Fuji_None:          "none",
		ip.FujiCommandDialMode(4): "",
	}

	for code, want := range check {
		got := FujiCommandDialModeAsString(code)
		if got != want {
			t.Errorf("FujiCommandDialModeAsString() return = '%s', want '%s'", got, want)
		}
	}
}

func TestFujiDeviceErrorAsString(t *testing.T) {
	check := map[ip.FujiDeviceError]string{
		ip.DE_Fuji_None:       "none",
		ip.FujiDeviceError(4): "",
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
		uint32(ip.EDX_Fuji_Auto): "auto",
		0x00000140:               "320",
		0x00001900:               "6400",
		0x40000064:               "L 100",
		0x40003200:               "H 12800",
		0x40006400:               "H 25600",
		0x4000C800:               "H 51200",
		0x800000C8:               "S 200",
		0x80000190:               "S 400",
		0x80000320:               "S 800",
		0x80000640:               "S 1600",
		0x80000C80:               "S 3200",
		0x80001900:               "S 6400",
	}

	for code, want := range check {
		got := FujiExposureIndexAsString(ip.FujiExposureIndex(code))
		if got != want {
			t.Errorf("FujiExposureIndexAsString() return = '%s', want '%s'", got, want)
		}
	}
}

func TestFujiFilmSimulationAsString(t *testing.T) {
	check := map[ip.FujiFilmSimulation]string{
		ip.FS_Fuji_Provia:             "PROVIA",
		ip.FS_Fuji_Velvia:             "Velvia",
		ip.FS_Fuji_Astia:              "ASTIA",
		ip.FS_Fuji_Monochrome:         "Monochrome",
		ip.FS_Fuji_Sepia:              "Sepia",
		ip.FS_Fuji_ProNegHigh:         "PRO Neg. Hi",
		ip.FS_Fuji_ProNegStandard:     "PRO Neg. Std",
		ip.FS_Fuji_MonochromeYeFilter: "Monochrome + Ye Filter",
		ip.FS_Fuji_MonochromeRFilter:  "Monochrome + R Filter",
		ip.FS_Fuji_MonochromeGFilter:  "Monochrome + G Filter",
		ip.FS_Fuji_ClassicChrome:      "Classic Chrome",
		ip.FS_Fuji_ACROS:              "ACROS",
		ip.FS_Fuji_ACROSYe:            "ACROS Ye",
		ip.FS_Fuji_ACROSR:             "ACROS R",
		ip.FS_Fuji_ACROSG:             "ACROS G",
		ip.FS_Fuji_ETERNA:             "ETERNA",
		ip.FujiFilmSimulation(0):      "",
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
		ip.FM_Fuji_On:         "on",
		ip.FM_Fuji_RedEye:     "red eye",
		ip.FM_Fuji_RedEyeOn:   "red eye on",
		ip.FM_Fuji_RedEyeSync: "red eye sync",
		ip.FM_Fuji_RedEyeRear: "red eye rear",
		ip.FM_Fuji_SlowSync:   "slow sync",
		ip.FM_Fuji_RearSync:   "rear sync",
		ip.FM_Fuji_Commander:  "commander",
		ip.FM_Fuji_Disabled:   "disabled",
		ip.FM_Fuji_Enabled:    "enabled",
		ptp.FlashMode(0):      "",
	}

	for code, want := range check {
		got := FujiFlashModeAsString(code)
		if got != want {
			t.Errorf("FujiFlashModeAsString() return = '%s', want '%s'", got, want)
		}
	}
}

func TestFujiFocusLockAsString(t *testing.T) {
	check := map[ip.FujiFocusLock]string{
		ip.FL_Fuji_On:       "on",
		ip.FL_Fuji_Off:      "off",
		ip.FujiFocusLock(2): "",
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
		ip.FCM_Fuji_Single_Auto:     "single auto",
		ip.FCM_Fuji_Continuous_Auto: "continuous auto",
		ptp.FocusMode(2):            "",
	}

	for code, want := range check {
		got := FujiFocusModeAsString(code)
		if got != want {
			t.Errorf("FujiFocusModeAsString() return = '%s', want '%s'", got, want)
		}
	}
}

func TestFujiImageAspectRatioAsString(t *testing.T) {
	check := map[ip.FujiImageSize]string{
		ip.IS_Fuji_Small_3x2:   "S 3:2",
		ip.IS_Fuji_Small_16x9:  "S 16:9",
		ip.IS_Fuji_Small_1x1:   "S 1:1",
		ip.IS_Fuji_Medium_3x2:  "M 3:2",
		ip.IS_Fuji_Medium_16x9: "M 16:9",
		ip.IS_Fuji_Medium_1x1:  "M 1:1",
		ip.IS_Fuji_Large_3x2:   "L 3:2",
		ip.IS_Fuji_Large_16x9:  "L 16:9",
		ip.IS_Fuji_Large_1x1:   "L 1:1",
		ip.FujiImageSize(0):    "",
	}

	for code, want := range check {
		got := FujiImageAspectRatioAsString(code)
		if got != want {
			t.Errorf("FujiImageAspectRatioAsString() return = '%s', want '%s'", got, want)
		}
	}
}

func TestFujiImageQualityAsString(t *testing.T) {
	check := map[ip.FujiImageQuality]string{
		ip.IQ_Fuji_Fine:         "fine",
		ip.IQ_Fuji_Normal:       "normal",
		ip.IQ_Fuji_FineAndRAW:   "fine + RAW",
		ip.IQ_Fuji_NormalAndRAW: "normal + RAW",
		ip.IQ_Fuji_RAW:          "RAW",
		ip.FujiImageQuality(0):  "",
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
		ip.WB_Fuji_Fluorescent1: "fluorescent 1",
		ip.WB_Fuji_Fluorescent2: "fluorescent 2",
		ip.WB_Fuji_Fluorescent3: "fluorescent 3",
		ip.WB_Fuji_Shade:        "shade",
		ip.WB_Fuji_Underwater:   "underwater",
		ip.WB_Fuji_Temperature:  "temprerature",
		ip.WB_Fuji_Custom:       "custom",
		ptp.WhiteBalance(0):     "",
	}

	for code, want := range check {
		got := FujiWhiteBalanceAsString(code)
		if got != want {
			t.Errorf("FujiWhiteBalanceAsString() return = '%s', want '%s'", got, want)
		}
	}
}

func TestFujiSelfTimerAsString(t *testing.T) {
	check := map[ip.FujiSelfTimer]string{
		ip.ST_Fuji_1Sec:     "1 second",
		ip.ST_Fuji_2Sec:     "2 seconds",
		ip.ST_Fuji_5Sec:     "5 seconds",
		ip.ST_Fuji_10Sec:    "10 seconds",
		ip.ST_Fuji_Off:      "off",
		ip.FujiSelfTimer(5): "",
	}

	for code, want := range check {
		got := FujiSelfTimerAsString(code)
		if got != want {
			t.Errorf("FujiSelfTimerAsString() return = '%s', want '%s'", got, want)
		}
	}
}
