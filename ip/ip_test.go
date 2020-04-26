package ip

import (
	"github.com/google/uuid"
	"testing"
)

func TestNewInitiator(t *testing.T) {
	got := NewInitiator()
	if got.GUID == uuid.Nil {
		t.Errorf("NewInitiator() GUID = %s; want valid non-empty UUID", got.GUID)
	}
	if got.FriendlyName != InitiatorFriendlyName {
		t.Errorf("NewInitiator() Frendlyname = %s; want %s", got.FriendlyName, InitiatorFriendlyName)
	}
}

func TestNewResponder(t *testing.T) {
	got := NewResponder(DefaultIpAddress, DefaultPort)
	if got.IpAddress != DefaultIpAddress {
		t.Errorf("NewResponder() IpAddress = %s; want %s", got.IpAddress, DefaultIpAddress)
	}
	if got.Port != DefaultPort {
		t.Errorf("NewResponder() IpAddress = %d; want %d", got.Port, DefaultPort)
	}
	if got.GUID != uuid.Nil {
		t.Errorf("NewResponder() FriendlyName = %s; want nil", got.GUID)
	}
	if got.FriendlyName != "" {
		t.Errorf("NewResponder() FriendlyName = %s; want nil", got.FriendlyName)
	}
}
