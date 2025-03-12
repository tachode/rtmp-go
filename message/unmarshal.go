package message

import (
	"fmt"
	"reflect"
)

// typeRegistry is a map of message types to prototypical instances of the message.
var typeRegistry map[Type]Message

func RegisterType(v Message) {
	if typeRegistry == nil {
		typeRegistry = make(map[Type]Message)
	}
	typeRegistry[v.Type()] = v
}

func Unmarshal(timestamp uint32, typ Type, streamId uint32, payload []byte) (Message, error) {
	prototype, ok := typeRegistry[typ]
	if !ok {
		return nil, fmt.Errorf("unknown RTMP message %v", typ)
	}
	copy := reflect.New(reflect.Indirect(reflect.ValueOf(prototype)).Type()).Interface()
	message, ok := copy.(Message)
	if !ok {
		return nil, fmt.Errorf("invalid registered type %T does not implement Message interface", prototype)
	}
	message.Metadata().Length = uint32(len(payload))
	message.Metadata().Timestamp = timestamp
	message.Metadata().StreamId = streamId
	err := message.Unmarshal(payload)
	if err != nil {
		return nil, err
	}
	return message, nil
}
