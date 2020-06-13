package ip

import (
	"encoding/binary"
	"errors"
	"github.com/google/uuid"
	ipInternal "github.com/malc0mn/ptp-ip/ip/internal"
	"github.com/malc0mn/ptp-ip/ptp"
	"strings"
)

const (
	// DPC_Fuji_UseInitSequence indicates the initialisation sequence being used. It MUST be set by the Initiator during
	// the initialisation sequence and depending on it's value, will require a different init sequence to be used.
	// See PM_Fuji_InitSequence for further info.
	DPC_Fuji_UseInitSequence ptp.DevicePropCode = 0xDF01
	// DPC_Fuji_AppVersion indicates the minium application version the camera will accept. It MUST be set by the
	// Initiator during the initialisation sequence. As soon as this is done, the camera will acknowledge the client and
	// store the client's friendly name to allow future connections without the need for a confirmation.
	DPC_Fuji_AppVersion ptp.DevicePropCode = 0xDF24

	DPC_Fuji_FilmSimulation     ptp.DevicePropCode = 0xD001
	DPC_Fuji_ImageFormat        ptp.DevicePropCode = 0xD018
	DPC_Fuji_RecmodeEnable      ptp.DevicePropCode = 0xD019
	DPC_Fuji_CommandDial        ptp.DevicePropCode = 0xD028
	DPC_Fuji_Iso                ptp.DevicePropCode = 0xD02A
	DPC_Fuji_MovieIso           ptp.DevicePropCode = 0xD02B
	DPC_Fuji_FocusPoint         ptp.DevicePropCode = 0xD17C
	DPC_Fuji_FocusLock          ptp.DevicePropCode = 0xD209
	DPC_Fuji_DeviceError        ptp.DevicePropCode = 0xD21B
	DPC_Fuji_ImageSpaceSD       ptp.DevicePropCode = 0xD229
	DPC_Fuji_MovieRemainingRime ptp.DevicePropCode = 0xD22A
	DPC_Fuji_ShutterSpeed       ptp.DevicePropCode = 0xD240
	DPC_Fuji_ImageAspect        ptp.DevicePropCode = 0xD241
	DPC_Fuji_BatteryLevel       ptp.DevicePropCode = 0xD242

	// FR_Fuji_DeviceBusy is returned in the following cases:
	//   - The FriendlyName stored in the camera does not match the FriendlyName being sent. Set the camera to 'change'
	//     so that it will accept a new FriendlyName
	//   - The camera side has timed out waiting for a connection and displays 'not found'. Set the camera to 'retry' to
	//     allow resending the InitCommandRequestPacket.
	// Seems to be an own version of RC_DeviceBusy.
	FR_Fuji_DeviceBusy FailReason = 0x00002019
	// FR_Fuji_InvalidParameter is returned when the InitCommandRequestPacket has the wrong protocol version.
	// Seems to be an own version of RC_InvalidParameter.
	FR_Fuji_InvalidParameter FailReason = 0x0000201D

	// OC_Fuji_GetDeviceInfo is not really the same as OC_GetDeviceInfo in that it returns a list of device
	OC_Fuji_GetDeviceInfo ptp.OperationCode = 0x902B

	// For convenience: to be used with FujiSendOperationRequest() when no parameter is required for the operation.
	PM_Fuji_NoParam = 0x00000000
	// When this parameter is 'too low', the camera will complain about the application version being 'the previous
	// version' and requests to 'upgrade the app'.
	// After multiple experiments, this parameter will affect the initialisation sequence being used.
	// On the X-T1 (firmware v5.51):
	//   - 0x00000003 IS accepted but the init sequence used here does not seem to work and probably needs some
	//     unknown command(s) to 'finish it off' correctly.
	//   - 0x00000004 SEEMS not to be accepted, i.e. a client confirmation prompt is never displayed by the camera. So
	//     it MIGHT work, but could expect a different set of commands.
	//   - 0x00000005 hits the sweet spot and the init sequence we use completes nicely.
	PM_Fuji_InitSequence = 0x00000005
	// When this parameter is 'too low', the camera will also complain about the application version being 'the previous
	// version' and requests to 'upgrade the app'. However, it does NOT affect the initialisation sequence at all.
	// The value here is that of the X-T1 on firmware version v5.51. We're not using it through this fixed value
	// anymore, but we now get it from the camera and confirm it by setting it to what the camera reports in the hope
	// that this will be future proof and we do not need to to adjust it ever again.
	PM_Fuji_AppVersion = 0x00020001

	// This is the Fuji Protocol Version required to construct a valid InitCommandRequestPacket.
	PV_Fuji ProtocolVersion = 0x8F53E4F2

	// RC_Fuji_DevicePropValue is the response code to a OC_GetDevicePropValue. The first parameter in the packet will
	// hold the property value.
	RC_Fuji_DevicePropValue = ptp.OperationResponseCode(ptp.OC_GetDevicePropValue)
	// RC_Fuji_DeviceInfo is the response code to OC_Fuji_GetDeviceInfo.
	RC_Fuji_DeviceInfo      = ptp.OperationResponseCode(OC_Fuji_GetDeviceInfo)
)

