package ip

import (
	"errors"
	"github.com/google/uuid"
	ipInternal "github.com/malc0mn/ptp-ip/ip/internal"
	"github.com/malc0mn/ptp-ip/ptp"
	"log"
)

const (
	// The application version property code for Fuji.
	DPC_Fuji_AppVersion    ptp.DevicePropCode = 0xDF01
	DPC_Fuji_CameraVersion ptp.DevicePropCode = 0xDF24

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

	OC_Fuji_GetDeviceInfo ptp.OperationCode = 0x902B

	// This is the Fuji Protocol Version required to construct a valid InitCommandRequestPacket.
	PV_Fuji ProtocolVersion = 0x8F53E4F2

	// The response code to a OC_GetDevicePropValue. The first parameter in the packet will hold the property value.
	RC_Fuji_DevicePropValue ptp.OperationResponseCode = 0x1015
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
}

func (forp *FujiOperationResponsePacket) PacketType() PacketType {
	return PKT_Invalid
}

func (forp *FujiOperationResponsePacket) TotalFixedFieldSize() int {
	return ipInternal.TotalSizeOfFixedFields(forp)
}

func (forp *FujiOperationResponsePacket) WasSuccessfull() bool {
	return forp.OperationResponseCode == ptp.RC_OK || forp.OperationResponseCode == RC_Fuji_DevicePropValue
}

func (forp *FujiOperationResponsePacket) ReasonAsError() error {
	return errors.New(ptp.OperationResponseCodeAsString(forp.OperationResponseCode))
}

// TODO: these is just a quick fix so we can establish what the exact handshake is.
//  We will need to rework the data reading process entirely to make it nice and dynamic.
type FujiOperationResponsePacketOne struct {
	DataPhase uint16
	OperationResponseCode ptp.OperationResponseCode
	TransactionID ptp.TransactionID
	Parameter1 uint32
}

func (forp *FujiOperationResponsePacketOne) PacketType() PacketType {
	return PKT_Invalid
}

func (forp *FujiOperationResponsePacketOne) TotalFixedFieldSize() int {
	return ipInternal.TotalSizeOfFixedFields(forp)
}

func (forp *FujiOperationResponsePacketOne) WasSuccessfull() bool {
	return forp.OperationResponseCode == ptp.RC_OK || forp.OperationResponseCode == RC_Fuji_DevicePropValue
}

func (forp *FujiOperationResponsePacketOne) ReasonAsError() error {
	return errors.New(ptp.OperationResponseCodeAsString(forp.OperationResponseCode))
}
// TODO: end

