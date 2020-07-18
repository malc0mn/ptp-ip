package ip

import (
	"github.com/google/uuid"
	"io"
	"log"
	"net"
)

func handleGenericMessages(conn net.Conn, lmp string) {
	// NO defer conn.Close() here since we need to mock a real responder and thus need to keep the connections open when
	// established and continuously listen for messages in a loop.
	for {
		h, pkt, err := readMessage(conn, lmp)
		if err == io.EOF {
			conn.Close()
			break
		}
		if pkt == nil {
			continue
		}

		var msg string
		var res PacketIn
		switch h.PacketType {
		case PKT_InitCommandRequest:
			msg, res = genericInitCommandRequestResponse(lmp, PV_VersionOnePointZero)
		case PKT_InitEventRequest:
			msg, res = genericInitEventRequestResponse()
		case PKT_OperationRequest:
			msg, res = genericOperationRequestResponse()
		default:
			log.Printf("%s unknown packet type %#x", lmp, h.PacketType)
			continue
		}

		if res != nil {
			if msg != "" {
				log.Printf("%s responding to %s", lmp, msg)
			}
			sendMessage(conn, res, lmp)
		}
	}
}

func genericInitCommandRequestResponse(friendlyName string, pv ProtocolVersion) (string, PacketIn) {
	uuid, _ := uuid.Parse(MockResponderGUID)
	return "InitCommandRequest",
		&InitCommandAckPacket{
			ConnectionNumber:         1,
			ResponderGUID:            uuid,
			ResponderFriendlyName:    friendlyName,
			ResponderProtocolVersion: uint32(pv),
		}
}

func genericInitEventRequestResponse() (string, PacketIn) {
	return "InitEventRequest", &InitEventAckPacket{}
}

func genericOperationRequestResponse() (string, PacketIn) {
	return "OperationRequest", &OperationResponsePacket{}
}