package ptp

import "encoding/binary"

// byteArrayToInt64 converts a byte array to an int64 where l is the number of significant bytes in the byte array.
// Setting l to 0 will cause l to be set to the length of the byte array passed in.
func byteArrayToInt64(b []byte, l int) int64 {
	if l == 0 || l > len(b) {
		l = len(b)
	}

	if l < 8 {
		pad := make([]byte, 8-l)
		b = append(b, pad...)
	}

	// Converting between uint64 and int64 does not change the sign bit, only the way it is interpreted.
	return int64(binary.LittleEndian.Uint64(b))
}
