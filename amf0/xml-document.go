package amf0

import (
	"io"
)

type XmlDocument string

func init() { RegisterType(new(XmlDocument)) }

func (v XmlDocument) Type() Type { return XmlDocumentMarker }

func (v *XmlDocument) Read(r io.Reader) error {
	value, err := readString[uint32](r)
	*v = XmlDocument(value)
	return err
}

func (v XmlDocument) Write(w io.Writer) error {
	return writeString[uint32](w, string(v))
}
