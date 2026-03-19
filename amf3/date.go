package amf3

import (
	"encoding/binary"
	"time"
)

// Date represents the AMF 3 date type (§3.10).
// In AMF 3, a Date is the number of milliseconds elapsed since the epoch
// of midnight, 1st Jan 1970 in UTC. No timezone information is sent.
// Dates can be sent as a reference to a previously occurring Date instance.
type Date time.Time

func init() { RegisterType(new(Date)) }

func (v *Date) isObjectRefType() {}

func (v Date) Type() Type { return DateMarker }

func (v *Date) Read(r *Reader) error {
	ref, header, isRef, err := r.readObjectRef()
	if err != nil {
		return err
	}
	if isRef {
		if d, ok := ref.(Date); ok {
			*v = d
			return nil
		}
		return &UnexpectedRefTypeError{"Date", ref}
	}

	// header bits after ref flag are not used for Date (§3.10)
	_ = header

	// Add to object reference table before reading payload
	r.addObjectRef(v)

	var millis float64
	err = binary.Read(r.r, binary.BigEndian, &millis)
	if err != nil {
		return err
	}
	*v = Date(time.UnixMilli(int64(millis)))

	// Update the reference table entry now that we have the value
	r.objects[len(r.objects)-1] = *v

	return nil
}

func (v Date) Write(w *Writer) error {
	// Dates are not identity-comparable in Go, so we don't attempt
	// object deduplication for Date values. Always write inline.
	// U29D-value: low bit = 1, remaining bits unused (send 0x01)
	err := writeU29(w.w, 0x01)
	if err != nil {
		return err
	}
	return binary.Write(w.w, binary.BigEndian, float64(time.Time(v).UnixMilli()))
}
