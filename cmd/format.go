package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	ptpfmt "github.com/malc0mn/ptp-ip/fmt"
	"github.com/malc0mn/ptp-ip/ip"
	"github.com/malc0mn/ptp-ip/ptp"
	"strconv"
	"strings"
	"text/tabwriter"
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
	switch vendor {
	case ptp.VE_FujiPhotoFilmCoLtd:
		return fujiFormatDeviceInfo(data.([]*ptp.DevicePropDesc), f)
	default:
		// TODO: add generic device info formatting.
		return ""
	}
}

func fujiFormatDeviceProperty(dpd *ptp.DevicePropDesc, f []string) string {
	if len(f) >= 1 && f[0] == "json" {
		var opt string
		if len(f) > 1 {
			opt = f[1]
		}

		return fujiFormatJson(&ptpfmt.DevicePropDescJSON{
			DevicePropDesc: dpd,
		}, opt)
	}

	return fujiFormatTable(dpd)
}

func fujiFormatDeviceInfo(list []*ptp.DevicePropDesc, f []string) string {
	if len(f) >= 1 && f[0] == "json" {
		var opt string
		if len(f) > 1 {
			opt = f[1]
		}

		return fujiFormatJsonList(list, opt)
	}

	return fujiFormatListAsTable(list)
}

func fujiFormatJsonList(list []*ptp.DevicePropDesc, opt string) string {
	lj := make([]*ptpfmt.DevicePropDescJSON, len(list))
	for i := 0; i < len(list); i++ {
		lj[i] = &ptpfmt.DevicePropDescJSON{
			DevicePropDesc: list[i],
		}
	}

	return fujiFormatJson(lj, opt)
}

func fujiFormatJson(v interface{}, opt string) string {
	var err error
	var res []byte
	if opt == "pretty" {
		res, err = json.MarshalIndent(v, "", "    ")
	} else {
		res, err = json.Marshal(v)
	}
	if err != nil {
		return err.Error()
	}

	return string(res)
}

func fujiFormatTable(dpd *ptp.DevicePropDesc) string {
	w, buf := newTabWriter()
	rows := longHeader()
	rows = append(rows, longPropDescFormat(dpd))
	formatRows(w, rows)

	return "\n" + buf.String()
}

func fujiFormatListAsTable(list []*ptp.DevicePropDesc) string {
	w, buf := newTabWriter()
	rows := shortHeader()
	for _, dpd := range list {
		rows = append(rows, shortPropDescFormat(dpd))
	}
	formatRows(w, rows)

	return "\n" + buf.String()
}

func newTabWriter() (*tabwriter.Writer, *bytes.Buffer) {
	buf := new(bytes.Buffer)

	return tabwriter.NewWriter(buf, 8, 4, 2, ' ', tabwriter.TabIndent), buf
}

func shortHeader() [][]string {
	return [][]string{
		{"DevicePropCode", "Property name", "Value as string", "Value as int64", "Value in hex"},
		{"--------------", "-------------", "---------------", "--------------", "------------"},
	}
}

func shortPropDescFormat(dpd *ptp.DevicePropDesc) []string {
	return []string{
		fmt.Sprintf("%0#4x", dpd.DevicePropertyCode),
		ptpfmt.FujiDevicePropCodeAsString(dpd.DevicePropertyCode),
		ptpfmt.FujiDevicePropValueAsString(dpd.DevicePropertyCode, dpd.CurrentValueAsInt64()),
		strconv.FormatInt(dpd.CurrentValueAsInt64(), 10),
		fmt.Sprintf("%0#8x", dpd.CurrentValueAsInt64()),
	}
}

func longHeader() [][]string {
	return [][]string{
		{"DevicePropCode", "Prop name", "Dflt val as str", "Dflt val as int64", "Dflt val in hex", "Cur val as str", "Cur val as int64", "Cur val in hex", "Vals allowed"},
		{"--------------", "---------", "---------------", "-----------------", "---------------", "--------------", "----------------", "--------------", "------------"},
	}
}

func longPropDescFormat(dpd *ptp.DevicePropDesc) []string {
	var allowed string

	switch form := dpd.Form.(type) {
	case *ptp.RangeForm:
		allowed = fmt.Sprintf(
			"Min: %#x, max: %#x, stepszie: %#x",
			form.MinimumValueAsInt64(), form.MaximumValueAsInt64(), form.StepSizeAsInt64(),
		)
	case *ptp.EnumerationForm:
		vals := form.SupportedValuesAsInt64Array()
		str := make([]string, len(vals))
		for i, val := range vals {
			str[i] = ptpfmt.ConvertToHexString(val)
		}
		allowed = strings.Join(str, ", ")
	}

	return []string{
		fmt.Sprintf("%0#4x", dpd.DevicePropertyCode),
		ptpfmt.FujiDevicePropCodeAsString(dpd.DevicePropertyCode),
		ptpfmt.FujiDevicePropValueAsString(dpd.DevicePropertyCode, dpd.FactoryDefaultValueAsInt64()),
		strconv.FormatInt(dpd.FactoryDefaultValueAsInt64(), 10),
		fmt.Sprintf("%0#8x", dpd.FactoryDefaultValueAsInt64()),
		ptpfmt.FujiDevicePropValueAsString(dpd.DevicePropertyCode, dpd.CurrentValueAsInt64()),
		strconv.FormatInt(dpd.CurrentValueAsInt64(), 10),
		fmt.Sprintf("%0#8x", dpd.CurrentValueAsInt64()),
		allowed,
	}
}

func formatRows(w *tabwriter.Writer, rows [][]string) {
	for _, row := range rows {
		fmt.Fprintln(w, strings.Join(row, "\t"))
	}
	w.Flush()
}
