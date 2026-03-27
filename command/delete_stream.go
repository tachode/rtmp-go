package command

import "github.com/tachode/rtmp-go/message"

// NetStream.deleteStream() command

func init() { RegisterCommand(new(DeleteStream)) }

type DeleteStream struct {
	Transaction    int
	DeleteStreamId int `amfParameter:"0"` // The ID of the stream that is destroyed on the server.
}

func (d DeleteStream) CommandName() string { return "deleteStream" }

func (d *DeleteStream) FromMessageCommand(cmd message.Command) error {
	return message.ReadFromCommand(cmd, d)
}

func (d *DeleteStream) ToMessageCommand() (message.Command, error) {
	return message.BuildCommand(d.CommandName(), d), nil
}

func (d *DeleteStream) MakeResponse(status Status) message.Command {
	return resultResponse(d.Transaction, status)
}
