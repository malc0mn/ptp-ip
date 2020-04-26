package ip

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/malc0mn/ptp-ip/internal"
	ipInternal "github.com/malc0mn/ptp-ip/ip/internal"
	"net"
)

const (
	DefaultPort           int    = 15740
	DefaultIpAddress      string = "192.168.0.1"
	InitiatorFriendlyName string = "Golang PTP/IP client"
)

var (
	BytesWrittenMismatch = errors.New("Bytes written mismatch: written %d wanted %d")
)

type Initiator struct {
	GUID         uuid.UUID
	FriendlyName string
}

func NewInitiator(friendlyName string) *Initiator {
	if friendlyName == "" {
		friendlyName = InitiatorFriendlyName
	}
	guid, err := uuid.NewRandom()
	internal.FailOnError(err)
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
	c.initCommandDataConn()
	c.initEventConn()
}

func (c *Client) DialWithStreamer() {
	c.initCommandDataConn()
	c.initEventConn()
	c.initStreamerConn()
}

func (c *Client) SendPacket(packet interface{}) error {
	payload := ipInternal.ToBytesLittleEndian(packet)
	// The packet length MUST include the header, so we add 4 bytes for the length field!
	lenBytes := 4
	length := ipInternal.ToBytesLittleEndian(uint32(len(payload) + lenBytes))

	n, err := c.commandDataConn.Write(length)
	internal.FailOnError(err)
	if n != lenBytes {
		return fmt.Errorf(BytesWrittenMismatch.Error(), n, lenBytes)
	}
	internal.LogDebug(fmt.Errorf("bytes written %d", n))

	n, err = c.commandDataConn.Write(payload)
	if n != len(payload) {
		return fmt.Errorf(BytesWrittenMismatch.Error(), n, len(payload))
	}
	internal.FailOnError(err)
	internal.LogDebug(fmt.Errorf("bytes written %d", n))

	return nil
}

func (c *Client) initCommandDataConn() {
	conn, err := net.Dial(c.Network(), c.String())
	internal.FailOnError(err)
	c.commandDataConn = conn

	icrp := NewInitCommandRequestPacket(c.InitiatorGUID(), c.InitiatorFriendlyName())
	c.SendPacket(icrp)
}

func (c *Client) initEventConn() {
	conn, err := net.Dial(c.Network(), c.String())
	internal.FailOnError(err)
	c.commandDataConn = conn

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

func NewClient(ip string, port int, friendlyName string) *Client {
	r := NewResponder(ip, port)
	i := NewInitiator(friendlyName)

	c := Client{nil, nil, nil, i, r}
	return &c
}

/*
func InitCommandRequest() {

}

func InitCommandAck() {

}
*/