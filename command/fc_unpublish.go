package command

import "github.com/tachode/rtmp-go/message"

// FCUnpublish() command

func init() { RegisterCommand(new(FCUnpublish)) }

type FCUnpublish struct {
	Transaction int
}

func (f FCUnpublish) CommandName() string { return "FCUnpublish" }

func (f *FCUnpublish) FromMessageCommand(cmd message.Command) error {
	f.Transaction = int(cmd.GetTransactionId())
	return nil
}

func (f *FCUnpublish) ToMessageCommand() (message.Command, error) {
	cmd := &message.Amf0CommandMessage{
		Command:       f.CommandName(),
		TransactionId: float64(f.Transaction),
	}
	return cmd, nil
}

func (f *FCUnpublish) MakeResponse(status Status) message.Command {
	cmd := &message.Amf0CommandMessage{
		Command:       "onFCUnpublish",
		TransactionId: float64(f.Transaction),
		Object:        nil,
		Parameters:    []any{status.ToObject()},
	}
	return cmd
}
