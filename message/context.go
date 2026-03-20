package message

import (
	"fmt"
	"reflect"

	"github.com/tachode/rtmp-go/amf3"
)

// typeRegistry is a map of message types to prototypical instances of the message.
var typeRegistry map[Type]Message

func RegisterType(v Message) {
	if typeRegistry == nil {
		typeRegistry = make(map[Type]Message)
	}
	typeRegistry[v.Type()] = v
}

type Context struct {
	amf3Writer *amf3.Writer
	amf3Reader *amf3.Reader
}

func NewContext() *Context {
	return &Context{
		amf3Writer: amf3.NewWriter(nil),
		amf3Reader: amf3.NewReader(nil),
	}
}

func (c *Context) Unmarshal(timestamp uint32, typ Type, streamId uint32, payload []byte) (Message, error) {
	prototype, ok := typeRegistry[typ]
	var copy any
	if ok {
		copy = reflect.New(reflect.Indirect(reflect.ValueOf(prototype)).Type()).Interface()
	} else {
		copy = &UnimplementedMessage{messageType: typ}
	}
	message, ok := copy.(Message)
	if !ok {
		return nil, fmt.Errorf("invalid registered type %T does not implement Message interface", prototype)
	}
	message.Metadata().Length = uint32(len(payload))
	message.Metadata().Timestamp = timestamp
	message.Metadata().StreamId = streamId
	message.Metadata().context = c
	err := message.Unmarshal(payload)
	if err != nil {
		return nil, err
	}
	return message, nil
}

func (c *Context) Marshal(m Message) ([]byte, error) {
	m.Metadata().context = c
	return m.Marshal()
}
