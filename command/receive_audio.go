package command

import "github.com/tachode/rtmp-go/message"

// NetStream.receiveAudio() command
// Informs the server whether to send audio to the client.

func init() { RegisterCommand(new(ReceiveAudio)) }

type ReceiveAudio struct {
	StreamId    int
	Transaction int
	ReceiveFlag bool `amfParameter:"0"` // true to receive audio, false to stop.
}

func (r ReceiveAudio) CommandName() string { return "receiveAudio" }

func (r *ReceiveAudio) FromMessageCommand(cmd message.Command) error {
	message.ReadFromCommand(cmd, r)
	return nil
}

func (r *ReceiveAudio) ToMessageCommand() (message.Command, error) {
	return message.BuildCommand(r.CommandName(), r), nil
}
