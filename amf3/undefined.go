package amf3

// Undefined represents the AMF 3 undefined type (§3.2).
// No further information is encoded for this value.
type Undefined struct{}

func init() { RegisterType(new(Undefined)) }

func (v Undefined) Type() Type { return UndefinedMarker }

func (v *Undefined) Read(r *Reader) error {
	return nil
}

func (v Undefined) Write(w *Writer) error {
	return nil
}
