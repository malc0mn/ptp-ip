package json

import (
	"encoding/json"
	"fmt"
	"github.com/malc0mn/ptp-ip/ip"
	"github.com/malc0mn/ptp-ip/ptp"
)

type DevicePropDescJSON struct {
	*ptp.DevicePropDesc
}

func (dpdj *DevicePropDescJSON) MarshalJSON() ([]byte, error) {
	var form interface{}
	switch dpdj.FormFlag {
	case ptp.DPF_FormFlag_Range:
		form = &RangeFormJSON{
			RangeForm: dpdj.Form.(*ptp.RangeForm),
		}
	case ptp.DPF_FormFlag_Enum:
		form = &EnumerationFormJSON{
			EnumerationForm: dpdj.Form.(*ptp.EnumerationForm),
		}
	}

	return json.Marshal(&struct {
		DevicePropertyCode  string      `json:"devicePropertyCode"`
		DevicePropertyName  string      `json:"deviceProprtyName"`
		DataType            string      `json:"dataType"`
		GetSet              bool        `json:"readOnly"`
		FactoryDefaultValue string      `json:"factoryDefaultValue"`
		CurrentValue        string      `json:"currentValue"`
		FormFlag            string      `json:"formType"`
		Form                interface{} `json:"form"`
	}{
		DevicePropertyCode:  ConvertToHexString(dpdj.DevicePropertyCode),
		DevicePropertyName:  DevicePropertyName(dpdj.DevicePropertyCode),
		DataType:            ptp.DataTypeCodeAsString(dpdj.DataType),
		GetSet:              dpdj.GetSet != ptp.DPD_GetSet,
		FactoryDefaultValue: ConvertToHexString(dpdj.FactoryDefaultValue),
		CurrentValue:        ConvertToHexString(dpdj.CurrentValue),
		FormFlag:            ptp.FormFlagAsString(dpdj.FormFlag),
		Form:                form,
	})
}

type RangeFormJSON struct {
	*ptp.RangeForm
}

func (rfj *RangeFormJSON) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		MinimumValue string `json:"min"`
		MaximumValue string `json:"max"`
		StepSize     string `json:"step"`
	}{
		MinimumValue: ConvertToHexString(rfj.MinimumValueAsInt64()),
		MaximumValue: ConvertToHexString(rfj.MaximumValueAsInt64()),
		StepSize:     ConvertToHexString(rfj.StepSizeAsInt64()),
	})
}

type EnumerationFormJSON struct {
	*ptp.EnumerationForm
}

func (ef *EnumerationFormJSON) MarshalJSON() ([]byte, error) {
	values := ef.SupportedValuesAsInt64Array()
	hex := make([]string, len(values))
	for i := 0; i < len(values); i++ {
		hex[i] = ConvertToHexString(values[i])
	}

	return json.Marshal(&struct {
		SupportedValues []string `json:"values"`
	}{
		SupportedValues: hex,
	})
}

func ConvertToHexString(v interface{}) string {
	return fmt.Sprintf("%#x", v)
}

// TODO: how to do this better without the need to pass in a vendor?
func DevicePropertyName(code ptp.DevicePropCode) string {
	res := ptp.DevicePropCodeAsString(code)
	if res == "" {
		res = ip.FujiDevicePropCodeAsString(code)
	}

	return res
}
