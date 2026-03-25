package command

import (
	"github.com/tachode/rtmp-go/amf0"
	"github.com/tachode/rtmp-go/message"
)

// NetStream.play2() command
// Unlike play, play2 can switch to a different bit rate stream without
// changing the timeline of the content played. The parameters are encoded
// as an AMF object with properties from the NetStreamPlayOptions class.

func init() { RegisterCommand(new(Play2)) }

// PlayTransition specifies how play2 transitions between streams.
type PlayTransition string

const (
	PlayTransitionAppend        PlayTransition = "append"        // Adds stream to playlist, begins playback with first stream.
	PlayTransitionAppendAndWait PlayTransition = "appendAndWait" // Builds playlist without starting playback.
	PlayTransitionReset         PlayTransition = "reset"         // Clears previous play calls, plays stream immediately.
	PlayTransitionResume        PlayTransition = "resume"        // Requests data from new connection starting where previous ended.
	PlayTransitionStop          PlayTransition = "stop"          // Stops playing the streams in a playlist.
	PlayTransitionSwap          PlayTransition = "swap"          // Replaces a content stream with different content.
	PlayTransitionSwitch        PlayTransition = "switch"        // Switches to a different bitrate of the same content.
)

type Play2 struct {
	StreamId      int
	Transaction   int
	StreamName    string         // The name of the new stream to transition to or to play.
	OldStreamName string         // The name of the old stream to transition from (empty if not transitioning).
	Start         float64        // The start time in seconds (-2 = live|recorded, -1 = live only, >= 0 = seek).
	Len           float64        // Duration of playback in seconds (-1 = until end).
	Offset        float64        // Absolute stream time for bitrate switch (-1 = fast switch).
	Transition    PlayTransition // The transition mode (see PlayTransition constants).
}

func (p Play2) CommandName() string { return "play2" }

func (p *Play2) FromMessageCommand(cmd message.Command) error {
	p.StreamId = int(cmd.Metadata().StreamId)
	p.Transaction = int(cmd.GetTransactionId())
	params := cmd.GetParameters()
	if len(params) > 0 {
		if obj, ok := params[0].(amf0.Object); ok {
			p.fromObject(obj)
		}
	}
	return nil
}

func (p *Play2) fromObject(obj amf0.Object) {
	if v, ok := obj["streamName"]; ok {
		if s, ok := message.ToString(v); ok {
			p.StreamName = s
		}
	}
	if v, ok := obj["oldStreamName"]; ok {
		if s, ok := message.ToString(v); ok {
			p.OldStreamName = s
		}
	}
	if v, ok := obj["start"]; ok {
		if n, ok := message.ToFloat64(v); ok {
			p.Start = n
		}
	}
	if v, ok := obj["len"]; ok {
		if n, ok := message.ToFloat64(v); ok {
			p.Len = n
		}
	}
	if v, ok := obj["offset"]; ok {
		if n, ok := message.ToFloat64(v); ok {
			p.Offset = n
		}
	}
	if v, ok := obj["transition"]; ok {
		if s, ok := message.ToString(v); ok {
			p.Transition = PlayTransition(s)
		}
	}
}

func (p *Play2) ToMessageCommand() (message.Command, error) {
	obj := amf0.Object{
		"streamName":    p.StreamName,
		"oldStreamName": p.OldStreamName,
		"start":         p.Start,
		"len":           p.Len,
		"offset":        p.Offset,
		"transition":    string(p.Transition),
	}
	cmd := &message.Amf0CommandMessage{
		MetadataFields: message.MetadataFields{
			StreamId: uint32(p.StreamId),
		},
		Command:       p.CommandName(),
		TransactionId: float64(p.Transaction),
		Parameters:    []any{obj},
	}
	return cmd, nil
}
