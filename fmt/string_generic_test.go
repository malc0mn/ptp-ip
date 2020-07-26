package fmt

import (
	"fmt"
	"github.com/malc0mn/ptp-ip/ptp"
	"testing"
)

func TestGenericDevicePropCodeAsString(t *testing.T) {
	check := map[ptp.DevicePropCode]string{
		ptp.DPC_BatteryLevel:             "battery level",
		ptp.DPC_FunctionalMode:           "functional mode",
		ptp.DPC_ImageSize:                "image size",
		ptp.DPC_CompressionSetting:       "compression setting",
		ptp.DPC_WhiteBalance:             "white balance",
		ptp.DPC_RGBGain:                  "RGB gain",
		ptp.DPC_FNumber:                  "F number",
		ptp.DPC_FocalLength:              "focal length",
		ptp.DPC_FocusDistance:            "focus distance",
		ptp.DPC_FocusMode:                "focus mode",
		ptp.DPC_ExposureMeteringMode:     "exposure metering mode",
		ptp.DPC_FlashMode:                "flash mode",
		ptp.DPC_ExposureTime:             "exposure time",
		ptp.DPC_ExposureProgramMode:      "exposure program mode",
		ptp.DPC_ExposureIndex:            "ISO",
		ptp.DPC_ExposureBiasCompensation: "exposure bias compensation",
		ptp.DPC_DateTime:                 "date time",
		ptp.DPC_CaptureDelay:             "capture delay",
		ptp.DPC_StillCaptureMode:         "still capture mode",
		ptp.DPC_Contrast:                 "contrast",
		ptp.DPC_Sharpness:                "sharpness",
		ptp.DPC_DigitalZoom:              "digital zoom",
		ptp.DPC_EffectMode:               "effect mode",
		ptp.DPC_BurstNumber:              "burst number",
		ptp.DPC_BurstInterval:            "burst interval",
		ptp.DPC_TimelapseNumber:          "timelapse number",
		ptp.DPC_TimelapseInterval:        "timelapse interval",
		ptp.DPC_FocusMeteringMode:        "focus metering mode",
		ptp.DPC_UploadURL:                "upload URL",
		ptp.DPC_Artist:                   "artist",
		ptp.DPC_CopyrightInfo:            "copyright info",
		ptp.DevicePropCode(0):            "",
	}

	for code, want := range check {
		got := GenericDevicePropCodeAsString(code)
		if got != want {
			t.Errorf("GenericDevicePropCodeAsString() return = '%s', want '%s'", got, want)
		}
	}
}

func TestPropToDevicePropCode(t *testing.T) {
	check := map[string]ptp.DevicePropCode{
		PRP_Delay:             ptp.DPC_CaptureDelay,
		PRP_Effect:            ptp.DPC_EffectMode,
		PRP_Exposure:          ptp.DPC_ExposureTime,
		PRP_ExpBias:           ptp.DPC_ExposureBiasCompensation,
		PRP_FlashMode:         ptp.DPC_FlashMode,
		PRP_FocusMeteringMode: ptp.DPC_FocusMeteringMode,
		PRP_ISO:               ptp.DPC_ExposureIndex,
		PRP_WhiteBalance:      ptp.DPC_WhiteBalance,
	}

	for prop, want := range check {
		got, err := GenericPropToDevicePropCode(prop)
		if err != nil {
			t.Errorf("GenericPropToDevicePropCode() error = %s, want <nil>", err)
		}
		if got != want {
			t.Errorf("GenericPropToDevicePropCode() return = '%#x', want '%#x'", got, want)
		}
	}

	prop := "test"
	got, err := GenericPropToDevicePropCode(prop)
	wantE := fmt.Sprintf("unknown field name '%s'", prop)
	if err.Error() != wantE {
		t.Errorf("GenericPropToDevicePropCode() error = %s, want %s", err, wantE)
	}
	wantC := ptp.DevicePropCode(0)
	if got != wantC {
		t.Errorf("GenericPropToDevicePropCode() return = %d, want %d", got, wantC)
	}
}

func TestFormFlagAsString(t *testing.T) {
	check := map[ptp.DevicePropFormFlag]string{
		ptp.DPF_FormFlag_None:     "none",
		ptp.DPF_FormFlag_Range:    "range",
		ptp.DPF_FormFlag_Enum:     "enum",
		ptp.DevicePropFormFlag(3): "",
	}

	for code, want := range check {
		got := FormFlagAsString(code)
		if got != want {
			t.Errorf("FormFlagAsString() return = '%s', want '%s'", got, want)
		}
	}
}

