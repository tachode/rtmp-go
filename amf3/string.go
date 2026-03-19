package amf3

// String represents the AMF 3 string type (§3.8).
// In AMF 3, there is a single string type (no separate long-string).
// Strings are encoded using UTF-8-vr, which supports string references
// via the implicit string reference table managed by Reader/Writer.
type String string

func init() { RegisterType(new(String)) }

func (v String) Type() Type { return StringMarker }

func (v *String) Read(r *Reader) error {
	s, err := r.readString()
	if err != nil {
		return err
	}
	*v = String(s)
	return nil
}

func (v String) Write(w *Writer) error {
	return w.writeString(string(v))
}
