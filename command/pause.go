package command

import "github.com/tachode/rtmp-go/message"

// NetStream.pause() command
// Tells the server to pause or resume playing.

func init() { RegisterCommand(new(Pause)) }

type Pause struct {
	StreamId     int
	Transaction  int
	PauseFlag    bool    // true to pause, false to resume play.
	MilliSeconds float64 // Stream time (in ms) at which the client paused or resumed.
}

func (p Pause) CommandName() string { return "pause" }

func (p *Pause) FromMessageCommand(cmd message.Command) error {
	p.StreamId = int(cmd.Metadata().StreamId)
	p.Transaction = int(cmd.GetTransactionId())
	params := cmd.GetParameters()
	if len(params) > 0 {
		if b, ok := message.ToBool(params[0]); ok {
			p.PauseFlag = b
		}
	}
	if len(params) > 1 {
		if n, ok := message.ToFloat64(params[1]); ok {
			p.MilliSeconds = n
		}
	}
	return nil
}

func (p *Pause) ToMessageCommand() (message.Command, error) {
	cmd := &message.Amf0CommandMessage{
		MetadataFields: message.MetadataFields{
			StreamId: uint32(p.StreamId),
		},
		Command:       p.CommandName(),
		TransactionId: float64(p.Transaction),
		Parameters:    []any{p.PauseFlag, p.MilliSeconds},
	}
	return cmd, nil
}

func (p *Pause) MakeStatus(status Status) message.Command {
	command := "onStatus"
	if status.Level == LevelError {
		command = "_error"
	}

	cmd := &message.Amf0CommandMessage{
		MetadataFields: message.MetadataFields{
			StreamId: uint32(p.StreamId),
		},
		Command:       command,
		TransactionId: float64(p.Transaction),
		Object:        nil,
		Parameters:    []any{status.ToObject()},
	}
	return cmd
}
