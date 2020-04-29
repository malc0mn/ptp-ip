package internal

import (
	"bytes"
	"encoding/binary"
	"reflect"
	"unicode/utf8"
)

func ToBytesLittleEndian(s interface{}) []byte {
	var b bytes.Buffer

	// binary.Write can only cope with fixed length values so we'll need to handle those ourselves.
	if binary.Size(s) < 0 {
		v := reflect.Indirect(reflect.ValueOf(s))

		for i := 0; i < v.NumField(); i++ {
			switch v.Field(i).Kind() {
			case reflect.String:
				// Add one to account for the null char.
				l := utf8.RuneCountInString(v.Field(i).String()) + 1
				r := make([]rune, l)
				// Convert string to runes and add null char for termination.
				copy(r[:], []rune(v.Field(i).String() + "\x00")[:])
				binary.Write(&b, binary.LittleEndian, r)
			default:
				binary.Write(&b, binary.LittleEndian, v.Field(i))
			}
		}
	} else {
		binary.Write(&b, binary.LittleEndian, s)
	}

	return b.Bytes()
}
