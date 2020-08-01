package fmt

import (
	"fmt"
	"github.com/malc0mn/ptp-ip/ptp"
	"math"
	"strconv"
)

// GenericDevicePropCodeAsString returns the DevicePropCode as string. When the DevicePropCode is unknown, it returns an empty
// string.
func GenericDevicePropCodeAsString(code ptp.DevicePropCode) string {
	switch code {
	case ptp.DPC_BatteryLevel:
		return "battery level"
	case ptp.DPC_FunctionalMode:
		return "functional mode"
	case ptp.DPC_ImageSize:
		return "image size"
	case ptp.DPC_CompressionSetting:
		return "compression setting"
	case ptp.DPC_WhiteBalance:
		return "white balance"
	case ptp.DPC_RGBGain:
		return "RGB gain"
	case ptp.DPC_FNumber:
		return "F-number"
	case ptp.DPC_FocalLength:
		return "focal length"
	case ptp.DPC_FocusDistance:
		return "focus distance"
	case ptp.DPC_FocusMode:
		return "focus mode"
	case ptp.DPC_ExposureMeteringMode:
		return "exposure metering mode"
	case ptp.DPC_FlashMode:
		return "flash mode"
	case ptp.DPC_ExposureTime:
		return "exposure time"
	case ptp.DPC_ExposureProgramMode:
		return "exposure program mode"
	case ptp.DPC_ExposureIndex:
		return "ISO"
	case ptp.DPC_ExposureBiasCompensation:
		return "exposure bias compensation"
	case ptp.DPC_DateTime:
		return "date time"
	case ptp.DPC_CaptureDelay:
		return "capture delay"
	case ptp.DPC_StillCaptureMode:
		return "still capture mode"
	case ptp.DPC_Contrast:
		return "contrast"
	case ptp.DPC_Sharpness:
		return "sharpness"
	case ptp.DPC_DigitalZoom:
		return "digital zoom"
	case ptp.DPC_EffectMode:
		return "effect mode"
	case ptp.DPC_BurstNumber:
		return "burst number"
	case ptp.DPC_BurstInterval:
		return "burst interval"
	case ptp.DPC_TimelapseNumber:
		return "timelapse number"
	case ptp.DPC_TimelapseInterval:
		return "timelapse interval"
	case ptp.DPC_FocusMeteringMode:
		return "focus metering mode"
	case ptp.DPC_UploadURL:
		return "upload URL"
	case ptp.DPC_Artist:
		return "artist"
	case ptp.DPC_CopyrightInfo:
		return "copyright info"
	default:
		return ""
	}
}

// GenericPropToDevicePropCode converts a standardised property string to a valid DevicePropertyCode.
func GenericPropToDevicePropCode(field string) (ptp.DevicePropCode, error) {
	switch field {
	case PRP_Delay:
		return ptp.DPC_CaptureDelay, nil
	case PRP_Effect:
		return ptp.DPC_EffectMode, nil
	case PRP_Exposure:
		return ptp.DPC_ExposureTime, nil
	case PRP_ExpBias:
		return ptp.DPC_ExposureBiasCompensation, nil
	case PRP_FlashMode:
		return ptp.DPC_FlashMode, nil
	case PRP_ISO:
		return ptp.DPC_ExposureIndex, nil
	case PRP_WhiteBalance:
		return ptp.DPC_WhiteBalance, nil
	case PRP_FocusMeteringMode:
		return ptp.DPC_FocusMeteringMode, nil
	default:
		return 0, fmt.Errorf("unknown field name '%s'", field)
	}
}

func FormFlagAsString(flag ptp.DevicePropFormFlag) string {
	switch flag {
	case ptp.DPF_FormFlag_None:
		return "none"
	case ptp.DPF_FormFlag_Range:
		return "range"
	case ptp.DPF_FormFlag_Enum:
		return "enum"
	default:
		return ""
	}
}

func DataTypeCodeAsString(code ptp.DataTypeCode) string {
	switch code {
	case ptp.DTC_UNDEF:
		return "undefined"
	case ptp.DTC_INT8:
		return "int8"
	case ptp.DTC_UINT8:
		return "uint8"
	case ptp.DTC_INT16:
		return "int16"
	case ptp.DTC_UINT16:
		return "uint16"
	case ptp.DTC_INT32:
		return "int32"
	case ptp.DTC_UINT32:
		return "uint32"
	case ptp.DTC_INT64:
		return "int64"
	case ptp.DTC_UINT64:
		return "uint64"
	case ptp.DTC_INT128:
		return "int128"
	case ptp.DTC_UINT128:
		return "uint128"
	case ptp.DTC_AINT8:
		return "aint8"
	case ptp.DTC_AUINT8:
		return "auint8"
	case ptp.DTC_AINT16:
		return "aint16"
	case ptp.DTC_AUINT16:
		return "auint16"
	case ptp.DTC_AINT32:
		return "aint32"
	case ptp.DTC_AUINT32:
		return "auint32"
	case ptp.DTC_AINT64:
		return "aint64"
	case ptp.DTC_AUINT64:
		return "auint64"
	case ptp.DTC_AINT128:
		return "aint128"
	case ptp.DTC_AUINT128:
		return "auint128"
	case ptp.DTC_STR:
		return "string"
	default:
		return ""
	}
}

