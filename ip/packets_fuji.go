package ip

import (
	"github.com/google/uuid"
	ipInternal "github.com/malc0mn/ptp-ip/ip/internal"
)

const (
	// This error happens when the InitCommandRequestPacket contains a wrong identifier and/or protocol version.
	FR_FUJI_WRONG_IDENTIFIER FailReason = 0x0000201d
	// This fail reason can happen in these situations:
	//   - The FriendlyName stored in the camera does not match the FriendlyName being sent. Set the camera to 'change'
	//     so that it will accept a new FriendlyName
	//   - The camera side has timed out waiting for a connection and displays 'not found'. Set the camera to 'retry' to
	//     allow resending the InitCommandRequestPacket.
	FR_FUJI_CLIENT_STATE FailReason = 0x00002019

	PV_FUJI ProtocolVersion = 0xc0a87802
)

// The Fuji version of the PTP/IP InitCommandRequestPacket deviates from the standard. Looking at what is sent 'over the
// wire', we see this sequence in little endian format as the START of the packet, right after the header fields (being
// the Length [4 bytes] and PacketType [4 bytes] fields):
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
// Fuji seems to have swapped the FriendlyName and ProtocolVersion fields as it is sending 20 bytes BEFORE sending the
// FriendlyName field. From this the order is concluded to be:
//    - GUID [16]byte{0xf2, 0xe4, 0x53, 0x8f, 0xad, 0xa5, 0x48, 0x5d, 0x87, 0xb2, 0x7f, 0x0b, 0xd3, 0xd5, 0xde, 0xd0}
//    - ProtocolVersion [4]byte{0x02, 0x78, 0xa8, 0xc0}
type FujiInitCommandRequestPacket struct {
	GUID            uuid.UUID
	ProtocolVersion ProtocolVersion
	FriendlyName    string
}

func (ficrp *FujiInitCommandRequestPacket) PacketType() PacketType {
	return PKT_InitCommandRequest
}

func (ficrp *FujiInitCommandRequestPacket) Payload() []byte {
	return ipInternal.MarshalLittleEndian(ficrp)
}

func NewFujiInitCommandRequestPacket(friendlyName string) *FujiInitCommandRequestPacket {
	guid, _ := uuid.Parse("f2e4538f-ada5-485d-87b2-7f0bd3d5ded0")
	fa := &FujiInitCommandRequestPacket{
		GUID:            guid,
		ProtocolVersion: PV_FUJI,
		FriendlyName:    friendlyName,
	}

	return fa
}
