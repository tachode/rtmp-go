package usercontrol

import (
	"github.com/tachode/rtmp-go/message"
)

// Stream EOF: the server sends this event to notify the client that the
// playback of data is over as requested on this stream.

func init() { RegisterEvent(new(StreamEOF)) }

type StreamEOF struct {
	StreamID uint32
}

func (e StreamEOF) EventType() message.UserControlMessageEvent {
	return message.UserControlStreamEOF
}

func (e *StreamEOF) FromMessage(msg *message.UserControlMessage) error {
	if len(msg.Parameters) < 1 {
		return message.ErrShortMessage
	}
	e.StreamID = msg.Parameters[0]
	return nil
}

func (e *StreamEOF) ToMessage() (*message.UserControlMessage, error) {
	return &message.UserControlMessage{
		Event:      message.UserControlStreamEOF,
		Parameters: []uint32{e.StreamID},
	}, nil
}
