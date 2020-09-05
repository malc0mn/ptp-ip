package ip

import (
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/malc0mn/ptp-ip/ip/internal"
	"github.com/malc0mn/ptp-ip/ptp"
)

// TODO: This solution is not OK, vendors can differ massively so it seems. Should this become an interface that all
//  vendors need to implement...? It would turn out to be a huge interface, so there will no doubt be a better solution?
//  Embedding ip.Client in a struct ip.FujiClient wasn't that good either. Take the Dial() method for example, this
//  calls the initCommandDataConn() and initEventConn() methods but when using embedding the methods on ip.Client get
//  called and not the ones on ip.FujiClient so you would also have to "override" the Dial() as well.
type VendorExtensions struct {
	cmdDataInit            func(*Client) error
	eventInit              func(*Client) error
	processStreamData      func(*Client) error
	newCmdDataInitPacket   func(uuid.UUID, string) InitCommandRequestPacket
	newEventInitPacket     func(uint32) InitEventRequestPacket
	newEventPacket         func() EventPacket
	extractTransactionId   func([]byte, connectionType) (ptp.TransactionID, error)
	getDeviceInfo          func(*Client) (interface{}, error)
	getDeviceState         func(*Client) (interface{}, error)
	getDevicePropertyDesc  func(*Client, ptp.DevicePropCode) (*ptp.DevicePropDesc, error)
	getDevicePropertyValue func(*Client, ptp.DevicePropCode) (uint32, error)
	setDeviceProperty      func(*Client, ptp.DevicePropCode, uint32) error
	operationRequestRaw    func(*Client, ptp.OperationCode, []uint32) ([][]byte, error)
	initiateCapture        func(*Client) ([]byte, error)
}

func (c *Client) loadVendorExtensions() {
	c.vendorExtensions = &VendorExtensions{
		cmdDataInit:            GenericInitCommandDataConn,
		eventInit:              GenericInitEventConn,
		processStreamData:      GenericProcessStreamData,
		newCmdDataInitPacket:   NewInitCommandRequestPacket,
		newEventInitPacket:     NewInitEventRequestPacket,
		newEventPacket:         NewEventPacket,
		extractTransactionId:   GenericExtractTransactionId,
		getDeviceInfo:          GenericGetDeviceInfo,
		getDeviceState:         GenericGetDeviceState,
		getDevicePropertyDesc:  GenericGetDevicePropertyDesc,
		getDevicePropertyValue: GenericGetDevicePropertyValue,
		setDeviceProperty:      GenericSetDeviceProperty,
		operationRequestRaw:    GenericOperationRequestRaw,
		initiateCapture:        GenericInitiateCapture,
	}

	switch c.ResponderVendor() {
	case ptp.VE_FujiPhotoFilmCoLtd:
		c.vendorExtensions.cmdDataInit = FujiInitCommandDataConn
		c.vendorExtensions.processStreamData = FujiProcessStreamData
		c.vendorExtensions.newCmdDataInitPacket = NewFujiInitCommandRequestPacket
		c.vendorExtensions.newEventInitPacket = NewFujiInitEventRequestPacket
		c.vendorExtensions.newEventPacket = NewFujiEventPacket
		c.vendorExtensions.extractTransactionId = FujiExtractTransactionId
		c.vendorExtensions.getDeviceInfo = FujiGetDeviceInfo
		c.vendorExtensions.getDeviceState = FujiGetDeviceState
		c.vendorExtensions.getDevicePropertyDesc = FujiGetDevicePropertyDesc
		c.vendorExtensions.getDevicePropertyValue = FujiGetDevicePropertyValue
		c.vendorExtensions.setDeviceProperty = FujiSetDeviceProperty
		c.vendorExtensions.operationRequestRaw = FujiSendOperationRequestAndGetRawResponse
		c.vendorExtensions.initiateCapture = FujiInitiateCapture
	}
}

// GenericInitCommandDataConn initiates the command/data connection. It expects an open TCP connection to the
// command/data port to be present.
func GenericInitCommandDataConn(c *Client) error {
	err := c.SendPacketToCmdDataConn(c.newCmdDataInitPacket())
	if err != nil {
		return err
	}

	res, _, err := c.waitForPacketFromCmdDataConn(nil)
	if err != nil {
		return err
	}

	switch pkt := res.(type) {
	case *InitFailPacket:
		err = pkt.ReasonAsError()
	case *InitCommandAckPacket:
		c.connectionNumber = pkt.ConnectionNumber
		c.responder.GUID = pkt.ResponderGUID
		c.responder.FriendlyName = pkt.ResponderFriendlyName
		c.responder.ProtocolVersion = pkt.ResponderProtocolVersion
		go c.responseListener()
		return nil
	default:
		err = fmt.Errorf("unexpected packet received %T", res)
	}

	c.Infoln("Closing Command/Data connection!")
	c.commandDataConn.Close()
	return err
}

