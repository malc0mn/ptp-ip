package ip

import (
	"fmt"
	"github.com/google/uuid"
	ipInternal "github.com/malc0mn/ptp-ip/ip/internal"
	"github.com/malc0mn/ptp-ip/ptp"
)

type VendorExtensions struct {
	cmdDataInit          func(c *Client) error
	eventInit            func(c *Client) error
	streamerInit         func(c *Client) error
	newCmdDataInitPacket func(guid uuid.UUID, friendlyName string) InitCommandRequestPacket
}

func (c *Client) loadVendorExtensions() {
	c.vendorExtensions = &VendorExtensions{
		cmdDataInit:          GenericInitCommandDataConn,
		eventInit:            GenericInitEventConn,
		streamerInit:         GenericInitStreamerConn,
		newCmdDataInitPacket: NewInitCommandRequestPacket,
	}

	switch c.ResponderVendor() {
	case ptp.VE_FujiPhotoFilmCoLtd:
		c.vendorExtensions.cmdDataInit = FujiInitCommandDataConn
		c.vendorExtensions.newCmdDataInitPacket = NewFujiInitCommandRequestPacket
	}
}

func GenericInitCommandDataConn(c *Client) error {
	var err error

	c.commandDataConn, err = ipInternal.RetryDialer(c.Network(), c.CommandDataAddres(), DefaultDialTimeout)
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

	ierp := NewInitEventRequestPacket(c.connectionNumber)
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
