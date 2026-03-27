package command

import (
	"github.com/tachode/rtmp-go/message"
)

// NetConnection.onStatus() command

func init() { RegisterCommand(new(OnStatus)) }

type OnStatus struct {
	Transaction int
	Status      `amfParameter:"0"`
}

func (o OnStatus) CommandName() string { return "onStatus" }

func (o *OnStatus) FromMessageCommand(cmd message.Command) error {
	return message.ReadFromCommand(cmd, o)
}

func (o *OnStatus) ToMessageCommand() (message.Command, error) {
	return message.BuildCommand(o.CommandName(), o), nil
}

func NewOnStatus(transaction int, status Status) *OnStatus {
	return &OnStatus{
		Transaction: transaction,
		Status:      status,
	}
}
