package ptp

import "testing"

var modes = map[DevicePropCode]map[int64]string{
	DPC_EffectMode: {
		int64(FXM_Undefined): "undefined",
		int64(FXM_Standard): "standard",
		int64(FXM_BlackWhite): "black and white",
		int64(FXM_Sepia): "sepia",
		int64(EffectMode(4)): "",
	},
	DPC_ExposureBiasCompensation: {
		int64(-3000): "-3",
		int64(-2667): "-2 2/3",
		int64(-2334): "-2 1/3",
		int64(-2000): "-2",
		int64(-1667): "-1 2/3",
		int64(-1334): "-1 1/3",
		int64(-1000): "-1",
		int64(0): "0",
		int64(334): "1/3",
		int64(667): "2/3",
		int64(1000): "1",
		int64(1334): "1 1/3",
		int64(1667): "1 2/3",
		int64(2000): "2",
		int64(2334): "2 1/3",
		int64(2667): "2 2/3",
		int64(3000): "3",
	},
	DPC_ExposureMeteringMode: {
		int64(EMM_Undefined): "undefined",
		int64(EMM_Avarage): "average",
		int64(EMM_CenterWeightedAvarage): "center weighted average",
		int64(EMM_MultiSpot): "multi spot",
		int64(EMM_CenterSpot): "center spot",
		int64(ExposureMeteringMode(5)): "",
	},
	DPC_ExposureProgramMode: {
		int64(EPM_Undefined): "undefined",
		int64(EPM_Manual): "manual",
		int64(EPM_Automatic): "automatic",
		int64(EPM_AperturePriority): "aperture priority",
		int64(EPM_ShutterPriority): "shutter priority",
		int64(EPM_ProgramCreative): "program creative",
		int64(EPM_ProgramAction): "program action",
		int64(EPM_Portrait): "portrait",
		int64(ExposureProgramMode(8)): "",
	},
	DPC_FlashMode: {
		int64(FLM_Undefined): "undefined",
		int64(FLM_AutoFlash): "auto flash",
		int64(FLM_FlashOff): "off",
		int64(FLM_FillFlash): "fill",
		int64(FLM_RedEyeAuto): "red eye auto",
		int64(FLM_RedEyeFill): "red eye fill",
		int64(FLM_ExternalSync): "external sync",
		int64(FlashMode(7)): "",
	},
	DPC_FocusMeteringMode: {
		int64(FMM_Undefined): "undefined",
		int64(FMM_CenterSpot): "center spot",
		int64(FMM_MultiSpot): "multi spot",
		int64(FocusMeteringMode(3)): "",
	},
	DPC_FocusMode: {
		int64(FCM_Undefined): "undefined",
		int64(FCM_Manual): "manual",
		int64(FCM_Automatic): "automatic",
		int64(FCM_AutomaticMacro): "automatic macro",
		int64(FocusMode(4)): "",
	},
	DPC_FunctionalMode: {
		int64(FUM_StandardMode): "standard",
		int64(FUM_SleepState): "sleep",
		int64(FunctionalMode(2)): "",
	},
	DPC_StillCaptureMode: {
		int64(SCM_Undefined): "undefined",
		int64(SCM_Normal): "normal",
		int64(SCM_Burst): "burst",
		int64(SCM_Timelapse): "timelapse",
		int64(StillCaptureMode(4)): "",
	},
	DPC_WhiteBalance: {
		int64(WB_Undefined): "undefined",
		int64(WB_Manual): "manual",
		int64(WB_Automatic): "automatic",
		int64(WB_OnePushAutomatic): "one push automatic",
		int64(WB_Daylight): "daylight",
		int64(WB_Fluorescent): "fluorescent",
		int64(WB_Tungsten): "tungsten",
		int64(WB_Flash): "flash",
		int64(WhiteBalance(8)): "",
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
	for code, want := range modes[DPC_EffectMode] {
		got := EffectModeAsString(EffectMode(code))
		if got != want {
			t.Errorf("EffectModeAsString() return = '%s', want '%s'", got, want)
		}
	}
}

func TestExposureBiasCompensationAsString(t *testing.T) {
	for ebv, want := range modes[DPC_ExposureBiasCompensation] {
		got := ExposureBiasCompensationAsString(int16(ebv))
		if got != want {
			t.Errorf("ExposureBiasCompensationAsString() return = '%s', want '%s'", got, want)
		}
	}
}

func TestExposureMeteringModeAsString(t *testing.T) {
	for code, want := range modes[DPC_ExposureMeteringMode] {
		got := ExposureMeteringModeAsString(ExposureMeteringMode(code))
		if got != want {
			t.Errorf("ExposureMeteringModeAsString() return = '%s', want '%s'", got, want)
		}
	}
}

func TestExposureProgramModeAsString(t *testing.T) {
	for code, want := range modes[DPC_ExposureProgramMode] {
		got := ExposureProgramModeAsString(ExposureProgramMode(code))
		if got != want {
			t.Errorf("ExposureProgramModeAsString() return = '%s', want '%s'", got, want)
		}
	}
}

func TestFlashModeAsString(t *testing.T) {
	for code, want := range modes[DPC_FlashMode] {
		got := FlashModeAsString(FlashMode(code))
		if got != want {
			t.Errorf("FlashModeAsString() return = '%s', want '%s'", got, want)
		}
	}
}

func TestFocusMeteringModeAsString(t *testing.T) {
	for code, want := range modes[DPC_FocusMeteringMode] {
		got := FocusMeteringModeAsString(FocusMeteringMode(code))
		if got != want {
			t.Errorf("FocusMeteringModeAsString() return = '%s', want '%s'", got, want)
		}
	}
}

func TestFocusModeAsString(t *testing.T) {
	for code, want := range modes[DPC_FocusMode] {
		got := FocusModeAsString(FocusMode(code))
		if got != want {
			t.Errorf("FocusModeAsString() return = '%s', want '%s'", got, want)
		}
	}
}

func TestFunctionalModeAsString(t *testing.T) {
	for code, want := range modes[DPC_FunctionalMode] {
		got := FunctionalModeAsString(FunctionalMode(code))
		if got != want {
			t.Errorf("FunctionalModeAsString() return = '%s', want '%s'", got, want)
		}
	}
}

func TestSelfTestTypeAsString(t *testing.T) {
	check := map[SelfTestType]string{
		STT_Default:     "default",
		SelfTestType(1): "",
	}

	for code, want := range check {
		got := SelfTestTypeAsString(code)
		if got != want {
			t.Errorf("SelfTestTypeAsString() return = '%s', want '%s'", got, want)
		}
	}
}

func TestStillCaptureModeAsString(t *testing.T) {
	for code, want := range modes[DPC_StillCaptureMode] {
		got := StillCaptureModeAsString(StillCaptureMode(code))
		if got != want {
			t.Errorf("StillCaptureModeAsString() return = '%s', want '%s'", got, want)
		}
	}
}

func TestWhiteBalanceAsString(t *testing.T) {
	for code, want := range modes[DPC_WhiteBalance] {
		got := WhiteBalanceAsString(WhiteBalance(code))
		if got != want {
			t.Errorf("WhiteBalanceAsString() return = '%s', want '%s'", got, want)
		}
	}
}
