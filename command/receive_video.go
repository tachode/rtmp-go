package command

import "github.com/tachode/rtmp-go/message"

// NetStream.receiveVideo() command
// Informs the server whether to send video to the client.

func init() { RegisterCommand(new(ReceiveVideo)) }

type ReceiveVideo struct {
	StreamId    int
	Transaction int
	ReceiveFlag bool // true to receive video, false to stop.
}

func (r ReceiveVideo) CommandName() string { return "receiveVideo" }

func (r *ReceiveVideo) FromMessageCommand(cmd message.Command) error {
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

func (r *ReceiveVideo) ToMessageCommand() (message.Command, error) {
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
