package ip

import (
	"github.com/google/uuid"
	ipInternal "github.com/malc0mn/ptp-ip/ip/internal"
	"github.com/malc0mn/ptp-ip/ptp"
)

const (
	// This fail reason is returned in the following cases:
	//   - The FriendlyName stored in the camera does not match the FriendlyName being sent. Set the camera to 'change'
	//     so that it will accept a new FriendlyName
	//   - The camera side has timed out waiting for a connection and displays 'not found'. Set the camera to 'retry' to
	//     allow resending the InitCommandRequestPacket.
	// Seems to be an own version of RC_DeviceBusy.
	FR_Fuji_DeviceBusy FailReason = 0x00002019
	// This error is returned when the InitCommandRequestPacket has the wrong protocol version.
	// Seems to be an own version of RC_InvalidParameter.
	FR_Fuji_InvalidParameter FailReason = 0x0000201D

	// This is the Fuji Protocol Version required to construct a valid InitCommandRequestPacket.
	PV_Fuji ProtocolVersion = 0x8F53E4F2
)

// The Fuji version of the PTP/IP InitCommandRequestPacket deviates from the standard. Looking at what is sent 'over the
// wire', we see this sequence in little endian format as the START of the packet, right after the header fields, being
// the Length (4 bytes) and PacketType (4 bytes) fields:
//   [20]byte{
//       0xf2, 0xe4, 0x53, 0x8f,
//       0xad, 0xa5, 0x48, 0x5d,
//       0x87, 0xb2, 0x7f, 0x0b,
//       0xd3, 0xd5, 0xde, 0xd0,
//       0x02, 0x78, 0xa8, 0xc0
//   }
// Referring to the PTP/IP standard which specifies the following:
//    - GUID (16 bytes)
//    - FriendlyName (variable)
//    - ProtocolVersion (4 bytes)
// Fuji seems to have moved the ProtocolVersion field to the top so it is sent first, followed by the GUID and finally
// the FriendlyName field.
// After several attempts, this conclusion stands: the first field is the 4 byte ProtocolVersion field which MUST be set
// to 0x8f53e4f2 for the InitCommandRequestPacket to be considered valid.
type FujiInitCommandRequestPacket struct {
	ProtocolVersion ProtocolVersion
	GUID            uuid.UUID
	FriendlyName    string
}

func (ficrp *FujiInitCommandRequestPacket) PacketType() PacketType {
	return PKT_InitCommandRequest
}

func (ficrp *FujiInitCommandRequestPacket) Payload() []byte {
	return ipInternal.MarshalLittleEndian(ficrp)
}

func (ficrp *FujiInitCommandRequestPacket) GetGUID() uuid.UUID {
	return ficrp.GUID
}

func (ficrp *FujiInitCommandRequestPacket) GetFriendlyName() string {
	return ficrp.FriendlyName
}

func (ficrp *FujiInitCommandRequestPacket) GetProtocolVersion() ProtocolVersion {
	return ficrp.ProtocolVersion
}

func (ficrp *FujiInitCommandRequestPacket) SetProtocolVersion(pv ProtocolVersion) {
	ficrp.ProtocolVersion = pv
}

func NewFujiInitCommandRequestPacket(guid uuid.UUID, friendlyName string) *FujiInitCommandRequestPacket {
	fa := &FujiInitCommandRequestPacket{
		ProtocolVersion: PV_Fuji,
		GUID:            guid,
		FriendlyName:    friendlyName,
	}

	return fa
}

// The Fuji OperationRequestPacket deviates from the PTP/IP standard in several ways:
//   - the packet type should be PKT_OperationRequest, but there is NO packet type sent out in the packet header!
//   - the DataPhase should be uint32 but Fuji uses uint16
type FujiOperationRequestPacket struct {
	DataPhaseInfo uint16
	OperationCode ptp.OperationCode
	TransactionID ptp.TransactionID
	Parameter1 uint32
	Parameter2 uint32
	Parameter3 uint32
	Parameter4 uint32
	Parameter5 uint32
}

func (forp *FujiOperationRequestPacket) PacketType() PacketType {
	return PKT_Invalid
}

func (forp *FujiOperationRequestPacket) Payload() []byte {
	return ipInternal.MarshalLittleEndian(forp)
}

func NewFujiOpenSessionCommand(tid ptp.TransactionID, sid ptp.SessionID) *FujiOperationRequestPacket {
	return &FujiOperationRequestPacket{
		DataPhaseInfo: uint16(DP_NoDataOrDataIn),
		OperationCode: ptp.OC_OpenSession,
		TransactionID: tid,
		Parameter1: uint32(sid),
	}
}

// The Fuji OperationResponsePacket deviates from the PTP/IP standard similarly to the Fuji OperationRequestPacket:
//   - the packet type should be PKT_OperationResponse, but there is NO packet type sent out in the packet header which
//     is, as one can imagine, extremely annoying when parsing the TCP/IP data coming in.
//   - the DataPhase should be uint32 but Fuji uses uint16
type FujiOperationResponsePacket struct {
	DataPhase uint16
	OperationResponseCode ptp.OperationResponseCode
	TransactionID ptp.TransactionID
	Parameter1 uint32
	Parameter2 uint32
	Parameter3 uint32
	Parameter4 uint32
	Parameter5 uint32
}

func (forp *FujiOperationResponsePacket) PacketType() PacketType {
	return PKT_Invalid
}

func (forp *FujiOperationResponsePacket) TotalFixedFieldSize() int {
	return ipInternal.TotalSizeOfFixedFields(forp)
}

// The PTP/IP protocol specifies how to set up the Command/Data connection which should immediately be followed by
// setting up the Event connection. However Fuji wants additional communications before it is satisfied that the
// Command/Data connection is properly setup. This additional initialisation is performed here.
func FujiInitSequence(c *Client) error {
	var err error

	err = c.SendPacketToCmdDataConn(NewFujiOpenSessionCommand(c.TransactionId(), 0x00000001))
	if err != nil {
		return err
	}

	p := new(FujiOperationResponsePacket)
	res, err := c.WaitForPacketFromCmdDataConn(p)
	if err != nil {
		return err
	}

	log.Printf("%v", res)

/*	err = c.SendPacketToCmdDataConn(&FujiOpenSessionPacket{
		DataPhaseInfo: DP_DataOut,
		TransactionID: c.TransactionId(),
	})
	if err != nil {
		return err
	}

	res, err = c.WaitForPacketFromCmdDataConn()
	if err != nil {
		return err
	}

	log.Printf("%v", res)

	c.incrementTransactionId()
*/
	return nil
}