package command

import "github.com/tachode/rtmp-go/message"

// NetConnection.close() command

func init() { RegisterCommand(new(Close)) }

type Close struct {
	Transaction int
}

func (c Close) CommandName() string { return "close" }

func (c *Close) FromMessageCommand(cmd message.Command) error {
	return message.ReadFromCommand(cmd, c)
}

func (c *Close) ToMessageCommand() (message.Command, error) {
	return message.BuildCommand(c.CommandName(), c), nil
}
