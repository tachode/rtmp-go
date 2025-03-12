package amf0

import "io"

type MutableValue interface {
	Read(r io.Reader) error
	Value
}

type Value interface {
	Type() Type
	Write(w io.Writer) error
}
