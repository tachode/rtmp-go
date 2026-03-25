package command

import "github.com/tachode/rtmp-go/message"

// NetStream.deleteStream() command

func init() { RegisterCommand(new(DeleteStream)) }

type DeleteStream struct {
	Transaction    int
	DeleteStreamId int // The ID of the stream that is destroyed on the server.
}

func (d DeleteStream) CommandName() string { return "deleteStream" }

func (d *DeleteStream) FromMessageCommand(cmd message.Command) error {
	d.Transaction = int(cmd.GetTransactionId())
	params := cmd.GetParameters()
	if len(params) > 0 {
		if n, ok := message.ToFloat64(params[0]); ok {
			d.DeleteStreamId = int(n)
		}
	}
	return nil
}

func (d *DeleteStream) ToMessageCommand() (message.Command, error) {
	cmd := &message.Amf0CommandMessage{
		Command:       d.CommandName(),
		TransactionId: float64(d.Transaction),
		Parameters:    []any{float64(d.DeleteStreamId)},
	}
	return cmd, nil
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
