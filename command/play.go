package command

import "github.com/tachode/rtmp-go/message"

// NetStream.play() command
// TODO: fill in parameters once encoding is determined

func init() { RegisterCommand(new(Play)) }

type Play struct {
	Transaction int
}

func (p Play) CommandName() string { return "play" }

func (p *Play) FromMessageCommand(cmd message.Command) error {
	p.Transaction = int(cmd.GetTransactionId())
	return nil
}

func (p *Play) ToMessageCommand() (message.Command, error) {
	cmd := &message.Amf0CommandMessage{
		Command:       p.CommandName(),
		TransactionId: float64(p.Transaction),
	}
	return cmd, nil
}