func DevicePropValueAsString(code ptp.DevicePropCode, v int64) string {
	switch code {
	case ptp.DPC_EffectMode:
		return EffectModeAsString(ptp.EffectMode(v))
	case ptp.DPC_ExposureBiasCompensation:
		return ExposureBiasCompensationAsString(int16(v))
	case ptp.DPC_ExposureMeteringMode:
		return ExposureMeteringModeAsString(ptp.ExposureMeteringMode(v))
	case ptp.DPC_ExposureProgramMode:
		return ExposureProgramModeAsString(ptp.ExposureProgramMode(v))
	case ptp.DPC_FlashMode:
		return FlashModeAsString(ptp.FlashMode(v))
	case ptp.DPC_FNumber:
		return FNumberAsString(uint16(v))
	case ptp.DPC_FocusMeteringMode:
		return FocusMeteringModeAsString(ptp.FocusMeteringMode(v))
	case ptp.DPC_FocusMode:
		return FocusModeAsString(ptp.FocusMode(v))
	case ptp.DPC_FunctionalMode:
		return FunctionalModeAsString(ptp.FunctionalMode(v))
	case ptp.DPC_StillCaptureMode:
		return StillCaptureModeAsString(ptp.StillCaptureMode(v))
	case ptp.DPC_WhiteBalance:
		return WhiteBalanceAsString(ptp.WhiteBalance(v))
	default:
		return ""
	}
}

func FNumberAsString(fn uint16) string {
	if fn == 0xffff {
		return "automatic"
	}

	return fmt.Sprintf("f/%.1f", float32(fn)/100)
}

func EffectModeAsString(fxm ptp.EffectMode) string {
	switch fxm {
	case ptp.FXM_Undefined:
		return "undefined"
	case ptp.FXM_Standard:
		return "standard"
	case ptp.FXM_BlackWhite:
		return "black and white"
	case ptp.FXM_Sepia:
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
		return sign + frac
	}

	return fmt.Sprintf("%d %s", int(i), frac)
}

func ExposureMeteringModeAsString(emm ptp.ExposureMeteringMode) string {
	switch emm {
	case ptp.EMM_Undefined:
		return "undefined"
	case ptp.EMM_Avarage:
		return "average"
	case ptp.EMM_CenterWeightedAvarage:
		return "center weighted average"
	case ptp.EMM_MultiSpot:
		return "multi spot"
	case ptp.EMM_CenterSpot:
		return "center spot"
	default:
		return ""
	}
}

func ExposureProgramModeAsString(epm ptp.ExposureProgramMode) string {
	switch epm {
	case ptp.EPM_Undefined:
		return "undefined"
	case ptp.EPM_Manual:
		return "manual"
	case ptp.EPM_Automatic:
		return "automatic"
	case ptp.EPM_AperturePriority:
		return "aperture priority"
	case ptp.EPM_ShutterPriority:
		return "shutter priority"
	case ptp.EPM_ProgramCreative:
		return "program creative"
	case ptp.EPM_ProgramAction:
		return "program action"
	case ptp.EPM_Portrait:
		return "portrait"
	default:
		return ""
	}
}

func FlashModeAsString(flm ptp.FlashMode) string {
	switch flm {
	case ptp.FLM_Undefined:
		return "undefined"
	case ptp.FLM_AutoFlash:
		return "auto flash"
	case ptp.FLM_FlashOff:
		return "off"
	case ptp.FLM_FillFlash:
		return "fill"
	case ptp.FLM_RedEyeAuto:
		return "red eye auto"
	case ptp.FLM_RedEyeFill:
		return "red eye fill"
	case ptp.FLM_ExternalSync:
		return "external sync"
	default:
		return ""
	}
}

func FocusMeteringModeAsString(fmm ptp.FocusMeteringMode) string {
	switch fmm {
	case ptp.FMM_Undefined:
		return "undefined"
	case ptp.FMM_CenterSpot:
		return "center spot"
	case ptp.FMM_MultiSpot:
		return "multi spot"
	default:
		return ""
	}
}

func FocusModeAsString(fcm ptp.FocusMode) string {
	switch fcm {
	case ptp.FCM_Undefined:
		return "undefined"
	case ptp.FCM_Manual:
		return "manual"
	case ptp.FCM_Automatic:
		return "automatic"
	case ptp.FCM_AutomaticMacro:
		return "automatic macro"
	default:
		return ""
	}
}

func FunctionalModeAsString(fum ptp.FunctionalMode) string {
	switch fum {
	case ptp.FUM_StandardMode:
		return "standard"
	case ptp.FUM_SleepState:
		return "sleep"
	default:
		return ""
	}
}

func SelfTestTypeAsString(stt ptp.SelfTestType) string {
	switch stt {
	case ptp.STT_Default:
		return "default"
	default:
		return ""
	}
}

func StillCaptureModeAsString(scm ptp.StillCaptureMode) string {
	switch scm {
	case ptp.SCM_Undefined:
		return "undefined"
	case ptp.SCM_Normal:
		return "normal"
	case ptp.SCM_Burst:
		return "burst"
	case ptp.SCM_Timelapse:
		return "timelapse"
	default:
		return ""
	}
}

func WhiteBalanceAsString(wb ptp.WhiteBalance) string {
	switch wb {
	case ptp.WB_Undefined:
		return "undefined"
	case ptp.WB_Manual:
		return "manual"
	case ptp.WB_Automatic:
		return "automatic"
	case ptp.WB_OnePushAutomatic:
		return "one push automatic"
	case ptp.WB_Daylight:
		return "daylight"
	case ptp.WB_Fluorescent:
		return "fluorescent"
	case ptp.WB_Tungsten:
		return "tungsten"
	case ptp.WB_Flash:
		return "flash"
	default:
		return ""
	}
}
