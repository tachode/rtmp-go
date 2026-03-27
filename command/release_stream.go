package command

import "github.com/tachode/rtmp-go/message"

// NetConnection.releaseStream() command

func init() { RegisterCommand(new(ReleaseStream)) }

type ReleaseStream struct {
	Transaction int
	StreamKey   string `amfParameter:"0"`
}

func (r ReleaseStream) CommandName() string { return "releaseStream" }

func (r *ReleaseStream) FromMessageCommand(cmd message.Command) error {
	return message.ReadFromCommand(cmd, r)
}

func (r *ReleaseStream) ToMessageCommand() (message.Command, error) {
	return message.BuildCommand(r.CommandName(), r), nil
}

func (r *ReleaseStream) MakeResponse(status Status) message.Command {
	return resultResponse(r.Transaction, status)
}
