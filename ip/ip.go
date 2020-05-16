package ip

import (
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/malc0mn/ptp-ip/internal"
	ipInternal "github.com/malc0mn/ptp-ip/ip/internal"
	"github.com/malc0mn/ptp-ip/ptp"
	"io"
	"net"
	"time"
)

const (
	DefaultDialTimeout           = 10 * time.Second
	DefaultReadTimeout           = 30 * time.Second
	DefaultPort           uint16 = 15740
	DefaultIpAddress             = "192.168.0.1"
	InitiatorFriendlyName        = "Golang PTP/IP client"
)

var (
	BytesWrittenMismatch = errors.New("bytes written mismatch: written %d wanted %d")
	ReadResponseError    = errors.New("unable to read response packet")
	WaitForResponseError = errors.New("timeout reached when waiting for response")
)

type Initiator struct {
	GUID         uuid.UUID
	FriendlyName string
}

// Creates a new Initiator using InitiatorFriendlyName as name and a randomly generated GUID.
// This is the same as calling NewInitiator with two empty strings as arguments.
func NewDefaultInitiator() (*Initiator, error) {
	return NewInitiator("", "")
}

// Creates a new Initiator with a friendlyName and GUID of your choosing.
// Passing an empt
func NewInitiator(friendlyName, guid string) (*Initiator, error) {
	var (
		err error
		id  uuid.UUID
	)

	if friendlyName == "" {
		friendlyName = InitiatorFriendlyName
	}

	if guid == "" {
		id, err = uuid.NewRandom()
		if err != nil {
			return nil, err
		}
	} else {
		id, err = uuid.Parse(guid)
		if err != nil {
			return nil, err
		}
	}

	i := &Initiator{
		GUID:         id,
		FriendlyName: friendlyName,
	}

	return i, nil
}

type Responder struct {
	IpAddress    string
	Port         uint16
	GUID         uuid.UUID
	FriendlyName string
}

// Implement the net.Addr interface
func (r Responder) Network() string {
	return "tcp"
}
func (r Responder) String() string {
	return fmt.Sprintf("%s:%d", r.IpAddress, r.Port)
}

func NewResponder(ip string, port uint16) *Responder {
	return &Responder{
		IpAddress: ip,
		Port:      port,
	}
}

type Client struct {
	connectionNumber uint32
	transactionId    ptp.TransactionID
	commandDataConn  net.Conn
	eventConn        net.Conn
	streamConn       net.Conn
	initiator        *Initiator
	responder        *Responder
}

func (c *Client) ConnectionNumber() uint32 {
	return c.connectionNumber
}

func (c *Client) TransactionId() ptp.TransactionID {
	return c.transactionId
}

func (c *Client) incrementTransactionId() {
	c.transactionId++
	// The min and max values are considered 'invalid', so we roll over to 1 when we reach the max value.
	if c.transactionId == 0xFFFFFFFF {
		c.transactionId = 0x00000001
	}
}

// Implement the net.Addr interface
func (c *Client) Network() string {
	return c.responder.Network()
}
func (c *Client) String() string {
	return c.responder.String()
}

func (c *Client) ResponderFriendlyName() string {
	return c.responder.FriendlyName
}

func (c *Client) InitiatorFriendlyName() string {
	return c.initiator.FriendlyName
}

func (c *Client) ResponderGUID() uuid.UUID {
	return c.responder.GUID
}

func (c *Client) ResponderGUIDAsString() string {
	return c.responder.GUID.String()
}

func (c *Client) InitiatorGUID() uuid.UUID {
	return c.initiator.GUID
}

func (c *Client) InitiatorGUIDAsString() string {
	return c.initiator.GUID.String()
}

func (c *Client) Dial() error {
	var err error

	err = c.initCommandDataConn()
	if err != nil {
		return err
	}

	err = c.initEventConn()
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) DialWithStreamer() error {
	var err error

	err = c.initCommandDataConn()
	if err != nil {
		return err
	}

	err = c.initEventConn()
	if err != nil {
		return err
	}

	err = c.initStreamerConn()
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) Close() error {
	var err error

	if c.commandDataConn != nil {
		err = c.commandDataConn.Close()
		if err != nil {
			return err
		}
	}

	if c.eventConn != nil {
		err = c.eventConn.Close()
		if err != nil {
			return err
		}
	}

	if c.streamConn != nil {
		c.streamConn.Close()
		if err != nil {
			return err
		}
	}

	return nil
}

// Sends a packet to the Command/Data connection.
func (c *Client) SendPacketToCmdDataConn(p PacketOut) error {
	return c.sendPacket(c.commandDataConn, p)
}

// Send a packet to the Event connection.
func (c *Client) SendPacketToEventConn(p PacketOut) error {
	return c.sendPacket(c.eventConn, p)
}

func (c *Client) sendPacket(w io.Writer, p PacketOut) error {
	internal.LogDebug(fmt.Errorf("[sendPacket] sending %T", p))

	pl := p.Payload()
	pll := len(pl)

	// The packet length MUST include the header, so we add 8 bytes for that!
	h := ipInternal.MarshalLittleEndian(Header{uint32(pll + HeaderSize), p.PacketType()})

	// Send header.
	n, err := w.Write(h)
	if err != nil {
		return err
	}
	if n != HeaderSize {
		return fmt.Errorf(BytesWrittenMismatch.Error(), n, HeaderSize)
	}
	internal.LogDebug(fmt.Errorf("[sendPacket] header bytes written %d", n))

	// Send payload.
	if pll == 0 {
		internal.LogDebug(errors.New("[sendPacket] packet has no payload"))
		return nil
	}

	n, err = w.Write(pl)
	if err != nil {
		return err
	}
	if n != pll {
		return fmt.Errorf(BytesWrittenMismatch.Error(), n, pll)
	}
	internal.LogDebug(fmt.Errorf("[sendPacket] payload bytes written %d", n))

	return nil
}

