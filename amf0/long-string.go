package amf0

import (
	"io"
)

type LongString string

func init() { RegisterType(new(LongString)) }

func (v LongString) Type() Type { return LongStringMarker }

func (v *LongString) Read(r io.Reader) error {
	value, err := readString[uint32](r)
	*v = LongString(value)
	return err
}

func (v LongString) Write(w io.Writer) error {
	return writeString[uint32](w, string(v))
}
