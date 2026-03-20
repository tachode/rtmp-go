package message

import (
	"bytes"
	"fmt"

	"github.com/tachode/rtmp-go/amf0"
	"github.com/tachode/rtmp-go/amf3"
)

// amf3Writer returns the AMF3 writer from the context, creating a temporary one if needed.
func (m *MetadataFields) amf3Writer() *amf3.Writer {
	if m.context != nil && m.context.amf3Writer != nil {
		return m.context.amf3Writer
	}
	return amf3.NewWriter(nil)
}

// amf3Reader returns the AMF3 reader from the context, creating a temporary one if needed.
func (m *MetadataFields) amf3Reader() *amf3.Reader {
	if m.context != nil && m.context.amf3Reader != nil {
		return m.context.amf3Reader
	}
	return amf3.NewReader(nil)
}

// readAmf0OrAmf3 reads a single value from buf. If the next byte is the AMF0
// avmplus-object-marker (0x11), the value is read as AMF3 using r. Otherwise,
// the byte is put back and the value is read as AMF0. This allows senders to
// freely choose AMF0 or AMF3 encoding for each value in AMF3 message types.
func readAmf0OrAmf3(buf *bytes.Buffer, r *amf3.Reader) (any, error) {
	marker, err := buf.ReadByte()
	if err != nil {
		return nil, err
	}
	if marker == byte(amf0.AvmplusObjectMarker) {
		return r.ReadValue()
	}
	buf.UnreadByte()
	return amf0.Read(buf)
}

// readAmf0OrAmf3String reads a value using readAmf0OrAmf3 and coerces it to a
// Go string. Accepts amf0.String and amf3.String.
func readAmf0OrAmf3String(buf *bytes.Buffer, r *amf3.Reader) (string, error) {
	v, err := readAmf0OrAmf3(buf, r)
	if err != nil {
		return "", err
	}
	switch s := v.(type) {
	case amf0.String:
		return string(s), nil
	case amf3.String:
		return string(s), nil
	default:
		return "", fmt.Errorf("expected string, got %T", v)
	}
}

// readAmf0OrAmf3Number reads a value using readAmf0OrAmf3 and coerces it to a
// Go float64. Accepts amf0.Number, amf3.Double, and amf3.Integer.
func readAmf0OrAmf3Number(buf *bytes.Buffer, r *amf3.Reader) (float64, error) {
	v, err := readAmf0OrAmf3(buf, r)
	if err != nil {
		return 0, err
	}
	switch n := v.(type) {
	case amf0.Number:
		return float64(n), nil
	case amf3.Double:
		return float64(n), nil
	case amf3.Integer:
		return float64(n), nil
	default:
		return 0, fmt.Errorf("expected number, got %T", v)
	}
}