func TestDataTypeCodeAsString(t *testing.T) {
	check := map[ptp.DataTypeCode]string{
		ptp.DTC_UNDEF:            "undefined",
		ptp.DTC_INT8:             "int8",
		ptp.DTC_UINT8:            "uint8",
		ptp.DTC_INT16:            "int16",
		ptp.DTC_UINT16:           "uint16",
		ptp.DTC_INT32:            "int32",
		ptp.DTC_UINT32:           "uint32",
		ptp.DTC_INT64:            "int64",
		ptp.DTC_UINT64:           "uint64",
		ptp.DTC_INT128:           "int128",
		ptp.DTC_UINT128:          "uint128",
		ptp.DTC_AINT8:            "aint8",
		ptp.DTC_AUINT8:           "auint8",
		ptp.DTC_AINT16:           "aint16",
		ptp.DTC_AUINT16:          "auint16",
		ptp.DTC_AINT32:           "aint32",
		ptp.DTC_AUINT32:          "auint32",
		ptp.DTC_AINT64:           "aint64",
		ptp.DTC_AUINT64:          "auint64",
		ptp.DTC_AINT128:          "aint128",
		ptp.DTC_AUINT128:         "auint128",
		ptp.DTC_STR:              "string",
		ptp.DataTypeCode(0xF000): "",
	}

	for code, want := range check {
		got := DataTypeCodeAsString(code)
		if got != want {
			t.Errorf("DataTypeCodeAsString() return = '%s', want '%s'", got, want)
		}
	}
}

var modes = map[ptp.DevicePropCode]map[int64]string{
	ptp.DPC_EffectMode: {
		int64(ptp.FXM_Undefined):  "undefined",
		int64(ptp.FXM_Standard):   "standard",
		int64(ptp.FXM_BlackWhite): "black and white",
		int64(ptp.FXM_Sepia):      "sepia",
		int64(ptp.EffectMode(4)):  "",
	},
	ptp.DPC_ExposureBiasCompensation: {
		int64(-3000): "-3",
		int64(-2667): "-2 2/3",
		int64(-2334): "-2 1/3",
		int64(-2000): "-2",
		int64(-1667): "-1 2/3",
		int64(-1334): "-1 1/3",
		int64(-1000): "-1",
		int64(0):     "0",
		int64(334):   "1/3",
		int64(667):   "2/3",
		int64(1000):  "1",
		int64(1334):  "1 1/3",
		int64(1667):  "1 2/3",
		int64(2000):  "2",
		int64(2334):  "2 1/3",
		int64(2667):  "2 2/3",
		int64(3000):  "3",
	},
	ptp.DPC_ExposureMeteringMode: {
		int64(ptp.EMM_Undefined):             "undefined",
		int64(ptp.EMM_Avarage):               "average",
		int64(ptp.EMM_CenterWeightedAvarage): "center weighted average",
		int64(ptp.EMM_MultiSpot):             "multi spot",
		int64(ptp.EMM_CenterSpot):            "center spot",
		int64(ptp.ExposureMeteringMode(5)):   "",
	},
	ptp.DPC_ExposureProgramMode: {
		int64(ptp.EPM_Undefined):          "undefined",
		int64(ptp.EPM_Manual):             "manual",
		int64(ptp.EPM_Automatic):          "automatic",
		int64(ptp.EPM_AperturePriority):   "aperture priority",
		int64(ptp.EPM_ShutterPriority):    "shutter priority",
		int64(ptp.EPM_ProgramCreative):    "program creative",
		int64(ptp.EPM_ProgramAction):      "program action",
		int64(ptp.EPM_Portrait):           "portrait",
		int64(ptp.ExposureProgramMode(8)): "",
	},
	ptp.DPC_FlashMode: {
		int64(ptp.FLM_Undefined):    "undefined",
		int64(ptp.FLM_AutoFlash):    "auto flash",
		int64(ptp.FLM_FlashOff):     "off",
		int64(ptp.FLM_FillFlash):    "fill",
		int64(ptp.FLM_RedEyeAuto):   "red eye auto",
		int64(ptp.FLM_RedEyeFill):   "red eye fill",
		int64(ptp.FLM_ExternalSync): "external sync",
		int64(ptp.FlashMode(7)):     "",
	},
	ptp.DPC_FocusMeteringMode: {
		int64(ptp.FMM_Undefined):        "undefined",
		int64(ptp.FMM_CenterSpot):       "center spot",
		int64(ptp.FMM_MultiSpot):        "multi spot",
		int64(ptp.FocusMeteringMode(3)): "",
	},
	ptp.DPC_FocusMode: {
		int64(ptp.FCM_Undefined):      "undefined",
		int64(ptp.FCM_Manual):         "manual",
		int64(ptp.FCM_Automatic):      "automatic",
		int64(ptp.FCM_AutomaticMacro): "automatic macro",
		int64(ptp.FocusMode(4)):       "",
	},
	ptp.DPC_FunctionalMode: {
		int64(ptp.FUM_StandardMode):  "standard",
		int64(ptp.FUM_SleepState):    "sleep",
		int64(ptp.FunctionalMode(2)): "",
	},
	ptp.DPC_StillCaptureMode: {
		int64(ptp.SCM_Undefined):       "undefined",
		int64(ptp.SCM_Normal):          "normal",
		int64(ptp.SCM_Burst):           "burst",
		int64(ptp.SCM_Timelapse):       "timelapse",
		int64(ptp.StillCaptureMode(4)): "",
	},
	ptp.DPC_WhiteBalance: {
		int64(ptp.WB_Undefined):        "undefined",
		int64(ptp.WB_Manual):           "manual",
		int64(ptp.WB_Automatic):        "automatic",
		int64(ptp.WB_OnePushAutomatic): "one push automatic",
		int64(ptp.WB_Daylight):         "daylight",
		int64(ptp.WB_Fluorescent):      "fluorescent",
		int64(ptp.WB_Tungsten):         "tungsten",
		int64(ptp.WB_Flash):            "flash",
		int64(ptp.WhiteBalance(8)):     "",
	},
}

