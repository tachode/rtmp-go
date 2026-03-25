package usercontrol

import (
	"github.com/tachode/rtmp-go/message"
)

// Ping Request: the server sends this event to test whether the client is
// reachable. Event data is a 4-byte timestamp representing the local server
// time when the server dispatched the command.

func init() { RegisterEvent(new(PingRequest)) }

type PingRequest struct {
	Timestamp uint32
}

func (e PingRequest) EventType() message.UserControlMessageEvent {
	return message.UserControlPingRequest
}

func (e *PingRequest) FromMessage(msg *message.UserControlMessage) error {
	if len(msg.Parameters) < 1 {
		return message.ErrShortMessage
	}
	e.Timestamp = msg.Parameters[0]
	return nil
}

func (e *PingRequest) ToMessage() (*message.UserControlMessage, error) {
	return &message.UserControlMessage{
		Event:      message.UserControlPingRequest,
		Parameters: []uint32{e.Timestamp},
	}, nil
}
