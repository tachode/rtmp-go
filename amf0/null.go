package amf0

import (
	"io"
)

type Null struct{}

func init() { RegisterType(new(Null)) }

func (v Null) Type() Type { return NullMarker }

func (v *Null) Read(r io.Reader) error {
	return nil
}

func (v Null) Write(w io.Writer) error {
	return nil
}
