package main

import (
	"github.com/malc0mn/ptp-ip/ip"
	"github.com/malc0mn/ptp-ip/ptp"
	"testing"
)

func TestFormatDeviceProperty(t *testing.T) {
	c, err := ip.NewClient(ip.DefaultVendor, ip.DefaultIpAddress, ip.DefaultPort, "", "", ip.LevelSilent)
	if err != nil {
		t.Fatal(err)
	}

	want := ptp.DevicePropCode(0)
	wantE := "error converting: strconv.ParseUint: parsing \"test\": invalid syntax or unknown field name 'test'"
	got, err := formatDeviceProperty(c, "test")
	if err.Error() != wantE {
		t.Errorf("formatDeviceProperty() error = %s; want %s", err, wantE)
	}
	if got != want {
		t.Errorf("formatDeviceProperty() got %#x; want %#x", got, want)
	}

	want = ptp.DevicePropCode(0x5005)
	got, err = formatDeviceProperty(c, "0x5005")
	if err != nil {
		t.Errorf("formatDeviceProperty() error = %s; want <nil>", err)
	}
	if got != want {
		t.Errorf("formatDeviceProperty() got %#x; want %#x", got, want)
	}

	want = ptp.DevicePropCode(0x500f)
	got, err = formatDeviceProperty(c, "iso")
	if err != nil {
		t.Errorf("formatDeviceProperty() error = %s; want <nil>", err)
	}
	if got != want {
		t.Errorf("formatDeviceProperty() got %#x; want %#x", got, want)
	}

	c, err = ip.NewClient("fuji", ip.DefaultIpAddress, ip.DefaultPort, "", "", ip.LevelSilent)
	if err != nil {
		t.Fatal(err)
	}

	want = ptp.DevicePropCode(0xd02a)
	got, err = formatDeviceProperty(c, "iso")
	if err != nil {
		t.Errorf("formatDeviceProperty() error = %s; want <nil>", err)
	}
	if got != want {
		t.Errorf("formatDeviceProperty() got %#x; want %#x", got, want)
	}
}
