package internal

import (
	"bytes"
	"encoding/binary"
	"github.com/malc0mn/ptp-ip/ptp"
	"io"
	"net"
	"reflect"
	"strings"
	"time"
	"unicode/utf8"
)

func marshal(s interface{}, bo binary.ByteOrder) []byte {
	var b bytes.Buffer

	_, hasSession := s.(ptp.Session)

	// binary.Write() can only cope with fixed length values so we'll need to handle anything else ourselves.
	// When a packet has a SessionID, we must skip sending it in the PTP/IP protocol.
	if binary.Size(s) < 0 || hasSession {
		v := reflect.Indirect(reflect.ValueOf(s))

		for i := 0; i < v.NumField(); i++ {
			if v.Type().Field(i).Name == "SessionID" {
				continue
			}

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
	// binary.Read() can only cope with fixed length values so we'll need to handle anything else ourselves.
	if binary.Size(s) < 0 {
		v := reflect.Indirect(reflect.ValueOf(s))

		for i := 0; i < v.NumField(); i++ {
			f := v.Field(i)
			switch f.Kind() {
			case reflect.String:
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
// We need a reader, a destination container and a "variable size" integer indicating the variable sized portion
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

// A wrapper around net.Dial() that will retry dialing 10 times on a "connection refused" error with a 500ms delay
// between retries.
func RetryDialer(network, address string, timeout time.Duration) (net.Conn, error) {
	var err error
	var retries = 10
	var wait = 500 * time.Millisecond
	var conn net.Conn

	for {
		conn, err = net.DialTimeout(network, address, timeout)
		// Insane isn't it? No "typed errors" from net.Dial()!
		if err != nil && strings.Contains(err.Error(), "connection refused") && retries > 0 {
			retries--
			time.Sleep(wait)
			continue
		}
		break
	}
	if err != nil {
		return nil, err
	}

	return conn, nil
}