// FujiInitCommandRequestPacket is the Fuji version of the PTP/IP InitCommandRequestPacket which deviates from the
// standard. Looking at what is sent 'over the wire', we see this sequence in little endian format as the START of the
// packet, right after the header fields, being the Length (4 bytes) and PacketType (4 bytes) fields:
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

func NewFujiInitCommandRequestPacket(guid uuid.UUID, friendlyName string) InitCommandRequestPacket {
	return &FujiInitCommandRequestPacket{
		ProtocolVersion: PV_Fuji,
		GUID:            guid,
		FriendlyName:    friendlyName,
	}
}

// NewFujiInitEventRequestPacket returns nil because Fuji does not require the Event channel to be initialised. This
// will skip any further event channel initialisation.
func NewFujiInitEventRequestPacket(connNum uint32) InitEventRequestPacket {
	return nil
}

// FujiOperationRequestPacket deviates from the PTP/IP standard in several ways:
//   - the packet type should be PKT_OperationRequest, but there is NO packet type sent out in the packet header (which
//     is really annoying)!
//   - the DataPhase should be uint32 but Fuji uses uint16
type FujiOperationRequestPacket struct {
	DataPhaseInfo uint16
	OperationCode ptp.OperationCode
	TransactionID ptp.TransactionID
	Parameter1    uint32
	Parameter2    uint32
	Parameter3    uint32
	Parameter4    uint32
	Parameter5    uint32
}

func (forp *FujiOperationRequestPacket) PacketType() PacketType {
	return PKT_Invalid
}

func (forp *FujiOperationRequestPacket) Payload() []byte {
	return ipInternal.MarshalLittleEndian(forp)
}

// FujiOperationResponsePacket deviates from the PTP/IP standard similarly to FujiOperationRequestPacket:
//   - the packet type should be PKT_OperationResponse, but there is NO packet type sent out in the packet header which
//     is, as one can imagine, extremely annoying when parsing the TCP/IP data coming in.
//   - the DataPhase should be uint32 but Fuji uses uint16
type FujiOperationResponsePacket struct {
	DataPhase             uint16
	OperationResponseCode ptp.OperationResponseCode
	TransactionID         ptp.TransactionID
	Parameter1            uint32
}

func (forp *FujiOperationResponsePacket) PacketType() PacketType {
	return PKT_Invalid
}

func (forp *FujiOperationResponsePacket) TotalFixedFieldSize() int {
	return ipInternal.TotalSizeOfFixedFields(forp)
}

// TODO: make this better, obviously...
func (forp *FujiOperationResponsePacket) WasSuccessfull() bool {
	return forp.OperationResponseCode == ptp.RC_OK ||
		forp.OperationResponseCode == ptp.RC_SessionAlreadyOpen ||
		forp.OperationResponseCode == RC_Fuji_DevicePropValue ||
		forp.OperationResponseCode == RC_Fuji_DeviceInfo
}

func (forp *FujiOperationResponsePacket) ReasonAsError() error {
	return errors.New(ptp.OperationResponseCodeAsString(forp.OperationResponseCode))
}

