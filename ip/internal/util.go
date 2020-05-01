package internal

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"io"
	"reflect"
	"unicode/utf8"
)

func marshal(s interface{}, bo binary.ByteOrder) []byte {
	var b bytes.Buffer

	// binary.Write can only cope with fixed length values so we'll need to handle anything else ourselves.
	if binary.Size(s) < 0 {
		v := reflect.Indirect(reflect.ValueOf(s))

		for i := 0; i < v.NumField(); i++ {
			f := v.Field(i)
			switch f.Kind() {
			case reflect.String:
				// Add one to account for the null char.
				l := utf8.RuneCountInString(f.String()) + 1
				r := make([]byte, l)
				// Convert string to runes.
				copy(r, f.String())
				binary.Write(&b, bo, r)
			default:
				binary.Write(&b, bo, f.Addr().Interface())
			}
		}
	} else {
		binary.Write(&b, bo, s)
	}

	return b.Bytes()
}

// Marshal data to a byte array, Little Endian formant, for transport.
func MarshalLittleEndian(s interface{}) []byte {
	return marshal(s, binary.LittleEndian)
}

func unmarshal(r io.Reader, s interface{}, bo binary.ByteOrder) error {
	// binary.Read can only cope with fixed length values so we'll need to handle anything else ourselves.
	if binary.Size(s) < 0 {
		v := reflect.Indirect(reflect.ValueOf(s))

		for i := 0; i < v.NumField(); i++ {
			f := v.Field(i)
			switch f.Kind() {
			case reflect.String:
				// Strings are null terminated!
				br := bufio.NewReader(r)
				b, err := br.ReadString(0)
				if err != nil {
					return err
				}
				f.SetString(b[:len(b)-1]) // -1 to drop the null termination!
			default:
				if err := binary.Read(r, bo, f.Addr().Interface()); err != nil {
					return err
				}
			}
		}
	} else {
		if err := binary.Read(r, bo, s); err != nil {
			return err
		}
	}

	return nil
}

// Unmarshal a byte array, Little Endian formant, upon reception.
func UnmarshalLittleEndian(r io.Reader, s interface{}) error {
	return unmarshal(r, s, binary.LittleEndian)
}
