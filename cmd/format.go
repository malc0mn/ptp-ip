package main

import (
	"encoding/json"
	"fmt"
	"github.com/malc0mn/ptp-ip/ip"
	ptpjson "github.com/malc0mn/ptp-ip/json"
	"github.com/malc0mn/ptp-ip/ptp"
)

func formatDevicePropVal(vendor ptp.VendorExtension, code ptp.DevicePropCode, v int64) string {
	switch vendor {
	case ptp.VE_FujiPhotoFilmCoLtd:
		return ip.FujiDevicePropValueAsString(code, v)
	default:
		return ptp.DevicePropValueAsString(code, v)
	}
}

// TODO: add generic device info formatting.
func formatDeviceInfo(vendor ptp.VendorExtension, data interface{}, f []string) string {
	switch vendor {
	case ptp.VE_FujiPhotoFilmCoLtd:
		return fujiFormatDeviceInfo(data.([]*ptp.DevicePropDesc), f)
	}

	return ""
}

func fujiFormatDeviceInfo(list []*ptp.DevicePropDesc, f []string) string {
	if len(f) >= 1 && f[0] == "json" {
		var opt string
		if len(f) > 1 {
			opt = f[1]
		}
		return fujiFormatJson(list, opt)
	}

	return fujiFormatTable(list)
}

func fujiFormatJson(list []*ptp.DevicePropDesc, opt string) string {
	lj := make([]*ptpjson.DevicePropDescJSON, len(list))
	for i := 0; i < len(list); i++ {
		lj[i] = &ptpjson.DevicePropDescJSON{
			DevicePropDesc: list[i],
		}
	}

	var err error
	var res []byte
	if opt == "pretty" {
		res, err = json.MarshalIndent(lj, "", "    ")
	} else {
		res, err = json.Marshal(lj)
	}
	if err != nil {
		return err.Error()
	}

	return string(res)
}

// TODO: write as ASCII table.
func fujiFormatTable(list []*ptp.DevicePropDesc) string {
	var s string

	for _, dpd := range list {
		s += fmt.Sprintf("%s: %s || %d - %v - %#x - %#x",
			ip.FujiDevicePropCodeAsString(dpd.DevicePropertyCode),
			ip.FujiDevicePropValueAsString(dpd.DevicePropertyCode, dpd.CurrentValueAsInt64()),
			uint16(dpd.CurrentValueAsInt64()),
			dpd.CurrentValue,
			dpd.CurrentValue,
			uint16(dpd.CurrentValueAsInt64()),
		)
		s += "\n"
	}

	return s
}
