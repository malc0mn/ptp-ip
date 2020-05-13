package main

import (
	"testing"
)

func TestUint16Value(t *testing.T) {
	u := uint16Value(55)

	want := "55"
	got := u.String()
	if got != want {
		t.Errorf("uint16Value String() = %s; want %s", got, want)
	}

	err := u.Set("0")
	if err == nil {
		t.Errorf("uint16Value Set() = %s; want value out of range", err)
	}

	err = u.Set("65536")
	if err == nil {
		t.Errorf("uint16Value Set() = %s; want value out of range", err)
	}
}
