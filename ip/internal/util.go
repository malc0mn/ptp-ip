package internal

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"io"
	"reflect"
	"unicode/utf8"
)

func ToBytesLittleEndian(s interface{}) []byte {
	var b bytes.Buffer

	// binary.Write can only cope with fixed length values so we'll need to handle anything else ourselves.
	if binary.Size(s) < 0 {
		v := reflect.Indirect(reflect.ValueOf(s))

		for i := 0; i < v.NumField(); i++ {
			switch v.Field(i).Kind() {
			case reflect.String:
				// Add one to account for the null char.
				l := utf8.RuneCountInString(v.Field(i).String()) + 1
				r := make([]byte, l)
				// Convert string to runes.
				copy(r, v.Field(i).String())
				binary.Write(&b, binary.LittleEndian, r)
			default:
				binary.Write(&b, binary.LittleEndian, v.Field(i).Addr().Interface())
			}
		}
	} else {
		binary.Write(&b, binary.LittleEndian, s)
	}

	return b.Bytes()
}

func FromBytesLittleEndian(r io.Reader, s interface{}) error {
	// binary.Read can only cope with fixed length values so we'll need to handle anything else ourselves.
	if binary.Size(s) < 0 {
		v := reflect.Indirect(reflect.ValueOf(s))

		for i := 0; i < v.NumField(); i++ {
			switch v.Field(i).Kind() {
			case reflect.String:
				// Strings are null terminated!
				br := bufio.NewReader(r)
				b, err := br.ReadString(0)
				if err != nil {
					return err
				}
				v.Field(i).SetString(b[:len(b)-1]) // -1 to drop the null termination!
			default:
				if err := binary.Read(r, binary.LittleEndian, v.Field(i).Addr().Interface()); err != nil {
					return err
				}
			}
		}
	} else {
		if err := binary.Read(r, binary.LittleEndian, s); err != nil {
			return err
		}
	}

	return nil
}
