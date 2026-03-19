package amf3

// Null represents the AMF 3 null type (§3.3).
// No further information is encoded for this value.
type Null struct{}

func init() { RegisterType(new(Null)) }

func (v Null) Type() Type { return NullMarker }

func (v *Null) Read(r *Reader) error {
	return nil
}

func (v Null) Write(w *Writer) error {
	return nil
}
