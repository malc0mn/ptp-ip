package ip

import (
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/google/uuid"
	ipInternal "github.com/malc0mn/ptp-ip/ip/internal"
	"github.com/malc0mn/ptp-ip/ptp"
	"io"
	"log"
	"net"
	"os"
	"time"
)

const (
	DefaultVendor         string         = "generic"
	DefaultDialTimeout                   = 10 * time.Second
	DefaultReadTimeout                   = 30 * time.Second
	DefaultPort           uint16         = 15740
	DefaultIpAddress      string         = "192.168.0.1"
	InitiatorFriendlyName string         = "Golang PTP/IP client"
	cmdDataConnection     connectionType = "cmd"
	eventConnection       connectionType = "event"
	streamConnection      connectionType = "stream"
)

var (
	BytesWrittenMismatch = errors.New("bytes written mismatch: written %d wanted %d")
	ReadResponseError    = errors.New("unable to read response packet")
	WaitForResponseError = errors.New("timeout reached when waiting for response")
	InvalidPacketError   = errors.New("invalid packet")
	NotConnectedError    = errors.New("not connected")
)

type connectionType string

// Initiator holds the identity of "ourselves".
type Initiator struct {
	GUID         uuid.UUID
	FriendlyName string
}

// NewDefaultInitiator creates a new Initiator using InitiatorFriendlyName as name and a randomly generated GUID.
// This is the same as calling NewInitiator() with two empty strings as arguments.
func NewDefaultInitiator() (*Initiator, error) {
	return NewInitiator("", "")
}

// NewInitiator creates a new Initiator with a friendlyName and GUID of your choosing.
// Passing an empty string as friendlyName will result in InitiatorFriendlyName being used.
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

// Responder holds the information of our responder, i.e. the camera. The PTP/IP protocol is designed to work perfectly
// fine with a single port to handle the command/data and event channels. However some vendors have chosen to work with
// a separate port for each channel, which is why there are three possible ports in this struct.
type Responder struct {
	Vendor          ptp.VendorExtension
	IpAddress       string
	CommandDataPort uint16
	EventPort       uint16
	StreamerPort    uint16
	GUID            uuid.UUID
	FriendlyName    string
	ProtocolVersion uint32
}

// Network returns a fixed value: "tcp".
func (r Responder) Network() string {
	return "tcp"
}

// CommandDataAddress returns the address of the command/data channel as string in the form of host:port.
func (r Responder) CommandDataAddress() string {
	return fmt.Sprintf("%s:%d", r.IpAddress, r.CommandDataPort)
}

// EventAddress returns the address of the event channel as string in the form of host:port.
func (r Responder) EventAddress() string {
	return fmt.Sprintf("%s:%d", r.IpAddress, r.EventPort)
}

// StreamerAddress returns the address streamer channel as string in the form of host:port.
func (r Responder) StreamerAddress() string {
	return fmt.Sprintf("%s:%d", r.IpAddress, r.StreamerPort)
}

// NewResponder creates a new responder struct.
func NewResponder(vendor string, ip string, cport uint16, eport uint16, sport uint16) *Responder {
	return &Responder{
		Vendor:          ptp.VendorStringToType(vendor),
		IpAddress:       ip,
		CommandDataPort: cport,
		EventPort:       eport,
		StreamerPort:    sport,
	}
}

// Client holds all parts needed to build our PTP/IP client:
//   - the connection number
//   - the current transaction ID
//   - the command/data channel connection
//   - the event channel connection
//   - the streamer channel connection
//   - the initiator info, i.e. us
//   - the responder info, i.e. camera
//   - the loaded vendor extensions
//   - a logger
type Client struct {
	connectionNumber uint32
	transactionId    ptp.TransactionID
	commandDataConn  net.Conn
	eventConn        net.Conn
	streamConn       net.Conn
	initiator        *Initiator
	responder        *Responder
	vendorExtensions *VendorExtensions
	EventChan        chan EventPacket
	Logger
}

// ConnectionNumber returns the connection number received from the responder after initialising the command/data
// connection.
func (c *Client) ConnectionNumber() uint32 {
	return c.connectionNumber
}

