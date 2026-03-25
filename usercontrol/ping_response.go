package usercontrol

import (
	"github.com/tachode/rtmp-go/message"
)

// Ping Response: the client sends this event to the server in response to the
// ping request. Event data is the 4-byte timestamp received with the
// PingRequest.

func init() { RegisterEvent(new(PingResponse)) }

type PingResponse struct {
	Timestamp uint32
}

func (e PingResponse) EventType() message.UserControlMessageEvent {
	return message.UserControlPingResponse
}

func (e *PingResponse) FromMessage(msg *message.UserControlMessage) error {
	if len(msg.Parameters) < 1 {
		return message.ErrShortMessage
	}
	e.Timestamp = msg.Parameters[0]
	return nil
}

func (e *PingResponse) ToMessage() (*message.UserControlMessage, error) {
	return &message.UserControlMessage{
		Event:      message.UserControlPingResponse,
		Parameters: []uint32{e.Timestamp},
	}, nil
}
