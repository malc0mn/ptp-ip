package fmt

import (
	"github.com/malc0mn/ptp-ip/ip"
	"github.com/malc0mn/ptp-ip/ptp"
	"testing"
)

func TestConvertToHexString(t *testing.T) {
	want := "0x501f"
	got := ConvertToHexString(ptp.DPC_CopyrightInfo)

	if got != want {
		t.Errorf("HexStringToUint64() got = %s; want %s", got, want)
	}
}

func TestHexStringToUint64(t *testing.T) {
	got, err := HexStringToUint64("0x10G4", 13)
	wantE := "error converting: strconv.ParseUint: parsing \"10G4\": invalid syntax"
	if err.Error() != wantE {
		t.Errorf("HexStringToUint64() error = %s; want %s", err, wantE)
	}

	got, err = HexStringToUint64("0x1004", 12)
	wantE = "error converting: strconv.ParseUint: parsing \"1004\": value out of range"
	if err.Error() != wantE {
		t.Errorf("HexStringToUint64() error = %s; want %s", err, wantE)
	}

	got, err = HexStringToUint64("0x1004", 13)
	if err != nil {
		t.Fatal(err)
	}

	wantI := uint64(4100)
	if got != wantI {
		t.Errorf("HexStringToUint64() return = %d; want %d", got, wantI)
	}

	got, err = HexStringToUint64("1004", 13)
	if err != nil {
		t.Errorf("HexStringToUint64() error = %s; want <nil>", err)
	}

	wantI = uint64(4100)
	if got != wantI {
		t.Errorf("HexStringToUint64() return = %d; want %d", got, wantI)
	}
}

func TestDevicePropCodeAsString(t *testing.T) {
	want := "ISO"
	got := DevicePropCodeAsString(ip.DPC_Fuji_ExposureIndex)
	if got != want {
		t.Errorf("DevicePropCodeAsString() got = %s; want %s", got, want)
	}

	got = DevicePropCodeAsString(ptp.DPC_ExposureIndex)
	if got != want {
		t.Errorf("DevicePropCodeAsString() got = %s; want %s", got, want)
	}
}

func TestPropNameToDevicePropCode(t *testing.T) {
	want := ip.DPC_Fuji_ExposureIndex
	got, err := PropNameToDevicePropCode(ptp.VE_FujiPhotoFilmCoLtd, "iso")
	if err != nil {
		t.Fatal(err)
	}
	if got != want {
		t.Errorf("PropNameToDevicePropCode() got = %#x; want %#x", got, want)
	}

	want = ptp.DPC_ExposureIndex
	got, err = PropNameToDevicePropCode(ptp.VendorExtension(0), "iso")
	if err != nil {
		t.Fatal(err)
	}
	if got != want {
		t.Errorf("PropNameToDevicePropCode() got = %#x; want %#x", got, want)
	}

	wantE := "unknown field name 'testing123'"
	got, err = PropNameToDevicePropCode(ptp.VendorExtension(0), "testing123")
	if err.Error() != wantE {
		t.Errorf("PropNameToDevicePropCode() err = %s; want %s", err, wantE)
	}
	if got != 0 {
		t.Errorf("PropNameToDevicePropCode() got = %#x; want %#x", got, want)
	}
}

func TestDevicePropValAsString(t *testing.T) {
	want := "PRO Neg. Hi"
	got := DevicePropValAsString(ptp.VE_FujiPhotoFilmCoLtd, ip.DPC_Fuji_FilmSimulation, 6)

	if got != want {
		t.Errorf("DevicePropValAsString() got = %s; want %s", got, want)
	}

	want = "center spot"
	got = DevicePropValAsString(ptp.VendorExtension(0), ptp.DPC_FocusMeteringMode, 1)

	if got != want {
		t.Errorf("DevicePropValAsString() got = %s; want %s", got, want)
	}
}
