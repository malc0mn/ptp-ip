package ptp

// The response phase consists of the transport-specific transmission of a 30-byte response dataset from the Responder
// to the Initiator.
type OperationResponse struct {
	// Indicates the interpretation of the response.
	ResponseCode OperationResponseCode

	// The identifier for the session within which this operation is being responded to. This value is assigned by the
	// Initiator using the OpenSession operation, and should be copied from the OperationRequest dataset that is
	// received by the Responder prior to responding.
	SessionID SessionID

	// The identifier of the particular transaction. This field should be copied from the OperationRequest dataset that
	// is received by the Responder prior to responding.
	TransactionID TransactionID

	// These fields hold the operation-specific nth response parameter. Response datasets may have at most five
	// parameters. The interpretation of any parameter is dependent upon the OperationCode for which the response has
	// been generated, and secondarily may be a function of the particular ResponseCode itself. Any unused parameter
	// fields should be set to 0x00000000. If a parameter holds a value that is less than 32 bits, the lowest
	// significant bits shall be used to store the value, with the most significant bits being set to zeros.
	Parameter1 interface{}
	Parameter2 interface{}
	Parameter3 interface{}
	Parameter4 interface{}
	Parameter5 interface{}
}
