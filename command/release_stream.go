package command

import "github.com/tachode/rtmp-go/message"

// NetConnection.releaseStream() command

func init() { RegisterCommand(new(ReleaseStream)) }

type ReleaseStream struct {
	Transaction int
	StreamKey   string
}

func (r ReleaseStream) CommandName() string { return "releaseStream" }

func (r *ReleaseStream) FromMessageCommand(cmd message.Command) error {
	r.Transaction = int(cmd.GetTransactionId())
	params := cmd.GetParameters()
	if len(params) > 0 {
		if s, ok := message.ToString(params[0]); ok {
			r.StreamKey = s
		}
	}
	return nil
}

func (r *ReleaseStream) ToMessageCommand() (message.Command, error) {
	cmd := &message.Amf0CommandMessage{
		Command:       r.CommandName(),
		TransactionId: float64(r.Transaction),
		Parameters:    []any{r.StreamKey},
	}
	return cmd, nil
}

func (r *ReleaseStream) MakeResponse(status Status) message.Command {
	cmd := &message.Amf0CommandMessage{
		Command:       "_result",
		TransactionId: float64(r.Transaction),
		Object:        nil,
		Parameters:    []any{status.ToObject()},
	}
	return cmd
}
