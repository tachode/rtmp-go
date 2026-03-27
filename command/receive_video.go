package command

import "github.com/tachode/rtmp-go/message"

// NetStream.receiveVideo() command
// Informs the server whether to send video to the client.

func init() { RegisterCommand(new(ReceiveVideo)) }

type ReceiveVideo struct {
	StreamId    int
	Transaction int
	ReceiveFlag bool `amfParameter:"0"` // true to receive video, false to stop.
}

func (r ReceiveVideo) CommandName() string { return "receiveVideo" }

func (r *ReceiveVideo) FromMessageCommand(cmd message.Command) error {
	return message.ReadFromCommand(cmd, r)
}

func (r *ReceiveVideo) ToMessageCommand() (message.Command, error) {
	return message.BuildCommand(r.CommandName(), r), nil
}
