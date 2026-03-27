package data

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/tachode/rtmp-go/message"
)

var UnknownHandlerError = errors.New("unknown data handler")

// Handler is the common interface for typed data message handlers.
type Handler interface {
	FromDataMessage(message.Data) error
	ToDataMessage() (message.Data, error)
	HandlerName() string
}

// handlerRegistry maps handler names to prototypical instances.
var handlerRegistry map[string]Handler

func RegisterHandler(v Handler) {
	if handlerRegistry == nil {
		handlerRegistry = make(map[string]Handler)
	}
	handlerRegistry[v.HandlerName()] = v
}

// FromDataMessage converts a message.Data into a typed Handler.
func FromDataMessage(msg message.Data) (Handler, error) {
	handlerName := msg.GetHandler()

	// @setDataFrame is a publisher-side wrapper; unwrap it.
	if handlerName == "@setDataFrame" {
		params := msg.GetParameters()
		if len(params) > 0 {
			if name, ok := message.ToString(params[0]); ok {
				handlerName = name
			}
		}
	}

	prototype, found := handlerRegistry[handlerName]
	if !found {
		return nil, fmt.Errorf("%w %s", UnknownHandlerError, handlerName)
	}
	copy := reflect.New(reflect.Indirect(reflect.ValueOf(prototype)).Type()).Interface()
	handler, ok := copy.(Handler)
	if !ok {
		return nil, fmt.Errorf("invalid registered type %T does not implement Handler interface", prototype)
	}
	err := handler.FromDataMessage(msg)
	return handler, err
}
