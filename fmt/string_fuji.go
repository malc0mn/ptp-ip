package fmt

import (
	"fmt"
	"github.com/malc0mn/ptp-ip/ip"
	"github.com/malc0mn/ptp-ip/ptp"
	"strconv"
)

func FujiDevicePropCodeAsString(code ptp.DevicePropCode) string {
	switch code {
	case ip.DPC_Fuji_FilmSimulation:
		return "film simulation"
	case ip.DPC_Fuji_ImageQuality:
		return "image quality"
	case ip.DPC_Fuji_RecMode:
		return "rec mode"
	case ip.DPC_Fuji_CommandDialMode:
		return "command dial mode"
	case ip.DPC_Fuji_ExposureIndex:
		return "ISO"
	case ip.DPC_Fuji_MovieISO:
		return "movie ISO"
	case ip.DPC_Fuji_FocusMeteringMode:
		return "focus point"
	case ip.DPC_Fuji_FocusLock:
		return "focus lock"
	case ip.DPC_Fuji_DeviceError:
		return "device error"
	case ip.DPC_Fuji_ImageSpaceSD:
		return "image space SD"
	case ip.DPC_Fuji_MovieRemainingTime:
		return "movie remaining time"
	case ip.DPC_Fuji_ShutterSpeed:
		return "shutter speed"
	case ip.DPC_Fuji_ImageAspectRatio:
		return "image size"
	case ip.DPC_Fuji_BatteryLevel:
		return "battery level"
	case ip.DPC_Fuji_InitSequence:
		return "init sequence"
	case ip.DPC_Fuji_AppVersion:
		return "app version"
	default:
		return GenericDevicePropCodeAsString(code)
	}
}

// FujiPropToDevicePropCode converts a standardised property string to a valid ptp.DevicePropertyCode.
func FujiPropToDevicePropCode(field string) (ptp.DevicePropCode, error) {
	switch field {
	case PRP_ISO:
		return ip.DPC_Fuji_ExposureIndex, nil
	case PRP_Effect:
		return ip.DPC_Fuji_FilmSimulation, nil
	case "recmode":
		return ip.DPC_Fuji_RecMode, nil
	case PRP_FocusMeteringMode:
		return ip.DPC_Fuji_FocusMeteringMode, nil
	default:
		return GenericPropToDevicePropCode(field)
	}
}

func FujiDevicePropValueAsString(code ptp.DevicePropCode, v int64) string {
	switch code {
	case ptp.DPC_BatteryLevel, ip.DPC_Fuji_BatteryLevel:
		return FujiBatteryLevelAsString(ip.FujiBatteryLevel(v))
	case ip.DPC_Fuji_CommandDialMode:
		return FujiCommandDialModeAsString(ip.FujiCommandDialMode(v))
	case ip.DPC_Fuji_DeviceError:
		return FujiDeviceErrorAsString(ip.FujiDeviceError(v))
	case ip.DPC_Fuji_ExposureIndex:
		return FujiExposureIndexAsString(ip.FujiExposureIndex(v))
	case ip.DPC_Fuji_FilmSimulation:
		return FujiFilmSimulationAsString(ip.FujiFilmSimulation(v))
	case ptp.DPC_FlashMode:
		return FujiFlashModeAsString(ptp.FlashMode(v))
	case ip.DPC_Fuji_FocusLock:
		return FujiFocusLockAsString(ip.FujiFocusLock(v))
	case ip.DPC_Fuji_FocusMeteringMode:
		return FujiFocusMeteringModeAsString(uint32(v))
	case ptp.DPC_FocusMode:
		return FujiFocusModeAsString(ptp.FocusMode(v))
	case ip.DPC_Fuji_ImageAspectRatio:
		return FujiImageAspectRatioAsString(ip.FujiImageSize(v))
	case ip.DPC_Fuji_ImageQuality:
		return FujiImageQualityAsString(ip.FujiImageQuality(v))
	case ptp.DPC_WhiteBalance:
		return FujiWhiteBalanceAsString(ptp.WhiteBalance(v))
	case ptp.DPC_CaptureDelay:
		return FujiSelfTimerAsString(ip.FujiSelfTimer(v))
	default:
		return DevicePropValueAsString(code, v)
	}
}

