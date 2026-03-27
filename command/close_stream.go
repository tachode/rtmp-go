package command

import "github.com/tachode/rtmp-go/message"

// NetStream.closeStream() command

func init() { RegisterCommand(new(CloseStream)) }

type CloseStream struct {
	StreamId    int
	Transaction int
}

func (c CloseStream) CommandName() string { return "closeStream" }

func (c *CloseStream) FromMessageCommand(cmd message.Command) error {
	return message.ReadFromCommand(cmd, c)
}

func (c *CloseStream) ToMessageCommand() (message.Command, error) {
	return message.BuildCommand(c.CommandName(), c), nil
}
