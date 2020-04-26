package internal

import (
	"bytes"
	"encoding/binary"
)

func ToBytesLittleEndian(s interface{}) []byte {
	var b bytes.Buffer
	// TODO: add error handling!
	binary.Write(&b, binary.LittleEndian, s)
	return b.Bytes()
}
