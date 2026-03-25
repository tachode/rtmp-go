package usercontrol

import (
	"github.com/tachode/rtmp-go/message"
)

// Stream Dry: the server sends this event to notify the client that there is
// no more data on the stream.

func init() { RegisterEvent(new(StreamDry)) }

type StreamDry struct {
	StreamID uint32
}

func (e StreamDry) EventType() message.UserControlMessageEvent {
	return message.UserControlStreamDry
}

func (e *StreamDry) FromMessage(msg *message.UserControlMessage) error {
	if len(msg.Parameters) < 1 {
		return message.ErrShortMessage
	}
	e.StreamID = msg.Parameters[0]
	return nil
}

func (e *StreamDry) ToMessage() (*message.UserControlMessage, error) {
	return &message.UserControlMessage{
		Event:      message.UserControlStreamDry,
		Parameters: []uint32{e.StreamID},
	}, nil
}
