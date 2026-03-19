package amf3_test

import (
	"bytes"
)

// u29Encode is a test helper to encode a U29 value for building test data.
func u29Encode(value uint32) []byte {
	var buf bytes.Buffer
	switch {
	case value < 0x80:
		buf.WriteByte(byte(value))
	case value < 0x4000:
		buf.WriteByte(byte(value>>7) | 0x80)
		buf.WriteByte(byte(value & 0x7F))
	case value < 0x200000:
		buf.WriteByte(byte(value>>14) | 0x80)
		buf.WriteByte(byte(value>>7) | 0x80)
		buf.WriteByte(byte(value & 0x7F))
	default:
		buf.WriteByte(byte(value>>22) | 0x80)
		buf.WriteByte(byte(value>>15) | 0x80)
		buf.WriteByte(byte(value>>8) | 0x80)
		buf.WriteByte(byte(value))
	}
	return buf.Bytes()
}

// utf8vr is a test helper to encode a UTF-8-vr string literal.
func utf8vr(s string) []byte {
	var buf bytes.Buffer
	length := uint32(len(s))
	buf.Write(u29Encode((length << 1) | 1))
	buf.WriteString(s)
	return buf.Bytes()
}

// utf8vrEmpty returns the encoding of an empty string.
func utf8vrEmpty() []byte {
	return []byte{0x01}
}

// utf8vrRef returns the encoding of a string reference.
func utf8vrRef(index uint32) []byte {
	return u29Encode(index << 1)
}
