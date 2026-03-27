package command

import "github.com/tachode/rtmp-go/message"

// NetStream.pause() command
// Tells the server to pause or resume playing.

func init() { RegisterCommand(new(Pause)) }

type Pause struct {
	StreamId     int
	Transaction  int
	PauseFlag    bool    `amfParameter:"0"` // true to pause, false to resume play.
	MilliSeconds float64 `amfParameter:"1"` // Stream time (in ms) at which the client paused or resumed.
}

func (p Pause) CommandName() string { return "pause" }

func (p *Pause) FromMessageCommand(cmd message.Command) error {
	return message.ReadFromCommand(cmd, p)
}

func (p *Pause) ToMessageCommand() (message.Command, error) {
	return message.BuildCommand(p.CommandName(), p), nil
}

func (p *Pause) MakeStatus(status Status) message.Command {
	return streamStatusResponse(p.StreamId, p.Transaction, status)
}
