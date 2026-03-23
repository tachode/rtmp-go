package command

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/tachode/rtmp-go/message"
)

var UnknownCommandError = errors.New("unknown command")

type Command interface {
	FromMessageCommand(message.Command) error
	ToMessageCommand() (message.Command, error)
	CommandName() string
}

// commandRegistry is a map of message types to prototypical instances of the message.
var commandRegistry map[string]Command

func RegisterCommand(v Command) {
	if commandRegistry == nil {
		commandRegistry = make(map[string]Command)
	}
	commandRegistry[v.CommandName()] = v
}

func FromMessageCommand(cmd message.Command) (Command, error) {
	commandName := cmd.GetCommand()
	prototype, found := commandRegistry[commandName]
	if !found {
		return nil, fmt.Errorf("%w %s", UnknownCommandError, commandName)
	}
	copy := reflect.New(reflect.Indirect(reflect.ValueOf(prototype)).Type()).Interface()
	command, ok := copy.(Command)
	if !ok {
		return nil, fmt.Errorf("invalid registered type %T does not implement Command interface", prototype)
	}
	err := command.FromMessageCommand(cmd)
	return command, err
}
