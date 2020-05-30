package main

import (
	"github.com/malc0mn/ptp-ip/ip"
	"os"
	"os/exec"
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	want := "generic"
	if conf.vendor != want {
		t.Errorf("conf.vendor = %s; want %s", conf.vendor, want)
	}

	want = "192.168.0.1"
	if conf.host != want {
		t.Errorf("conf.host = %s; want %s", conf.host, want)
	}

	wantPort := uint16Value(15740)
	if conf.port != wantPort {
		t.Errorf("conf.port = %d; want %d", conf.port, wantPort)
	}

	want = "127.0.0.1"
	if conf.srvAddr != want {
		t.Errorf("conf.srvAddr = %s; want %s", conf.srvAddr, want)
	}

	if conf.srvPort != wantPort {
		t.Errorf("conf.srvPort = %d; want %d", conf.srvPort, wantPort)
	}
}

func TestLoadconfigOk1(t *testing.T) {
	file = "testdata/test_ok1.conf"
	loadConfig()

	want := "Golang test OK1 client"
	if conf.fname != want {
		t.Errorf("loadConfig() fname = %s; want %s", conf.fname, want)
	}

	want = "cca455de-79ac-4b12-9731-91e433a899cf"
	if conf.guid != want {
		t.Errorf("loadConfig() guid = %s; want %s", conf.guid, want)
	}

	want = "fuji"
	if conf.vendor != want {
		t.Errorf("loadConfig() vendor = %s; want %s", conf.host, want)
	}

	want = "192.168.0.2"
	if conf.host != want {
		t.Errorf("loadConfig() host = %s; want %s", conf.host, want)
	}

	wantPort := uint16Value(35740)
	if conf.port != wantPort {
		t.Errorf("loadConfig() port = %d; want %d", conf.port, wantPort)
	}

	wantEnabled := true
	if server != wantEnabled {
		t.Errorf("loadConfig() server = %v; want %v", server, wantEnabled)
	}

	want = "127.0.0.2"
	if conf.srvAddr != want {
		t.Errorf("loadConfig() saddr = %s; want %s", conf.srvAddr, want)
	}

	wantPort = uint16Value(25740)
	if conf.srvPort != wantPort {
		t.Errorf("loadConfig() sport = %d; want %d", conf.srvPort, wantPort)
	}
}

func TestLoadconfigOk2(t *testing.T) {
	conf = &config{
		vendor:  ip.DefaultVendor,
		host:    ip.DefaultIpAddress,
		port:    uint16Value(ip.DefaultPort),
		srvAddr: defaultIp,
		srvPort: uint16Value(ip.DefaultPort),
	}

	file = "testdata/test_ok2.conf"
	loadConfig()

	want := "Golang test OK2 client"
	if conf.fname != want {
		t.Errorf("loadConfig() fname = %s; want %s", conf.fname, want)
	}

	want = "9fe5160c-4951-404d-9505-10baaf725606"
	if conf.guid != want {
		t.Errorf("loadConfig() guid = %s; want %s", conf.guid, want)
	}

	want = "fuji"
	if conf.vendor != want {
		t.Errorf("loadConfig() vendor = %s; want %s", conf.host, want)
	}

	want = "192.168.0.2"
	if conf.host != want {
		t.Errorf("loadConfig() host = %s; want %s", conf.host, want)
	}

	wantPort := uint16Value(15740)
	if conf.port != wantPort {
		t.Errorf("loadConfig() port = %d; want %d", conf.port, wantPort)
	}

	wantPort = uint16Value(55740)
	if conf.cport != wantPort {
		t.Errorf("loadConfig() cport = %d; want %d", conf.cport, wantPort)
	}

	wantPort = uint16Value(55741)
	if conf.eport != wantPort {
		t.Errorf("loadConfig() eport = %d; want %d", conf.eport, wantPort)
	}

	wantPort = uint16Value(55742)
	if conf.sport != wantPort {
		t.Errorf("loadConfig() sport = %d; want %d", conf.sport, wantPort)
	}

	wantEnabled := true
	if server != wantEnabled {
		t.Errorf("loadConfig() server = %v; want %v", server, wantEnabled)
	}

	want = "127.0.0.3"
	if conf.srvAddr != want {
		t.Errorf("loadConfig() saddr = %s; want %s", conf.srvAddr, want)
	}

	wantPort = uint16Value(35740)
	if conf.srvPort != wantPort {
		t.Errorf("loadConfig() sport = %d; want %d", conf.srvPort, wantPort)
	}
}

func TestLoadConfigWrongPath(t *testing.T) {
	if os.Getenv("CONF_FAIL") == "1" {
		file = "does-not-exist.conf"
		loadConfig()
		return
	}

	want := 102
	cmd := exec.Command(os.Args[0], "-test.run=TestLoadConfigWrongPath")
	cmd.Env = append(os.Environ(), "CONF_FAIL=1")
	err := cmd.Run()
	if e, ok := err.(*exec.ExitError); ok && !e.Success() && e.ExitCode() != want {
		t.Fatalf("loadConfig() ran with err %v, want exit status %d", err, want)
	}
}

func TestLoadConfigFail1(t *testing.T) {
	if os.Getenv("CONF_FAIL") == "1" {
		file = "testdata/test_fail1.conf"
		loadConfig()
		return
	}

	want := 1
	cmd := exec.Command(os.Args[0], "-test.run=TestLoadConfigFail")
	cmd.Env = append(os.Environ(), "CONF_FAIL=1")
	err := cmd.Run()
	if e, ok := err.(*exec.ExitError); ok && !e.Success() && e.ExitCode() != want {
		t.Fatalf("loadConfig() ran with err %v, want exit status %d", err, want)
	}
}

func TestLoadConfigFail2(t *testing.T) {
	if os.Getenv("CONF_FAIL") == "1" {
		file = "testdata/test_fail2.conf"
		loadConfig()
		return
	}

	want := 1
	cmd := exec.Command(os.Args[0], "-test.run=TestLoadConfigFail")
	cmd.Env = append(os.Environ(), "CONF_FAIL=1")
	err := cmd.Run()
	if e, ok := err.(*exec.ExitError); ok && !e.Success() && e.ExitCode() != want {
		t.Fatalf("loadConfig() ran with err %v, want exit status %d", err, want)
	}
}
