package ip

import (
	"encoding/binary"
	"errors"
	"github.com/google/uuid"
	ipInternal "github.com/malc0mn/ptp-ip/ip/internal"
	"github.com/malc0mn/ptp-ip/ptp"
)

type FujiBatteryLevel uint16
type FujiCommandDialMode uint16
type FujiDeviceError uint16
type FujiFilmSimulation uint16
type FujiFocusLock uint16
type FujiImageSize uint16
type FujiImageQuality uint16
type FujiMovieMode uint16
type FujiSelfTimer uint16

const (
	BAT_Fuji_3bCritical FujiBatteryLevel = 0x0001
	BAT_Fuji_3bOne      FujiBatteryLevel = 0x0002
	BAT_Fuji_3bTwo      FujiBatteryLevel = 0x0003
	BAT_Fuji_3bFull     FujiBatteryLevel = 0x0004
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

	IQ_Fuji_Fine         FujiImageQuality = 0x0002
	IQ_Fuji_Normal       FujiImageQuality = 0x0003
	IQ_Fuji_FineAndRAW   FujiImageQuality = 0x0004
	IQ_Fuji_NormalAndRAW FujiImageQuality = 0x0005
	IQ_Fuji_RAW          FujiImageQuality = 0x0006

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
	DPC_Fuji_Recmode         ptp.DevicePropCode = 0xD019
	DPC_Fuji_CommandDialMode ptp.DevicePropCode = 0xD028
	DPC_Fuji_ISO             ptp.DevicePropCode = 0xD02A
	DPC_Fuji_MovieISO        ptp.DevicePropCode = 0xD02B
	DPC_Fuji_FocusPoint      ptp.DevicePropCode = 0xD17C
	DPC_Fuji_FocusLock       ptp.DevicePropCode = 0xD209
	// DPC_Fuji_CurrentState is a property code that will return a list of properties with their current value.
	DPC_Fuji_CurrentState       ptp.DevicePropCode = 0xD212
	DPC_Fuji_DeviceError        ptp.DevicePropCode = 0xD21B
	DPC_Fuji_ImageSpaceSD       ptp.DevicePropCode = 0xD229
	DPC_Fuji_MovieRemainingTime ptp.DevicePropCode = 0xD22A
	DPC_Fuji_ShutterSpeed       ptp.DevicePropCode = 0xD240
	DPC_Fuji_ImageSize          ptp.DevicePropCode = 0xD241
	DPC_Fuji_BatteryLevel       ptp.DevicePropCode = 0xD242
	// DPC_Fuji_InitSequence indicates the initialisation sequence being used. It MUST be set by the Initiator during
	// the initialisation sequence and depending on it's value, will require a different init sequence to be used.
	// See PM_Fuji_InitSequence for further info.
	DPC_Fuji_InitSequence ptp.DevicePropCode = 0xDF01
	// DPC_Fuji_AppVersion indicates the minium application version the camera will accept. It MUST be set by the
	// Initiator during the initialisation sequence. As soon as this is done, the camera will acknowledge the client and
	// store the client's friendly name to allow future connections without the need for a confirmation.
	DPC_Fuji_AppVersion ptp.DevicePropCode = 0xDF24

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

	OC_Fuji_GetLastImage    ptp.OperationCode = 0x9022
	OC_Fuji_SetFocusPoint   ptp.OperationCode = 0x9026
	OC_Fuji_ResetFocusPoint ptp.OperationCode = 0x9027

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

	// RC_Fuji_DevicePropValue is the response code to a OC_GetDevicePropValue. The first parameter in the packet will
	// hold the property value.
	RC_Fuji_DevicePropValue = ptp.OperationResponseCode(ptp.OC_GetDevicePropValue)
	// RC_Fuji_DeviceInfo is the response code to OC_Fuji_GetDeviceInfo.
	RC_Fuji_DeviceInfo = ptp.OperationResponseCode(OC_Fuji_GetDeviceInfo)
)

