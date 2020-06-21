package ip

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	ipInternal "github.com/malc0mn/ptp-ip/ip/internal"
	"github.com/malc0mn/ptp-ip/ptp"
)

// TODO: This solution is not OK, vendors can differ massively so it seems. Should this become an interface that all
//  vendors need to implement...? It would turn out to be a huge interface, so there will no doubt be a better solution?
//  Embedding ip.Client in a stuct ip.FujiClient wasn't that good either. Take the Dial() method for example, this calls
//  the initCommandDataConn() and initEventConn() methods but when using embedding the methods on ip.Client get called
//  and not the ones on ip.FujiClient so you would also have to "override" the Dial() as well.
type VendorExtensions struct {
	cmdDataInit          func(c *Client) error
	eventInit            func(c *Client) error
	streamerInit         func(c *Client) error
	newCmdDataInitPacket func(guid uuid.UUID, friendlyName string) InitCommandRequestPacket
	newEventInitPacket   func(connNum uint32) InitEventRequestPacket
	getDeviceInfo        func(c *Client) (interface{}, error)
	getDeviceState       func(c *Client) (interface{}, error)
	operationRequestRaw  func(c *Client, code ptp.OperationCode, params []uint32) ([][]byte, error)
}

func (c *Client) loadVendorExtensions() {
	c.vendorExtensions = &VendorExtensions{
		cmdDataInit:          GenericInitCommandDataConn,
		eventInit:            GenericInitEventConn,
		streamerInit:         GenericInitStreamerConn,
		newCmdDataInitPacket: NewInitCommandRequestPacket,
		newEventInitPacket:   NewInitEventRequestPacket,
		getDeviceInfo:        GenericGetDeviceInfo,
		getDeviceState:       GenericGetDeviceState,
		operationRequestRaw:  GenericOperationRequestRaw,
	}

	switch c.ResponderVendor() {
	case ptp.VE_FujiPhotoFilmCoLtd:
		c.vendorExtensions.cmdDataInit = FujiInitCommandDataConn
		c.vendorExtensions.newCmdDataInitPacket = NewFujiInitCommandRequestPacket
		c.vendorExtensions.newEventInitPacket = NewFujiInitEventRequestPacket
		c.vendorExtensions.getDeviceInfo = FujiGetDeviceInfo
		c.vendorExtensions.getDeviceState = FujiGetDeviceState
		c.vendorExtensions.operationRequestRaw = FujiOperationRequestRaw
	}
}

func GenericInitCommandDataConn(c *Client) error {
	var err error

	c.commandDataConn, err = ipInternal.RetryDialer(c.Network(), c.CommandDataAddress(), DefaultDialTimeout)
	if err != nil {
		return err
	}

	c.configureTcpConn(cmdDataConnection)

	err = c.SendPacketToCmdDataConn(c.newCmdDataInitPacket())
	if err != nil {
		return err
	}

	res, err := c.WaitForPacketFromCmdDataConn(nil)
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
		return nil
	default:
		err = fmt.Errorf("unexpected packet received %T", res)
	}

	c.Infoln("Closing Command/Data connection!")
	c.commandDataConn.Close()
	return err
}

func GenericInitEventConn(c *Client) error {
	var err error

	c.eventConn, err = ipInternal.RetryDialer(c.Network(), c.EventAddress(), DefaultDialTimeout)
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

	res, err := c.WaitForPacketFromEventConn()
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

func GenericInitStreamerConn(c *Client) error {
	var err error

	c.streamConn, err = ipInternal.RetryDialer(c.Network(), c.StreamerAddress(), DefaultDialTimeout)
	if err != nil {
		return err
	}

	c.configureTcpConn(streamConnection)

	return nil
}

// Request the Responder's device information.
func GenericGetDeviceInfo(c *Client) (interface{}, error) {
	err := c.SendPacketToCmdDataConn(&OperationRequestPacket{
		DataPhaseInfo:    DP_NoDataOrDataIn,
		OperationRequest: ptp.GetDeviceInfo(),
	})

	if err != nil {
		return nil, err
	}

	res, err := c.WaitForPacketFromCmdDataConn(nil)
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

// Request the Responder's device status.
func GenericGetDeviceState(c *Client) (interface{}, error) {
	return nil, errors.New("Command not supported!")
}

func GenericOperationRequestRaw(c *Client, code ptp.OperationCode, params []uint32) ([][]byte, error) {
	or := ptp.OperationRequest{
		OperationCode: code,
	}

	// TODO: how to eliminate this crazyness? Rework the OperationRequest struct perhaps with a [5]interface{} instead
	//  of 5 separate fields...?
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
		or.Parameter5 = params[5]
	}

	err := c.SendPacketToCmdDataConn(&OperationRequestPacket{
		DataPhaseInfo:    DP_NoDataOrDataIn,
		OperationRequest: or,
	})

	if err != nil {
		return nil, err
	}

	var raw [][]byte
	raw[0], err = c.ReadRawFromCmdDataConn()
	// TODO: handle possible followup packets depending on the data phase returned.

	return raw, err
}