// Reads a packet from the Command/Data connection.
func (c *Client) ReadPacketFromCmdDataConn() (PacketIn, error) {
	c.commandDataConn.SetReadDeadline(time.Now().Add(DefaultReadTimeout))
	return c.readResponse(c.commandDataConn)
}

// Waits 30 seconds for a packet on the Command/Data connection.
func (c *Client) WaitForPacketFromCmdDataConn() (PacketIn, error) {
	var (
		res PacketIn
		err error
	)

	for wait, timeout := true, time.After(DefaultReadTimeout); wait; {
		select {
		case <-timeout:
			wait = false
			err = WaitForResponseError
		default:
			res, err = c.ReadPacketFromCmdDataConn()
			if err != io.EOF || res != nil {
				wait = false
			}
			time.Sleep(20 * time.Millisecond)
		}
	}
	if err != nil {
		return nil, err
	}

	return res, nil
}

// Reads a packet from the Event connection.
func (c *Client) ReadPacketFromEventConn() (PacketIn, error) {
	c.eventConn.SetReadDeadline(time.Now().Add(DefaultReadTimeout))
	return c.readResponse(c.eventConn)
}

// Waits 30 seconds for a packet on the Event connection.
func (c *Client) WaitForPacketFromEventConn() (PacketIn, error) {
	var (
		res PacketIn
		err error
	)

	for wait, timeout := true, time.After(DefaultReadTimeout); wait; {
		select {
		case <-timeout:
			wait = false
			err = WaitForResponseError
		default:
			res, err = c.ReadPacketFromEventConn()
			if err == nil {
				wait = false
			}
			time.Sleep(20 * time.Millisecond)
		}
	}

	return res, nil
}

func (c *Client) readResponse(r io.Reader) (PacketIn, error) {
	var h Header
	if err := binary.Read(r, binary.LittleEndian, &h); err != nil {
		return nil, err
	}

	if h.Length == 0 {
		return nil, ReadResponseError
	}

	p, err := NewPacketInFromPacketType(h.PacketType)
	if err != nil {
		return nil, err
	}

	// TODO: this variable string calculation works for now, but there MUST be a better way to handle this!
	// We calculate the size of the variable portion of the packet here!
	// If there is no variable portion, vs will be 0.
	vs := int(h.Length) - HeaderSize - p.TotalFixedFieldSize()
	if err := ipInternal.UnmarshalLittleEndian(r, p, vs); err != nil && err != io.EOF {
		return nil, err
	}

	return p, nil
}

func (c *Client) initCommandDataConn() error {
	var err error

	c.commandDataConn, err = ipInternal.RetryDialer(c.Network(), c.String(), DefaultDialTimeout)
	if err != nil {
		return err
	}

	err = c.commandDataConn.(*net.TCPConn).SetKeepAlive(true)
	if err != nil {
		internal.LogError(fmt.Errorf("TCP keepalive not enabled for command/data connection: %s", err))
		return err
	}

	icrp := NewInitCommandRequestPacket(c.InitiatorGUID(), c.InitiatorFriendlyName())
	err = c.SendPacketToCmdDataConn(icrp)
	if err != nil {
		return err
	}

	res, err := c.WaitForPacketFromCmdDataConn()
	if err != nil {
		return err
	}

	switch pkt := res.(type) {
	case *InitFailPacket:
		err = pkt.ReasonAsError()
	case *InitCommandAckPacket:
		c.connectionNumber = pkt.ConnectionNumber
		return nil
	default:
		err = fmt.Errorf("unexpected packet received %T", res)
	}

	c.commandDataConn.Close()
	return err
}

func (c *Client) initEventConn() error {
	var err error

	c.eventConn, err = ipInternal.RetryDialer(c.Network(), c.String(), DefaultDialTimeout)
	if err != nil {
		return err
	}

	err = c.eventConn.(*net.TCPConn).SetKeepAlive(true)
	if err != nil {
		internal.LogError(fmt.Errorf("TCP keepalive not enabled for event connection: %s", err))
		return err
	}

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

// Not all devices will have a streamer service. When this connection fails, we will fail silently.
func (c *Client) initStreamerConn() error {
	var err error

	c.streamConn, err = ipInternal.RetryDialer(c.Network(), c.String(), DefaultDialTimeout)
	if err != nil {
		return err
	}

	err = c.streamConn.(*net.TCPConn).SetKeepAlive(true)
	if err != nil {
		internal.LogError(fmt.Errorf("TCP keepalive not enabled for streamer connection: %s", err))
	}

	return nil
}

// Creates a new PTP/IP client.
// Passing an empty string to friendlyName will use the default friendly name.
// Passing an empty string as guid will generate a random V4 UUID upon initialisation.
func NewClient(ip string, port uint16, friendlyName string, guid string) (*Client, error) {
	i, err := NewInitiator(friendlyName, guid)
	if err != nil {
		return nil, err
	}

	c := &Client{
		initiator: i,
		responder: NewResponder(ip, port),
	}

	return c, nil
}
