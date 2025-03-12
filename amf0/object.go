package amf0

import (
	"io"
	"reflect"
)

type Object map[string]any

func init() { RegisterType(new(Object)) }

func (v Object) Type() Type { return ObjectMarker }

func (v *Object) Read(r io.Reader) error {
	*v = make(Object)
	var end ObjectEnd

	for {
		var key String
		err := key.Read(r)
		if err != nil {
			return err
		}
		value, err := Read(r)
		if err != nil {
			return err
		}
		if reflect.TypeOf(value) == reflect.TypeOf(end) {
			return nil
		}
		(*v)[string(key)] = value
	}
}

func (v Object) Write(w io.Writer) error {
	for key, val := range v {
		err := writeString[uint16](w, key)
		if err != nil {
			return err
		}
		err = Write(w, val)
		if err != nil {
			return err
		}

	}
	// Add an "ObjectEnd" entry to mark the end of the object
	err := writeString[uint16](w, "")
	if err != nil {
		return err
	}
	err = Write(w, &ObjectEnd{})
	if err != nil {
		return err
	}
	return nil
}
