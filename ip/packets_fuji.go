package ip

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/malc0mn/ptp-ip/ip/internal"
	"github.com/malc0mn/ptp-ip/ptp"
	"io"
	"time"
)

type FujiBatteryLevel uint16
type FujiCommandDialMode uint16
type FujiDeviceError uint16
type FujiExposureIndex uint32
type FujiFilmSimulation uint16
type FujiFocusLock uint16
type FujiImageSize uint16
type FujiImageQuality uint16
type FujiMovieMode uint16
type FujiSelfTimer uint16

const (
	BAT_Fuji_3bOne      FujiBatteryLevel = 0x0001
	BAT_Fuji_3bTwo      FujiBatteryLevel = 0x0002
	BAT_Fuji_3bFull     FujiBatteryLevel = 0x0003
	BAT_Fuji_5bCritical FujiBatteryLevel = 0x0006
	BAT_Fuji_5bOne      FujiBatteryLevel = 0x0007
	BAT_Fuji_5bTwo      FujiBatteryLevel = 0x0008
	BAT_Fuji_5bThree    FujiBatteryLevel = 0x0009
	BAT_Fuji_5bFour     FujiBatteryLevel = 0x000A
	BAT_Fuji_5bFull     FujiBatteryLevel = 0x000B

	CMD_Fuji_Both         FujiCommandDialMode = 0x0000
	CMD_Fuji_Aperture     FujiCommandDialMode = 0x0001
	CMD_Fuji_ShutterSpeed FujiCommandDialMode = 0x0002
	CMD_Fuji_None         FujiCommandDialMode = 0x0003

	// EDX_Fuji_Extended indicates if the ISO setting is an 'emulated' one. These are typically the extreme lows and
	// extreme highs such as 'L100' or 'H25600'. The 4 LSBs indicate the value.
	// Note: changing the ISO setting from an 'automatic value' to a 'manual value' will cause the fields output by
	// FujiGetDeviceState to change!
	EDX_Fuji_Extended uint16 = 0x4000
	// EDX_Fuji_MaxSensitivity indicates if the ISO setting is 'automatic' with a maximum sensitivity indicated by the 4
	// LSBs.
	// Note: changing the ISO setting from an 'automatic value' to a 'manual value' will cause the fields output by
	// FujiGetDeviceState to change!
	EDX_Fuji_MaxSensitivity uint16            = 0x8000
	EDX_Fuji_Auto           FujiExposureIndex = 0xFFFFFFFF

	FCM_Fuji_Single_Auto     ptp.FocusMode = 0x8001
	FCM_Fuji_Continuous_Auto ptp.FocusMode = 0x8002

	FM_Fuji_On         ptp.FlashMode = 0x8001
	FM_Fuji_RedEye     ptp.FlashMode = 0x8002
	FM_Fuji_RedEyeOn   ptp.FlashMode = 0x8003
	FM_Fuji_RedEyeSync ptp.FlashMode = 0x8004
	FM_Fuji_RedEyeRear ptp.FlashMode = 0x8005
	FM_Fuji_SlowSync   ptp.FlashMode = 0x8006
	FM_Fuji_RearSync   ptp.FlashMode = 0x8007
	FM_Fuji_Commander  ptp.FlashMode = 0x8008
	FM_Fuji_Disabled   ptp.FlashMode = 0x8009
	FM_Fuji_Enabled    ptp.FlashMode = 0x800A

	FL_Fuji_Off FujiFocusLock = 0x0000
	FL_Fuji_On  FujiFocusLock = 0x0001

	DE_Fuji_None FujiDeviceError = 0x0000

	FS_Fuji_Provia             FujiFilmSimulation = 0x0001
	FS_Fuji_Velvia             FujiFilmSimulation = 0x0002
	FS_Fuji_Astia              FujiFilmSimulation = 0x0003
	FS_Fuji_Monochrome         FujiFilmSimulation = 0x0004
	FS_Fuji_Sepia              FujiFilmSimulation = 0x0005
	FS_Fuji_ProNegHigh         FujiFilmSimulation = 0x0006
	FS_Fuji_ProNegStandard     FujiFilmSimulation = 0x0007
	FS_Fuji_MonochromeYeFilter FujiFilmSimulation = 0x0008
	FS_Fuji_MonochromeRFilter  FujiFilmSimulation = 0x0009
	FS_Fuji_MonochromeGFilter  FujiFilmSimulation = 0x000A
	FS_Fuji_ClassicChrome      FujiFilmSimulation = 0x000B
	FS_Fuji_ACROS              FujiFilmSimulation = 0x000C
	FS_Fuji_ACROSYe            FujiFilmSimulation = 0x000D
	FS_Fuji_ACROSR             FujiFilmSimulation = 0x000E
	FS_Fuji_ACROSG             FujiFilmSimulation = 0x000F
	FS_Fuji_ETERNA             FujiFilmSimulation = 0x0010

	IS_Fuji_Small_3x2   FujiImageSize = 0x0002
	IS_Fuji_Small_16x9  FujiImageSize = 0x0003
	IS_Fuji_Small_1x1   FujiImageSize = 0x0004
	IS_Fuji_Medium_3x2  FujiImageSize = 0x0006
	IS_Fuji_Medium_16x9 FujiImageSize = 0x0007
	IS_Fuji_Medium_1x1  FujiImageSize = 0x0008
	IS_Fuji_Large_3x2   FujiImageSize = 0x000A
	IS_Fuji_Large_16x9  FujiImageSize = 0x000B
	IS_Fuji_Large_1x1   FujiImageSize = 0x000C

	// IQ_Fuji_Fine indicates jpeg only shooting in 'fine' quality.
	IQ_Fuji_Fine FujiImageQuality = 0x0002
	// IQ_Fuji_Normal indicates jpeg only shooting in 'normal' quality.
	IQ_Fuji_Normal FujiImageQuality = 0x0003
	// IQ_Fuji_FineAndRAW indicates RAW + jpeg shooting in 'fine' quality. This is the highest possible quality setting
	// the camera will accept when in remote control mode. There is no 'RAW only' mode!
	IQ_Fuji_FineAndRAW FujiImageQuality = 0x0004
	// IQ_Fuji_NormalAndRAW indicates RAW + jpeg shooting in 'normal' quality.
	IQ_Fuji_NormalAndRAW FujiImageQuality = 0x0005

	MM_Fuji_None    FujiMovieMode = 0x0000
	MM_Fuji_Present FujiMovieMode = 0x0001

	ST_Fuji_Off   FujiSelfTimer = 0x0000
	ST_Fuji_1Sec  FujiSelfTimer = 0x0001
	ST_Fuji_2Sec  FujiSelfTimer = 0x0002
	ST_Fuji_5Sec  FujiSelfTimer = 0x0003
	ST_Fuji_10Sec FujiSelfTimer = 0x0004

	WB_Fuji_Fluorescent1 ptp.WhiteBalance = 0x8001
	WB_Fuji_Fluorescent2 ptp.WhiteBalance = 0x8002
	WB_Fuji_Fluorescent3 ptp.WhiteBalance = 0x8003
	WB_Fuji_Shade        ptp.WhiteBalance = 0x8006
	WB_Fuji_Underwater   ptp.WhiteBalance = 0x800A
	WB_Fuji_Temperature  ptp.WhiteBalance = 0x800B
	WB_Fuji_Custom       ptp.WhiteBalance = 0x800C

	DPC_Fuji_FilmSimulation  ptp.DevicePropCode = 0xD001
	DPC_Fuji_ImageQuality    ptp.DevicePropCode = 0xD018
	DPC_Fuji_RecMode         ptp.DevicePropCode = 0xD019
	DPC_Fuji_CommandDialMode ptp.DevicePropCode = 0xD028
	DPC_Fuji_ExposureIndex   ptp.DevicePropCode = 0xD02A
	DPC_Fuji_MovieISO        ptp.DevicePropCode = 0xD02B
	// DPC_Fuji_ImageSize is the Fuji equivalent of ptp.DPC_ImageSize. However ptp.DPC_ImageSize is directly supported
	// as well.
	DPC_Fuji_ImageSize         ptp.DevicePropCode = 0xD174
	DPC_Fuji_FocusMeteringMode ptp.DevicePropCode = 0xD17C
	DPC_Fuji_FocusLock         ptp.DevicePropCode = 0xD209
	// DPC_Fuji_CurrentState is a property code that will return a list of properties with their current value.
	DPC_Fuji_CurrentState ptp.DevicePropCode = 0xD212
	DPC_Fuji_DeviceError  ptp.DevicePropCode = 0xD21B
	// DPC_Fuji_CapturesRemaining indicates the amount of still image captures the internal storage can hold based on
	// the current capture quality and resolution settings.
	DPC_Fuji_CapturesRemaining  ptp.DevicePropCode = 0xD229
	DPC_Fuji_MovieRemainingTime ptp.DevicePropCode = 0xD22A
	DPC_Fuji_ShutterSpeed       ptp.DevicePropCode = 0xD240
	DPC_Fuji_ImageAspectRatio   ptp.DevicePropCode = 0xD241
	DPC_Fuji_BatteryLevel       ptp.DevicePropCode = 0xD242
	// DPC_Fuji_InitSequence indicates the initialisation sequence being used. It MUST be set by the Initiator during
	// the initialisation sequence and depending on it's value, will require a different init sequence to be used.
	// See PM_Fuji_InitSequence for further info.
	DPC_Fuji_InitSequence ptp.DevicePropCode = 0xDF01
	// DPC_Fuji_AppVersion indicates the minium application version the camera will accept. It MUST be set by the
	// Initiator during the initialisation sequence. As soon as this is done, the camera will acknowledge the client and
	// store the client's friendly name to allow future connections without the need for a confirmation.
	DPC_Fuji_AppVersion ptp.DevicePropCode = 0xDF24

	// EC_Fuji_PreviewAvailable is sent out as the second event during the ptp.OC_InitiateCapture operation indicating
	// the preview buffer is filled with a preview of the captured image. The client MUST empty this buffer by executing
	// the OC_Fuji_GetCapturePreview operation to make the camera send out a ptp.EC_CaptureComplete event which will
	// round up the ptp.OC_InitiateCapture operation allowing for a new capture to be taken.
	// Parameter2 of the event object will hold the size in bytes of the image preview data.
	EC_Fuji_PreviewAvailable ptp.EventCode = 0xC001
	// EC_Fuji_ObjectAdded is the first event sent during the ptp.OC_InitiateCapture operation informing the initiator
	// of a new object having been added to the device. Sadly none of the parameters hold the object handle allowing
	// the initiator to retrieve the full object.
	EC_Fuji_ObjectAdded ptp.EventCode = 0xC004

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

	// OC_Fuji_GetCapturePreview requests the contents of the preview buffer of the image captured by the
	// ptp.OC_InitiateCapture operation. This operation MUST be called immediately after receiving the
	// EC_Fuji_PreviewAvailable event to empty the preview buffer, thereby triggering the camera to sent the
	// ptp.EC_CaptureComplete event after which a new capture can be executed.
	OC_Fuji_GetCapturePreview ptp.OperationCode = 0x9022
	OC_Fuji_SetFocusPoint     ptp.OperationCode = 0x9026
	OC_Fuji_ResetFocusPoint   ptp.OperationCode = 0x9027

	// OC_Fuji_GetDeviceInfo returns a list of DevicePropDesc structs so it is not at all the same as OC_GetDeviceInfo.
	OC_Fuji_GetDeviceInfo ptp.OperationCode = 0x902B

	OC_Fuji_SetShutterSpeed         ptp.OperationCode = 0x902C
	OC_Fuji_SetAperture             ptp.OperationCode = 0x902D
	OC_Fuji_SetExposureCompensation ptp.OperationCode = 0x902E

	// PM_Fuji_NoParam is for convenience: to be used with FujiSendOperationRequestAndGetResponse() when no parameter is required for
	// the operation.
	PM_Fuji_NoParam = 0x00000000
	// PM_Fuji_InitSequence defines the init sequence to be used.
	// When this parameter is 'too low', the camera will complain about the application version being 'the previous
	// version' and requests to 'upgrade the app'.
	// After multiple experiments, this parameter will affect the initialisation sequence being used.
	// On the X-T1 (firmware v5.51):
	//   - 0x00000003 IS accepted but the init sequence used here does not seem to work and probably needs some
	//     unknown command(s) to 'finish it off' correctly.
	//   - 0x00000004 SEEMS not to be accepted, i.e. a client confirmation prompt is never displayed by the camera. So
	//     it MIGHT work, but could expect a different set of commands.
	//   - 0x00000005 hits the sweet spot and the init sequence we use completes nicely.
	// Could it maybe be that this is not so much 'init sequence' as "operation mode'? There is a mode for image
	// transfers over Wi-Fi as well, but this was not investigated deeper just yet...
	PM_Fuji_InitSequence = 0x00000005
	// PM_Fuji_AppVersion defines the minimal supported app version by the Responder.
	// When this parameter is 'too low', the camera will also complain about the application version being 'the previous
	// version' and requests to 'upgrade the app'. However, it does NOT affect the initialisation sequence at all.
	// The value here is that of the X-T1 on firmware version v5.51. We're not using it through this fixed value
	// anymore, but we now get it from the camera and confirm it by setting it to what the camera reports in the hope
	// that this will be future proof and we do not need to to adjust it ever again.
	PM_Fuji_AppVersion = 0x00020001

	// PV_Fuji is the Fuji Protocol Version required to construct a valid InitCommandRequestPacket.
	PV_Fuji ProtocolVersion = 0x8F53E4F2

	// RC_Fuji_GetDevicePropValue is the response code to a OC_GetDevicePropValue. The first parameter in the packet will
	// hold the property value.
	RC_Fuji_GetDevicePropValue = ptp.OperationResponseCode(ptp.OC_GetDevicePropValue)
	// RC_Fuji_GetDevicePropDesc is the response code to OC_GetDevicePropDesc
	RC_Fuji_GetDevicePropDesc = ptp.OperationResponseCode(ptp.OC_GetDevicePropDesc)
	// RC_Fuji_GetDeviceInfo is the response code to OC_Fuji_GetDeviceInfo.
	RC_Fuji_GetDeviceInfo = ptp.OperationResponseCode(OC_Fuji_GetDeviceInfo)
	// RC_Fuji_GetCapturePreview is the response code to OC_Fuji_GetCapturePreview
	RC_Fuji_GetCapturePreview = ptp.OperationResponseCode(OC_Fuji_GetCapturePreview)
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
	return internal.MarshalLittleEndian(ficrp)
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

func NewFujiInitCommandRequestPacketForClient(c *Client) InitCommandRequestPacket {
	return NewFujiInitCommandRequestPacket(c.InitiatorGUID(), c.InitiatorFriendlyName())
}

func NewFujiInitCommandRequestPacketWithVersion(guid uuid.UUID, friendlyName string, protocolVersion ProtocolVersion) InitCommandRequestPacket {
	icrp := NewFujiInitCommandRequestPacket(guid, friendlyName)
	icrp.SetProtocolVersion(protocolVersion)

	return icrp
}

// NewFujiInitEventRequestPacket returns nil because Fuji does not require the Event channel to be initialised. This
// will skip any further event channel initialisation.
func NewFujiInitEventRequestPacket(_ uint32) InitEventRequestPacket {
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
	return internal.MarshalLittleEndian(forp)
}

// FujiOperationResponsePacket deviates from the PTP/IP standard similarly to FujiOperationRequestPacket:
//   - the packet type should be PKT_OperationResponse, but there is NO packet type sent out in the packet header which
//     is, as one can imagine, extremely annoying when parsing the TCP/IP data coming in
//   - the DataPhase should be uint32 but Fuji uses uint16
//   - the parameters returned vary in size but the PTP/IP spec defines them as UINT32, extremely annoying
// To solve the varying parameter length, no parameters have been added to the struct but are handled separately.
type FujiOperationResponsePacket struct {
	DataPhase             uint16
	OperationResponseCode ptp.OperationResponseCode
	TransactionID         ptp.TransactionID
}

func (forp *FujiOperationResponsePacket) PacketType() PacketType {
	return PKT_Invalid
}

func (forp *FujiOperationResponsePacket) TotalFixedFieldSize() int {
	return internal.TotalSizeOfFixedFields(forp)
}

func operationCodeToOKResponseCode(oc ptp.OperationCode) ptp.OperationResponseCode {
	switch oc {
	case OC_Fuji_GetDeviceInfo:
		return RC_Fuji_GetDeviceInfo
	case ptp.OC_GetDevicePropDesc:
		return RC_Fuji_GetDevicePropDesc
	case ptp.OC_GetDevicePropValue:
		return RC_Fuji_GetDevicePropValue
	case ptp.OC_OpenSession:
		return ptp.RC_SessionAlreadyOpen
	default:
		return 0
	}
}

// WasSuccessful indicates if the operation request was successful by investigating the operation response code. By
// default it will check for ptp.RC_OK. If you expect another valid response code, you can pass it in.
func (forp *FujiOperationResponsePacket) WasSuccessful(rc ptp.OperationResponseCode) bool {
	return forp.OperationResponseCode == ptp.RC_OK || (rc != 0 && forp.OperationResponseCode == rc)
}

// ReasonAsError returns an error based on the operation response code.
func (forp *FujiOperationResponsePacket) ReasonAsError() error {
	return ptp.OperationResponseCodeAsError(forp.OperationResponseCode)
}

// FujiEventPacket is the Fuji version of the PTP/IP EventPacket which again deviates from the standard. 'Over the wire'
// we see these sequences, triggered by ptp.OC_InitiateCapute, in little endian format right after the the Length field
// (4 bytes) (the PacketType field is missing):
//   [24]byte{
//       0x04, 0x00, 0x04, 0xc0,
//       0x01, 0x00, 0x00, 0x00,
//       0x06, 0x00, 0x00, 0x00,
//       0x06, 0x00, 0x00, 0x00,
//       0x00, 0x00, 0x00, 0x00,
//       0x00, 0x00, 0x00, 0x00,
//   }
//   [24]byte{
//       0x04, 0x00, 0x01, 0xc0,
//       0x01, 0x00, 0x00, 0x00,
//       0x06, 0x00, 0x00, 0x00,
//       0x06, 0x00, 0x00, 0x00,
//       0x29, 0xf1, 0x00, 0x00,
//       0x00, 0x00, 0x00, 0x00,
//   }
// Referring to the PTP/IP standard which specifies the following:
//    - EventCode (16 bytes)
//    - TransactionID (32 bytes)
//    - Parameter1-3 (4 bytes each)
// There is an additional unknown field before the EventCode it always seems to be set to 0x004 for events so it was
// kept in line with the FujiOperationRequestPacket and hence dubbed DataPhase; although that makes no sense whatsoever.
// EventCode seems to adhere to the PTP standard concerning vendor extensions in that it starts with 0xC making the MSN
// 1100.
// There is another additional unknown field right after EventCode. It could be a SessionID? But experimenting with the
// OpenSession command and a different parameter did not seem to change the value of this field. So no clue what this
// field really is, seems to be set to 0x00000001; named it Amount for now.
// TransactionID is clear, but Parameter1 always seems to be set to the current TransactionID value as well!
type FujiEventPacket struct {
	DataPhase     uint16
	EventCode     ptp.EventCode
	Amount        uint32
	TransactionID ptp.TransactionID
	Parameter1    uint32
	Parameter2    uint32
	Parameter3    uint32
}

func (fep *FujiEventPacket) GetEventCode() ptp.EventCode {
	return fep.EventCode
}

func (fep *FujiEventPacket) PacketType() PacketType {
	return PKT_Invalid
}

func (fep *FujiEventPacket) TotalFixedFieldSize() int {
	return internal.TotalSizeOfFixedFields(fep)
}

func NewFujiEventPacket() EventPacket {
	return &FujiEventPacket{}
}

// FujiExtractTransactionId extracts the transaction ID from a full raw inbound packet. This packet must include the
// full header containing length and packet type.
func FujiExtractTransactionId(p []byte, ct connectionType) (ptp.TransactionID, error) {
	errFmt := "packet too small: got length %d"

	var data []byte
	switch ct {
	case cmdDataConnection:
		if len(p) < 8 {
			return 0, fmt.Errorf(errFmt, len(p))
		}

		data = p[8:12]
	case eventConnection:
		if len(p) < 12 {
			return 0, fmt.Errorf(errFmt, len(p))
		}

		data = p[12:16]
	}

	return ptp.TransactionID(binary.LittleEndian.Uint32(data)), nil
}

// FujiInitCommandDataConn initialises the Fuji command/data connection. It expects an open TCP connection to the
// command/data port to be present.
// The PTP/IP protocol specifies how to set up the command/data connection which should immediately be followed by
// setting up the event connection. However Fuji wants additional communications before it is satisfied that the
// command/data connection is properly setup. This additional initialisation is performed here.
// The sequence is as follows:
//   1. Open a session.
//   2. Set device property DPC_Fuji_InitSequence to the correct number of the init sequence being used by the
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
//      Initiator. This also opens up the event connection port 55741 used by Fuji so we can connect to it and complete
//      the init sequence there.
func FujiInitCommandDataConn(c *Client) error {
	// The first part of the sequence is according to the PTP/IP standard, save for the different packet format.
	if err := GenericInitCommandDataConn(c); err != nil {
		return err
	}

	c.Info("Opening a session...")
	if err := FujiSendOperationRequestIgnoreResponse(c, ptp.OC_OpenSession, 0x00000001, 0); err != nil {
		return err
	}

	c.Info("Setting correct init sequence number...")
	c.Infof("Should you be prompted, please accept the new connection request on the %s.", c.ResponderFriendlyName())
	if err := FujiSetDeviceProperty(c, DPC_Fuji_InitSequence, PM_Fuji_InitSequence); err != nil {
		return err
	}

	c.Info("Getting current minimum application version...")
	val, err := FujiGetDevicePropertyValue(c, DPC_Fuji_AppVersion)
	if err != nil {
		return err
	}
	c.Infof("Acknowledging current minimal application version as communicated by the %s: %#x", c.ResponderFriendlyName(), val)
	if err := FujiSetDeviceProperty(c, DPC_Fuji_AppVersion, val); err != nil {
		return err
	}

	c.Info("Initiating open capture...")
	if err := FujiSendOperationRequestIgnoreResponse(c, ptp.OC_InitiateOpenCapture, PM_Fuji_NoParam, 0); err != nil {
		return err
	}

	return nil
}

// FujiProcessStreamData reads raw image data from the incoming stream and sends them to the streamer channel.
func FujiProcessStreamData(c *Client) error {
	go func() {
		c.Info("[fujiStreamListener] subscribing stream listener to streamer connection...")
		for {
			select {
			case <-c.closeStreamChan:
				c.Info("[fujiStreamListener] stopping stream listener.")
				close(c.StreamChan)
				c.StreamChan = nil
				return
			default:
				if data, err := c.ReadRawFromStreamConn(); err == nil {
					// As always, packet length first.
					l := binary.LittleEndian.Uint32(data[:4])
					c.Debugf("[fujiStreamListener] Packet length %d", l)

					// Four bytes always zero followed by what is clearly a counter which resets on 0xff, so one byte
					// only.
					count := binary.LittleEndian.Uint16(append(data[8:9], 0))
					c.Debugf("[fujiStreamListener] Image number %d", count)

					// Unknown what the next 9 bytes are, but they always END in two bytes with unknown significance
					// (seen 0xff, 0xff as well as 0x5e, 0x49 and 0x4b, 0xbf) after which the image data begins, filling
					// the rest of the packet.
					c.StreamChan <- data[18:]
				}
			}
		}
	}()

	return nil
}

// FujiSetDeviceProperty sets a device property to the given value.
func FujiSetDeviceProperty(c *Client, code ptp.DevicePropCode, val uint32) error {
	tid := c.incrementTransactionId()

	resCh := make(chan []byte, 2)
	if err := c.subscribe(tid, resCh); err != nil {
		return err
	}
	defer c.unsubscribe(tid)

	if err := c.SendPacketToCmdDataConn(&FujiOperationRequestPacket{
		DataPhaseInfo: uint16(DP_NoDataOrDataIn),
		OperationCode: ptp.OC_SetDevicePropValue,
		TransactionID: tid,
		Parameter1:    uint32(code),
	}); err != nil {
		return err
	}

	if err := c.SendPacketToCmdDataConn(&FujiOperationRequestPacket{
		DataPhaseInfo: uint16(DP_DataOut),
		OperationCode: ptp.OC_SetDevicePropValue,
		TransactionID: tid,
		Parameter1:    val,
	}); err != nil {
		return err
	}

	p := new(FujiOperationResponsePacket)
	if _, _, err := c.WaitForPacketFromCommandDataSubscriber(resCh, p); err != nil {
		return err
	}

	if !p.WasSuccessful(0) {
		return p.ReasonAsError()
	}

	return nil
}

// FujiGetDevicePropertyValue gets the value for the given device property.
// TODO: add third parameter to indicate how many parameters from the response object are expected?
func FujiGetDevicePropertyValue(c *Client, dpc ptp.DevicePropCode) (uint32, error) {
	var val uint32
	var err error

	// First we get the actual value from the Responder.
	if val, _, err = FujiSendOperationRequestAndGetResponse(c, ptp.OC_GetDevicePropValue, uint32(dpc), 4); err != nil {
		return 0, err
	}

	return val, nil
}

// FujiSendOperationRequest sends an operation request to the camera and returns a channel that will receive the
// response messages as a raw byte array.
// If a parameter is not required, simply pass in PM_Fuji_NoParam!
func FujiSendOperationRequest(c *Client, code ptp.OperationCode, param uint32) (chan []byte, error) {
	resCh := make(chan []byte, 2)

	_, err := FujiSendOperationRequestWithChan(c, code, param, resCh)

	return resCh, err
}

// FujiSendOperationRequestWithChan sends an operation request to the camera and returns the current transaction ID
// the given response channel has been subscribed to.
// If a parameter is not required, simply pass in PM_Fuji_NoParam!
func FujiSendOperationRequestWithChan(c *Client, code ptp.OperationCode, param uint32, resCh chan []byte) (ptp.TransactionID, error) {
	tid := c.incrementTransactionId()

	if err := c.subscribe(tid, resCh); err != nil {
		return 0, err
	}

	return tid, c.SendPacketToCmdDataConn(&FujiOperationRequestPacket{
		DataPhaseInfo: uint16(DP_NoDataOrDataIn),
		OperationCode: code,
		TransactionID: tid,
		Parameter1:    param,
	})
}

// FujiSendOperationRequestIgnoreResponse sends an operation request to the camera. If a parameter is not required,
// simply pass in PM_Fuji_NoParam!
// Use this wrapper function if you do not care about the actual response value but just want to know if it was
// successful.
func FujiSendOperationRequestIgnoreResponse(c *Client, code ptp.OperationCode, param uint32, pSize int) error {
	_, _, err := FujiSendOperationRequestAndGetResponse(c, code, param, pSize)

	return err
}

// FujiSendOperationRequestAndGetResponse sends an operation request to the camera. If a parameter is not required,
// simply pass in PM_Fuji_NoParam!
// Sometimes, the response will have an additional variable sized parameter. Use pSize to indicate you are expecting one
// by passing the size in bytes of the expected data. Pass 0 when not expecting anything.
// The byte array being returned may contain excess dat that could not be unmarshalled. This will often be the case so
// check this data to see if it is not nil and handle it accordingly.
func FujiSendOperationRequestAndGetResponse(c *Client, code ptp.OperationCode, param uint32, pSize int) (uint32, []byte, error) {
	resCh, err := FujiSendOperationRequest(c, code, param)
	if err != nil {
		return 0, nil, err
	}

	p := new(FujiOperationResponsePacket)
	_, xs, err := c.WaitForPacketFromCommandDataSubscriber(resCh, p)
	if err != nil {
		return 0, nil, err
	}

	// Make sure we also grab the end of data packet should it be there...
	if p.DataPhase == uint16(DP_DataOut) {
		eodp := new(FujiOperationResponsePacket)
		if _, _, err := c.WaitForPacketFromCommandDataSubscriber(resCh, eodp); err != nil {
			return 0, nil, err
		}

		if eodp != nil && !eodp.WasSuccessful(0) {
			return 0, nil, eodp.ReasonAsError()
		}
	}

	if xs == nil && pSize > 0 {
		return 0, nil, errors.New("expected additional value but none was returned")
	}

	r := bytes.NewReader(xs)
	var parameter uint32
	if pSize > 0 {
		// This here is done to remove complexity in the variable parameter size being returned depending on the
		// requested property. The PTP/IP standard states the parameters should always be uint32, but Fuji...
		// The size adjustment here would mean that when expecting 4 bytes but only 2 are returned, you will still get
		// that data. Be careful though...
		if len(xs) < pSize {
			pSize = len(xs)
		}

		v := make([]byte, pSize)
		if err := binary.Read(r, binary.LittleEndian, &v); err != nil {
			return 0, nil, err
		}
		if pSize < 4 {
			pad := make([]byte, 4-pSize)
			v = append(v, pad...)
		}
		parameter = binary.LittleEndian.Uint32(v)

		// Remove bytes read from unmarshalled data.
		xs = xs[pSize:]
	}

	if !p.WasSuccessful(operationCodeToOKResponseCode(code)) {
		return 0, nil, p.ReasonAsError()
	}

	return parameter, xs, nil
}

// FujiSendOperationRequestAndGetRawResponse wraps FujiSendOperationRequest and returns the raw camera response data.
func FujiSendOperationRequestAndGetRawResponse(c *Client, code ptp.OperationCode, params []uint32) ([][]byte, error) {
	var err error

	field := uint32(PM_Fuji_NoParam)
	if len(params) != 0 {
		field = params[0]
	}

	resCh := make(chan []byte, 2)
	tid, err := FujiSendOperationRequestWithChan(c, code, field, resCh)
	if err != nil {
		return nil, err
	}
	defer c.unsubscribe(tid)

	var raw [][]byte
	for {
		r, err := c.WaitForRawPacketFromCommandDataSubscriber(resCh)
		if err == nil {
			raw = append(raw, r)
			// Keep reading as long as the Responder tells us there is more data.
			if len(r) > 7 && binary.LittleEndian.Uint16(r[4:6]) == uint16(DP_DataOut) {
				continue
			}
		}
		break
	}

	// TODO: find out why no error is returned when receiving a timeout!

	// TODO: check if there is data on the event connection and read that as well!

	return raw, err
}

// FujiGetDevicePropDesc retrieves the description for the given device property code. Beware that this method can
// return no error and at the same time return nil for *ptp.DevicePropDesc! This means that the requested device
// property cannot be described: the camera gave a response but returned no property data.
// With the Fuji implementation one cannot be sure if the property does not exist or cannot be described as there is no
// clear error being returned.
func FujiGetDevicePropertyDesc(c *Client, code ptp.DevicePropCode) (*ptp.DevicePropDesc, error) {
	c.Infof("Requesting %s device property description for %#x...", c.ResponderFriendlyName(), code)
	_, xs, err := FujiSendOperationRequestAndGetResponse(c, ptp.OC_GetDevicePropDesc, uint32(code), 0)
	if err != nil {
		return nil, err
	}

	r := bytes.NewReader(xs)

	dpd, err := fujiReadDevicePropDesc(c, r)
	// When requesting the description of a non-existing device property, the camera does not return an error code, it
	// just does not return any data. Another annoying complexity we need to handle here...
	if err != nil && err != io.EOF {
		return nil, err
	}

	return dpd, nil
}

// FujiGetDeviceInfo retrieves the current settings of a Fuji device. It is not at all a GetDeviceInfo call as specified
// in the PTP/IP specification, but it is more of a GetDevicePropDescList call that simply does not exist in the PTP/IP
// specification.
func FujiGetDeviceInfo(c *Client) (interface{}, error) {
	c.Infof("Requesting %s device info...", c.ResponderFriendlyName())
	numProps, xs, err := FujiSendOperationRequestAndGetResponse(c, OC_Fuji_GetDeviceInfo, PM_Fuji_NoParam, 4)
	if err != nil {
		return nil, err
	}

	c.Debugf("Number of properties returned: %d", numProps)

	r := bytes.NewReader(xs)
	list := make([]*ptp.DevicePropDesc, numProps)

	for i := 0; i < int(numProps); i++ {
		var l uint32
		if err := binary.Read(r, binary.LittleEndian, &l); err != nil {
			return nil, err
		}

		c.Debugf("Property length: %d", l)

		dpd, err := fujiReadDevicePropDesc(c, r)
		if err != nil {
			return nil, err
		}

		list[i] = dpd
	}

	return list, nil
}

func fujiReadDevicePropDesc(c *Client, r io.Reader) (*ptp.DevicePropDesc, error) {
	dpd := new(ptp.DevicePropDesc)
	if err := binary.Read(r, binary.LittleEndian, &dpd.DevicePropertyCode); err != nil {
		return nil, err
	}
	if err := binary.Read(r, binary.LittleEndian, &dpd.DataType); err != nil {
		return nil, err
	}
	if err := binary.Read(r, binary.LittleEndian, &dpd.GetSet); err != nil {
		return nil, err
	}

	c.Debugf("Size of property values in bytes: %d", dpd.SizeOfValueInBytes())

	// We now know the DataTypeCode so we know what to expect next.
	dpd.FactoryDefaultValue = make([]byte, dpd.SizeOfValueInBytes())
	if err := binary.Read(r, binary.LittleEndian, dpd.FactoryDefaultValue); err != nil {
		return nil, err
	}

	dpd.CurrentValue = make([]byte, dpd.SizeOfValueInBytes())
	if err := binary.Read(r, binary.LittleEndian, dpd.CurrentValue); err != nil {
		return nil, err
	}

	// Read the type of form that will follow.
	if err := binary.Read(r, binary.LittleEndian, &dpd.FormFlag); err != nil {
		return nil, err
	}

	switch dpd.FormFlag {
	case ptp.DPF_FormFlag_Range:
		form := new(ptp.RangeForm)
		c.Debug("Property is a range type, filling range form...")

		form.SetDevicePropDesc(dpd)

		// Minimum possible value.
		form.MinimumValue = make([]byte, dpd.SizeOfValueInBytes())
		if err := binary.Read(r, binary.LittleEndian, form.MinimumValue); err != nil {
			return nil, err
		}

		// Maximum possible value.
		form.MaximumValue = make([]byte, dpd.SizeOfValueInBytes())
		if err := binary.Read(r, binary.LittleEndian, form.MaximumValue); err != nil {
			return nil, err
		}

		// Stepper value.
		form.StepSize = make([]byte, dpd.SizeOfValueInBytes())
		if err := binary.Read(r, binary.LittleEndian, form.StepSize); err != nil {
			return nil, err
		}

		dpd.Form = form
	case ptp.DPF_FormFlag_Enum:
		form := new(ptp.EnumerationForm)
		c.Debug("Property is an enum type, filling enum form...")

		form.SetDevicePropDesc(dpd)

		// First read the number of values that will follow.
		var num uint16
		if err := binary.Read(r, binary.LittleEndian, &num); err != nil {
			return nil, err
		}
		form.NumberOfValues = int(num)

		// Now fill the enumeration form with the actual values.
		for i := 0; i < form.NumberOfValues; i++ {
			v := make([]byte, dpd.SizeOfValueInBytes())
			if err := binary.Read(r, binary.LittleEndian, v); err != nil {
				return nil, err
			}
			form.SupportedValues = append(form.SupportedValues, v)
		}
		dpd.Form = form
	}

	return dpd, nil
}

// FujiGetDeviceState returns a list of properties with their current values. The values being returned will depend on
// the exposure program mode of the camera: it will change if the camera is in aperture priority, shutter priority,
// manual or auto.
func FujiGetDeviceState(c *Client) (interface{}, error) {
	c.Infof("Requesting %s device state...", c.ResponderFriendlyName())
	numProps, xs, err := FujiSendOperationRequestAndGetResponse(c, ptp.OC_GetDevicePropValue, uint32(DPC_Fuji_CurrentState), 2)
	if err != nil {
		return nil, err
	}

	c.Debugf("Number of properties returned: %d", numProps)

	r := bytes.NewReader(xs)
	list := make([]*ptp.DevicePropDesc, numProps)

	for i := 0; i < int(numProps); i++ {
		dpd := new(ptp.DevicePropDesc)
		if err := binary.Read(r, binary.LittleEndian, &dpd.DevicePropertyCode); err != nil {
			return nil, err
		}
		c.Debugf("Property code: %#x", dpd.DevicePropertyCode)

		dpd.DataType = ptp.DTC_UINT32
		dpd.CurrentValue = make([]byte, dpd.SizeOfValueInBytes())
		if err := binary.Read(r, binary.LittleEndian, dpd.CurrentValue); err != nil {
			return nil, err
		}
		c.Debugf("Property value: %#x", dpd.CurrentValue)

		list[i] = dpd
	}

	return list, nil
}

// FujiInitiateCapture releases the shutter and returns a byte array containing the raw image data representing a preview
// of the image taken.
// The sequence is a bit odd: it partly follows the PTP/IP spec but expects the client to request the preview buffer
// from the camera in order for the ptp.EC_CaptureComplete to be sent out.
// Failing to do this, will not allow the client to release the shutter again. The operation request will be accepted
// but no further actions will be taken by the camera.
func FujiInitiateCapture(c *Client) ([]byte, error) {
	c.Infof("Releasing %s shutter...", c.ResponderFriendlyName())
	if err := FujiSendOperationRequestIgnoreResponse(c, ptp.OC_InitiateCapture, PM_Fuji_NoParam, 0); err != nil {
		return nil, err
	}

	var pvSize int
	invalidEvent := "invalid event received, expected '%#x' got '%#x'"
	for _, ec := range []ptp.EventCode{EC_Fuji_ObjectAdded, EC_Fuji_PreviewAvailable} {
		select {
		case msg := <-c.eventChan:
			if msg.GetEventCode() != ec {
				return nil, fmt.Errorf(invalidEvent, ec, msg.GetEventCode())
			}
			var txt string
			var extra string
			switch ec {
			case EC_Fuji_ObjectAdded:
				txt = "object added"
			case EC_Fuji_PreviewAvailable:
				txt = "preview available"
				pvSize = int(msg.(*FujiEventPacket).Parameter2)
				extra = fmt.Sprintf(": preview size is %d bytes", pvSize)
			}
			c.Debugf("Received %s event (%#x)%s.", txt, msg.GetEventCode(), extra)
		case <-time.After(DefaultReadTimeout):
			return nil, WaitForEventError
		}
	}

	raw, err := FujiSendOperationRequestAndGetRawResponse(c, OC_Fuji_GetCapturePreview, nil)
	if err != nil {
		return nil, err
	}

	select {
	case msg := <-c.eventChan:
		if msg.GetEventCode() != ptp.EC_CaptureComplete {
			return nil, fmt.Errorf("invalid event received, expected '%#x' got '%#x'", ptp.EC_CaptureComplete, msg.GetEventCode())
		}
		c.Debugf("Received capture complete event (%#x).", msg.GetEventCode())
	case <-time.After(DefaultReadTimeout):
		return nil, WaitForEventError
	}

	var img []byte
	for _, pkt := range raw {
		code := binary.LittleEndian.Uint16(pkt[6:8])
		switch {
		case code == uint16(ptp.RC_OK):
			break
		case code != uint16(OC_Fuji_GetCapturePreview):
			return nil, errors.New("failed reading image data")
		}
		img = append(img, pkt[12:]...)
	}

	if len(img) != pvSize {
		c.Warnf("Preview size mismatch: expected %d, got %d. Returning possibly malformed data nonetheless.", pvSize, len(img))
	}

	return img, nil
}
