package internal

import (
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

func unmarshal(r io.Reader, s interface{}, vs int, bo binary.ByteOrder) error {
	// binary.Read can only cope with fixed length values so we'll need to handle anything else ourselves.
	if binary.Size(s) < 0 {
		v := reflect.Indirect(reflect.ValueOf(s))

		for i := 0; i < v.NumField(); i++ {
			f := v.Field(i)
			switch f.Kind() {
			case reflect.String:
				// Strings are null terminated!
				b := make([]byte, vs)
				if err := binary.Read(r, bo, b); err != nil {
					return err
				}
				f.SetString(string(b[:]))
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
// We will need a reader, a destination container and a "variable size" integer indicating the variable sized portion
// of the packet.
func UnmarshalLittleEndian(r io.Reader, s interface{}, vs int) error {
	return unmarshal(r, s, vs, binary.LittleEndian)
}

func TotalSizeOfFixedFields(s interface{}) int {
	tfs := binary.Size(s)
	if tfs < 0 {
		tfs = 0
		v := reflect.Indirect(reflect.ValueOf(s))
		for i := 0; i < v.NumField(); i++ {
			f := v.Field(i)
			switch f.Kind() {
			case reflect.String:
				// Skip string fields, we do not calculate their size.
				continue
			default:
				tfs += binary.Size(f.Addr().Interface())
			}
		}
	}

	return tfs
}
