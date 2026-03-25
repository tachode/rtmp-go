package command

import (
	"github.com/tachode/rtmp-go/message"
)

// getStreamLength command

func init() { RegisterCommand(new(GetStreamLength)) }

type GetStreamLength struct {
	StreamId    int
	Transaction int
	StreamKey   string
}

func (g GetStreamLength) CommandName() string { return "getStreamLength" }

func (g *GetStreamLength) FromMessageCommand(cmd message.Command) error {
	g.StreamId = int(cmd.Metadata().StreamId)
	g.Transaction = int(cmd.GetTransactionId())
	params := cmd.GetParameters()
	if len(params) > 0 {
		if s, ok := message.ToString(params[0]); ok {
			g.StreamKey = s
		}
	}
	return nil
}

func (g *GetStreamLength) ToMessageCommand() (message.Command, error) {
	cmd := &message.Amf0CommandMessage{
		MetadataFields: message.MetadataFields{
			StreamId: uint32(g.StreamId),
		},
		Command:       g.CommandName(),
		TransactionId: float64(g.Transaction),
		Parameters:    []any{g.StreamKey},
	}
	return cmd, nil
}

func (g *GetStreamLength) MakeResponse(duration float64) message.Command {
	cmd := &message.Amf0CommandMessage{
		MetadataFields: message.MetadataFields{
			StreamId: uint32(g.StreamId),
		},
		Command:       "onResult",
		TransactionId: float64(g.Transaction),
		Object:        nil,
		Parameters:    []any{duration},
	}
	return cmd
}
