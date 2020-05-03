package ip

import (
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/malc0mn/ptp-ip/internal"
	ipInternal "github.com/malc0mn/ptp-ip/ip/internal"
	"io"
	"net"
	"time"
)

const (
	DefaultDialTimeout           = 10 * time.Second
	DefaultReadTimeout           = 30 * time.Second
	DefaultPort           int    = 15740
	DefaultIpAddress      string = "192.168.0.1"
	InitiatorFriendlyName string = "Golang PTP/IP client"
)

var (
	BytesWrittenMismatch = errors.New("bytes written mismatch: written %d wanted %d")
	ReadResponseError    = errors.New("unable to read response packet")
)

type Initiator struct {
	GUID         uuid.UUID
	FriendlyName string
}

func NewDefaultInitiator() (*Initiator, error) {
	return NewInitiator("", "")
}

func NewInitiator(friendlyName string, guid string) (*Initiator, error) {
	var err error
	var id uuid.UUID

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

	i := Initiator{
		GUID:         id,
		FriendlyName: friendlyName,
	}

	return &i, nil
}

type Responder struct {
	IpAddress    string
	Port         int
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

func NewResponder(ip string, port int) *Responder {
	r := Responder{
		IpAddress: ip,
		Port:      port,
	}
	return &r
}

type Client struct {
	connectionNumber uint32
	commandDataConn  net.Conn
	eventConn        net.Conn
	streamConn       net.Conn
	initiator        *Initiator
	responder        *Responder
}

func (c *Client) ConnectionNumber() uint32 {
	return c.connectionNumber
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
	internal.LogDebug(fmt.Errorf("[sendPacket] bytes written %d", n))

	// Send payload.
	n, err = w.Write(pl)
	if n != pll {
		return fmt.Errorf(BytesWrittenMismatch.Error(), n, pll)
	}
	if err != nil {
		return err
	}

	internal.LogDebug(fmt.Errorf("[sendPacket] bytes written %d", n))

	return nil
}

func (c *Client) ReadPacketFromCmdDataConn() (PacketIn, error) {
	c.commandDataConn.SetReadDeadline(time.Now().Add(DefaultReadTimeout))
	return c.readResponse(c.commandDataConn)
}

func (c *Client) ReadPacketFromEventConn() (PacketIn, error) {
	c.eventConn.SetReadDeadline(time.Now().Add(DefaultReadTimeout))
	return c.readResponse(c.eventConn)
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

	// TODO: this vs calculation works for now, but there must be a better way to handle this!
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

	icrp := NewInitCommandRequestPacket(c.InitiatorGUID(), c.InitiatorFriendlyName())
	err = c.SendPacketToCmdDataConn(icrp)
	if err != nil {
		return err
	}

	var res PacketIn
	res, err = c.ReadPacketFromCmdDataConn()
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

	ierp := NewInitEventRequestPacket(c.connectionNumber)
	err = c.SendPacketToEventConn(ierp)
	if err != nil {
		return err
	}

	var res PacketIn
	res, err = c.ReadPacketFromEventConn()
	if err != nil {
		return err
	}

	switch pkt := res.(type) {
	case *InitFailPacket:
		err = pkt.ReasonAsError()
	case *InitEventAckPacket:
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

	return nil
}

// Creates a new PTP/IP client.
// Passing an empty string to friendlyName will use the default friendly name.
// Passing an empty string as guid will generate a random V4 UUID upon initialisation.
func NewClient(ip string, port int, friendlyName string, guid string) (*Client, error) {
	i, err := NewInitiator(friendlyName, guid)
	if err != nil {
		return nil, err
	}

	c := Client{
		initiator: i,
		responder: NewResponder(ip, port),
	}

	return &c, nil
}

/*
func InitCommandRequest() {

}

func InitCommandAck() {

}
*/
