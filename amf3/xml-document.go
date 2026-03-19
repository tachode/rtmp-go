package amf3

import (
	"io"
)

// XmlDocument represents the AMF 3 XMLDocument type (§3.9).
// This is the legacy XML type (flash.xml.XMLDocument).
// The XML content is flattened to a UTF-8 string.
// XMLDocuments can be sent as a reference to a previously occurring instance.
type XmlDocument string

func init() { RegisterType(new(XmlDocument)) }

func (v *XmlDocument) isObjectRefType() {}

func (v XmlDocument) Type() Type { return XmlDocumentMarker }

func (v *XmlDocument) Read(r *Reader) error {
	ref, header, isRef, err := r.readObjectRef()
	if err != nil {
		return err
	}
	if isRef {
		if x, ok := ref.(XmlDocument); ok {
			*v = x
			return nil
		}
		return &UnexpectedRefTypeError{"XmlDocument", ref}
	}

	// header = byte-length of UTF-8 encoded XML
	length := header

	// Add to object reference table before reading payload
	r.addObjectRef(v)

	buf := make([]byte, length)
	_, err = io.ReadFull(r.r, buf)
	if err != nil {
		return err
	}
	*v = XmlDocument(buf)

	// Update the reference now that we have the value
	r.objects[len(r.objects)-1] = *v

	return nil
}

func (v XmlDocument) Write(w *Writer) error {
	length := uint32(len(v))
	// U29X-value: low bit = 1, remaining bits = byte-length
	err := writeU29(w.w, (length<<1)|1)
	if err != nil {
		return err
	}
	_, err = w.Write([]byte(v))
	return err
}