func FujiBatteryLevelAsString(bat ip.FujiBatteryLevel) string {
	switch bat {
	case ip.BAT_Fuji_3bCritical:
		return "critical"
	case ip.BAT_Fuji_3bOne:
		return "1/3"
	case ip.BAT_Fuji_3bTwo:
		return "2/3"
	case ip.BAT_Fuji_3bFull:
		return "3/3"
	case ip.BAT_Fuji_5bCritical:
		return "critical"
	case ip.BAT_Fuji_5bOne:
		return "1/5"
	case ip.BAT_Fuji_5bTwo:
		return "2/5"
	case ip.BAT_Fuji_5bThree:
		return "3/5"
	case ip.BAT_Fuji_5bFour:
		return "4/5"
	case ip.BAT_Fuji_5bFull:
		return "5/5"
	default:
		return ""
	}
}

func FujiCommandDialModeAsString(cmd ip.FujiCommandDialMode) string {
	switch cmd {
	case ip.CMD_Fuji_Both:
		return "both"
	case ip.CMD_Fuji_Aperture:
		return "aperture"
	case ip.CMD_Fuji_ShutterSpeed:
		return "shutter speed"
	case ip.CMD_Fuji_None:
		return "none"
	default:
		return ""
	}
}

func FujiDeviceErrorAsString(de ip.FujiDeviceError) string {
	switch de {
	case ip.DE_Fuji_None:
		return "none"
	default:
		return ""
	}
}

func FujiExposureIndexAsString(edx ip.FujiExposureIndex) string {
	if edx == ip.EDX_Fuji_Auto {
		return "auto"
	}

	prefix := ""
	val := int64(edx & 0x0000FFFF)

	switch uint16(edx >> 16) {
	case ip.EDX_Fuji_Extended:
		prefix = "H"
		if val < 200 {
			prefix = "L"
		}
	case ip.EDX_Fuji_MaxSensitivity:
		prefix = "S"
	}

	if prefix == "" {
		return strconv.FormatInt(val, 10)
	}

	return fmt.Sprintf("%s%d", prefix, val)
}

func FujiFilmSimulationAsString(fs ip.FujiFilmSimulation) string {
	switch fs {
	case ip.FS_Fuji_Provia:
		return "PROVIA"
	case ip.FS_Fuji_Velvia:
		return "Velvia"
	case ip.FS_Fuji_Astia:
		return "ASTIA"
	case ip.FS_Fuji_Monochrome:
		return "Monochrome"
	case ip.FS_Fuji_Sepia:
		return "Sepia"
	case ip.FS_Fuji_ProNegHigh:
		return "PRO Neg. Hi"
	case ip.FS_Fuji_ProNegStandard:
		return "PRO Neg. Std"
	case ip.FS_Fuji_MonochromeYeFilter:
		return "Monochrome + Ye Filter"
	case ip.FS_Fuji_MonochromeRFilter:
		return "Monochrome + R Filter"
	case ip.FS_Fuji_MonochromeGFilter:
		return "Monochrome + G Filter"
	case ip.FS_Fuji_ClassicChrome:
		return "Classic Chrome"
	case ip.FS_Fuji_ACROS:
		return "ACROS"
	case ip.FS_Fuji_ACROSYe:
		return "ACROS Ye"
	case ip.FS_Fuji_ACROSR:
		return "ACROS R"
	case ip.FS_Fuji_ACROSG:
		return "ACROS G"
	case ip.FS_Fuji_ETERNA:
		return "ETERNA"
	default:
		return ""
	}
}

