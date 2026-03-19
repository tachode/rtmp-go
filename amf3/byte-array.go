package amf3

import (
	"io"
)

// ByteArray represents the AMF 3 ByteArray type (§3.14).
// An array of raw bytes. ByteArray instances can be sent as a reference
// to a previously occurring instance.
type ByteArray []byte

func init() { RegisterType(new(ByteArray)) }

func (v *ByteArray) isObjectRefType() {}

func (v ByteArray) Type() Type { return ByteArrayMarker }

func (v *ByteArray) Read(r *Reader) error {
	ref, header, isRef, err := r.readObjectRef()
	if err != nil {
		return err
	}
	if isRef {
		if b, ok := ref.(ByteArray); ok {
			*v = b
			return nil
		}
		return &UnexpectedRefTypeError{"ByteArray", ref}
	}

	// header = byte-length of the ByteArray
	length := header

	// Add to object reference table before reading payload
	r.addObjectRef(v)

	*v = make([]byte, length)
	_, err = io.ReadFull(r.r, *v)
	if err != nil {
		return err
	}

	// Update the reference now that we have the value
	r.objects[len(r.objects)-1] = *v

	return nil
}

func (v ByteArray) Write(w *Writer) error {
	length := uint32(len(v))
	// U29B-value: low bit = 1, remaining bits = byte-length
	err := writeU29(w.w, (length<<1)|1)
	if err != nil {
		return err
	}
	_, err = w.Write(v)
	return err
}
