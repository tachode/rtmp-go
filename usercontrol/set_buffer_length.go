package usercontrol

import (
	"github.com/tachode/rtmp-go/message"
)

// Set Buffer Length: the client sends this event to inform the server of the
// buffer size (in milliseconds) that is used to buffer any data coming over
// a stream.

func init() { RegisterEvent(new(SetBufferLength)) }

type SetBufferLength struct {
	StreamID     uint32
	BufferLength uint32 // in milliseconds
}

func (e SetBufferLength) EventType() message.UserControlMessageEvent {
	return message.UserControlSetBufferLength
}

func (e *SetBufferLength) FromMessage(msg *message.UserControlMessage) error {
	if len(msg.Parameters) < 2 {
		return message.ErrShortMessage
	}
	e.StreamID = msg.Parameters[0]
	e.BufferLength = msg.Parameters[1]
	return nil
}

func (e *SetBufferLength) ToMessage() (*message.UserControlMessage, error) {
	return &message.UserControlMessage{
		Event:      message.UserControlSetBufferLength,
		Parameters: []uint32{e.StreamID, e.BufferLength},
	}, nil
}