func FujiFlashModeAsString(mode ptp.FlashMode) string {
	switch mode {
	case ip.FM_Fuji_On:
		return "on"
	case ip.FM_Fuji_RedEye:
		return "red eye"
	case ip.FM_Fuji_RedEyeOn:
		return "red eye on"
	case ip.FM_Fuji_RedEyeSync:
		return "red eye sync"
	case ip.FM_Fuji_RedEyeRear:
		return "red eye rear"
	case ip.FM_Fuji_SlowSync:
		return "slow sync"
	case ip.FM_Fuji_RearSync:
		return "rear sync"
	case ip.FM_Fuji_Commander:
		return "commander"
	case ip.FM_Fuji_Disabled:
		return "disabled"
	case ip.FM_Fuji_Enabled:
		return "enabled"
	default:
		return FlashModeAsString(mode)
	}
}

func FujiFocusLockAsString(fl ip.FujiFocusLock) string {
	switch fl {
	case ip.FL_Fuji_On:
		return "on"
	case ip.FL_Fuji_Off:
		return "off"
	default:
		return ""
	}
}

func FujiFocusMeteringModeAsString(fmm uint32) string {
	// TODO: what are the 4 msb?
	mask := uint32(0x000000FF)
	x := fmm >> 8 & mask
	y := fmm & mask

	return fmt.Sprintf("%dx%d", x, y)
}

func FujiFocusModeAsString(fm ptp.FocusMode) string {
	switch fm {
	case ip.FCM_Fuji_Single_Auto:
		return "single auto"
	case ip.FCM_Fuji_Continuous_Auto:
		return "continuous auto"
	default:
		return ""
	}
}

func FujiImageAspectRatioAsString(is ip.FujiImageSize) string {
	switch is {
	case ip.IS_Fuji_Small_3x2:
		return "S 3:2"
	case ip.IS_Fuji_Small_16x9:
		return "S 16:9"
	case ip.IS_Fuji_Small_1x1:
		return "S 1:1"
	case ip.IS_Fuji_Medium_3x2:
		return "M 3:2"
	case ip.IS_Fuji_Medium_16x9:
		return "M 16:9"
	case ip.IS_Fuji_Medium_1x1:
		return "M 1:1"
	case ip.IS_Fuji_Large_3x2:
		return "L 3:2"
	case ip.IS_Fuji_Large_16x9:
		return "L 16:9"
	case ip.IS_Fuji_Large_1x1:
		return "L 1:1"
	default:
		return ""
	}
}

func FujiImageQualityAsString(iq ip.FujiImageQuality) string {
	switch iq {
	case ip.IQ_Fuji_Fine:
		return "fine"
	case ip.IQ_Fuji_Normal:
		return "normal"
	case ip.IQ_Fuji_FineAndRAW:
		return "fine + RAW"
	case ip.IQ_Fuji_NormalAndRAW:
		return "normal + RAW"
	case ip.IQ_Fuji_RAW:
		return "RAW"
	default:
		return ""
	}
}

func FujiWhiteBalanceAsString(wb ptp.WhiteBalance) string {
	switch wb {
	case ip.WB_Fuji_Fluorescent1:
		return "fluorescent 1"
	case ip.WB_Fuji_Fluorescent2:
		return "fluorescent 2"
	case ip.WB_Fuji_Fluorescent3:
		return "fluorescent 3"
	case ip.WB_Fuji_Shade:
		return "shade"
	case ip.WB_Fuji_Underwater:
		return "underwater"
	case ip.WB_Fuji_Temperature:
		return "temprerature"
	case ip.WB_Fuji_Custom:
		return "custom"
	default:
		return WhiteBalanceAsString(wb)
	}
}

// TODO: FujiRecModeAsString(rm ip.FujiRecMode)

func FujiSelfTimerAsString(st ip.FujiSelfTimer) string {
	switch st {
	case ip.ST_Fuji_1Sec:
		return "1 second"
	case ip.ST_Fuji_2Sec:
		return "2 seconds"
	case ip.ST_Fuji_5Sec:
		return "5 seconds"
	case ip.ST_Fuji_10Sec:
		return "10 seconds"
	case ip.ST_Fuji_Off:
		return "off"
	default:
		return ""
	}
}