// TransactionId returns the current transaction ID.
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

// Network returns a fixed value: "tcp".
func (c *Client) Network() string {
	return c.responder.Network()
}

// CommandDataAddress returns the address from the responder's command/data channel as string in the form of host:port.
func (c *Client) CommandDataAddress() string {
	return c.responder.CommandDataAddress()
}

// EventAddress returns the address from the responder's event channel as string in the form of host:port.
func (c *Client) EventAddress() string {
	return c.responder.EventAddress()
}

// StreamerAddress returns the address from the responder's streamer channel as string in the form of host:port.
func (c *Client) StreamerAddress() string {
	return c.responder.StreamerAddress()
}

// ResponderFriendlyName returns the responder's friendly name.
func (c *Client) ResponderFriendlyName() string {
	return c.responder.FriendlyName
}

// InitiatorFriendlyName returns the initiator's friendly name.
func (c *Client) InitiatorFriendlyName() string {
	return c.initiator.FriendlyName
}

// ResponderVendor returns the vendor code from the responder.
func (c *Client) ResponderVendor() ptp.VendorExtension {
	return c.responder.Vendor
}

// ResponderGUID returns the GUID from the responder as uuid.UUID.
func (c *Client) ResponderGUID() uuid.UUID {
	return c.responder.GUID
}

// ResponderGUIDAsString returns the GUID from the responder formatted as a string.
func (c *Client) ResponderGUIDAsString() string {
	return c.responder.GUID.String()
}

// InitiatorGUID returns the GUID from the initiator as uuid.UUID.
func (c *Client) InitiatorGUID() uuid.UUID {
	return c.initiator.GUID
}

// InitiatorGUIDAsString returns the GUID from the initiator formatted as a string.
func (c *Client) InitiatorGUIDAsString() string {
	return c.initiator.GUID.String()
}

// SetCommandDataPort allows setting the command/data channel port.
func (c *Client) SetCommandDataPort(port uint16) {
	c.responder.CommandDataPort = port
}

// SetEventPort allows setting the event channel port.
func (c *Client) SetEventPort(port uint16) {
	c.responder.EventPort = port
}

// SetStreamerPort allows setting the streamer channel port.
func (c *Client) SetStreamerPort(port uint16) {
	c.responder.StreamerPort = port
}

// SetLogger allows setting a custom logger. This defaults to the Go log package.
func (c *Client) SetLogger(log Logger) {
	c.Logger = log
}

// Dial will initialise the command/data and Event connections.
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

// DialWithStreamer will call Dial and also attempt to open the steamer channel used for live preview. Not all devices
// have such a channel.
func (c *Client) DialWithStreamer() error {
	var err error

	err = c.Dial()
	if err != nil {
		return err
	}

	err = c.initStreamerConn()
	if err != nil {
		return err
	}

	return nil
}

// Close closes all open connections for the client.
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

// SendPacketToCmdDataConn sends a packet to the command/data connection.
func (c *Client) SendPacketToCmdDataConn(p PacketOut) error {
	return c.sendPacket(c.commandDataConn, p)
}

// SendPacketToEventConn sends a packet to the Event connection.
func (c *Client) SendPacketToEventConn(p PacketOut) error {
	return c.sendPacket(c.eventConn, p)
}

// We write directly to the connection here without using bufio. The Payload() method and marshaling functions are
// already writing to a bytes buffer before we write to the connection.
func (c *Client) sendPacket(w io.Writer, p PacketOut) error {
	if w == nil {
		return NotConnectedError
	}
	if p == nil {
		return InvalidPacketError
	}
	c.Debugf("[sendPacket] sending %T", p)

	pl := p.Payload()
	pll := len(pl)

	// An invalid packet type means it does not adhere to the PTP/IP standard, so we only send the length field here.
	if p.PacketType() == PKT_Invalid {
		// Send length only. The length must include the size of the length field, so we add 4 bytes for that!
		if _, err := w.Write(ipInternal.MarshalLittleEndian(uint32(pll + 4))); err != nil {
			return err
		}
	} else {
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
		c.Debugf("[sendPacket] header bytes written %d", n)
	}

	// Send payload.
	if pll == 0 {
		c.Debugf("[sendPacket] packet has no payload")
		return nil
	}

	n, err := w.Write(pl)
	if err != nil {
		return err
	}
	if n != pll {
		return fmt.Errorf(BytesWrittenMismatch.Error(), n, pll)
	}
	c.Debugf("[sendPacket] payload bytes written %d", n)

	return nil
}

