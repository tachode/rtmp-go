package amf0

import (
	"fmt"
	"io"
	"reflect"
)

// This typeRegistry contains a mapping from AMF0 type markers to
// prototypical instances
var typeRegistry map[Type]MutableValue

func RegisterType(v MutableValue) {
	if typeRegistry == nil {
		typeRegistry = make(map[Type]MutableValue)
	}
	typeRegistry[v.Type()] = v
}

// Reads a single value from a reader
func Read(r io.Reader) (out any, err error) {
	var marker [1]byte

	_, err = io.ReadFull(r, marker[:])
	if err != nil {
		return
	}

	prototype, ok := typeRegistry[Type(marker[0])]
	if !ok {
		return nil, fmt.Errorf("unknown AMF 0 type marker %v", int(marker[0]))
	}
	prototypeReadableCopy := reflect.New(reflect.Indirect(reflect.ValueOf(prototype)).Type()).Interface()
	readableOut, ok := prototypeReadableCopy.(MutableValue)
	if !ok {
		return nil, fmt.Errorf("invalid registered type %T does not implement Value interface", prototype)
	}

	err = readableOut.Read(r)

	// Cast to WritableValue to double-check that the type implements the Value interface correctly
	out, ok = reflect.Indirect(reflect.ValueOf(readableOut)).Interface().(Value)

	if !ok {
		return nil, fmt.Errorf("invalid registered type %T does not implement WritableValue interface", prototype)
	}

	return
}

// Convenience function when the type is known to be a string of some kind
func ReadString(r io.Reader) (string, error) {
	val, err := Read(r)
	if err != nil {
		return "", err
	}
	switch val := val.(type) {
	case String:
		return string(val), nil
	case LongString:
		return string(val), nil
	case XmlDocument:
		return string(val), nil
	default:
		return "", fmt.Errorf("expected string, got %T", val)
	}
}

// Convenience function when the type is known to be a Number
func ReadNumber(r io.Reader) (float64, error) {
	val, err := Read(r)
	if err != nil {
		return 0, err
	}
	switch val := val.(type) {
	case Number:
		return float64(val), nil
	default:
		return 0, fmt.Errorf("expected number, got %T", val)
	}
}

// Convenience function when the type is known to be an Object
func ReadObject(r io.Reader) (Object, error) {
	val, err := Read(r)
	if err != nil {
		return nil, err
	}
	switch val := val.(type) {
	case Object:
		return val, nil
	case Null:
		return nil, nil
	default:
		return nil, fmt.Errorf("expected object, got %T", val)
	}
}
