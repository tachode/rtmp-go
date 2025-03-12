package amf0

import (
	"encoding/binary"
	"io"
)

type Number float64

func init() { RegisterType(new(Number)) }

func (v Number) Type() Type { return NumberMarker }

func (v *Number) Read(r io.Reader) error {
	return binary.Read(r, binary.BigEndian, v)
}

func (v Number) Write(w io.Writer) error {
	return binary.Write(w, binary.BigEndian, v)
}
