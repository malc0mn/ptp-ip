package main

import (
	"os"
	"os/exec"
	"testing"
)

func TestLoadconfigOk(t *testing.T) {
	file = "testdata/test.conf"
	loadConfig()

	want := "Golang test client"
	if conf.fname != want {
		t.Errorf("loadConfig() fname = %s; want %s", conf.fname, want)
	}

	want = "cca455de-79ac-4b12-9731-91e433a899cf"
	if conf.guid != want {
		t.Errorf("loadConfig() guid = %s; want %s", conf.guid, want)
	}

	want = "192.168.0.2"
	if conf.host != want {
		t.Errorf("loadConfig() host = %s; want %s", conf.host, want)
	}

	wantPort := uint16Value(1740)
	if conf.port != wantPort {
		t.Errorf("loadConfig() port = %d; want %d", conf.port, wantPort)
	}

	wantEnabled := true
	if server != wantEnabled {
		t.Errorf("loadConfig() server = %v; want %v", server, wantEnabled)
	}

	want = "127.0.0.2"
	if conf.saddr != want {
		t.Errorf("loadConfig() saddr = %s; want %s", conf.saddr, want)
	}

	wantPort = uint16Value(25740)
	if conf.sport != wantPort {
		t.Errorf("loadConfig() sport = %d; want %d", conf.sport, wantPort)
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

func TestLoadConfigFail(t *testing.T) {
	if os.Getenv("CONF_FAIL") == "1" {
		file = "testdata/test_fail.conf"
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