// GenericInitEventConn initiates the event connection.
func GenericInitEventConn(c *Client) error {
	var err error

	c.eventConn, err = internal.RetryDialer(c.Network(), c.EventAddress(), DefaultDialTimeout)
	if err != nil {
		return err
	}

	c.configureTcpConn(eventConnection)

	ierp := c.newEventInitPacket()
	if ierp == nil {
		c.Info("No further event channel init required.")
		return nil
	}
	err = c.SendPacketToEventConn(ierp)
	if err != nil {
		return err
	}

	res, _, err := c.waitForPacketFromEventConn(nil)
	if err != nil {
		return err
	}

	switch pkt := res.(type) {
	case *InitFailPacket:
		err = pkt.ReasonAsError()
	case *InitEventAckPacket:
		c.incrementTransactionId()
		return nil
	default:
		err = fmt.Errorf("unexpected packet received %T", res)
	}

	c.Infoln("Closing Event connection!")
	c.eventConn.Close()
	return err
}

// GenericProcessStreamData does absolutely nothing since the standard PTP/IP protocol does not have a streamer
// connection.
func GenericProcessStreamData(_ *Client) error {
	return nil
}

// GenericExtractTransactionId extracts the transaction ID from a full raw inbound packet. This packet must include the
// full header containing length and packet type.
func GenericExtractTransactionId(p []byte, _ connectionType) (ptp.TransactionID, error) {
	if len(p) < 13 {
		return 0, fmt.Errorf("packet too small: got length %d", len(p))
	}

	var data []byte
	pt := PacketType(binary.LittleEndian.Uint32(p[4:8]))
	switch pt {
	case PKT_OperationResponse, PKT_Event:
		data = p[10:14]
	case PKT_StartData, PKT_Data, PKT_EndData, PKT_Cancel:
		data = p[8:12]
		// TODO: PKT_ProbeRequest and PKT_ProbeResponse do not have a transaction ID, how to handle those?
	}

	return ptp.TransactionID(binary.LittleEndian.Uint32(data)), nil
}

// Request the Responder's device information.
func GenericGetDeviceInfo(c *Client) (interface{}, error) {
	tid := c.incrementTransactionId()

	resCh := make(chan []byte, 2)
	defer close(resCh)
	if err := c.subscribe(tid, resCh); err != nil {
		return nil, err
	}

	err := c.SendPacketToCmdDataConn(&OperationRequestPacket{
		DataPhaseInfo:    DP_NoDataOrDataIn,
		OperationRequest: ptp.GetDeviceInfo(tid),
	})

	if err != nil {
		return nil, err
	}

	res, _, err := c.WaitForPacketFromCommandDataSubscriber(resCh, nil)
	if err != nil {
		return nil, err
	}

	switch pkt := res.(type) {
	case *OperationResponsePacket:
		return pkt, nil
	default:
		err = fmt.Errorf("unexpected packet received %T", res)
	}

	return nil, err
}

// GenericGetDeviceState requests the Responder's device status.
func GenericGetDeviceState(_ *Client) (interface{}, error) {
	return nil, errors.New("command not supported")
}

// GenericGetDevicePropertyValue requests the value for the given property from the Responder's.
func GenericGetDevicePropertyDesc(c *Client, dpc ptp.DevicePropCode) (*ptp.DevicePropDesc, error) {
	return nil, errors.New("command not YET supported")
}

// GenericGetDevicePropertyValue requests the value for the given property from the Responder's.
func GenericGetDevicePropertyValue(c *Client, dpc ptp.DevicePropCode) (uint32, error) {
	return 0, errors.New("command not YET supported")
}

// GenericSetDeviceProperty sets the value for the given property on the Responder.
func GenericSetDeviceProperty(c *Client, dpc ptp.DevicePropCode, val uint32) error {
	return errors.New("command not YET supported")
}

func GenericOperationRequestRaw(c *Client, code ptp.OperationCode, params []uint32) ([][]byte, error) {
	tid := c.incrementTransactionId()

	or := ptp.OperationRequest{
		OperationCode: code,
		TransactionID: tid,
	}

	// TODO: how to eliminate this crazyness WITHOUT reflection? Rework the OperationRequest struct perhaps with a
	//  [5]interface{} instead of 5 separate fields...?
	if len(params) >= 1 {
		or.Parameter1 = params[0]
	}
	if len(params) >= 2 {
		or.Parameter2 = params[1]
	}
	if len(params) >= 3 {
		or.Parameter3 = params[2]
	}
	if len(params) >= 4 {
		or.Parameter4 = params[3]
	}
	if len(params) == 5 {
		or.Parameter5 = params[4]
	}
	resCh := make(chan []byte, 2)
	if err := c.subscribe(tid, resCh); err != nil {
		return nil, err
	}
	defer c.unsubscribe(tid)

	err := c.SendPacketToCmdDataConn(&OperationRequestPacket{
		DataPhaseInfo:    DP_NoDataOrDataIn,
		OperationRequest: or,
	})

	if err != nil {
		return nil, err
	}

	var raw [][]byte
	raw[0], err = c.WaitForRawPacketFromCommandDataSubscriber(resCh)
	if err != nil {
		return nil, err
	}

	// TODO: handle possible followup packets depending on the data phase returned.

	return raw, err
}

func GenericInitiateCapture(c *Client) ([]byte, error) {
	return nil, errors.New("command not YET supported")
}
