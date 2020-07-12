package ptp

import (
	"encoding/binary"
	"testing"
)

func TestDevicePropCodeAsString(t *testing.T) {
	check := map[DevicePropCode]string{
		DPC_BatteryLevel:             "battery level",
		DPC_FunctionalMode:           "functional mode",
		DPC_ImageSize:                "image size",
		DPC_CompressionSetting:       "compression setting",
		DPC_WhiteBalance:             "white balance",
		DPC_RGBGain:                  "RGB gain",
		DPC_FNumber:                  "F number",
		DPC_FocalLength:              "focal length",
		DPC_FocusDistance:            "focus distance",
		DPC_FocusMode:                "focus mode",
		DPC_ExposureMeteringMode:     "exposure metering mode",
		DPC_FlashMode:                "flash mode",
		DPC_ExposureTime:             "exposure time",
		DPC_ExposureProgramMode:      "exposure program mode",
		DPC_ExposureIndex:            "exposure index",
		DPC_ExposureBiasCompensation: "exposure bias compensation",
		DPC_DateTime:                 "date time",
		DPC_CaptureDelay:             "capture delay",
		DPC_StillCaptureMode:         "still capture mode",
		DPC_Contrast:                 "contrast",
		DPC_Sharpness:                "sharpness",
		DPC_DigitalZoom:              "digital zoom",
		DPC_EffectMode:               "effect mode",
		DPC_BurstNumber:              "burst number",
		DPC_BurstInterval:            "burst interval",
		DPC_TimelapseNumber:          "timelapse number",
		DPC_TimelapseInterval:        "timelapse interval",
		DPC_FocusMeteringMode:        "focus metering mode",
		DPC_UploadURL:                "upload URL",
		DPC_Artist:                   "artist",
		DPC_CopyrightInfo:            "copyright info",
		DevicePropCode(0):            "",
	}

	for code, want := range check {
		got := DevicePropCodeAsString(code)
		if got != want {
			t.Errorf("DevicePropCodeAsString() return = '%s', want '%s'", got, want)
		}
	}
}

func TestDevicePropDesc_SizeOfValueInBytes(t *testing.T) {
	check := map[DataTypeCode]int{
		DTC_INT8:   1,
		DTC_UINT8:  1,
		DTC_INT16:  2,
		DTC_UINT16: 2,
		DTC_INT32:  4,
		DTC_UINT32: 4,
		DTC_INT64:  8,
		DTC_UINT64: 8,
	}

	for code, want := range check {
		dpd := DevicePropDesc{
			DataType: code,
		}
		got := dpd.SizeOfValueInBytes()
		if got != want {
			t.Errorf("SizeOfValueInBytes() return = %d, want %d", got, want)
		}
	}
}

func TestDevicePropDesc_FactoryDefaultValueAsInt64(t *testing.T) {
	b := make([]byte, 2)
	binary.LittleEndian.PutUint16(b, 0x0140)

	dpd := DevicePropDesc{
		FactoryDefaultValue: b,
	}

	got := dpd.FactoryDefaultValueAsInt64()
	want := int64(0x0140)
	if got != want {
		t.Errorf("FactoryDefaultValueAsInt64() return = %d, want %d", got, want)
	}
}

func TestDevicePropDesc_CurrentValueAsInt64(t *testing.T) {
	b := make([]byte, 2)
	binary.LittleEndian.PutUint16(b, 0x0340)

	dpd := DevicePropDesc{
		CurrentValue: b,
	}

	got := dpd.CurrentValueAsInt64()
	want := int64(0x0340)
	if got != want {
		t.Errorf("CurrentValueAsInt64() return = %d, want %d", got, want)
	}
}

func TestRangeForm_MinimumValueAsInt64(t *testing.T) {
	b := make([]byte, 2)
	binary.LittleEndian.PutUint16(b, 0x0040)

	rf := RangeForm{
		MinimumValue: b,
	}

	got := rf.MinimumValueAsInt64()
	want := int64(0x0040)
	if got != want {
		t.Errorf("MinimumValueAsInt64() return = %d, want %d", got, want)
	}
}

func TestRangeForm_MaximumValueAsInt64(t *testing.T) {
	b := make([]byte, 2)
	binary.LittleEndian.PutUint16(b, 0x9040)

	rf := RangeForm{
		MaximumValue: b,
	}

	got := rf.MaximumValueAsInt64()
	want := int64(0x9040)
	if got != want {
		t.Errorf("MaximumValueAsInt64() return = %d, want %d", got, want)
	}
}

func TestRangeForm_StepSizeAsInt64(t *testing.T) {
	b := make([]byte, 2)
	binary.LittleEndian.PutUint16(b, 0x0010)

	rf := RangeForm{
		StepSize: b,
	}

	got := rf.StepSizeAsInt64()
	want := int64(0x0010)
	if got != want {
		t.Errorf("StepSizeAsInt64() return = %d, want %d", got, want)
	}
}

func TestEnumerationForm_SupportedValuesAsInt64Array(t *testing.T) {
	total := 5
	b := make([][]byte, total)
	for i := 0; i < total; i++ {
		b[i] = make([]byte, 2)
		binary.LittleEndian.PutUint16(b[i], uint16(0x0010*i))
	}

	rf := EnumerationForm{
		NumberOfValues:  total,
		SupportedValues: b,
	}

	got := rf.SupportedValuesAsInt64Array()
	want := []int64{0x0000, 0x0010, 0x0020, 0x0030, 0x0040}
	for i := 0; i < total; i++ {
		if got[i] != want[i] {
			t.Errorf("StepSizeAsInt64() %d = %d, want %d", i, got[i], want[i])
		}
	}
}

func TestFormFlagAsString(t *testing.T) {
	check := map[DevicePropFormFlag]string{
		DPF_FormFlag_None:     "none",
		DPF_FormFlag_Range:    "range",
		DPF_FormFlag_Enum:     "enum",
		DevicePropFormFlag(3): "",
	}

	for code, want := range check {
		got := FormFlagAsString(code)
		if got != want {
			t.Errorf("FormFlagAsString() return = '%s', want '%s'", got, want)
		}
	}
}

func TestDataTypeCodeAsString(t *testing.T) {
	check := map[DataTypeCode]string{
		DTC_UNDEF:            "undefined",
		DTC_INT8:             "int8",
		DTC_UINT8:            "uint8",
		DTC_INT16:            "int16",
		DTC_UINT16:           "uint16",
		DTC_INT32:            "int32",
		DTC_UINT32:           "uint32",
		DTC_INT64:            "int64",
		DTC_UINT64:           "uint64",
		DTC_INT128:           "int128",
		DTC_UINT128:          "uint128",
		DTC_AINT8:            "aint8",
		DTC_AUINT8:           "auint8",
		DTC_AINT16:           "aint16",
		DTC_AUINT16:          "auint16",
		DTC_AINT32:           "aint32",
		DTC_AUINT32:          "auint32",
		DTC_AINT64:           "aint64",
		DTC_AUINT64:          "auint64",
		DTC_AINT128:          "aint128",
		DTC_AUINT128:         "auint128",
		DTC_STR:              "string",
		DataTypeCode(0xF000): "",
	}

	for code, want := range check {
		got := DataTypeCodeAsString(code)
		if got != want {
			t.Errorf("DataTypeCodeAsString() return = '%s', want '%s'", got, want)
		}
	}
}
