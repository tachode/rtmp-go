package command

import (
	"github.com/tachode/rtmp-go/message"
)

// getStreamLength command

func init() { RegisterCommand(new(GetStreamLength)) }

type GetStreamLength struct {
	StreamId    int
	Transaction int
	StreamKey   string `amfParameter:"0"`
}

func (g GetStreamLength) CommandName() string { return "getStreamLength" }

func (g *GetStreamLength) FromMessageCommand(cmd message.Command) error {
	message.ReadFromCommand(cmd, g)
	return nil
}

func (g *GetStreamLength) ToMessageCommand() (message.Command, error) {
	return message.BuildCommand(g.CommandName(), g), nil
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