func FujiDevicePropCodeAsString(code ptp.DevicePropCode) string {
	if msg := ptp.DevicePropCodeAsString(code); msg != "" {
		return msg
	}

	switch code {
	case DPC_Fuji_FilmSimulation:
		return "film simulation"
	case DPC_Fuji_ImageQuality:
		return "image quality"
	case DPC_Fuji_Recmode:
		return "rec mode"
	case DPC_Fuji_CommandDialMode:
		return "command dial mode"
	case DPC_Fuji_ISO:
		return "ISO"
	case DPC_Fuji_MovieISO:
		return "movie ISO"
	case DPC_Fuji_FocusPoint:
		return "focus point"
	case DPC_Fuji_FocusLock:
		return "focus lock"
	case DPC_Fuji_DeviceError:
		return "device error"
	case DPC_Fuji_ImageSpaceSD:
		return "image space SD"
	case DPC_Fuji_MovieRemainingTime:
		return "movie remaining time"
	case DPC_Fuji_ShutterSpeed:
		return "shutter speed"
	case DPC_Fuji_ImageSize:
		return "image size"
	case DPC_Fuji_BatteryLevel:
		return "battery level"
	case DPC_Fuji_InitSequence:
		return "init sequence"
	case DPC_Fuji_AppVersion:
		return "app version"
	default:
		return ""
	}
}

func FujiDevicePropValueAsString(code ptp.DevicePropCode, v uint16) string {
	if msg := ptp.DevicePropValueAsString(code, v); msg != "" {
		return msg
	}

	switch code {
	case ptp.DPC_BatteryLevel, DPC_Fuji_BatteryLevel:
		return FujiBatteryLevelAsString(FujiBatteryLevel(v))
	case DPC_Fuji_CommandDialMode:
		return FujiCommandDialModeAsString(FujiCommandDialMode(v))
	case DPC_Fuji_DeviceError:
		return FujiDeviceErrorAsString(FujiDeviceError(v))
	case DPC_Fuji_FilmSimulation:
		return FujiFilmSimulationAsString(FujiFilmSimulation(v))
	case ptp.DPC_FlashMode:
		return FujiFlashModeAsString(ptp.FlashMode(v))
	case DPC_Fuji_FocusLock:
		return FujiFocusLockAsString(FujiFocusLock(v))
	case ptp.DPC_FocusMode:
		return FujiFocusModeAsString(ptp.FocusMode(v))
	case DPC_Fuji_ImageSize:
		return FujiImageSizeAsString(FujiImageSize(v))
	case DPC_Fuji_ImageQuality:
		return FujiImageQualityAsString(FujiImageQuality(v))
	case ptp.DPC_WhiteBalance:
		return FujiWhiteBalanceAsString(ptp.WhiteBalance(v))
	case ptp.DPC_CaptureDelay:
		return FujiSelfTimerAsString(FujiSelfTimer(v))
	default:
		return ""
	}
}

func FujiBatteryLevelAsString(bat FujiBatteryLevel) string {
	switch bat {
	case BAT_Fuji_3bCritical:
		return "critical"
	case BAT_Fuji_3bOne:
		return "1/3"
	case BAT_Fuji_3bTwo:
		return "2/3"
	case BAT_Fuji_3bFull:
		return "3/3"
	case BAT_Fuji_5bCritical:
		return "critical"
	case BAT_Fuji_5bOne:
		return "1/5"
	case BAT_Fuji_5bTwo:
		return "2/5"
	case BAT_Fuji_5bThree:
		return "3/5"
	case BAT_Fuji_5bFour:
		return "4/5"
	case BAT_Fuji_5bFull:
		return "5/5"
	default:
		return ""
	}
}

func FujiCommandDialModeAsString(cmd FujiCommandDialMode) string {
	switch cmd {
	case CMD_Fuji_Both:
		return "both"
	case CMD_Fuji_Aperture:
		return "aperture"
	case CMD_Fuji_ShutterSpeed:
		return "shutter speed"
	case CMD_Fuji_None:
		return "none"
	default:
		return ""
	}
}

func FujiDeviceErrorAsString(de FujiDeviceError) string {
	switch de {
	case DE_Fuji_None:
		return "none"
	default:
		return ""
	}
}

func FujiFilmSimulationAsString(fs FujiFilmSimulation) string {
	switch fs {
	case FS_Fuji_Provia:
		return "Provia"
	case FS_Fuji_Velvia:
		return "Velvia"
	case FS_Fuji_Astia:
		return "Astia"
	case FS_Fuji_Monochrome:
		return "Monochrome"
	case FS_Fuji_Sepia:
		return "Sepia"
	case FS_Fuji_ProNegHigh:
		return "Pro. Neg. Hi"
	case FS_Fuji_ProNegStandard:
		return "Pro Neg. Std."
	case FS_Fuji_MonochromeYeFilter:
		return "Monochrome + Ye Filter"
	case FS_Fuji_MonochromeRFilter:
		return "Monochrome + R Filter"
	case FS_Fuji_MonochromeGFilter:
		return "Monochrome + G Filter"
	case FS_Fuji_ClassicChrome:
		return "Classic Chrome"
	case FS_Fuji_ACROS:
		return "ACROS"
	case FS_Fuji_ACROSYe:
		return "ACROS Ye"
	case FS_Fuji_ACROSR:
		return "ACROS R"
	case FS_Fuji_ACROSG:
		return "ACROS G"
	case FS_Fuji_ETERNA:
		return "ETERNA"
	default:
		return ""
	}
}

