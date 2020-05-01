package ip

import (
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/malc0mn/ptp-ip/internal"
	ipInternal "github.com/malc0mn/ptp-ip/ip/internal"
	"io"
	"log"
	"net"
)

const (
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

func NewDefaultInitiator() *Initiator {
	return NewInitiator("", uuid.Nil)
}

func NewInitiator(friendlyName string, guid uuid.UUID) *Initiator {
	if friendlyName == "" {
		friendlyName = InitiatorFriendlyName
	}
	if guid == uuid.Nil {
		var err error
		guid, err = uuid.NewRandom()
		internal.FailOnError(err)
	}
	i := Initiator{guid, friendlyName}
	return &i
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
	r := Responder{ip, port, uuid.Nil, ""}
	return &r
}

type Client struct {
	commandDataConn net.Conn
	eventConn       net.Conn
	streamConn      net.Conn
	initiator       *Initiator
	responder       *Responder
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

func (c *Client) InitiatorGUID() uuid.UUID {
	return c.initiator.GUID
}

func (c *Client) Dial() {
	icap := c.initCommandDataConn()
	c.initEventConn(icap)
}

func (c *Client) DialWithStreamer() {
	c.initCommandDataConn()
	c.initEventConn(1)
	c.initStreamerConn()
}

// Sends a packet to the Command/Data connection.
func (c *Client) SendPacketToCmdDataConn(p Packet) error {
	return c.sendPacket(c.commandDataConn, p)
}

// Send a packet to the Event connection.
func (c *Client) SendPacketToEventConn(p Packet) error {
	return c.sendPacket(c.commandDataConn, p)
}

func (c *Client) sendPacket(w io.Writer, p Packet) error {
	pl := p.Payload()

	// The packet length MUST include the header, so we add 8 bytes for that!
	lenBytes := 8
	h := ipInternal.MarshalLittleEndian(Header{uint32(len(pl) + lenBytes), p.PacketType()})

	// Send header.
	n, err := w.Write(h)
	internal.FailOnError(err)
	if n != lenBytes {
		return fmt.Errorf(BytesWrittenMismatch.Error(), n, lenBytes)
	}
	internal.LogDebug(fmt.Errorf("[sendPacket] bytes written %d", n))

	// Send payload.
	n, err = w.Write(pl)
	if n != len(pl) {
		return fmt.Errorf(BytesWrittenMismatch.Error(), n, len(pl))
	}
	internal.FailOnError(err)
	internal.LogDebug(fmt.Errorf("[sendPacket] bytes written %d", n))

	return nil
}

func (c *Client) ReadPacketFromCmdDataConn() (Packet, error) {
	return c.readResponse(c.commandDataConn)
}

func (c *Client) ReadPacketFromEventConn() (Packet, error) {
	return c.readResponse(c.eventConn)
}

func (c *Client) readResponse(r io.Reader) (Packet, error) {
	var h Header
	if err := binary.Read(r, binary.LittleEndian, &h); err != nil {
		return nil, err
	}

	if h.Length != 0 {
		p, err := NewPacketFromPacketType(h.PacketType)
		if err != nil {
			return nil, err
		}

		if err := ipInternal.UnmarshalLittleEndian(r, p); err != nil && err != io.EOF {
			return nil, err
		}

		return p, nil
	}

	return nil, ReadResponseError
}

func (c *Client) initCommandDataConn() uint32 {
	conn, err := net.Dial(c.Network(), c.String())
	internal.FailOnError(err)
	c.commandDataConn = conn

	icrp := NewInitCommandRequestPacket(c.InitiatorGUID(), c.InitiatorFriendlyName())
	c.SendPacketToCmdDataConn(icrp)
	r, err := c.ReadPacketFromCmdDataConn()
	if err != nil {
		log.Fatal(err)
	}
	return r.(*InitCommandAckPacket).ConnectionNumber
}

func (c *Client) initEventConn(connNum uint32) {
	conn, err := net.Dial(c.Network(), c.String())
	internal.FailOnError(err)
	c.eventConn = conn

	ierp := NewInitEventRequestPacket(connNum)
	c.SendPacketToEventConn(ierp)
}

// Not all devices will have a streamer service. When this connection fails, we will fail silently.
func (c *Client) initStreamerConn() {
	conn, err := net.Dial(c.Network(), c.String())
	if err != nil {
		internal.LogInfo(err)
		return
	}
	c.streamConn = conn
}

func NewClient(ip string, port int, friendlyName string, guid uuid.UUID) *Client {
	r := NewResponder(ip, port)
	i := NewInitiator(friendlyName, guid)

	c := Client{nil, nil, nil, i, r}
	return &c
}

/*
func InitCommandRequest() {

}

func InitCommandAck() {

}
*/
