package ip

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/malc0mn/ptp-ip/ip/internal"
	"github.com/malc0mn/ptp-ip/ptp"
	"io"
	"log"
	"net"
	"os"
	"strings"
	"sync"
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
	BytesWrittenMismatch = "bytes written mismatch: written %d wanted %d"
	ConnectionLostError  = errors.New("connection lost")
	ReadResponseError    = errors.New("unable to read response packet")
	WaitForResponseError = errors.New("timeout reached when waiting for response")
	WaitForEventError    = errors.New("timeout reached when waiting for event")
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
//   - an async event channel receiving events from the Responder's event connection
//   - an async streamer channel receiving raw image data from the Responder's streaming connection if there is one
//   - a channel to request the streamer to close down
//   - a logger
type Client struct {
	connectionNumber uint32
	transactionId    ptp.TransactionID
	transactionIdMu  sync.Mutex
	commandDataConn  net.Conn
	eventConn        net.Conn
	streamConn       net.Conn
	initiator        *Initiator
	responder        *Responder
	vendorExtensions *VendorExtensions
	cmdDataChan      chan []byte
	cmdDataSubs      map[ptp.TransactionID]chan<- []byte
	cmdDataSubsMu    sync.Mutex
	eventChan        chan EventPacket
	StreamChan       chan []byte
	closeStreamChan  chan struct{}
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

// incrementTransactionId increments the transaction ID in a thread safe way.
func (c *Client) incrementTransactionId() ptp.TransactionID {
	c.transactionIdMu.Lock()
	c.transactionId++
	// The min and max values are considered 'invalid', so we roll over to 1 when we reach the max value.
	if c.transactionId == 0xFFFFFFFF {
		c.transactionId = 0x00000001
	}
	tid := c.transactionId // must copy the value before releasing the lock to reliably return it!
	c.transactionIdMu.Unlock()

	return tid
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

	err = c.initStreamConn()
	if err != nil {
		return err
	}

	return nil
}

// Close closes all open connections for the client.
func (c *Client) Close() error {
	var err error

	// streamConn must be closed first so we can do it cleanly, otherwise the camera might terminate it for us causing
	// any possible listeners to panic.
	if c.streamConn != nil {
		err = c.closeStreamConn()
		if err != nil {
			return err
		}
	}

	// TODO: add a closeEventConn() method so we can properly shut down the event channel like we do with the streamer.
	if c.eventConn != nil {
		err = c.eventConn.Close()
		c.eventConn = nil
		if err != nil {
			return err
		}
	}

	if c.commandDataConn != nil {
		err = c.commandDataConn.Close()
		c.commandDataConn = nil
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
		if _, err := w.Write(internal.MarshalLittleEndian(uint32(pll + 4))); err != nil {
			return err
		}
	} else {
		// The packet length MUST include the header, so we add 8 bytes for that!
		h := internal.MarshalLittleEndian(Header{uint32(pll + HeaderSize), p.PacketType()})

		// Send header.
		n, err := w.Write(h)
		if err != nil {
			return err
		}
		if n != HeaderSize {
			return fmt.Errorf(BytesWrittenMismatch, n, HeaderSize)
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
		return fmt.Errorf(BytesWrittenMismatch, n, pll)
	}
	c.Debugf("[sendPacket] payload bytes written %d", n)

	return nil
}

// readRawFromCmdDataConn reads raw data from the command/data connection with a read timout of 30 seconds.
func (c *Client) readRawFromCmdDataConn() ([]byte, error) {
	if c.commandDataConn == nil {
		return nil, fmt.Errorf("connection lost")
	}
	c.commandDataConn.SetReadDeadline(time.Now().Add(DefaultReadTimeout))
	return c.readRawResponse(c.commandDataConn)
}

// waitForRawFromCmdDataConn waits 30 seconds for a packet on the command/data connection.
func (c *Client) waitForRawFromCmdDataConn() ([]byte, error) {
	var (
		res []byte
		err error
	)

	for wait, timeout := true, time.After(DefaultReadTimeout); wait; {
		select {
		case <-timeout:
			wait = false
			err = WaitForResponseError
		default:
			res, err = c.readRawFromCmdDataConn()
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

// readPacketFromCmdDataConn reads a packet from the command/data connection with a read timout of 30 seconds.
// When expecting a specific packet, you can pass it in, otherwise pass nil.
// The byte array that is returned will contain any excess data that was not unmarshalled, empty otherwise.
func (c *Client) readPacketFromCmdDataConn(p PacketIn) (PacketIn, []byte, error) {
	if c.commandDataConn == nil {
		return nil, nil, ConnectionLostError
	}
	c.commandDataConn.SetReadDeadline(time.Now().Add(DefaultReadTimeout))
	return c.readResponse(c.commandDataConn, p)
}

// waitForPacketFromCmdDataConn waits 30 seconds for a packet on the command/data connection.
// This function will return a packet satisfying PacketIn together with any excess data that was not unmarshalled as a
// byte array. The excess data will be empty if there was none.
func (c *Client) waitForPacketFromCmdDataConn(p PacketIn) (PacketIn, []byte, error) {
	var (
		res PacketIn
		xs  []byte
		err error
	)

	for wait, timeout := true, time.After(DefaultReadTimeout); wait; {
		select {
		case <-timeout:
			wait = false
			err = WaitForResponseError
		default:
			res, xs, err = c.readPacketFromCmdDataConn(p)
			if err != io.EOF || res != nil {
				wait = false
			}
			time.Sleep(20 * time.Millisecond)
		}
	}
	if err != nil {
		return nil, nil, err
	}

	return res, xs, nil
}

// readPacketFromEventConn reads a packet from the Event connection.
// The byte array that is returned will contain any excess data that was not unmarshalled, empty otherwise.
func (c *Client) readPacketFromEventConn(p PacketIn) (PacketIn, []byte, error) {
	if c.eventConn == nil {
		return nil, nil, ConnectionLostError
	}
	c.eventConn.SetReadDeadline(time.Now().Add(DefaultReadTimeout))
	return c.readResponse(c.eventConn, p)
}

// waitForPacketFromEventConn waits for a packet on the Event connection.
// This function will return a packet satisfying EventPacket together with any excess data that was not unmarshalled as
// a byte array. The excess data will be empty if there was none.
func (c *Client) waitForPacketFromEventConn(p EventPacket) (PacketIn, []byte, error) {
	var (
		res PacketIn
		xs  []byte
		err error
	)

	for wait, timeout := true, time.After(DefaultReadTimeout); wait; {
		select {
		case <-timeout:
			wait = false
			err = WaitForEventError
		default:
			res, xs, err = c.readPacketFromEventConn(p)
			if err != io.EOF || res != nil {
				wait = false
			}
			time.Sleep(20 * time.Millisecond)
		}
	}
	if err != nil {
		return nil, nil, err
	}

	return res, xs, nil
}

// ReadRawFromStreamConn reads raw data from the streamer connection with a read timout of 30 seconds.
func (c *Client) ReadRawFromStreamConn() ([]byte, error) {
	c.commandDataConn.SetReadDeadline(time.Now().Add(DefaultReadTimeout))
	return c.readRawResponse(c.streamConn)
}

// TODO: this must be refactored to work like the events: continuously read and push to a channel in such a way that we
//  do not mix up packets (use transaction ID properly) like what's happening now with liveview polling the camera state
//  every second.
func (c *Client) readResponse(r io.Reader, p PacketIn) (PacketIn, []byte, error) {
	var err error
	var h Header
	var hl int

	// An invalid packet type means it does not adhere to the PTP/IP standard, so we only read the length field here.
	if p != nil && p.PacketType() == PKT_Invalid {
		var l uint32
		if err := binary.Read(r, binary.LittleEndian, &l); err != nil {
			return nil, nil, err
		}
		hl = int(l) - 4
	} else {
		if err := binary.Read(r, binary.LittleEndian, &h); err != nil {
			return nil, nil, err
		}

		if h.Length == 0 {
			return nil, nil, ReadResponseError
		}
		hl = int(h.Length) - HeaderSize
	}

	if p == nil {
		if p, err = NewPacketInFromPacketType(h.PacketType); err != nil {
			return nil, nil, err
		}
	}

	// TODO: this variable string calculation works for now, but there MUST be a better way to handle this!
	// We calculate the size of the variable portion of the packet here!
	// If there is no variable portion, vs will be 0.
	vs := hl - p.TotalFixedFieldSize()
	xs, err := internal.UnmarshalLittleEndian(r, p, hl, vs)
	if err != nil && err != io.EOF {
		return nil, nil, err
	}

	return p, xs, nil
}

// TODO: this must be refactored to work like the events: continuously read and push to a channel in such a way that we
//  do not mix up packets (use transaction ID properly) like what's happening now with liveview polling the camera state
//  every second.
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

// subscribe registers a channel to receive responses for a specific transaction ID.
func (c *Client) subscribe(tid ptp.TransactionID, ch chan<- []byte) error {
	c.cmdDataSubsMu.Lock()
	defer c.cmdDataSubsMu.Unlock()

	if _, ok := c.cmdDataSubs[tid]; ok {
		return fmt.Errorf("attempt to double subscribe transaction id %d", tid)
	}
	c.cmdDataSubs[tid] = ch

	return nil
}

// unsubscribe removes a subscription for a given transaction ID and closes the corresponding channel.
func (c *Client) unsubscribe(tid ptp.TransactionID) {
	c.cmdDataSubsMu.Lock()
	if ch, ok := c.cmdDataSubs[tid]; ok {
		close(ch)
		delete(c.cmdDataSubs, tid)
	}
	c.cmdDataSubsMu.Unlock()
}

// responseListener listens on the Command/Data connection for incoming packets and publishes them to a registered
// subscriber based on the transaction ID of the packet.
func (c *Client) responseListener() {
	c.cmdDataChan = make(chan []byte, 10)
	lmp := "[responseListener]"
	c.Infof("%s subscribing response listener to command/data connection...", lmp)
	for {
		p, err := c.waitForRawFromCmdDataConn()
		if err == nil {
			tid, err := c.vendorExtensions.extractTransactionId(p, cmdDataConnection)
			if err != nil {
				c.Error(err)
				continue
			}
			c.Debugf("%s publishing new response with length '%d' for transaction ID '%d'...", lmp, binary.LittleEndian.Uint32(p[0:4]), tid)

			if _, ok := c.cmdDataSubs[tid]; !ok {
				panic(fmt.Sprintf("No subscriber for transaction ID %d!", tid))
			}
			c.cmdDataSubs[tid] <- p
			continue
		} else if err == WaitForResponseError || strings.Contains(err.Error(), "i/o timeout") {
			continue
		}
		c.Errorf("%s message listener stopped: %s", lmp, err)
		return
	}
}

func (c *Client) initCommandDataConn() error {
	var err error

	c.commandDataConn, err = internal.RetryDialer(c.Network(), c.CommandDataAddress(), DefaultDialTimeout)
	if err != nil {
		return err
	}

	c.configureTcpConn(cmdDataConnection)

	if err := c.vendorExtensions.cmdDataInit(c); err != nil {
		return fmt.Errorf("command data connection: %s", err)
	}

	return nil
}

// WaitForRawPacketFromCommandDataSubscriber waits 30 seconds for a packet to be sent to a command/data channel
// subscriber registered using the subscribe method.
func (c *Client) WaitForRawPacketFromCommandDataSubscriber(ch <-chan []byte) ([]byte, error) {
	var (
		res []byte
		err error
	)

	for wait, timeout := true, time.After(DefaultReadTimeout); wait; {
		select {
		case <-timeout:
			wait = false
			err = WaitForResponseError
		case res = <-ch:
			wait = false
		}
	}
	if err != nil {
		return nil, err
	}

	return res, nil
}

// WaitForPacketFromCommandDataSubscriber waits 30 seconds for a packet to be sent to a command/data channel subscriber
// registered using the subscribe method.
// This function will return a packet satisfying PacketIn together with any excess data that was not unmarshalled as a
// byte array. The excess data will be empty if there was none.
func (c *Client) WaitForPacketFromCommandDataSubscriber(ch <-chan []byte, p PacketIn) (PacketIn, []byte, error) {
	res, err := c.WaitForRawPacketFromCommandDataSubscriber(ch)
	if err != nil {
		return nil, nil, err
	}

	return c.readResponse(bytes.NewReader(res), p)
}

func (c *Client) newCmdDataInitPacket() InitCommandRequestPacket {
	return c.vendorExtensions.newCmdDataInitPacket(c.InitiatorGUID(), c.InitiatorFriendlyName())
}

// TODO: refactor this one to work exactly like the responseListener!
func (c *Client) initEventConn() error {
	if err := c.vendorExtensions.eventInit(c); err != nil {
		return fmt.Errorf("event connection error: %s", err)
	}

	lmp := "[eventListener]"
	c.eventChan = make(chan EventPacket, 10)
	go func() {
		c.Infof("%s subscribing event listener to event connection...", lmp)
		for {
			p := c.vendorExtensions.newEventPacket()
			_, _, err := c.waitForPacketFromEventConn(p)
			if err == nil {
				c.Debugf("%s publishing new event '%#x' to event channel...", lmp, p.GetEventCode())
				c.eventChan <- p
				continue
			} else if err == WaitForEventError || strings.Contains(err.Error(), "i/o timeout") {
				continue
			}
			c.Errorf("%s message listener stopped: %s", lmp, err)
			return
		}
	}()

	return nil
}

func (c *Client) newEventInitPacket() InitEventRequestPacket {
	return c.vendorExtensions.newEventInitPacket(c.connectionNumber)
}

func (c *Client) initStreamConn() error {
	if c.streamConn == nil {
		var err error

		c.streamConn, err = internal.RetryDialer(c.Network(), c.StreamerAddress(), DefaultDialTimeout)
		if err != nil {
			return err
		}

		c.configureTcpConn(streamConnection)

		c.StreamChan = make(chan []byte, 50)
		c.closeStreamChan = make(chan struct{})

		return c.vendorExtensions.processStreamData(c)
	}

	return nil
}

func (c *Client) closeStreamConn() error {
	if c.StreamChan != nil {
		close(c.closeStreamChan)
	}

	err := c.streamConn.Close()
	c.streamConn = nil

	return err
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
		initiator:   i,
		responder:   NewResponder(vendor, ip, port, port, port),
		cmdDataSubs: make(map[ptp.TransactionID]chan<- []byte),
		Logger:      NewLogger(logLevel, os.Stderr, "", log.LstdFlags),
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

// GetDevicePropertyDescription gets the description of the given device property.
func (c *Client) GetDevicePropertyDescription(code ptp.DevicePropCode) (*ptp.DevicePropDesc, error) {
	return c.vendorExtensions.getDevicePropertyDesc(c, code)
}

// GetDevicePropertyValue gets the value of the given device property.
func (c *Client) GetDevicePropertyValue(code ptp.DevicePropCode) (uint32, error) {
	return c.vendorExtensions.getDevicePropertyValue(c, code)
}

// SetDeviceProperty sets the given device property to the specified value.
func (c *Client) SetDeviceProperty(code ptp.DevicePropCode, val uint32) error {
	return c.vendorExtensions.setDeviceProperty(c, code, val)
}

// OperationRequestRaw allows to perform any operation request and returns the raw result intended for reverse
// engineering purposes.
func (c *Client) OperationRequestRaw(code ptp.OperationCode, params []uint32) ([][]byte, error) {
	return c.vendorExtensions.operationRequestRaw(c, code, params)
}

// InitiateCapture releases the shutter and captures an image. If the responder supports it, a preview of the captured
// image is returned as a byte array.
func (c *Client) InitiateCapture() ([]byte, error) {
	return c.vendorExtensions.initiateCapture(c)
}

// ToggleLiveView opens or closes the streamer connection on the camera, if it has one, and initiates or closes the
// StreamChan on the client.
// StreamChan will receive raw image data that can be processed by the client.
func (c *Client) ToggleLiveView(en bool) error {
	if en {
		return c.initStreamConn()
	}

	return c.closeStreamConn()
}
