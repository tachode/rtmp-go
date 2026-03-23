package command

import (
	"errors"

	"github.com/tachode/rtmp-go/amf0"
	"github.com/tachode/rtmp-go/message"
)

// NetResultion.connect() command

func init() { RegisterCommand(new(Result)) }

type Result struct {
	Transaction int
	Status
}

func (r Result) CommandName() string { return "_result" }

func (r *Result) FromMessageCommand(cmd message.Command) error {
	r.Transaction = int(cmd.GetTransactionId())
	obj := cmd.GetObject()
	if obj == nil {
		return errors.New("_result command contains no command object")
	}
	r.Status.Level = Level(GetString(obj, "level"))
	r.Status.Code = StatusCode(GetString(obj, "code"))
	r.Status.Description = GetString(obj, "description")
	return nil
}

func (r *Result) ToMessageCommand() (message.Command, error) {
	obj := amf0.Object{
		"fmsVer":       "FMS/4,0,0,1121",
		"capabilities": 255,
		"mode":         1,
	}
	_ = obj
	cmd := &message.Amf0CommandMessage{
		Command:       r.CommandName(),
		TransactionId: float64(r.Transaction),
		Object:        nil,
		Parameters: []any{
			amf0.Object{
				"level":          r.Level,
				"code":           r.Code,
				"description":    r.Description,
				"objectEncoding": 0,
			},
		},
	}
	return cmd, nil
}
