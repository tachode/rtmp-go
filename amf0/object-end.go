package amf0

import "io"

type ObjectEnd struct{}

func init() { RegisterType(new(ObjectEnd)) }

func (v ObjectEnd) Type() Type { return ObjectEndMarker }

func (v *ObjectEnd) Read(r io.Reader) error {
	return nil
}

func (v ObjectEnd) Write(w io.Writer) error {
	return nil
}
