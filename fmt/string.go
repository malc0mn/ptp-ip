package fmt

import (
	"fmt"
	"github.com/malc0mn/ptp-ip/ptp"
	"strconv"
	"strings"
)

const (
	PRP_Delay             string = "delay"
	PRP_Effect            string = "effect"
	PRP_Exposure          string = "exposure"
	PRP_ExpBias           string = "exp-bias"
	PRP_FlashMode         string = "flashmode"
	PRP_FocusMeteringMode string = "focusmtr"
	PRP_ISO               string = "iso"
	PRP_WhiteBalance      string = "whitebalance"
)

func ConvertToHexString(v interface{}) string {
	return fmt.Sprintf("%#x", v)
}

// HexStringToUint64 will convert a string in hexadecimal notation to an unsigned 64 bit integer. String values may
// start with 0x but this is not mandatory.
func HexStringToUint64(code string, bitSize int) (uint64, error) {
	cod, err := strconv.ParseUint(strings.Replace(code, "0x", "", -1), 16, bitSize)
	if err != nil {
		return 0, fmt.Errorf("error converting: %s", err)
	}

	return cod, nil
}

// TODO: how to do this better without the need to pass in a vendor? This is called from MarshalJSON() which cannot
//  accept parameters.
func DevicePropCodeAsString(code ptp.DevicePropCode) string {
	res := GenericDevicePropCodeAsString(code)
	if res == "" {
		res = FujiDevicePropCodeAsString(code)
	}

	return res
}

// PropNameToDevicePropCode converts a string to a device property code.
func PropNameToDevicePropCode(vendor ptp.VendorExtension, param string) (ptp.DevicePropCode, error) {
	switch vendor {
	case ptp.VE_FujiPhotoFilmCoLtd:
		return FujiPropToDevicePropCode(param)
	default:
		return GenericPropToDevicePropCode(param)
	}
}

func DevicePropValAsString(vendor ptp.VendorExtension, code ptp.DevicePropCode, v int64) string {
	switch vendor {
	case ptp.VE_FujiPhotoFilmCoLtd:
		return FujiDevicePropValueAsString(code, v)
	default:
		return DevicePropValueAsString(code, v)
	}
}
