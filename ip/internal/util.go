package internal

import (
	"bytes"
	"encoding/binary"
)

func ToBytesLittleEndian(s interface{}) []byte {
	var b bytes.Buffer
	binary.Write(&b, binary.LittleEndian, s)
	return b.Bytes()
}