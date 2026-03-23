package message

import (
	"bytes"
	"errors"
	"fmt"
	"io"

	"github.com/tachode/rtmp-go/amf0"
)

// Amf3CommandMessage represents an AMF3 Command Message (type 17).
// Per the RTMP spec, type 17 messages begin with a format selector byte (0x00),
// followed by AMF-encoded values. AMF3-encoded values are prefixed with the
// AMF0 avmplus-object-marker (0x11); the switch to AMF3 is not sticky.
type Amf3CommandMessage struct {
	MetadataFields
	Command       string
	TransactionId float64
	Object        any // can be an amf0.Object or an amf3.Object
	Parameters    []any
}

func init() { RegisterType(new(Amf3CommandMessage)) }

func (m Amf3CommandMessage) Type() Type {
	return TypeAmf3CommandMessage
}

func (m Amf3CommandMessage) Marshal() ([]byte, error) {
	out := bytes.NewBuffer(nil)

	// Format selector byte
	out.WriteByte(0x00)

	// Command name and transaction ID are AMF0-encoded
	if err := amf0.Write(out, amf0.String(m.Command)); err != nil {
		return nil, fmt.Errorf("could not write command string: %w", err)
	}
	if err := amf0.Write(out, amf0.Number(m.TransactionId)); err != nil {
		return nil, fmt.Errorf("could not write transaction id: %w", err)
	}

	// Command object and parameters are AMF3-encoded (prefixed with 0x11)
	w := m.amf3Writer()
	w.SetWriter(out)
	defer w.SetWriter(nil)

	out.WriteByte(byte(amf0.AvmplusObjectMarker))
	if err := w.WriteValue(m.Object); err != nil {
		return nil, fmt.Errorf("could not write command object: %w", err)
	}

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

func (m *Amf3CommandMessage) Unmarshal(data []byte) error {
	if len(data) < 1 {
		return fmt.Errorf("amf3 command message: %w", ErrShortMessage)
	}

	// Skip format selector byte
	buf := bytes.NewBuffer(data[1:])

	r := m.amf3Reader()
	r.SetReader(buf)
	defer r.SetReader(nil)

	// Command name and transaction ID
	var err error
	if m.Command, err = readAmf0OrAmf3String(buf, r); err != nil {
		return fmt.Errorf("could not read command string: %w", err)
	}
	if m.TransactionId, err = readAmf0OrAmf3Number(buf, r); err != nil {
		return fmt.Errorf("could not read transaction id: %w", err)
	}

	// Read command object
	if buf.Len() == 0 {
		return nil
	}
	m.Object, err = readAmf0OrAmf3(buf, r)
	if err != nil {
		return fmt.Errorf("could not read command object: %w", err)
	}

	// Read parameters
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

func (m Amf3CommandMessage) String() string {
	return fmt.Sprintf("%v: %+v Command=%v(tid=%v, obj=%+v, %+v)", m.Type(),
		m.MetadataFields, m.Command, m.TransactionId, m.Object, m.Parameters)
}

func (m Amf3CommandMessage) GetCommand() string        { return m.Command }
func (m Amf3CommandMessage) GetTransactionId() float64 { return m.TransactionId }
func (m Amf3CommandMessage) GetParameters() []any      { return m.Parameters }

func (m Amf3CommandMessage) GetObject() Object {
	if obj, ok := m.Object.(Object); ok {
		return obj
	}
	return nil
}
