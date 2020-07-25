package ptp

import (
	"encoding/binary"
	"testing"
)

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
