package main

import (
	"encoding/json"
	"fmt"
	"github.com/malc0mn/ptp-ip/ip"
	ptpjson "github.com/malc0mn/ptp-ip/json"
	"github.com/malc0mn/ptp-ip/ptp"
	"log"
)

func formatDeviceInfo(vendor ptp.VendorExtension, data interface{}, f []string) string {
	switch vendor {
	case ptp.VE_FujiPhotoFilmCoLtd:
		return fujiFormatDeviceInfo(data.([]*ptp.DevicePropDesc), f)
	}

	return ""
}

func fujiFormatDeviceInfo(list []*ptp.DevicePropDesc, f []string) string {
log.Printf("%v - %T", list, list)
	switch f[0] {
	case "json":
		var opt string
		if len(f) > 1 {
			opt = f[1]
		}
		return fujiFormatJson(list, opt)
	default:
		return fujiFormatTable(list)
	}
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
			ip.FujiDevicePropValueAsString(dpd.DevicePropertyCode, uint16(dpd.CurrentValueAsInt64())),
			uint16(dpd.CurrentValueAsInt64()),
			dpd.CurrentValue,
			dpd.CurrentValue,
			uint16(dpd.CurrentValueAsInt64()),
		)
		s += "\n"
	}

	return s
}