func TestDevicePropValueAsString(t *testing.T) {
	for dpc, vals := range modes {
		for val, want := range vals {
			got := DevicePropValueAsString(dpc, int64(val))
			if got != want {
				t.Errorf("DevicePropValueAsString() return = '%s', want '%s'", got, want)
			}
		}
	}
}

func TestEffectModeAsString(t *testing.T) {
	for code, want := range modes[ptp.DPC_EffectMode] {
		got := EffectModeAsString(ptp.EffectMode(code))
		if got != want {
			t.Errorf("EffectModeAsString() return = '%s', want '%s'", got, want)
		}
	}
}

func TestExposureBiasCompensationAsString(t *testing.T) {
	for ebv, want := range modes[ptp.DPC_ExposureBiasCompensation] {
		got := ExposureBiasCompensationAsString(int16(ebv))
		if got != want {
			t.Errorf("ExposureBiasCompensationAsString() return = '%s', want '%s'", got, want)
		}
	}
}

func TestExposureMeteringModeAsString(t *testing.T) {
	for code, want := range modes[ptp.DPC_ExposureMeteringMode] {
		got := ExposureMeteringModeAsString(ptp.ExposureMeteringMode(code))
		if got != want {
			t.Errorf("ExposureMeteringModeAsString() return = '%s', want '%s'", got, want)
		}
	}
}

func TestExposureProgramModeAsString(t *testing.T) {
	for code, want := range modes[ptp.DPC_ExposureProgramMode] {
		got := ExposureProgramModeAsString(ptp.ExposureProgramMode(code))
		if got != want {
			t.Errorf("ExposureProgramModeAsString() return = '%s', want '%s'", got, want)
		}
	}
}

func TestFlashModeAsString(t *testing.T) {
	for code, want := range modes[ptp.DPC_FlashMode] {
		got := FlashModeAsString(ptp.FlashMode(code))
		if got != want {
			t.Errorf("FlashModeAsString() return = '%s', want '%s'", got, want)
		}
	}
}

func TestFocusMeteringModeAsString(t *testing.T) {
	for code, want := range modes[ptp.DPC_FocusMeteringMode] {
		got := FocusMeteringModeAsString(ptp.FocusMeteringMode(code))
		if got != want {
			t.Errorf("FocusMeteringModeAsString() return = '%s', want '%s'", got, want)
		}
	}
}

func TestFocusModeAsString(t *testing.T) {
	for code, want := range modes[ptp.DPC_FocusMode] {
		got := FocusModeAsString(ptp.FocusMode(code))
		if got != want {
			t.Errorf("FocusModeAsString() return = '%s', want '%s'", got, want)
		}
	}
}

func TestFunctionalModeAsString(t *testing.T) {
	for code, want := range modes[ptp.DPC_FunctionalMode] {
		got := FunctionalModeAsString(ptp.FunctionalMode(code))
		if got != want {
			t.Errorf("FunctionalModeAsString() return = '%s', want '%s'", got, want)
		}
	}
}

func TestSelfTestTypeAsString(t *testing.T) {
	check := map[ptp.SelfTestType]string{
		ptp.STT_Default:     "default",
		ptp.SelfTestType(1): "",
	}

	for code, want := range check {
		got := SelfTestTypeAsString(code)
		if got != want {
			t.Errorf("SelfTestTypeAsString() return = '%s', want '%s'", got, want)
		}
	}
}

func TestStillCaptureModeAsString(t *testing.T) {
	for code, want := range modes[ptp.DPC_StillCaptureMode] {
		got := StillCaptureModeAsString(ptp.StillCaptureMode(code))
		if got != want {
			t.Errorf("StillCaptureModeAsString() return = '%s', want '%s'", got, want)
		}
	}
}

func TestWhiteBalanceAsString(t *testing.T) {
	for code, want := range modes[ptp.DPC_WhiteBalance] {
		got := WhiteBalanceAsString(ptp.WhiteBalance(code))
		if got != want {
			t.Errorf("WhiteBalanceAsString() return = '%s', want '%s'", got, want)
		}
	}
}
