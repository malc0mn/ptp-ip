package ip

import (
	"encoding/binary"
	"fmt"
	"github.com/malc0mn/ptp-ip/ptp"
	"io"
	"log"
	"net"
)

func handleFujiMessages(conn net.Conn, lmp string) {
	// NO defer conn.Close() here since we need to mock a real Fuji responder and thus need to keep the connections open
	// when established and continuously listen for messages in a loop.
	for {
		l, raw, err := readMessageRaw(conn, lmp)
		if err == io.EOF {
			conn.Close()
			break
		}
		if raw == nil {
			continue
		}

		log.Printf("%s read %d raw bytes", lmp, l)

		var (
			msg  string
			resp PacketIn
			data []byte
		)
		eodp := false

		// This construction is thanks to the Fuji decision of not properly using packet types. Watch out for the caveat
		// here: we need to swap the order of the DataPhase and the OperationRequestCode because we are reading what are
		// actually two uint16 numbers as if they were a single uint32!
		switch binary.LittleEndian.Uint32(raw[0:4]) {
		case uint32(PKT_InitCommandRequest):
			msg, resp = genericInitCommandRequestResponse(lmp, ProtocolVersion(0))
		case constructPacketType(ptp.OC_GetDevicePropDesc):
			msg, resp, data = fujiGetDevicePropDescResponse(raw[4:8], raw[8:10])
			eodp = true
		case constructPacketType(ptp.OC_GetDevicePropValue):
			msg, resp, data = fujiGetDevicePropValueResponse(raw[4:8], raw[8:10])
			eodp = true
		case constructPacketType(ptp.OC_InitiateOpenCapture):
			msg, resp = fujiInitiateOpenCaptureResponse(raw[4:8])
		case constructPacketType(ptp.OC_OpenSession):
			msg, resp = fujiOpenSessionResponse(raw[4:8])
		case constructPacketTypeWithDataPhase(ptp.OC_SetDevicePropValue, DP_DataOut):
			// SetDevicePropValue involves two messages, only the second one needs a response from us!
			msg, resp = fujiSetDevicePropValue(raw[4:8])
		}

		if resp != nil {
			if msg != "" {
				log.Printf("%s responding to %s", lmp, msg)
			}
			sendMessage(conn, resp, data, lmp)
			if eodp {
				log.Printf("%s sending end of data packet", lmp)
				sendMessage(conn, fujiEndOfDataPacket(raw[4:8]), nil, lmp)
			}
		}
	}
}

func constructPacketType(code ptp.OperationCode) uint32 {
	return constructPacketTypeWithDataPhase(code, DP_NoDataOrDataIn)
}

func constructPacketTypeWithDataPhase(code ptp.OperationCode, dp DataPhase) uint32 {
	return uint32(code)<<16 | uint32(dp)
}

func fujiGetDevicePropDescResponse(tid []byte, prop []byte) (string, *FujiOperationResponsePacket, []byte) {
	var p []byte

	switch binary.LittleEndian.Uint16(prop) {
	case uint16(DPC_Fuji_FilmSimulation):
		p = []byte{0x01, 0xd0, 0x04, 0x00, 0x01, 0x01, 0x00, 0x01, 0x00, 0x02, 0x0b, 0x00, 0x01, 0x00, 0x02, 0x00, 0x03,
			0x00, 0x04, 0x00, 0x05, 0x00, 0x06, 0x00, 0x07, 0x00, 0x08, 0x00, 0x09, 0x00, 0x0a, 0x00, 0x0b, 0x00,
		}
	}

	return fmt.Sprintf("GetDevicePropDesc %#x", binary.LittleEndian.Uint16(prop)),
		fujiOperationResponsePacket(DP_DataOut, RC_Fuji_GetDevicePropDesc, tid),
		p
}

func fujiGetDevicePropValueResponse(tid []byte, prop []byte) (string, *FujiOperationResponsePacket, []byte) {
	var par uint32

	switch binary.LittleEndian.Uint16(prop) {
	case uint16(DPC_Fuji_AppVersion):
		par = PM_Fuji_AppVersion
	}

	p := make([]byte, 4)
	binary.LittleEndian.PutUint32(p, par)

	return fmt.Sprintf("GetDevicePropValue %#x", binary.LittleEndian.Uint16(prop)),
		fujiOperationResponsePacket(DP_DataOut, RC_Fuji_GetDevicePropValue, tid),
		p
}

func fujiInitiateOpenCaptureResponse(tid []byte) (string, *FujiOperationResponsePacket) {
	return "InitiateOpenCapture",
		fujiEndOfDataPacket(tid)
}

func fujiOpenSessionResponse(tid []byte) (string, *FujiOperationResponsePacket) {
	return "OpenSession",
		fujiEndOfDataPacket(tid)
}

func fujiSetDevicePropValue(tid []byte) (string, *FujiOperationResponsePacket) {
	return "SetDevicePropValue",
		fujiEndOfDataPacket(tid)
}

func fujiEndOfDataPacket(tid []byte) *FujiOperationResponsePacket {
	return fujiOperationResponsePacket(DP_Unknown, ptp.RC_OK, tid)
}

func fujiOperationResponsePacket(dp DataPhase, orc ptp.OperationResponseCode, tid []byte) *FujiOperationResponsePacket {
	return &FujiOperationResponsePacket{
		DataPhase:             uint16(dp),
		OperationResponseCode: orc,
		TransactionID:         ptp.TransactionID(binary.LittleEndian.Uint32(tid)),
	}
}
