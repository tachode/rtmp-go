package command

import (
	"github.com/tachode/rtmp-go/message"
)

// NetStream.publish() command

func init() { RegisterCommand(new(Publish)) }

type HowToPublish string

const (
	HowToPublishRecord        HowToPublish = "record"
	HowToPublishAppend        HowToPublish = "append"
	HowToPublishAppendWithGap HowToPublish = "appendWithGap"
	HowToPublishLive          HowToPublish = "live"
)

type Publish struct {
	StreamId     int
	Transaction  int
	StreamKey    string
	HowToPublish HowToPublish
}

func (p Publish) CommandName() string { return "publish" }

func (p *Publish) FromMessageCommand(cmd message.Command) error {
	p.StreamId = int(cmd.Metadata().StreamId)
	p.Transaction = int(cmd.GetTransactionId())
	params := cmd.GetParameters()
	if len(params) > 0 {
		if s, ok := message.ToString(params[0]); ok {
			p.StreamKey = s
		}
	}
	if len(params) > 1 {
		if s, ok := message.ToString(params[1]); ok {
			p.HowToPublish = HowToPublish(s)
		} else {
			p.HowToPublish = HowToPublishLive
		}
	}
	return nil
}

func (p *Publish) ToMessageCommand() (message.Command, error) {
	cmd := &message.Amf0CommandMessage{
		MetadataFields: message.MetadataFields{
			StreamId: uint32(p.StreamId),
		},
		Command:       p.CommandName(),
		TransactionId: float64(p.Transaction),
		Parameters:    []any{p.StreamKey},
	}
	return cmd, nil
}

func (p *Publish) MakeStatus(status Status, clientId int) message.Command {
	p0 := status.ToObject()
	p0["clientid"] = clientId

	cmd := &message.Amf0CommandMessage{
		MetadataFields: message.MetadataFields{
			StreamId: uint32(p.StreamId),
		},
		Command:       "onStatus",
		TransactionId: float64(p.Transaction),
		Object:        nil,
		Parameters:    []any{p0},
	}
	return cmd
}