// FujiInitCommandDataConn initialises the Fuji command/data connection.
// The PTP/IP protocol specifies how to set up the command/data connection which should immediately be followed by
// setting up the event connection. However Fuji wants additional communications before it is satisfied that the
// command/data connection is properly setup. This additional initialisation is performed here.
// The sequence is as follows:
//   1. Open a session.
//   2. Set device property DPC_Fuji_UseInitSequence to the correct number of the init sequence being used by the
//      Initiator.
//   3. If the client name differs from the one stored, the Responder will now prompt the user to acknowledge the client
//      connection, displaying the client name that was communicated using the InitCommandRequestPacket.
//   4. We will wait for 30 seconds for an acknowledgement from the Responder which means the user has pressed the 'OK'
//      button on the camera.
//   5. Next we will request the value of device property DPC_Fuji_AppVersion which holds the current minimal
//      application version supported by the Responder and we will simply acknowledge it by setting it to the same
//      value.
//      This way we will always support any future versions as required by the firmware; unless of course a newer init
//      sequence should be required.
//   6. Finally, we send the operation request OC_InitiateOpenCapture which makes the Responder hand over control to the
//      Initiator. This also opens up the event connection port used by Fuji on port 55741 so we can connect to it and
//      complete the init sequence there.
func FujiInitCommandDataConn(c *Client) error {
	// The first part of the sequence is according to the PTP/IP standard, save for the different packet format.
	if err := GenericInitCommandDataConn(c); err != nil {
		return err
	}

	c.log.Print("Opening a session...")
	if _, err := FujiSendOperationRequest(c, ptp.OC_OpenSession, 0x00000001); err != nil {
		return err
	}

	c.log.Print("Setting correct init sequence number...")
	c.log.Printf("Should you be prompted, please accept the new connection request on the %s.", c.ResponderFriendlyName())
	if err := FujiSetDeviceProperty(c, DPC_Fuji_UseInitSequence, PM_Fuji_InitSequence); err != nil {
		return err
	}

	c.log.Print("Getting current minimum application version...")
	val, err := FujiGetDevicePropertyValue(c, DPC_Fuji_AppVersion)
	if err != nil {
		return err
	}
	c.log.Printf("Acknowledging current minimal application version as communicated by the %s: %#x", c.ResponderFriendlyName(), val)
	if err := FujiSetDeviceProperty(c, DPC_Fuji_AppVersion, val); err != nil {
		return err
	}

	c.log.Print("Initiating open capture...")
	if _, err := FujiSendOperationRequest(c, ptp.OC_InitiateOpenCapture, PM_Fuji_NoParam); err != nil {
		return err
	}

	return nil
}

// Sets a device property to the given value.
func FujiSetDeviceProperty(c *Client, code ptp.DevicePropCode, val uint32) error {
	c.incrementTransactionId()

	if err := c.SendPacketToCmdDataConn(&FujiOperationRequestPacket{
		DataPhaseInfo: uint16(DP_NoDataOrDataIn),
		OperationCode: ptp.OC_SetDevicePropValue,
		TransactionID: c.TransactionId(),
		Parameter1:    uint32(code),
	}); err != nil {
		return err
	}

	if err := c.SendPacketToCmdDataConn(&FujiOperationRequestPacket{
		DataPhaseInfo: uint16(DP_DataOut),
		OperationCode: ptp.OC_SetDevicePropValue,
		TransactionID: c.TransactionId(),
		Parameter1:    val,
	}); err != nil {
		return err
	}

	p := new(FujiOperationResponsePacket)
	if _, err := c.WaitForPacketFromCmdDataConn(p); err != nil {
		return err
	}

	if !p.WasSuccessfull() {
		return p.ReasonAsError()
	}

	return nil
}

// Get the value for the given device property.
// TODO: add third parameter to indicate how many parameters from the response object are expected?
func FujiGetDevicePropertyValue(c *Client, dpc ptp.DevicePropCode) (uint32, error) {
	var val uint32
	var err error

	// First we get the actual value from the Responder.
	if val, err = FujiSendOperationRequest(c, ptp.OC_GetDevicePropValue, uint32(dpc)); err != nil {
		return 0, err
	}

	// Next we also get sort of an 'end of data' packet which is of no real use to us save for additional error
	// handling.
	p := new(FujiOperationResponsePacket)
	if _, err := c.WaitForPacketFromCmdDataConn(p); err != nil {
		return 0, err
	}

	if !p.WasSuccessfull() {
		return 0, p.ReasonAsError()
	}

	return val, nil
}

