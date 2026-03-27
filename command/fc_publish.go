package command

import "github.com/tachode/rtmp-go/message"

// FCPublish() command

func init() { RegisterCommand(new(FCPublish)) }

type FCPublish struct {
	Transaction int
}

func (f FCPublish) CommandName() string { return "FCPublish" }

func (f *FCPublish) FromMessageCommand(cmd message.Command) error {
	message.ReadFromCommand(cmd, f)
	return nil
}

func (f *FCPublish) ToMessageCommand() (message.Command, error) {
	return message.BuildCommand(f.CommandName(), f), nil
}

func (f *FCPublish) MakeResponse(status Status) message.Command {
	cmd := &message.Amf0CommandMessage{
		Command:       "onFCPublish",
		TransactionId: float64(f.Transaction),
		Object:        nil,
		Parameters:    []any{status.ToObject()},
	}
	return cmd
}
