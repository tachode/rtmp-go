package command

import "github.com/tachode/rtmp-go/message"

// NetStream.play() command

func init() { RegisterCommand(new(Play)) }

type Play struct {
	StreamId      int
	Transaction   int
	StreamKey     string  `amfParameter:"0"`           // Name of the stream to play.
	StartPosition float64 `amfParameter:"1"`           // Start position in seconds (-2 = live|recorded, -1 = live only, >= 0 = seek).
	Duration      float64 `amfParameter:"2,omitempty"` // Duration of playback in seconds (-1 = play until end).
	Reset         bool    `amfParameter:"3,omitempty"` // Whether to flush any previous playlist.
}

func (p Play) CommandName() string { return "play" }

func (p *Play) FromMessageCommand(cmd message.Command) error {
	return message.ReadFromCommand(cmd, p)
}

func (p *Play) ToMessageCommand() (message.Command, error) {
	return message.BuildCommand(p.CommandName(), p), nil
}

func (p *Play) MakeStatus(status Status) message.Command {
	return responseCommand("onStatus", p.StreamId, p.Transaction, status.ToObject())
}
