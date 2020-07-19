package ip

import "testing"

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
