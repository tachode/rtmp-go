package command

import (
	"errors"

	"github.com/tachode/rtmp-go/amf0"
	"github.com/tachode/rtmp-go/message"
)

// NetConnection.connect() command

func init() { RegisterCommand(new(Connect)) }

type Connect struct {
	Transaction int
	App         string
	TcUrl       string
}

func (c Connect) CommandName() string { return "connect" }

func (c *Connect) FromMessageCommand(cmd message.Command) error {
	c.Transaction = int(cmd.GetTransactionId())
	obj := cmd.GetObject()
	if obj == nil {
		return errors.New("connect command contains no command object")
	}
	c.App = GetString(obj, "app")
	c.TcUrl = GetString(obj, "tcUrl")
	return nil
}

func (c *Connect) ToMessageCommand() (message.Command, error) {
	cmd := &message.Amf0CommandMessage{
		Command:       c.CommandName(),
		TransactionId: float64(c.Transaction),
		Object: amf0.Object{
			"app":   c.App,
			"tcUrl": c.TcUrl,
		},
	}
	return cmd, nil
}

func (c *Connect) MakeResponse(status Status, amfLevel int) message.Command {
	p0 := status.ToObject()
	p0["objectEncoding"] = amfLevel

	command := "_result"
	if status.Level == LevelError {
		command = "_error"
	}

	cmd := &message.Amf0CommandMessage{
		Command:       command,
		TransactionId: float64(c.Transaction),
		Object:        nil,
		Parameters:    []any{p0},
	}
	return cmd
}
