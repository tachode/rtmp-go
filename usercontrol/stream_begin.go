package usercontrol

import (
	"github.com/tachode/rtmp-go/message"
)

// Stream Begin: the server sends this event to notify the client that a stream
// has become functional and can be used for communication.

func init() { RegisterEvent(new(StreamBegin)) }

type StreamBegin struct {
	StreamID uint32
}

func (e StreamBegin) EventType() message.UserControlMessageEvent {
	return message.UserControlStreamBegin
}

func (e *StreamBegin) FromMessage(msg *message.UserControlMessage) error {
	if len(msg.Parameters) < 1 {
		return message.ErrShortMessage
	}
	e.StreamID = msg.Parameters[0]
	return nil
}

func (e *StreamBegin) ToMessage() (*message.UserControlMessage, error) {
	return &message.UserControlMessage{
		Event:      message.UserControlStreamBegin,
		Parameters: []uint32{e.StreamID},
	}, nil
}
