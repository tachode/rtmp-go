package message

import (
	"bytes"
	"errors"
	"fmt"
	"io"

	"github.com/tachode/rtmp-go/amf0"
)

// Amf3DataMessage represents an AMF3 Data Message (type 15).
// Per the RTMP spec, type 15 messages begin with a format selector byte (0x00),
// followed by AMF-encoded values. AMF3-encoded values are prefixed with the
// AMF0 avmplus-object-marker (0x11); the switch to AMF3 is not sticky.
type Amf3DataMessage struct {
	MetadataFields
	Handler    string
	Parameters []any
}

func init() { RegisterType(new(Amf3DataMessage)) }

func (m Amf3DataMessage) Type() Type {
	return TypeAmf3DataMessage
}

func (m Amf3DataMessage) Marshal() ([]byte, error) {
	out := bytes.NewBuffer(nil)

	// Format selector byte
	out.WriteByte(0x00)

	// Handler name is AMF0-encoded
	if err := amf0.Write(out, amf0.String(m.Handler)); err != nil {
		return nil, fmt.Errorf("could not write handler string: %w", err)
	}

	// Parameters are AMF3-encoded (each prefixed with 0x11)
	w := m.amf3Writer()
	w.SetWriter(out)
	defer w.SetWriter(nil)

	if m.Parameters != nil {
		for _, param := range m.Parameters {
			out.WriteByte(byte(amf0.AvmplusObjectMarker))
			if err := w.WriteValue(param); err != nil {
				return nil, fmt.Errorf("could not write parameter: %w", err)
			}
		}
	}
	return out.Bytes(), nil
}

func (m *Amf3DataMessage) Unmarshal(data []byte) error {
	if len(data) < 1 {
		return fmt.Errorf("amf3 data message: %w", ErrShortMessage)
	}

	// Skip format selector byte
	buf := bytes.NewBuffer(data[1:])

	r := m.amf3Reader()
	r.SetReader(buf)
	defer r.SetReader(nil)

	// Handler name
	var err error
	if m.Handler, err = readAmf0OrAmf3String(buf, r); err != nil {
		return fmt.Errorf("could not read handler string: %w", err)
	}

	m.Parameters = nil
	for buf.Len() > 0 {
		param, err := readAmf0OrAmf3(buf, r)
		if err != nil {
			if errors.Is(err, io.EOF) || errors.Is(err, io.ErrUnexpectedEOF) {
				break
			}
			return fmt.Errorf("could not read parameter: %w", err)
		}
		m.Parameters = append(m.Parameters, param)
	}

	return nil
}

func (m Amf3DataMessage) String() string {
	return fmt.Sprintf("%v: %+v Handler=%v(%+v)", m.Type(),
		m.MetadataFields, m.Handler, m.Parameters)
}

func (m Amf3DataMessage) GetHandler() string   { return m.Handler }
func (m Amf3DataMessage) GetParameters() []any { return m.Parameters }
