package command

import "github.com/tachode/rtmp-go/message"

// NetStream.play() command

func init() { RegisterCommand(new(Play)) }

type Play struct {
	StreamId      int
	Transaction   int
	StreamKey     string
	StartPosition float64
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
	return nil
}

func (p *Play) ToMessageCommand() (message.Command, error) {
	cmd := &message.Amf0CommandMessage{
		MetadataFields: message.MetadataFields{
			StreamId: uint32(p.StreamId),
		},
		Command:       p.CommandName(),
		TransactionId: float64(p.Transaction),
		Parameters:    []any{p.StreamKey, p.StartPosition},
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
