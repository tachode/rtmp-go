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
	message.ReadFromCommand(cmd, r)
	return nil
}

func (r *ReleaseStream) ToMessageCommand() (message.Command, error) {
	return message.BuildCommand(r.CommandName(), r), nil
}

func (r *ReleaseStream) MakeResponse(status Status) message.Command {
	command := "_result"
	if status.Level == LevelError {
		command = "_error"
	}

	cmd := &message.Amf0CommandMessage{
		Command:       command,
		TransactionId: float64(r.Transaction),
		Object:        nil,
		Parameters:    []any{status.ToObject()},
	}
	return cmd
}
