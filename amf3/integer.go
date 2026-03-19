package amf3

import (
	"encoding/binary"
)

// Integer represents the AMF 3 integer type (§3.6).
// Integers are serialized using a variable length unsigned 29-bit integer (U29).
// Values are treated as signed 29-bit integers (-2^28 to 2^28-1).
// If a signed int >= 2^28 or unsigned uint >= 2^29, it should be sent as a Double.
type Integer int32

func init() { RegisterType(new(Integer)) }

func (v Integer) Type() Type { return IntegerMarker }

func (v *Integer) Read(r *Reader) error {
	val, err := readU29(r.r)
	if err != nil {
		return err
	}
	// Sign-extend from 29 bits to 32 bits
	if val&0x10000000 != 0 {
		*v = Integer(int32(val) | ^int32(0x1FFFFFFF))
	} else {
		*v = Integer(int32(val))
	}
	return nil
}

func (v Integer) Write(w *Writer) error {
	// Encode as unsigned 29-bit value; negative numbers use two's complement
	return writeU29(w.w, uint32(v)&0x1FFFFFFF)
}

// Double represents the AMF 3 double type (§3.7).
// Encoded as an 8-byte IEEE-754 double precision floating point value
// in network byte order.
type Double float64

func init() { RegisterType(new(Double)) }

func (v Double) Type() Type { return DoubleMarker }

func (v *Double) Read(r *Reader) error {
	return binary.Read(r.r, binary.BigEndian, v)
}

func (v Double) Write(w *Writer) error {
	return binary.Write(w.w, binary.BigEndian, v)
}
