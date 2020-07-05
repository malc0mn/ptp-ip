package ptp

import (
	"fmt"
	"math"
	"strconv"
)

type EffectMode uint16
type ExposureMeteringMode uint16
type ExposureProgramMode uint16
type FlashMode uint16
type FocusMeteringMode uint16
type FocusMode uint16
type FunctionalMode uint16
type SelfTestType uint16
type StillCaptureMode uint16
type WhiteBalance uint16

const (
	EMM_Undefined             ExposureMeteringMode = 0x0000
	EMM_Avarage               ExposureMeteringMode = 0x0001
	EMM_CenterWeightedAvarage ExposureMeteringMode = 0x0002
	EMM_MultiSpot             ExposureMeteringMode = 0x0003
	EMM_CenterSpot            ExposureMeteringMode = 0x0004

	EPM_Undefined        ExposureProgramMode = 0x0000
	EPM_Manual           ExposureProgramMode = 0x0001
	EPM_Automatic        ExposureProgramMode = 0x0002
	EPM_AperturePriority ExposureProgramMode = 0x0003
	EPM_ShutterPriority  ExposureProgramMode = 0x0004
	EPM_ProgramCreative  ExposureProgramMode = 0x0005
	EPM_ProgramAction    ExposureProgramMode = 0x0006
	EPM_Portrait         ExposureProgramMode = 0x0007

	FCM_Undefined      FocusMode = 0x0000
	FCM_Manual         FocusMode = 0x0001
	FCM_Automatic      FocusMode = 0x0002
	FCM_AutomaticMacro FocusMode = 0x0003

	FLM_Undefined    FlashMode = 0x0000
	FLM_AutoFlash    FlashMode = 0x0001
	FLM_FlashOff     FlashMode = 0x0002
	FLM_FillFlash    FlashMode = 0x0003
	FLM_RedEyeAuto   FlashMode = 0x0004
	FLM_RedEyeFill   FlashMode = 0x0005
	FLM_ExternalSync FlashMode = 0x0006

	FMM_Undefined  FocusMeteringMode = 0x0000
	FMM_CenterSpot FocusMeteringMode = 0x0001
	FMM_MultiSpot  FocusMeteringMode = 0x0002

	FUM_StandardMode FunctionalMode = 0x0000
	FUM_SleepState   FunctionalMode = 0x0001

	FXM_Undefined  EffectMode = 0x0000
	FXM_Standard   EffectMode = 0x0001
	FXM_BlackWhite EffectMode = 0x0002
	FXM_Sepia      EffectMode = 0x0003

	SCM_Undefined StillCaptureMode = 0x0000
	SCM_Normal    StillCaptureMode = 0x0001
	SCM_Burst     StillCaptureMode = 0x0002
	SCM_Timelapse StillCaptureMode = 0x0003

	// STT_Default is the default device-specific self-test.
	STT_Default SelfTestType = 0x0000

	WB_Undefined WhiteBalance = 0x0000
	// WB_Manual indicates the white balance is set directly using the RGB Gain property and is static until changed.
	WB_Manual WhiteBalance = 0x0001
	// WB_Automatic indicates the device attempts to set the white balance using some kind of automatic mechanism.
	WB_Automatic WhiteBalance = 0x0002
	// WB_OnePushAutomatic indicates the user must press the capture button while pointing the device at a white field,
	// at which time the device determines the white balance setting.
	WB_OnePushAutomatic WhiteBalance = 0x0003
	// WB_Daylight indicates the device attempts to set the white balance to a value that is appropriate for use in
	// daylight conditions.
	WB_Daylight WhiteBalance = 0x0004
	// WB_Fluorescent indicates the device attempts to set the white balance to a value that is appropriate for use in
	// with fluorescent lighting conditions.
	WB_Fluorescent WhiteBalance = 0x0005
	// WB_Tungsten indicates the device attempts to set the white balance to a value that is appropriate for use in
	// conditions with a tungsten light source.
	WB_Tungsten WhiteBalance = 0x0006
	// WB_Flash indicates the device attempts to set the white balance to a value that is appropriate for flash
	// conditions.
	WB_Flash WhiteBalance = 0x0007
)

func DevicePropValueAsString(code DevicePropCode, v int64) string {
	switch code {
	case DPC_EffectMode:
		return EffectModeAsString(EffectMode(v))
	case DPC_ExposureBiasCompensation:
		return ExposureBiasCompensationAsString(int16(v))
	case DPC_ExposureMeteringMode:
		return ExposureMeteringModeAsString(ExposureMeteringMode(v))
	case DPC_ExposureProgramMode:
		return ExposureProgramModeAsString(ExposureProgramMode(v))
	case DPC_FlashMode:
		return FlashModeAsString(FlashMode(v))
	case DPC_FocusMeteringMode:
		return FocusMeteringModeAsString(FocusMeteringMode(v))
	case DPC_FocusMode:
		return FocusModeAsString(FocusMode(v))
	case DPC_FunctionalMode:
		return FunctionalModeAsString(FunctionalMode(v))
	case DPC_StillCaptureMode:
		return StillCaptureModeAsString(StillCaptureMode(v))
	case DPC_WhiteBalance:
		return WhiteBalanceAsString(WhiteBalance(v))
	default:
		return ""
	}
}

