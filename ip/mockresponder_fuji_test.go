package ip

import (
	"encoding/binary"
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
			msg string
			res PacketIn
			par []byte
		)
		eodp := false

		// This construction is thanks to the Fuji decision of not properly using packet types. Watch out for the
		// caveat here: we need to swap the order of the data phase and the operation request code because we're
		// we're reading two uint16 numbers as a uint32!!!
		switch binary.LittleEndian.Uint32(raw[0:4]) {
		case uint32(PKT_InitCommandRequest):
			msg, res = genericInitCommandRequestResponse(lmp, ProtocolVersion(0))
		case uint32(ptp.OC_GetDevicePropDesc)<<16 | uint32(DP_NoDataOrDataIn):
			msg, res = fujiGetDevicePropDescResponse(raw[4:8])
			eodp = true
		case uint32(ptp.OC_GetDevicePropValue)<<16 | uint32(DP_NoDataOrDataIn):
			msg, res, par = fujiGetDevicePropValueResponse(raw[4:8], raw[8:10])
			eodp = true
		case uint32(ptp.OC_InitiateOpenCapture)<<16 | uint32(DP_NoDataOrDataIn):
			msg, res = fujiInitiateOpenCaptureResponse(raw[4:8])
		case uint32(ptp.OC_OpenSession)<<16 | uint32(DP_NoDataOrDataIn):
			msg, res = fujiOpenSessionResponse(raw[4:8])
		case uint32(ptp.OC_SetDevicePropValue)<<16 | uint32(DP_DataOut):
			// SetDevicePropValue involves two messages, only the second one needs a response from us!
			msg, res = fujiSetDevicePropValue(raw[4:8])
		}

		if res != nil {
			if msg != "" {
				log.Printf("%s responding to %s", lmp, msg)
			}
			sendMessage(conn, res, lmp)
			if par != nil {
				log.Printf("%s sending parameter %#v", lmp, par)
				conn.Write(par)
			}
			if eodp {
				log.Printf("%s sending end of data packet", lmp)
				sendMessage(conn, fujiEndOfDataPacket(raw[4:8]), lmp)
			}
		}
	}
}

func fujiGetDevicePropDescResponse(tid []byte) (string, *FujiOperationResponsePacket) {
	return "GetDevicePropDesc",
		fujiOperationResponsePacket(DP_DataOut, RC_Fuji_GetDevicePropDesc, tid)
}

func fujiGetDevicePropValueResponse(tid []byte, prop []byte) (string, *FujiOperationResponsePacket, []byte) {
	var par uint32

	switch binary.LittleEndian.Uint16(prop) {
	case uint16(DPC_Fuji_AppVersion):
		par = PM_Fuji_AppVersion
	}

	p := make([]byte, 4)
	binary.LittleEndian.PutUint32(p, par)

	return "GetDevicePropValue",
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
