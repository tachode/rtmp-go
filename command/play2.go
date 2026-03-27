package command

import (
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

// PlayOptions holds the NetStreamPlayOptions properties encoded as an AMF
// object in the first parameter of a play2 command.
type PlayOptions struct {
	StreamName    string         `amf:"streamName"`    // The name of the new stream to transition to or to play.
	OldStreamName string         `amf:"oldStreamName"` // The name of the old stream to transition from (empty if not transitioning).
	Start         float64        `amf:"start"`         // The start time in seconds (-2 = live|recorded, -1 = live only, >= 0 = seek).
	Len           float64        `amf:"len"`           // Duration of playback in seconds (-1 = until end).
	Offset        float64        `amf:"offset"`        // Absolute stream time for bitrate switch (-1 = fast switch).
	Transition    PlayTransition `amf:"transition"`    // The transition mode (see PlayTransition constants).
}

type Play2 struct {
	StreamId    int
	Transaction int
	PlayOptions `amfParameter:"0"`
}

func (p Play2) CommandName() string { return "play2" }

func (p *Play2) FromMessageCommand(cmd message.Command) error {
	message.ReadFromCommand(cmd, p)
	return nil
}

func (p *Play2) ToMessageCommand() (message.Command, error) {
	return message.BuildCommand(p.CommandName(), p), nil
}
