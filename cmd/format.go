package main

import (
	"encoding/json"
	"fmt"
	ptpfmt "github.com/malc0mn/ptp-ip/fmt"
	"github.com/malc0mn/ptp-ip/ip"
	"github.com/malc0mn/ptp-ip/ptp"
)

func formatDeviceProperty(c *ip.Client, param string) (ptp.DevicePropCode, error) {
	var cod ptp.DevicePropCode

	conv, errH := ptpfmt.HexStringToUint64(param, 16)
	if errH != nil {
		var errS error
		cod, errS = ptpfmt.PropNameToDevicePropCode(c.ResponderVendor(), param)
		if errS != nil {
			return 0, fmt.Errorf("%s or %s", errH, errS)
		} else {
			c.Debugf("Converted %s: %#x", param, cod)
		}
	} else {
		cod = ptp.DevicePropCode(conv)
		c.Debugf("Converted uint16: %#x", cod)
	}

	return cod, nil
}

func formatDeviceInfo(vendor ptp.VendorExtension, data interface{}, f []string) string {
	// TODO: how to cleanly eliminate this switch? Move it to the client?
	switch vendor {
	case ptp.VE_FujiPhotoFilmCoLtd:
		return fujiFormatDeviceInfo(data.([]*ptp.DevicePropDesc), f)
		//default:
		// TODO: add generic device info formatting.
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
	lj := make([]*ptpfmt.DevicePropDescJSON, len(list))
	for i := 0; i < len(list); i++ {
		lj[i] = &ptpfmt.DevicePropDescJSON{
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
			ptpfmt.FujiDevicePropCodeAsString(dpd.DevicePropertyCode),
			ptpfmt.FujiDevicePropValueAsString(dpd.DevicePropertyCode, dpd.CurrentValueAsInt64()),
			uint16(dpd.CurrentValueAsInt64()),
			dpd.CurrentValue,
			dpd.CurrentValue,
			uint16(dpd.CurrentValueAsInt64()),
		)
		s += "\n"
	}

	return s
}