// The PTP/IP protocol specifies how to set up the Command/Data connection which should immediately be followed by
// setting up the Event connection. However Fuji wants additional communications before it is satisfied that the
// Command/Data connection is properly setup. This additional initialisation is performed here.
// The sequence is as follows:
//   - Open a session
//   - Set device property DPC_Fuji_AppVersion to the correct version as expected by the camera
//   - The camera will now prompt the user to acknowledge the client connection, displaying the client name that was
//     set using the InitCommandRequestPacket.
//   - We will wait for 30 seconds for an acknowledgement from the camera which means the user has pressed the 'OK'
//     on the camera.
func FujiInitSequence(c *Client) error {
	var err error

	log.Print("Opening session...")
	err = c.SendPacketToCmdDataConn(NewFujiOpenSessionCommand(c.TransactionId(), 0x00000001))
	if err != nil {
		return err
	}

	p := new(FujiOperationResponsePacket)
	res, err := c.WaitForPacketFromCmdDataConn(p)
	if err != nil {
		return err
	}

	if !p.WasSuccessfull() {
		return p.ReasonAsError()
	}

	log.Printf("%v -- %s", res, p.ReasonAsError())

	c.incrementTransactionId()

	log.Print("Announcing our application version...")
	err = c.SendPacketToCmdDataConn(&FujiOperationRequestPacket{
		DataPhaseInfo: uint16(DP_NoDataOrDataIn),
		OperationCode: ptp.OC_SetDevicePropValue,
		TransactionID: c.TransactionId(),
		Parameter1: uint32(DPC_Fuji_AppVersion),
	})
	if err != nil {
		return err
	}

	err = c.SendPacketToCmdDataConn(&FujiOperationRequestPacket{
		DataPhaseInfo: uint16(DP_DataOut),
		OperationCode: ptp.OC_SetDevicePropValue,
		TransactionID: c.TransactionId(),
		// TODO: Parameter 1 is the application version, we need to make this dynamic so it can be passed in from the
		//   command line. This will differ from camera model to camera model, so we'll probably default this to a high
		//   value that, hopefully, will always work.
		Parameter1: 0x00000003,
	})
	if err != nil {
		return err
	}

	log.Print("Please accept the new connection request on the camera...")
	res, err = c.WaitForPacketFromCmdDataConn(p)
	if err != nil {
		return err
	}

	if !p.WasSuccessfull() {
		return p.ReasonAsError()
	}

	log.Printf("%v -- %s", res, p.ReasonAsError())

	c.incrementTransactionId()

	log.Print("Getting current 'other version'...")
	err = c.SendPacketToCmdDataConn(&FujiOperationRequestPacket{
		DataPhaseInfo: uint16(DP_NoDataOrDataIn),
		OperationCode: ptp.OC_GetDevicePropValue,
		TransactionID: c.TransactionId(),
		Parameter1: uint32(DPC_Fuji_CameraVersion),
	})
	if err != nil {
		return err
	}

	po := new(FujiOperationResponsePacketOne)
	res, err = c.WaitForPacketFromCmdDataConn(po)
	if err != nil {
		return err
	}

	if !po.WasSuccessfull() {
		return po.ReasonAsError()
	}

	log.Printf("%v -- %s", res, po.ReasonAsError())

	log.Printf("Current 'other version' as communicated by the camera: %#x", po.Parameter1)

	res, err = c.WaitForPacketFromCmdDataConn(p)
	if err != nil {
		return err
	}

	if !p.WasSuccessfull() {
		return p.ReasonAsError()
	}

	log.Printf("%v -- %s", res, p.ReasonAsError())

	c.incrementTransactionId()

	log.Printf("Setting current 'other version' to: %#x", 0x000200ff)
	err = c.SendPacketToCmdDataConn(&FujiOperationRequestPacket{
		DataPhaseInfo: uint16(DP_NoDataOrDataIn),
		OperationCode: ptp.OC_SetDevicePropValue,
		TransactionID: c.TransactionId(),
		Parameter1: uint32(DPC_Fuji_CameraVersion),
	})
	if err != nil {
		return err
	}

	err = c.SendPacketToCmdDataConn(&FujiOperationRequestPacket{
		DataPhaseInfo: uint16(DP_DataOut),
		OperationCode: ptp.OC_SetDevicePropValue,
		TransactionID: c.TransactionId(),
		// TODO: Parameter 1 is the exact same as we received from the camera in the previous request! It is also some
		//   sort of version that needs to be 'high enough' for the camera to accept us.
		Parameter1: 0x000200ff,
	})
	if err != nil {
		return err
	}

	res, err = c.WaitForPacketFromCmdDataConn(p)
	if err != nil {
		return err
	}

	if !p.WasSuccessfull() {
		return p.ReasonAsError()
	}

	log.Printf("%v -- %s", res, p.ReasonAsError())

	c.incrementTransactionId()

	log.Print("Requesting device info...")
	err = c.SendPacketToCmdDataConn(&FujiOperationRequestPacket{
		DataPhaseInfo: uint16(DP_NoDataOrDataIn),
		OperationCode: OC_Fuji_GetDeviceInfo,
		TransactionID: c.TransactionId(),
	})
	if err != nil {
		return err
	}

	res, err = c.WaitForPacketFromCmdDataConn(p)
	if err != nil {
		return err
	}

	if !p.WasSuccessfull() {
		return p.ReasonAsError()
	}

	log.Printf("%v -- %s", res, p.ReasonAsError())

	c.incrementTransactionId()

	log.Print("Initiating open capture...")
	err = c.SendPacketToCmdDataConn(&FujiOperationRequestPacket{
		DataPhaseInfo: uint16(DP_NoDataOrDataIn),
		OperationCode: ptp.OC_InitiateOpenCapture,
		TransactionID: c.TransactionId(),
	})
	if err != nil {
		return err
	}

	res, err = c.WaitForPacketFromCmdDataConn(p)
	if err != nil {
		return err
	}

	if !p.WasSuccessfull() {
		return p.ReasonAsError()
	}

	log.Printf("%v -- %s", res, p.ReasonAsError())

	return nil
}