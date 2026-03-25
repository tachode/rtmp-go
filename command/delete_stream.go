package command

import "github.com/tachode/rtmp-go/message"

// NetConnection.deleteStream() command

func init() { RegisterCommand(new(DeleteStream)) }

type DeleteStream struct {
	Transaction int
}

func (d DeleteStream) CommandName() string { return "deleteStream" }

func (d *DeleteStream) FromMessageCommand(cmd message.Command) error {
	d.Transaction = int(cmd.GetTransactionId())
	return nil
}

func (d *DeleteStream) ToMessageCommand() (message.Command, error) {
	cmd := &message.Amf0CommandMessage{
		Command:       d.CommandName(),
		TransactionId: float64(d.Transaction),
	}
	return cmd, nil
}

func (c *DeleteStream) MakeResponse(status Status) message.Command {
	command := "_result"
	if status.Level == LevelError {
		command = "_error"
	}

	cmd := &message.Amf0CommandMessage{
		Command:       command,
		TransactionId: float64(c.Transaction),
		Object:        nil,
		Parameters:    []any{status.ToObject()},
	}
	return cmd
}
