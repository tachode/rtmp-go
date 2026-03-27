package command

import "github.com/tachode/rtmp-go/message"

// FCPublish() command

func init() { RegisterCommand(new(FCPublish)) }

type FCPublish struct {
	Transaction int
}

func (f FCPublish) CommandName() string { return "FCPublish" }

func (f *FCPublish) FromMessageCommand(cmd message.Command) error {
	return message.ReadFromCommand(cmd, f)
}

func (f *FCPublish) ToMessageCommand() (message.Command, error) {
	return message.BuildCommand(f.CommandName(), f), nil
}

func (f *FCPublish) MakeResponse(status Status) message.Command {
	return responseCommand("onFCPublish", 0, f.Transaction, status.ToObject())
}
