package ptp

type Event struct {
	// Indicates the event.
	EventCode EventCode

	// Indicates the SessionID of the session for which the event is relevant. If the event is relevant to all open
	// sessions, this field should be set to 0xFFFFFFFF.
	SessionID SessionID

	// If the event corresponds to a previously initiated transaction, this field shall hold the TransactionID of that
	// operation. If the event is not specific to a particular transaction, this field shall be set to 0xFFFFFFFF.
	// Refer to Clause 9.3.1 for a description of TransactionID.
	TransactionID TransactionID

	// These fields hold the event-specific nth parameter. Events may have at most three parameters. The interpretation
	// of any parameter is dependent upon the EventCode. Any unused parameter fields should be set to 0x00000000. If a
	// parameter holds a value that is less than 32 bits, the lowest significant bits shall be used to store the value,
	// with the most significant bits being set to zeros.
	Parameter1 interface{}
	Parameter2 interface{}
	Parameter3 interface{}
}
