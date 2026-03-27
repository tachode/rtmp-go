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
	message.ReadFromCommand(cmd, d)
	return nil
}

func (d *DeleteStream) ToMessageCommand() (message.Command, error) {
	return message.BuildCommand(d.CommandName(), d), nil
}

func (d *DeleteStream) MakeResponse(status Status) message.Command {
	command := "_result"
	if status.Level == LevelError {
		command = "_error"
	}

	cmd := &message.Amf0CommandMessage{
		Command:       command,
		TransactionId: float64(d.Transaction),
		Object:        nil,
		Parameters:    []any{status.ToObject()},
	}
	return cmd
}