// Send an operation request to the camera. If a parameter is not required, simply pass in PM_Fuji_NoParam!
func FujiSendOperationRequest(c *Client, code ptp.OperationCode, param uint32) (uint32, error) {
	c.incrementTransactionId()

	if err := c.SendPacketToCmdDataConn(&FujiOperationRequestPacket{
		DataPhaseInfo: uint16(DP_NoDataOrDataIn),
		OperationCode: code,
		TransactionID: c.TransactionId(),
		Parameter1:    param,
	}); err != nil {
		return 0, err
	}

	p := new(FujiOperationResponsePacket)
	if _, err := c.WaitForPacketFromCmdDataConn(p); err != nil {
		return 0, err
	}

	if !p.WasSuccessfull() {
		return 0, p.ReasonAsError()
	}

	return p.Parameter1, nil
}

// FujiGetDeviceInfo retrieves the current settings of a Fuji device. It is not at all a GetDeviceInfo call as specified
// in the PTP/IP specification, but it is more of a GetDevicePropDescList call that simply does not exist in the PTP/IP
// specification.
func FujiGetDeviceInfo(c *Client) (PacketIn, error) {
	c.log.Printf("Requesting %s device info...", c.ResponderFriendlyName())
	numProps, err := FujiSendOperationRequest(c, OC_Fuji_GetDeviceInfo, PM_Fuji_NoParam)
	if err != nil {
		return nil, err
	}

	c.log.Printf("Number of properties returned: %d", numProps)

	list := make([]*ptp.DevicePropDesc, numProps)

	for i := 0; i < int(numProps); i++ {
		var l uint32
		if err := binary.Read(c.commandDataConn, binary.LittleEndian, &l); err != nil {
			return nil, err
		}

		c.log.Printf("Property length: %d", l)

		// First read DevicePropertyCode, DataType and GetSet; i.e. read until we encounter the first interface field.
		dpd := new(ptp.DevicePropDesc)
		if err := ipInternal.UnmarshalLittleEndian(c.commandDataConn, dpd, int(l)-4, 0); err != nil && !strings.Contains(err.Error(), "invalid type *interface") {
			return nil, err
		}

		c.log.Printf("Size of property values in bytes: %d", dpd.SizeOfValueInBytes())

		// We now know the DataTypeCode so we know what to expect next.
		fdv := make([]byte, dpd.SizeOfValueInBytes())
		if err := binary.Read(c.commandDataConn, binary.LittleEndian, &fdv); err != nil {
			return nil, err
		}
		dpd.FactoryDefaultValue = fdv

		cv := make([]byte, dpd.SizeOfValueInBytes())
		if err := binary.Read(c.commandDataConn, binary.LittleEndian, &cv); err != nil {
			return nil, err
		}
		dpd.CurrentValue = cv

		// Read the type of form that will follow.
		if err := binary.Read(c.commandDataConn, binary.LittleEndian, &dpd.FormFlag); err != nil {
			return nil, err
		}

		switch dpd.FormFlag {
		case ptp.DPF_FormFlag_Range:
			form := new(ptp.RangeForm)
			c.log.Printf("Property is a range type, filling range form...")

			// Minimum possible value.
			min := make([]byte, dpd.SizeOfValueInBytes())
			if err := binary.Read(c.commandDataConn, binary.LittleEndian, &min); err != nil {
				return nil, err
			}
			form.MinimumValue = min

			// Maximum possible value.
			max := make([]byte, dpd.SizeOfValueInBytes())
			if err := binary.Read(c.commandDataConn, binary.LittleEndian, &max); err != nil {
				return nil, err
			}
			form.MaximumValue = max

			// Stepper value.
			step := make([]byte, dpd.SizeOfValueInBytes())
			if err := binary.Read(c.commandDataConn, binary.LittleEndian, &step); err != nil {
				return nil, err
			}
			form.StepSize = step

			dpd.Form = form
		case ptp.DPF_FormFlag_Enum:
			form := new(ptp.EnumerationForm)
			c.log.Printf("Property is an enum type, filling enum form...")

			// First read the number of values that will follow.
			var num uint16
			if err := binary.Read(c.commandDataConn, binary.LittleEndian, &num); err != nil {
				return nil, err
			}
			form.NumberOfValues = int(num)

			// Now fill the enumeration form with the actual values.
			for i := 0; i < form.NumberOfValues; i++ {
				v := make([]byte, dpd.SizeOfValueInBytes())
				if err := binary.Read(c.commandDataConn, binary.LittleEndian, v); err != nil {
					return nil, err
				}
				form.SupportedValues = append(form.SupportedValues, v)
			}
			dpd.Form = form
		}

		list = append(list, dpd)
	}

	// TODO: what to return??
	return nil, nil
}
