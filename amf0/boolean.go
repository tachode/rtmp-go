package amf0

import (
	"encoding/binary"
	"io"
)

type Boolean bool

func init() { RegisterType(new(Boolean)) }

func (v Boolean) Type() Type { return BooleanMarker }

func (v *Boolean) Read(r io.Reader) error {
	return binary.Read(r, binary.BigEndian, v)
}

func (v Boolean) Write(w io.Writer) error {
	return binary.Write(w, binary.BigEndian, v)
}
