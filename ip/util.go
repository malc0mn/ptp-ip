package ip

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// HexStringToUint64 will convert a string in hexadecimal notation to an unsigned 64 bit integer. String values may
// start with 0x but this is not mandatory.
func HexStringToUint64(code string, bitSize int) (uint64, error) {
	cod, err := strconv.ParseUint(strings.Replace(code, "0x", "", -1), 16, bitSize)
	if err != nil {
		return 0, errors.New(fmt.Sprintf("error converting: %s", err))
	}

	return cod, nil
}