// ReadPacketFromCmdDataConn reads a packet from the command/data connection with a read timout of 30 seconds.
// When expecting a specific packet, you can pass it in, otherwise pass nil.
func (c *Client) ReadPacketFromCmdDataConn(p PacketIn) (PacketIn, error) {
	c.commandDataConn.SetReadDeadline(time.Now().Add(DefaultReadTimeout))
	return c.readResponse(c.commandDataConn, p)
}

// ReadRawFromCmdDataConn reads raw data from the command/data connection with a read timout of 5 seconds. It is
// intended primarily for debugging and/or reverse engineering purposes.
func (c *Client) ReadRawFromCmdDataConn() ([]byte, error) {
	c.commandDataConn.SetReadDeadline(time.Now().Add(5 * time.Second))
	return c.readRawResponse(c.commandDataConn)
}

// WaitForPacketFromCmdDataConn waits 30 seconds for a packet on the command/data connection.
func (c *Client) WaitForPacketFromCmdDataConn(p PacketIn) (PacketIn, error) {
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
			res, err = c.ReadPacketFromCmdDataConn(p)
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

// ReadPacketFromEventConn reads a packet from the Event connection.
func (c *Client) ReadPacketFromEventConn(p PacketIn) (PacketIn, error) {
	c.eventConn.SetReadDeadline(time.Now().Add(DefaultReadTimeout))
	return c.readResponse(c.eventConn, p)
}

// WaitForPacketFromEventConn waits for a packet on the Event connection.
func (c *Client) WaitForPacketFromEventConn(p EventPacket) (PacketIn, error) {
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
			res, err = c.ReadPacketFromEventConn(p)
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

func (c *Client) readResponse(r io.Reader, p PacketIn) (PacketIn, error) {
	var err error
	var h Header
	var hl int

	// An invalid packet type means it does not adhere to the PTP/IP standard, so we only read the length field here.
	if p != nil && p.PacketType() == PKT_Invalid {
		var l uint32
		if err := binary.Read(r, binary.LittleEndian, &l); err != nil {
			return nil, err
		}
		hl = int(l) - 4
	} else {
		if err := binary.Read(r, binary.LittleEndian, &h); err != nil {
			return nil, err
		}

		if h.Length == 0 {
			return nil, ReadResponseError
		}
		hl = int(h.Length) - HeaderSize
	}

	if p == nil {
		if p, err = NewPacketInFromPacketType(h.PacketType); err != nil {
			return nil, err
		}
	}

	// TODO: this variable string calculation works for now, but there MUST be a better way to handle this!
	// We calculate the size of the variable portion of the packet here!
	// If there is no variable portion, vs will be 0.
	vs := hl - p.TotalFixedFieldSize()
	if err := ipInternal.UnmarshalLittleEndian(r, p, hl, vs); err != nil && err != io.EOF {
		return nil, err
	}

	return p, nil
}

// The reading approach taken here is so that we can return the full raw data but still reliably read the complete
// expected data length.
func (c *Client) readRawResponse(r io.Reader) ([]byte, error) {
	l := make([]byte, 4)
	if err := binary.Read(r, binary.LittleEndian, &l); err != nil {
		return nil, err
	}

	len := binary.LittleEndian.Uint32(l)
	b := make([]byte, int(len)-4)
	if err := binary.Read(r, binary.LittleEndian, &b); err != nil {
		return nil, err
	}

	return append(l, b...), nil
}

func (c *Client) initCommandDataConn() error {
	err := c.vendorExtensions.cmdDataInit(c)
	if err != nil {
		return errors.New(fmt.Sprintf("command data connection: %s", err))
	}
	return nil
}

func (c *Client) newCmdDataInitPacket() InitCommandRequestPacket {
	return c.vendorExtensions.newCmdDataInitPacket(c.InitiatorGUID(), c.InitiatorFriendlyName())
}

func (c *Client) initEventConn() error {
	if err := c.vendorExtensions.eventInit(c); err != nil {
		return errors.New(fmt.Sprintf("event connection error: %s", err))
	}

	c.EventChan = make(chan EventPacket, 10)
	go func() {
		c.Info("Subscribing message listener to event connection...")
		for {
			p := c.vendorExtensions.newEventPacket()
			_, err := c.WaitForPacketFromEventConn(p)
			if err == nil {
				c.Debugln("Publishing new message to event channel...")
				c.EventChan <- p
				continue
			} else if err == WaitForResponseError {
				continue
			}
			c.Errorf("Message listener stopped: %s", err)
			return
		}
	}()

	return nil
}

func (c *Client) newEventInitPacket() InitEventRequestPacket {
	return c.vendorExtensions.newEventInitPacket(c.connectionNumber)
}

// Not all devices will have a streamer service. When this connection fails, we will fail silently.
func (c *Client) initStreamerConn() error {
	return c.vendorExtensions.streamerInit(c)
}

func (c *Client) configureTcpConn(t connectionType) {
	var conn net.Conn

	switch t {
	case cmdDataConnection:
		conn = c.commandDataConn
	case eventConnection:
		conn = c.eventConn
	case streamConnection:
		conn = c.streamConn
	}

	// The PTP/IP protocol specifically asks to enable keep alive.
	if err := conn.(*net.TCPConn).SetKeepAlive(true); err != nil {
		c.Warnf("TCP_KEEPALIVE not enabled for %s connection: %s", t, err)
	} else {
		c.Infof("TCP_KEEPALIVE enabled for %s connection", t)
	}

	// The PTP/IP protocol specifically asks to disable Nagle's algorithm. TCP_NODELAY SHOULD be enabled by default in
	// golang but there's no harm in making sure since performance here is negligible.
	if err := conn.(*net.TCPConn).SetNoDelay(true); err != nil {
		c.Warnf("TCP_NODELAY not enabled for %s connection: %s", t, err)
	} else {
		c.Infof("TCP_NODELAY enabled for %s connection", t)
	}
}

// NewClient creates a new PTP/IP client.
// Passing an empty string to friendlyName will use the default friendly name.
// Passing an empty string as guid will generate a random V4 UUID upon initialisation.
func NewClient(vendor string, ip string, port uint16, friendlyName string, guid string, logLevel LogLevel) (*Client, error) {
	i, err := NewInitiator(friendlyName, guid)
	if err != nil {
		return nil, err
	}

	c := &Client{
		initiator: i,
		responder: NewResponder(vendor, ip, port, port, port),
		Logger:    NewLogger(logLevel, os.Stderr, "", log.LstdFlags),
	}

	c.loadVendorExtensions()

	return c, nil
}

// GetDeviceInfo requests the Responder's device information. The data that should be returned is clearly specified by
// the PTP/IP protocol but will, alas, greatly differ from vendor to vendor.
func (c *Client) GetDeviceInfo() (interface{}, error) {
	return c.vendorExtensions.getDeviceInfo(c)
}

// GetDeviceState requests the Responder's device status. This is not part of the PTP/IP specification but is
// implemented by Fuji as a means to display the current camera settings in their mobile app.
func (c *Client) GetDeviceState() (interface{}, error) {
	return c.vendorExtensions.getDeviceState(c)
}

// GetDevicePropertyValue gets the value of the given device property.
func (c *Client) GetDevicePropertyValue(code ptp.DevicePropCode) (uint32, error) {
	return c.vendorExtensions.getDevicePropertyValue(c, code)
}

// OperationRequestRaw allows to perform any operation request and returns the raw result intended for reverse
// engineering purposes.
func (c *Client) OperationRequestRaw(code ptp.OperationCode, params []uint32) ([][]byte, error) {
	return c.vendorExtensions.operationRequestRaw(c, code, params)
}