func EffectModeAsString(fxm EffectMode) string {
	switch fxm {
	case FXM_Undefined:
		return "undefined"
	case FXM_Standard:
		return "standard"
	case FXM_BlackWhite:
		return "black and white"
	case FXM_Sepia:
		return "sepia"
	default:
		return ""
	}
}

func ExposureBiasCompensationAsString(ebv int16) string {
	i, f := math.Modf(float64(ebv) / float64(1000))

	if f == 0 {
		return strconv.FormatInt(int64(i), 10)
	}

	// Tried to use big.Rat to do the conversion, but it's trying to be "too precise" :/
	frac := "1/3"
	if math.Abs(f) > 0.4 {
		frac = "2/3"
	}

	if i == 0 {
		sign := ""
		if f < 0 {
			sign = "-"
		}
		return fmt.Sprintf("%s%s", sign, frac)
	}

	return fmt.Sprintf("%d %s", int(i), frac)
}

func ExposureMeteringModeAsString(emm ExposureMeteringMode) string {
	switch emm {
	case EMM_Undefined:
		return "undefined"
	case EMM_Avarage:
		return "average"
	case EMM_CenterWeightedAvarage:
		return "center weighted average"
	case EMM_MultiSpot:
		return "multi spot"
	case EMM_CenterSpot:
		return "center spot"
	default:
		return ""
	}
}

func ExposureProgramModeAsString(epm ExposureProgramMode) string {
	switch epm {
	case EPM_Undefined:
		return "undefined"
	case EPM_Manual:
		return "manual"
	case EPM_Automatic:
		return "automatic"
	case EPM_AperturePriority:
		return "aperture priority"
	case EPM_ShutterPriority:
		return "shutter priority"
	case EPM_ProgramCreative:
		return "program creative"
	case EPM_ProgramAction:
		return "program action"
	case EPM_Portrait:
		return "portrait"
	default:
		return ""
	}
}

func FlashModeAsString(flm FlashMode) string {
	switch flm {
	case FLM_Undefined:
		return "undefined"
	case FLM_AutoFlash:
		return "auto flash"
	case FLM_FlashOff:
		return "off"
	case FLM_FillFlash:
		return "fill"
	case FLM_RedEyeAuto:
		return "red eye auto"
	case FLM_RedEyeFill:
		return "red eye fill"
	case FLM_ExternalSync:
		return "external sync"
	default:
		return ""
	}
}

func FocusMeteringModeAsString(fmm FocusMeteringMode) string {
	switch fmm {
	case FMM_Undefined:
		return "undefined"
	case FMM_CenterSpot:
		return "center spot"
	case FMM_MultiSpot:
		return "multi spot"
	default:
		return ""
	}
}

func FocusModeAsString(fcm FocusMode) string {
	switch fcm {
	case FCM_Undefined:
		return "undefined"
	case FCM_Manual:
		return "manual"
	case FCM_Automatic:
		return "automatic"
	case FCM_AutomaticMacro:
		return "automatic macro"
	default:
		return ""
	}
}

func FunctionalModeAsString(fum FunctionalMode) string {
	switch fum {
	case FUM_StandardMode:
		return "standard"
	case FUM_SleepState:
		return "sleep"
	default:
		return ""
	}
}

func SelfTestTypeAsString(stt SelfTestType) string {
	switch stt {
	case STT_Default:
		return "default"
	default:
		return ""
	}
}

func StillCaptureModeAsString(scm StillCaptureMode) string {
	switch scm {
	case SCM_Undefined:
		return "undefined"
	case SCM_Normal:
		return "normal"
	case SCM_Burst:
		return "burst"
	case SCM_Timelapse:
		return "timelapse"
	default:
		return ""
	}
}

func WhiteBalanceAsString(wb WhiteBalance) string {
	switch wb {
	case WB_Undefined:
		return "undefined"
	case WB_Manual:
		return "manual"
	case WB_Automatic:
		return "automatic"
	case WB_OnePushAutomatic:
		return "one push automatic"
	case WB_Daylight:
		return "daylight"
	case WB_Fluorescent:
		return "fluorescent"
	case WB_Tungsten:
		return "tungsten"
	case WB_Flash:
		return "flash"
	default:
		return ""
	}
}
