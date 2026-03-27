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
	return message.ReadFromCommand(cmd, g)
}

func (g *GetStreamLength) ToMessageCommand() (message.Command, error) {
	return message.BuildCommand(g.CommandName(), g), nil
}

func (g *GetStreamLength) MakeResponse(duration float64) message.Command {
	return responseCommand("onResult", g.StreamId, g.Transaction, duration)
}
