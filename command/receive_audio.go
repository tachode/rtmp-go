package command

import "github.com/tachode/rtmp-go/message"

// NetStream.receiveAudio() command
// Informs the server whether to send audio to the client.

func init() { RegisterCommand(new(ReceiveAudio)) }

type ReceiveAudio struct {
	StreamId    int
	Transaction int
	ReceiveFlag bool // true to receive audio, false to stop.
}

func (r ReceiveAudio) CommandName() string { return "receiveAudio" }

func (r *ReceiveAudio) FromMessageCommand(cmd message.Command) error {
	r.StreamId = int(cmd.Metadata().StreamId)
	r.Transaction = int(cmd.GetTransactionId())
	params := cmd.GetParameters()
	if len(params) > 0 {
		if b, ok := message.ToBool(params[0]); ok {
			r.ReceiveFlag = b
		}
	}
	return nil
}

func (r *ReceiveAudio) ToMessageCommand() (message.Command, error) {
	cmd := &message.Amf0CommandMessage{
		MetadataFields: message.MetadataFields{
			StreamId: uint32(r.StreamId),
		},
		Command:       r.CommandName(),
		TransactionId: float64(r.Transaction),
		Parameters:    []any{r.ReceiveFlag},
	}
	return cmd, nil
}
