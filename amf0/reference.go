package amf0

import (
	"encoding/binary"
	"io"
)

type Reference uint16

func init() { RegisterType(new(Reference)) }

func (v Reference) Type() Type { return ReferenceMarker }

func (v *Reference) Read(r io.Reader) error {
	return binary.Read(r, binary.BigEndian, v)
}

func (v Reference) Write(w io.Writer) error {
	return binary.Write(w, binary.BigEndian, v)
}
