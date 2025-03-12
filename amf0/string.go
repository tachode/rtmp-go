package amf0

import (
	"io"
)

type String string

func init() { RegisterType(new(String)) }

func (v String) Type() Type { return StringMarker }

func (v *String) Read(r io.Reader) error {
	value, err := readString[uint16](r)
	*v = String(value)
	return err
}

func (v String) Write(w io.Writer) error {
	return writeString[uint16](w, string(v))
}
