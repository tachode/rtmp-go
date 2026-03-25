package usercontrol

import (
	"github.com/tachode/rtmp-go/message"
)

// Stream Is Recorded: the server sends this event to notify the client that
// the stream is a recorded stream.

func init() { RegisterEvent(new(StreamIsRecorded)) }

type StreamIsRecorded struct {
	StreamID uint32
}

func (e StreamIsRecorded) EventType() message.UserControlMessageEvent {
	return message.UserControlStreamIsRecorded
}

func (e *StreamIsRecorded) FromMessage(msg *message.UserControlMessage) error {
	if len(msg.Parameters) < 1 {
		return message.ErrShortMessage
	}
	e.StreamID = msg.Parameters[0]
	return nil
}

func (e *StreamIsRecorded) ToMessage() (*message.UserControlMessage, error) {
	return &message.UserControlMessage{
		Event:      message.UserControlStreamIsRecorded,
		Parameters: []uint32{e.StreamID},
	}, nil
}
