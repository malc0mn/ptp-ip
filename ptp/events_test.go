package ptp

import "testing"

func TestEvent_Session(t *testing.T) {
	event := &Event{
		SessionID: 5,
	}

	got := event.Session()
	want := SessionID(5)
	if got != want {
		t.Errorf("Session() return = %d, want %d", got, want)
	}
}
