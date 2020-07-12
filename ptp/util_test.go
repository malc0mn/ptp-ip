package ptp

import (
	"encoding/binary"
	"testing"
)

func TestByteArrayToInt64(t *testing.T) {
	b := make([]byte, 4)
	binary.LittleEndian.PutUint16(b, 0x6140)

	got := byteArrayToInt64(b, 2)
	want := int64(0x6140)
	if got != want {
		t.Errorf("byteArrayToInt64() return = %d, want %d", got, want)
	}
}
