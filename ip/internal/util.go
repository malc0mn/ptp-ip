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
	"unicode/utf16"
)

func marshal(s interface{}, bo binary.ByteOrder, b *bytes.Buffer) {
	// binary.Write() can only cope with fixed length values so we'll need to handle anything else ourselves.
	if _, hasSession := s.(ptp.Session); binary.Size(s) < 0 || hasSession {
		v := reflect.Indirect(reflect.ValueOf(s))

		for i := 0; i < v.NumField(); i++ {
			// When a dataset has a SessionID, we must skip sending it according to the PTP/IP protocol.
			if v.Type().Field(i).Name == "SessionID" {
				continue
			}

			f := v.Field(i)
			switch f.Kind() {
			case reflect.Struct:
				marshal(f.Addr().Interface(), bo, b)
			case reflect.String:
				// TODO: the PTP protocol sets a limit of 255 characters per string including the terminating null
				//  character. We must still enforce this limit here.
				// A rune in Go is an alias for uint32 but the PTP protocol expects 2 byte Unicode characters according
				// to the ISO10646 standard, so we convert them to utf16 (which is uint16) here.
				binary.Write(b, bo, utf16.Encode([]rune(f.String())))
				// Strings must be null terminated.
				binary.Write(b, bo, uint16(0))
			default:
				binary.Write(b, bo, f.Addr().Interface())
			}
		}
	} else {
		binary.Write(b, bo, s)
	}
}

// Marshal data to a byte array, Little Endian formant, for transport.
func MarshalLittleEndian(s interface{}) []byte {
	var b bytes.Buffer
	marshal(s, binary.LittleEndian, &b)

	return b.Bytes()
}

// We always read using reflection to fill each field of s as we go along. This way, we can fill structs like the
// ptp.OperationResponsePacket which does not necessarily receive all parameter fields 'over the wire'. According to the
// protocol we should, but unfortunately it depends on the vendor's implementation. So we need to make sure this
// unmarshal function is usable by all future implementations.
// The int returned is the left over length of the data that has NOT been unmarshalled. It is the responsibility of the
// caller to handle it.
func unmarshal(r io.Reader, s interface{}, l int, vs int, bo binary.ByteOrder) (int, error) {
	v := reflect.Indirect(reflect.ValueOf(s))

	for i := 0; i < v.NumField(); i++ {
		// When a dataset has a SessionID, we must skip it since the PTP/IP protocol does not send it.
		if v.Type().Field(i).Name == "SessionID" {
			continue
		}

		f := v.Field(i)
		switch f.Kind() {
		case reflect.Struct:
			var err error
			l, err = unmarshal(r, f.Addr().Interface(), l, vs, bo)
			if err != nil {
				return 0, err
			}
		case reflect.String:
			// The PTP protocol expects 2 byte Unicode characters according to the ISO10646 standard, so we convert
			// them to string here.
			b := make([]uint16, vs / 2)
			if err := binary.Read(r, bo, b); err != nil {
				return 0, err
			}
			// The slice operation happening here is to drop the null terminator.
			f.SetString(string(utf16.Decode(b[:len(b) - 1])))
			l -= vs
		default:
			if err := binary.Read(r, bo, f.Addr().Interface()); err != nil {
				return 0, err
			}
			l -= binary.Size(f.Addr().Interface())
		}

		if l == 0 {
			return l, nil
		}
	}

	return l, nil
}

// Unmarshal a byte array, Little Endian formant, upon reception.
// We need a reader, a destination container, the total expected length and a "variable size" integer indicating the
// variable sized portion of the packet.
// Any data that is left over after reading to s will be returned as as a byte array to be dealt with by the caller.
func UnmarshalLittleEndian(r io.Reader, s interface{}, l int, vs int) ([]byte, error) {
	var xs []byte

	left, err := unmarshal(r, s, l, vs, binary.LittleEndian)
	if left > 0 {
		xs = make([]byte, left)
		binary.Read(r, binary.LittleEndian, &xs)
	}

	return xs, err
}

func TotalSizeOfFixedFields(s interface{}) int {
	tfs := binary.Size(s)

	// The SessionID Field is dropped in the PTP/IP implementation.
	if _, hasSession := s.(ptp.Session); hasSession {
		tfs -= 4
	}

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
// TODO: make this loop cancelable!
func RetryDialer(network, address string, timeout time.Duration) (net.Conn, error) {
	var err error
	var retries = 10
	var wait = 500 * time.Millisecond
	var conn net.Conn

	for {
		conn, err = net.DialTimeout(network, address, timeout)
		// Insane isn't it? No sentinel errors from net.Dial()!
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
