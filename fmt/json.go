package fmt

import (
	"encoding/json"
	"github.com/malc0mn/ptp-ip/ptp"
)

type DevicePropDescJSON struct {
	*ptp.DevicePropDesc
}

type ValueLabel struct {
	Value string `json:"value"`
	Label string `json:"label"`
}

type CodeLabel struct {
	Code  string `json:"code"`
	Label string `json:"label"`
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
		DevicePropertyCode  CodeLabel
		DataType            string `json:"dataType"`
		GetSet              bool   `json:"readOnly"`
		FactoryDefaultValue ValueLabel
		CurrentValue        ValueLabel
		FormFlag            string      `json:"formType"`
		Form                interface{} `json:"form"`
	}{
		DevicePropertyCode: CodeLabel{
			Code:  ConvertToHexString(dpdj.DevicePropertyCode),
			Label: DevicePropCodeAsString(dpdj.DevicePropertyCode),
		},
		DataType: DataTypeCodeAsString(dpdj.DataType),
		GetSet:   dpdj.GetSet != ptp.DPD_GetSet,
		FactoryDefaultValue: ValueLabel{
			Value: ConvertToHexString(dpdj.FactoryDefaultValueAsInt64()),
			Label: FujiDevicePropValueAsString(dpdj.DevicePropertyCode, dpdj.FactoryDefaultValueAsInt64()),
		},
		CurrentValue: ValueLabel{
			Value: ConvertToHexString(dpdj.CurrentValueAsInt64()),
			Label: FujiDevicePropValueAsString(dpdj.DevicePropertyCode, dpdj.CurrentValueAsInt64()),
		},
		FormFlag: FormFlagAsString(dpdj.FormFlag),
		Form:     form,
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
	hex := make([]ValueLabel, len(values))
	for i := 0; i < len(values); i++ {
		hex[i] = ValueLabel{
			Value: ConvertToHexString(values[i]),
			Label: FujiDevicePropValueAsString(ef.DevicePropDesc.DevicePropertyCode, values[i]),
		}
	}

	return json.Marshal(&struct {
		SupportedValues []ValueLabel `json:"values"`
	}{
		SupportedValues: hex,
	})
}
