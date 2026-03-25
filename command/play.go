package command

import "github.com/tachode/rtmp-go/message"

// NetStream.play() command

func init() { RegisterCommand(new(Play)) }

type Play struct {
	StreamId      int
	Transaction   int
	StreamKey     string  // Name of the stream to play.
	StartPosition float64 // Start position in seconds (-2 = live|recorded, -1 = live only, >= 0 = seek).
	Duration      float64 // Duration of playback in seconds (-1 = play until end).
	Reset         bool    // Whether to flush any previous playlist.
}

func (p Play) CommandName() string { return "play" }

func (p *Play) FromMessageCommand(cmd message.Command) error {
	p.StreamId = int(cmd.Metadata().StreamId)
	p.Transaction = int(cmd.GetTransactionId())
	params := cmd.GetParameters()
	if len(params) > 0 {
		if s, ok := message.ToString(params[0]); ok {
			p.StreamKey = s
		}
	}
	if len(params) > 1 {
		if n, ok := message.ToFloat64(params[1]); ok {
			p.StartPosition = n
		}
	}
	if len(params) > 2 {
		if n, ok := message.ToFloat64(params[2]); ok {
			p.Duration = n
		}
	}
	if len(params) > 3 {
		if b, ok := message.ToBool(params[3]); ok {
			p.Reset = b
		}
	}
	return nil
}

func (p *Play) ToMessageCommand() (message.Command, error) {
	cmd := &message.Amf0CommandMessage{
		MetadataFields: message.MetadataFields{
			StreamId: uint32(p.StreamId),
		},
		Command:       p.CommandName(),
		TransactionId: float64(p.Transaction),
		Parameters:    []any{p.StreamKey, p.StartPosition, p.Duration, p.Reset},
	}
	return cmd, nil
}

func (p *Play) MakeStatus(status Status) message.Command {
	cmd := &message.Amf0CommandMessage{
		MetadataFields: message.MetadataFields{
			StreamId: uint32(p.StreamId),
		},
		Command:       "onStatus",
		TransactionId: float64(p.Transaction),
		Object:        nil,
		Parameters:    []any{status.ToObject()},
	}
	return cmd
}
