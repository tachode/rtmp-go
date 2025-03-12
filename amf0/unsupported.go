package amf0

import "io"

type Unsupported struct{}

func init() { RegisterType(new(Unsupported)) }

func (v Unsupported) Type() Type { return UnsupportedMarker }

func (v *Unsupported) Read(r io.Reader) error {
	return nil
}

func (v Unsupported) Write(w io.Writer) error {
	return nil
}
