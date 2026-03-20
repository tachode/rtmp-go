package message

import (
	"bytes"
	"fmt"

	"github.com/tachode/rtmp-go/amf0"
)

// Amf3SharedObjectMessage represents an AMF3 Shared Object Message (type 16).
// The wire format is the same as the AMF0 shared object message, but prefixed
// with a format selector byte (0x00). Event data containing AMF values uses
// AMF3 encoding instead of AMF0.
// See rtmp_specification_1.0.pdf §7.1.3.
type Amf3SharedObjectMessage struct {
	Amf0SharedObjectMessage
}

func init() { RegisterType(new(Amf3SharedObjectMessage)) }

func (m Amf3SharedObjectMessage) Type() Type {
	return TypeAmf3SharedObjectMessage
}

func (m Amf3SharedObjectMessage) Marshal() ([]byte, error) {
	inner, err := m.Amf0SharedObjectMessage.Marshal()
	if err != nil {
		return nil, err
	}
	return append([]byte{0x00}, inner...), nil
}

func (m *Amf3SharedObjectMessage) Unmarshal(data []byte) error {
	if len(data) < 1 {
		return fmt.Errorf("amf3 shared object message: %w", ErrShortMessage)
	}
	// Skip format selector byte, delegate to AMF0 implementation
	return m.Amf0SharedObjectMessage.Unmarshal(data[1:])
}

func (m Amf3SharedObjectMessage) String() string {
	return fmt.Sprintf("%v: %+v Name=%q Version=%d Flags=%x Events=%v",
		m.Type(), m.MetadataFields, m.Name, m.CurrentVersion, m.Flags, m.Events)
}

// DecodeEvent decodes the event at index i, returning its type and the
// name-value pair. The name is an AMF0 string body; the value may be AMF0 or
// AMF3 encoded (detected automatically).
func (m Amf3SharedObjectMessage) DecodeEvent(i int) (eventType SharedObjectEventType, name string, value any, err error) {
	eventType = m.Events[i].Type
	buf := bytes.NewBuffer(m.Events[i].Data)

	var nameStr amf0.String
	if err := nameStr.Read(buf); err != nil {
		return eventType, "", nil, fmt.Errorf("could not read shared object event name: %w", err)
	}
	name = string(nameStr)

	if buf.Len() > 0 {
		r := m.amf3Reader()
		r.SetReader(buf)
		defer r.SetReader(nil)

		value, err = readAmf0OrAmf3(buf, r)
		if err != nil {
			return eventType, name, nil, fmt.Errorf("could not read shared object event value: %w", err)
		}
	}
	return eventType, name, value, nil
}

// AddEvent encodes a name-value pair as event data using AMF3 and appends it
// to the Events slice.
func (m *Amf3SharedObjectMessage) AddEvent(eventType SharedObjectEventType, name string, value any) error {
	out := &bytes.Buffer{}

	if err := amf0.String(name).Write(out); err != nil {
		return fmt.Errorf("could not write shared object event name: %w", err)
	}

	if value != nil {
		w := m.amf3Writer()
		w.SetWriter(out)
		defer w.SetWriter(nil)

		out.WriteByte(byte(amf0.AvmplusObjectMarker))
		if err := w.WriteValue(value); err != nil {
			return fmt.Errorf("could not write shared object event value: %w", err)
		}
	}
	m.Events = append(m.Events, SharedObjectEvent{Type: eventType, Data: out.Bytes()})
	return nil
}
