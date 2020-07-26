package fmt

import (
	"bytes"
	"encoding/json"
	"github.com/malc0mn/ptp-ip/ip"
	"github.com/malc0mn/ptp-ip/ptp"
	"io/ioutil"
	"testing"
)

func TestMarshalJSON(t *testing.T) {
	list := []*ptp.DevicePropDesc{
		{
			DevicePropertyCode:  ptp.DPC_CaptureDelay,
			DataType:            ptp.DTC_UINT16,
			GetSet:              ptp.DPD_GetSet,
			FactoryDefaultValue: []uint8{0x0, 0x0},
			CurrentValue:        []uint8{0x0, 0x0},
			FormFlag:            ptp.DPF_FormFlag_Enum,
			Form: &ptp.EnumerationForm{
				NumberOfValues: 3,
				SupportedValues: [][]uint8{
					{0x00, 0x00}, {0x02, 0x00}, {0x04, 0x00},
				},
			},
		},
		{
			DevicePropertyCode:  ptp.DPC_FlashMode,
			DataType:            ptp.DTC_UINT16,
			GetSet:              ptp.DPD_GetSet,
			FactoryDefaultValue: []uint8{0x2, 0x0},
			CurrentValue:        []uint8{0x9, 0x80},
			FormFlag:            ptp.DPF_FormFlag_Enum,
			Form: &ptp.EnumerationForm{
				NumberOfValues: 2,
				SupportedValues: [][]uint8{
					{0x09, 0x80}, {0x0a, 0x80},
				},
			},
		},
		{
			DevicePropertyCode:  ptp.DPC_WhiteBalance,
			DataType:            ptp.DTC_UINT16,
			GetSet:              ptp.DPD_GetSet,
			FactoryDefaultValue: []uint8{0x2, 0x0},
			CurrentValue:        []uint8{0x2, 0x0},
			FormFlag:            ptp.DPF_FormFlag_Enum,
			Form: &ptp.EnumerationForm{
				NumberOfValues: 10,
				SupportedValues: [][]uint8{
					{0x02, 0x00}, {0x04, 0x00}, {0x06, 0x80}, {0x01, 0x80}, {0x02, 0x80}, {0x03, 0x80}, {0x06, 0x00},
					{0x0a, 0x80}, {0x0b, 0x80}, {0x0c, 0x80},
				},
			},
		},
		{
			DevicePropertyCode:  ptp.DPC_ExposureBiasCompensation,
			DataType:            ptp.DTC_INT16,
			GetSet:              ptp.DPD_GetSet,
			FactoryDefaultValue: []uint8{0x0, 0x0},
			CurrentValue:        []uint8{0x0, 0x0},
			FormFlag:            ptp.DPF_FormFlag_Enum,
			Form: &ptp.EnumerationForm{
				NumberOfValues: 19,
				SupportedValues: [][]uint8{
					{0x48, 0xf4}, {0x95, 0xf5}, {0xe3, 0xf6}, {0x30, 0xf8}, {0x7d, 0xf9}, {0xcb, 0xfa}, {0x18, 0xfc},
					{0x65, 0xfd}, {0xb3, 0xfe}, {0x00, 0x00}, {0x4d, 0x01}, {0x9b, 0x02}, {0xe8, 0x03}, {0x35, 0x05},
					{0x83, 0x06}, {0xd0, 0x07}, {0x1d, 0x09}, {0x6b, 0x0a}, {0xb8, 0x0b},
				},
			},
		},
		{
			DevicePropertyCode:  ip.DPC_Fuji_FilmSimulation,
			DataType:            ptp.DTC_UINT16,
			GetSet:              ptp.DPD_GetSet,
			FactoryDefaultValue: []uint8{0x1, 0x0},
			CurrentValue:        []uint8{0x2, 0x0},
			FormFlag:            ptp.DPF_FormFlag_Enum,
			Form: &ptp.EnumerationForm{
				NumberOfValues: 11,
				SupportedValues: [][]uint8{
					{0x01, 0x00}, {0x02, 0x00}, {0x03, 0x00}, {0x04, 0x00}, {0x05, 0x00}, {0x06, 0x00}, {0x07, 0x00}, {0x08, 0x00},
					{0x09, 0x00}, {0x0a, 0x00}, {0x0b, 0x00},
				},
			},
		},
		{
			DevicePropertyCode:  ip.DPC_Fuji_ExposureIndex,
			DataType:            ptp.DTC_UINT32,
			GetSet:              ptp.DPD_GetSet,
			FactoryDefaultValue: []uint8{0xff, 0xff, 0xff, 0xff},
			CurrentValue:        []uint8{0x0, 0x19, 0x0, 0x80},
			FormFlag:            ptp.DPF_FormFlag_Enum,
			Form: &ptp.EnumerationForm{
				NumberOfValues: 25,
				SupportedValues: [][]uint8{
					{0x90, 0x01, 0x00, 0x80}, {0x20, 0x03, 0x00, 0x80}, {0x40, 0x06, 0x00, 0x80}, {0x80, 0x0c, 0x00, 0x80},
					{0x00, 0x19, 0x00, 0x80}, {0x64, 0x00, 0x00, 0x40}, {0xc8, 0x00, 0x00, 0x00}, {0xfa, 0x00, 0x00, 0x00},
					{0x40, 0x01, 0x00, 0x00}, {0x90, 0x01, 0x00, 0x00}, {0xf4, 0x01, 0x00, 0x00}, {0x80, 0x02, 0x00, 0x00},
					{0x20, 0x03, 0x00, 0x00}, {0xe8, 0x03, 0x00, 0x00}, {0xe2, 0x04, 0x00, 0x00}, {0x40, 0x06, 0x00, 0x00},
					{0xd0, 0x07, 0x00, 0x00}, {0xc4, 0x09, 0x00, 0x00}, {0x80, 0x0c, 0x00, 0x00}, {0xa0, 0x0f, 0x00, 0x00},
					{0x88, 0x13, 0x00, 0x00}, {0x00, 0x19, 0x00, 0x00}, {0x00, 0x32, 0x00, 0x40}, {0x00, 0x64, 0x00, 0x40},
					{0x00, 0xc8, 0x00, 0x40},
				},
			},
		},
		{
			DevicePropertyCode:  ip.DPC_Fuji_RecMode,
			DataType:            ptp.DTC_UINT16,
			GetSet:              ptp.DPD_GetSet,
			FactoryDefaultValue: []uint8{0x1, 0x0},
			CurrentValue:        []uint8{0x1, 0x0},
			FormFlag:            ptp.DPF_FormFlag_Enum,
			Form: &ptp.EnumerationForm{
				NumberOfValues: 2,
				SupportedValues: [][]uint8{
					{0x0, 0x0}, {0x1, 0x0},
				},
			},
		},
		{
			DevicePropertyCode:  ip.DPC_Fuji_FocusMeteringMode,
			DataType:            ptp.DTC_UINT32,
			GetSet:              ptp.DPD_GetSet,
			FactoryDefaultValue: []uint8{0x0, 0x0, 0x0, 0x0},
			CurrentValue:        []uint8{0x2, 0x7, 0x2, 0x3},
			FormFlag:            ptp.DPF_FormFlag_Range,
			Form: &ptp.RangeForm{
				MinimumValue: []uint8{0x00, 0x00, 0x00, 0x00},
				MaximumValue: []uint8{0x07, 0x07, 0x09, 0x10},
				StepSize:     []uint8{0x01, 0x00, 0x00, 0x00},
			},
		},
	}

	for _, f := range list {
		f.Form.SetDevicePropDesc(f)
	}

	lj := make([]*DevicePropDescJSON, len(list))
	for i := 0; i < len(list); i++ {
		lj[i] = &DevicePropDescJSON{
			DevicePropDesc: list[i],
		}
	}

	want, err := ioutil.ReadFile("testdata/plain.json")
	if err != nil {
		t.Fatal(err)
	}

	got, err := json.Marshal(lj)
	if err != nil {
		t.Fatal(err)
	}
	if bytes.Compare(got, want) != 0 {
		t.Errorf("MarshalJSON() got = %s; want %s", got, want)
	}

	want, err = ioutil.ReadFile("testdata/pretty.json")
	if err != nil {
		t.Fatal(err)
	}
	got, err = json.MarshalIndent(lj, "", "    ")
	if err != nil {
		t.Fatal(err)
	}
	if bytes.Compare(got, want) != 0 {
		t.Errorf("MarshalJSON() got = %s; want %s", got, want)
	}
}
