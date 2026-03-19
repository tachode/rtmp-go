package amf3

import (
	"io"
)

// Xml represents the AMF 3 XML type (§3.13).
// This is the E4X XML type introduced in ActionScript 3.0.
// The XML content is flattened to a UTF-8 string.
// XML instances can be sent as a reference to a previously occurring instance.
type Xml string

func init() { RegisterType(new(Xml)) }

func (v *Xml) isObjectRefType() {}

func (v Xml) Type() Type { return XmlMarker }

func (v *Xml) Read(r *Reader) error {
	ref, header, isRef, err := r.readObjectRef()
	if err != nil {
		return err
	}
	if isRef {
		if x, ok := ref.(Xml); ok {
			*v = x
			return nil
		}
		return &UnexpectedRefTypeError{"Xml", ref}
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
	*v = Xml(buf)

	// Update the reference now that we have the value
	r.objects[len(r.objects)-1] = *v

	return nil
}

func (v Xml) Write(w *Writer) error {
	length := uint32(len(v))
	// U29X-value: low bit = 1, remaining bits = byte-length
	err := writeU29(w.w, (length<<1)|1)
	if err != nil {
		return err
	}
	_, err = w.Write([]byte(v))
	return err
}
