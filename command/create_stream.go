package command

import "github.com/tachode/rtmp-go/message"

// NetConnection.createStream() command

func init() { RegisterCommand(new(CreateStream)) }

type CreateStream struct {
	Transaction int
}

func (c CreateStream) CommandName() string { return "createStream" }

func (c *CreateStream) FromMessageCommand(cmd message.Command) error {
	return message.ReadFromCommand(cmd, c)
}

func (c *CreateStream) ToMessageCommand() (message.Command, error) {
	return message.BuildCommand(c.CommandName(), c), nil
}

func (c *CreateStream) MakeResponse(streamId int) message.Command {
	return responseCommand("_result", 0, c.Transaction, streamId)
}

func (c *CreateStream) MakeErrorResponse(status Status) message.Command {
	return responseCommand("_error", 0, c.Transaction, status.ToObject())
}
