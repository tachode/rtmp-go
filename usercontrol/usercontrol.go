package usercontrol

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/tachode/rtmp-go/message"
)

var UnknownEventError = errors.New("unknown user control event")

type Event interface {
	FromMessage(*message.UserControlMessage) error
	ToMessage() (*message.UserControlMessage, error)
	EventType() message.UserControlMessageEvent
}

// eventRegistry maps event types to prototypical instances.
var eventRegistry map[message.UserControlMessageEvent]Event

func RegisterEvent(v Event) {
	if eventRegistry == nil {
		eventRegistry = make(map[message.UserControlMessageEvent]Event)
	}
	eventRegistry[v.EventType()] = v
}

func FromMessage(msg *message.UserControlMessage) (Event, error) {
	prototype, found := eventRegistry[msg.Event]
	if !found {
		return nil, fmt.Errorf("%w %s (%d)", UnknownEventError, msg.Event, msg.Event)
	}
	cp := reflect.New(reflect.Indirect(reflect.ValueOf(prototype)).Type()).Interface()
	event, ok := cp.(Event)
	if !ok {
		return nil, fmt.Errorf("invalid registered type %T does not implement Event interface", prototype)
	}
	err := event.FromMessage(msg)
	return event, err
}
