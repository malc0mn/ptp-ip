package main

import (
	"fmt"
	"github.com/malc0mn/ptp-ip/ip"
	"github.com/malc0mn/ptp-ip/ptp"
	"log"
)

func printDeviceInfo(vendor ptp.VendorExtension, data interface{}) string {
	switch vendor {
	case ptp.VE_FujiPhotoFilmCoLtd:
		return fujiPrintDeviceInfo(data.([]*ptp.DevicePropDesc))
	}

	return ""
}

func fujiPrintDeviceInfo(list []*ptp.DevicePropDesc) string {
	var s string
log.Printf("%v - %T", list, list)
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
