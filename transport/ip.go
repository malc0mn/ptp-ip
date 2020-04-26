package transport

import (
	"fmt"
	"github.com/google/uuid"
)

const (
	DefaultPort           int    = 15740
	DefaultIpAddress      string = "192.168.0.1"
	InitiatorFriendlyName string = "Golang PTP/IP client"
)

type Inititor struct {
	GUID         uuid.UUID
	FriendlyName string
}

type Responder struct {
	IpAddress    string
	Port         int
	GUID         uuid.UUID
	FriendlyName string
}

// Implement the net.Addr interface
func (r *Responder) Network() string {
	return "tcp"
}
func (r *Responder) String() string {
	return fmt.Sprintf("%s:%d", r.IpAddress, r.Port)
}

func NewInitiator() *Inititor {
	// TODO: handle error from uuid.NewRandom.
	guid, _ := uuid.NewRandom()
	i := Inititor{guid, InitiatorFriendlyName}
	return &i
}

func NewResponder(ip string, port int) *Responder {
	r := Responder{ip, port, uuid.Nil, ""}
	return &r
}

/*
func InitCommandRequest() {

}

func InitCommandAck() {

}
*/
