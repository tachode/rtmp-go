package command

import (
	"github.com/tachode/rtmp-go/message"
)

// NetConnection.onStatus() command

func init() { RegisterCommand(new(OnStatus)) }

type OnStatus struct {
	Transaction int
	Status
}

func (o OnStatus) CommandName() string { return "onStatus" }

func (o *OnStatus) FromMessageCommand(cmd message.Command) error {
	o.Transaction = int(cmd.GetTransactionId())
	params := cmd.GetParameters()
	if len(params) > 0 {
		if obj, ok := params[0].(message.Object); ok {
			o.Status.FromObject(obj)
		}
	}
	return nil
}

func (o *OnStatus) ToMessageCommand() (message.Command, error) {
	cmd := &message.Amf0CommandMessage{
		Command:       o.CommandName(),
		TransactionId: float64(o.Transaction),
		Parameters:    []any{o.Status.ToObject()},
	}
	return cmd, nil
}

func NewOnStatus(transaction int, status Status) *OnStatus {
	return &OnStatus{
		Transaction: transaction,
		Status:      status,
	}
}
