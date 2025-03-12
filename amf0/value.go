package amf0

import "io"

type Value interface {
	Read(r io.Reader) error
	WritableValue
}

type WritableValue interface {
	Type() Type
	Write(w io.Writer) error
}
