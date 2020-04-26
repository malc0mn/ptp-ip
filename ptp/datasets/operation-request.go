package ptp

import ptp "github.com/malc0mn/ptp-ip/ptp/consts"

// Each session has a SessionID that consists of one device-unique 32-bit unsigned integer (UINT32). SessionIDs are
// assigned by the Initiator as a parameter to the OpenSession operation, and must be non-zero.
type SessionID uint32

// Each transaction within a session has a unique transaction identifier called TransactionID that is a session-unique
// 32-bit unsigned integer (UINT32). TransactionIDs are continuous sequences in numerical order starting from
// 0x00000001. The TransactionID used for the OpenSession operation shall be 0x00000000. The first operation issued by
// an Initiator after an OpenSession operation has a TransactionID of 0x00000001, the second operation has a
// TransactionID of 0x00000002, etc. The TransactionID of 0xFFFFFFFF is not valid, and is reserved for context-specific
// meanings. The presence of TransactionID allows asynchronous events to refer to specific previously initiated
// operations. If this field reaches its maximum value (0xFFFFFFFE), the device should "rollover" to 0x00000001.
// TransactionIDs allow events to refer to particular operation requests, allow correspondence between data objects and
// their describing datasets, and aid in debugging.
type TransactionID uint32

// The operation request phase consists of the transport-specific transmission of a 30-byte operation dataset from the
// Initiator to the Responder.
type OperationRequest struct {
	// The code indicating which operation is being initiated.
	OperationCode ptp.OperationCode

	// The identifier for the session within which this operation is being initiated. This value is assigned by the
	// Initiator using the OpenSession operation. This field should be set to 0x00000000 for operations that do not
	// occur within a session, and for the OpenSession OperationRequest dataset.
	SessionID SessionID

	// The identifier of this particular transaction. This value shall be a value that is unique within a particular
	// session, and shall increment by one for each subsequent transaction. Refer to Clause 9.3.1 for a description of
	// transaction identifiers. This field should be set to 0x00000000 for the OpenSession operation.
	TransactionID TransactionID

	// These fields hold the operation-specific nth parameter. Operations may have at most five parameters. The
	// interpretation of any parameter is dependent upon the OperationCode. Any unused parameter fields should be set to
	// 0x00000000. If a parameter holds a value that is less than 32 bits, the lowest significant bits shall be used to
	// store the value, with the most significant bits being set to zeros.
	Parameter1 interface{}
	Parameter2 interface{}
	Parameter3 interface{}
	Parameter4 interface{}
	Parameter5 interface{}
}