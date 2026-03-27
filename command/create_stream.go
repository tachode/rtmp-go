package command

import "github.com/tachode/rtmp-go/message"

// NetConnection.createStream() command

func init() { RegisterCommand(new(CreateStream)) }

type CreateStream struct {
	Transaction int
}

func (c CreateStream) CommandName() string { return "createStream" }

func (c *CreateStream) FromMessageCommand(cmd message.Command) error {
	message.ReadFromCommand(cmd, c)
	return nil
}

func (c *CreateStream) ToMessageCommand() (message.Command, error) {
	return message.BuildCommand(c.CommandName(), c), nil
}

func (c *CreateStream) MakeResponse(streamId int) message.Command {
	cmd := &message.Amf0CommandMessage{
		Command:       "_result",
		TransactionId: float64(c.Transaction),
		Object:        nil,
		Parameters:    []any{streamId},
	}
	return cmd
}

func (c *CreateStream) MakeErrorResponse(status Status) message.Command {
	cmd := &message.Amf0CommandMessage{
		Command:       "_error",
		TransactionId: float64(c.Transaction),
		Object:        nil,
		Parameters:    []any{status.ToObject()},
	}
	return cmd
}
