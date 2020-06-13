package ip

import (
	"fmt"
	"github.com/google/uuid"
	ipInternal "github.com/malc0mn/ptp-ip/ip/internal"
	"github.com/malc0mn/ptp-ip/ptp"
)

// TODO: This solution is not OK, vendors can differ massively so it seems. Should this become an interface that all
//  vendors need to implement...? It would turn out to be a huge interface, so there will no doubt be a better solution?
type VendorExtensions struct {
	cmdDataInit          func(c *Client) error
	eventInit            func(c *Client) error
	streamerInit         func(c *Client) error
	newCmdDataInitPacket func(guid uuid.UUID, friendlyName string) InitCommandRequestPacket
	newEventInitPacket   func(connNum uint32) InitEventRequestPacket
	getDeviceInfo        func(c *Client) (PacketIn, error)
	operationRequestRaw  func(c *Client, code uint32, params []uint32) ([]byte, error)
}

func (c *Client) loadVendorExtensions() {
	c.vendorExtensions = &VendorExtensions{
		cmdDataInit:          GenericInitCommandDataConn,
		eventInit:            GenericInitEventConn,
		streamerInit:         GenericInitStreamerConn,
		newCmdDataInitPacket: NewInitCommandRequestPacket,
		newEventInitPacket:   NewInitEventRequestPacket,
		getDeviceInfo:        GenericGetDeviceInfo,
		operationRequestRaw:  GenericOperationRequestRaw,
	}

	switch c.ResponderVendor() {
	case ptp.VE_FujiPhotoFilmCoLtd:
		c.vendorExtensions.cmdDataInit = FujiInitCommandDataConn
		c.vendorExtensions.newCmdDataInitPacket = NewFujiInitCommandRequestPacket
		c.vendorExtensions.newEventInitPacket = NewFujiInitEventRequestPacket
		c.vendorExtensions.getDeviceInfo = FujiGetDeviceInfo
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

	c.log.Println("Closing Command/Data connection!")
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
		c.log.Print("No further event channel init required.")
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

	c.log.Println("Closing Event connection!")
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
func GenericGetDeviceInfo(c *Client) (PacketIn, error) {
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

func GenericOperationRequestRaw(c *Client, code uint32, params []uint32) ([]byte, error) {
	return nil, nil
}