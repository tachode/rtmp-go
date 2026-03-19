package amf3

// Value represents a serializable AMF 3 value.
type Value interface {
	Type() Type
	Write(w *Writer) error
}

// MutableValue represents an AMF 3 value that can be deserialized.
type MutableValue interface {
	Value
	Read(r *Reader) error
}