func FujiFlashModeAsString(mode ptp.FlashMode) string {
	switch mode {
	case FM_Fuji_On:
		return "on"
	case FM_Fuji_RedEye:
		return "red eye"
	case FM_Fuji_RedEyeOn:
		return "red eye on"
	case FM_Fuji_RedEyeSync:
		return "red eye sync"
	case FM_Fuji_RedEyeRear:
		return "red eye rear"
	case FM_Fuji_SlowSync:
		return "slow sync"
	case FM_Fuji_RearSync:
		return "rear sync"
	case FM_Fuji_Commander:
		return "commander"
	case FM_Fuji_Disabled:
		return "disabled"
	case FM_Fuji_Enabled:
		return "enabled"
	default:
		return ""
	}
}

func FujiFocusLockAsString(fl FujiFocusLock) string {
	switch fl {
	case FL_Fuji_On:
		return "on"
	case FL_Fuji_Off:
		return "off"
	default:
		return ""
	}
}

func FujiFocusModeAsString(fm ptp.FocusMode) string {
	switch fm {
	case FCM_Fuji_Single_Auto:
		return "single auto"
	case FCM_Fuji_Continuous_Auto:
		return "continuous auto"
	default:
		return ""
	}
}

func FujiImageSizeAsString(is FujiImageSize) string {
	switch is {
	case IS_Fuji_Small_3x2:
		return "S 3:2"
	case IS_Fuji_Small_16x9:
		return "S 16:9"
	case IS_Fuji_Small_1x1:
		return "S 1:1"
	case IS_Fuji_Medium_3x2:
		return "M 3:2"
	case IS_Fuji_Medium_16x9:
		return "M 16:9"
	case IS_Fuji_Medium_1x1:
		return "M 1:1"
	case IS_Fuji_Large_3x2:
		return "L 3:2"
	case IS_Fuji_Large_16x9:
		return "L 16:9"
	case IS_Fuji_Large_1x1:
		return "L 1:1"
	default:
		return ""
	}
}

func FujiImageQualityAsString(iq FujiImageQuality) string {
	switch iq {
	case IQ_Fuji_Fine:
		return "fine"
	case IQ_Fuji_Normal:
		return "normal"
	case IQ_Fuji_FineAndRAW:
		return "fine + RAW"
	case IQ_Fuji_NormalAndRAW:
		return "normal + RAW"
	case IQ_Fuji_RAW:
		return "RAW"
	default:
		return ""
	}
}

func FujiWhiteBalanceAsString(wb ptp.WhiteBalance) string {
	switch wb {
	case WB_Fuji_Fluorescent1:
		return "fluorescent 1"
	case WB_Fuji_Fluorescent2:
		return "fluorescent 2"
	case WB_Fuji_Fluorescent3:
		return "fluorescent 3"
	case WB_Fuji_Shade:
		return "shade"
	case WB_Fuji_Underwater:
		return "underwater"
	case WB_Fuji_Temperature:
		return "temprerature"
	case WB_Fuji_Custom:
		return "custom"
	default:
		return ""
	}
}

func FujiSelfTimerAsString(st FujiSelfTimer) string {
	switch st {
	case ST_Fuji_1Sec:
		return "1 second"
	case ST_Fuji_2Sec:
		return "2 seconds"
	case ST_Fuji_5Sec:
		return "5 seconds"
	case ST_Fuji_10Sec:
		return "10 seconds"
	case ST_Fuji_Off:
		return "off"
	default:
		return ""
	}
}

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
	return ipInternal.TotalSizeOfFixedFields(forp)
}

// TODO: make this better, obviously...
func (forp *FujiOperationResponsePacket) WasSuccessful() bool {
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
	if _, _, err := FujiSendOperationRequestAndGetResponse(c, ptp.OC_OpenSession, 0x00000001, 0); err != nil {
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
	if _, _, err := FujiSendOperationRequestAndGetResponse(c, ptp.OC_InitiateOpenCapture, PM_Fuji_NoParam, 0); err != nil {
		return err
	}

	return nil
}

// FujiSetDeviceProperty sets a device property to the given value.
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

	if !p.WasSuccessful() {
		return p.ReasonAsError()
	}

	return nil
}

