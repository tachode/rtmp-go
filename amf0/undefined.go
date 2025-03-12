package amf0

import "io"

type Undefined struct{}

func init() { RegisterType(new(Undefined)) }

func (v Undefined) Type() Type { return UndefinedMarker }

func (v *Undefined) Read(r io.Reader) error {
	return nil
}

func (v Undefined) Write(w io.Writer) error {
	return nil
}