func FujiGetEndOfDataPacket(c *Client, orp *FujiOperationResponsePacket) (*FujiOperationResponsePacket, error) {
	if orp.DataPhase != uint16(DP_DataOut) {
		return nil, nil
	}

	eodp := new(FujiOperationResponsePacket)
	if _, err := c.WaitForPacketFromCmdDataConn(eodp); err != nil {
		return nil, err
	}

	return eodp, nil
}

// FujiGetDevicePropertyValue gets the value for the given device property.
// TODO: add third parameter to indicate how many parameters from the response object are expected?
func FujiGetDevicePropertyValue(c *Client, dpc ptp.DevicePropCode) (uint32, error) {
	var val uint32
	var err error
	var rp *FujiOperationResponsePacket

	// First we get the actual value from the Responder.
	if val, rp, err = FujiSendOperationRequestAndGetResponse(c, ptp.OC_GetDevicePropValue, uint32(dpc), 4); err != nil {
		return 0, err
	}

	eodp, err := FujiGetEndOfDataPacket(c, rp)
	if err != nil {
		return 0, err
	}

	if eodp != nil && !eodp.WasSuccessful() {
		return 0, eodp.ReasonAsError()
	}

	return val, nil
}

// FujiSendOperationRequest sends an operation request to the camera. If a parameter is not required, simply pass in
// PM_Fuji_NoParam!
func FujiSendOperationRequest(c *Client, code ptp.OperationCode, param uint32) error {
	c.incrementTransactionId()

	err := c.SendPacketToCmdDataConn(&FujiOperationRequestPacket{
		DataPhaseInfo: uint16(DP_NoDataOrDataIn),
		OperationCode: code,
		TransactionID: c.TransactionId(),
		Parameter1:    param,
	})

	return err
}

// FujiSendOperationRequestAndGetResponse sends an operation request to the camera. If a parameter is not required,
// simply pass in PM_Fuji_NoParam!
// Sometimes, the response will have an additional variable sized parameter. Use pSize to indicate you are expecting one
// by passing the size in bytes of the expected data. Pass 0 when not expecting anything.
func FujiSendOperationRequestAndGetResponse(c *Client, code ptp.OperationCode, param uint32, pSize int) (uint32, *FujiOperationResponsePacket, error) {
	if err := FujiSendOperationRequest(c, code, param); err != nil {
		return 0, nil, err
	}

	p := new(FujiOperationResponsePacket)
	if _, err := c.WaitForPacketFromCmdDataConn(p); err != nil {
		return 0, nil, err
	}

	var parameter uint32
	if pSize > 0 {
		b := make([]byte, pSize)
		if err := binary.Read(c.commandDataConn, binary.LittleEndian, &b); err != nil {
			return 0, nil, err
		}
		if pSize < 4 {
			pad := make([]byte, 4-pSize)
			b = append(b, pad...)
		}
		parameter = binary.LittleEndian.Uint32(b)
	}

	if !p.WasSuccessful() {
		return 0, nil, p.ReasonAsError()
	}

	return parameter, p, nil
}

// FujiOperationRequestRaw wraps FujiSendOperationRequest and returns the raw camera response data.
func FujiOperationRequestRaw(c *Client, code ptp.OperationCode, params []uint32) ([][]byte, error) {
	field := uint32(PM_Fuji_NoParam)
	if len(params) != 0 {
		field = params[0]
	}
	if err := FujiSendOperationRequest(c, code, field); err != nil {
		return nil, err
	}

	var raw [][]byte
	var err error
	r, err := c.ReadRawFromCmdDataConn()
	if err == nil {
		raw = append(raw, r)
		if len(raw[0]) > 7 && binary.LittleEndian.Uint16(raw[0][4:6]) == uint16(DP_DataOut) {
			r, err := c.ReadRawFromCmdDataConn()
			if err == nil {
				raw = append(raw, r)
			}
		}
	}

	return raw, err
}

// FujiGetDeviceInfo retrieves the current settings of a Fuji device. It is not at all a GetDeviceInfo call as specified
// in the PTP/IP specification, but it is more of a GetDevicePropDescList call that simply does not exist in the PTP/IP
// specification.
func FujiGetDeviceInfo(c *Client) (interface{}, error) {
	c.Infof("Requesting %s device info...", c.ResponderFriendlyName())
	numProps, rp, err := FujiSendOperationRequestAndGetResponse(c, OC_Fuji_GetDeviceInfo, PM_Fuji_NoParam, 4)
	if err != nil {
		return nil, err
	}

	c.Debugf("Number of properties returned: %d", numProps)

	list := make([]*ptp.DevicePropDesc, numProps)

	for i := 0; i < int(numProps); i++ {
		var l uint32
		if err := binary.Read(c.commandDataConn, binary.LittleEndian, &l); err != nil {
			return nil, err
		}

		c.Debugf("Property length: %d", l)

		dpd := new(ptp.DevicePropDesc)
		if err := binary.Read(c.commandDataConn, binary.LittleEndian, &dpd.DevicePropertyCode); err != nil {
			return nil, err
		}
		if err := binary.Read(c.commandDataConn, binary.LittleEndian, &dpd.DataType); err != nil {
			return nil, err
		}
		if err := binary.Read(c.commandDataConn, binary.LittleEndian, &dpd.GetSet); err != nil {
			return nil, err
		}

		c.Debugf("Size of property values in bytes: %d", dpd.SizeOfValueInBytes())

		// We now know the DataTypeCode so we know what to expect next.
		dpd.FactoryDefaultValue = make([]byte, dpd.SizeOfValueInBytes())
		if err := binary.Read(c.commandDataConn, binary.LittleEndian, dpd.FactoryDefaultValue); err != nil {
			return nil, err
		}

		dpd.CurrentValue = make([]byte, dpd.SizeOfValueInBytes())
		if err := binary.Read(c.commandDataConn, binary.LittleEndian, dpd.CurrentValue); err != nil {
			return nil, err
		}

		// Read the type of form that will follow.
		if err := binary.Read(c.commandDataConn, binary.LittleEndian, &dpd.FormFlag); err != nil {
			return nil, err
		}

		switch dpd.FormFlag {
		case ptp.DPF_FormFlag_Range:
			form := new(ptp.RangeForm)
			c.Debug("Property is a range type, filling range form...")

			// Minimum possible value.
			form.MinimumValue = make([]byte, dpd.SizeOfValueInBytes())
			if err := binary.Read(c.commandDataConn, binary.LittleEndian, form.MinimumValue); err != nil {
				return nil, err
			}

			// Maximum possible value.
			form.MaximumValue = make([]byte, dpd.SizeOfValueInBytes())
			if err := binary.Read(c.commandDataConn, binary.LittleEndian, form.MaximumValue); err != nil {
				return nil, err
			}

			// Stepper value.
			form.StepSize = make([]byte, dpd.SizeOfValueInBytes())
			if err := binary.Read(c.commandDataConn, binary.LittleEndian, form.StepSize); err != nil {
				return nil, err
			}

			dpd.Form = form
		case ptp.DPF_FormFlag_Enum:
			form := new(ptp.EnumerationForm)
			c.Debug("Property is an enum type, filling enum form...")

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

		list[i] = dpd
	}

	eodp, err := FujiGetEndOfDataPacket(c, rp)
	if err != nil {
		return nil, err
	}

	if eodp != nil && !eodp.WasSuccessful() {
		return nil, eodp.ReasonAsError()
	}

	return list, nil
}

func FujiGetDeviceState(c *Client) (interface{}, error) {
	c.Infof("Requesting %s device state...", c.ResponderFriendlyName())
	numProps, rp, err := FujiSendOperationRequestAndGetResponse(c, ptp.OC_GetDevicePropValue, uint32(DPC_Fuji_CurrentState), 2)
	if err != nil {
		return nil, err
	}

	c.Debugf("Number of properties returned: %d", numProps)

	list := make([]*ptp.DevicePropDesc, numProps)

	for i := 0; i < int(numProps); i++ {
		dpd := new(ptp.DevicePropDesc)
		if err := binary.Read(c.commandDataConn, binary.LittleEndian, &dpd.DevicePropertyCode); err != nil {
			return nil, err
		}
		c.Debugf("Property code: %#x", dpd.DevicePropertyCode)

		dpd.DataType = ptp.DTC_UINT32
		dpd.CurrentValue = make([]byte, dpd.SizeOfValueInBytes())
		if err := binary.Read(c.commandDataConn, binary.LittleEndian, dpd.CurrentValue); err != nil {
			return nil, err
		}
		c.Debugf("Property value: %#x", dpd.CurrentValue)

		list[i] = dpd
	}

	eodp, err := FujiGetEndOfDataPacket(c, rp)
	if err != nil {
		return nil, err
	}

	if eodp != nil && !eodp.WasSuccessful() {
		return nil, eodp.ReasonAsError()
	}

	return list, nil
}
