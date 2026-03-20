package message

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/tachode/rtmp-go/amf0"
)

// SharedObjectEventType represents the type of event in a shared object message.
// See rtmp_specification_1.0.pdf §7.1.3.
type SharedObjectEventType uint8

//go:generate stringer -type=SharedObjectEventType -trimprefix=SharedObject
const (
	SharedObjectUse           SharedObjectEventType = 1
	SharedObjectRelease       SharedObjectEventType = 2
	SharedObjectRequestChange SharedObjectEventType = 3
	SharedObjectChange        SharedObjectEventType = 4
	SharedObjectSuccess       SharedObjectEventType = 5
	SharedObjectSendMessage   SharedObjectEventType = 6
	SharedObjectStatus        SharedObjectEventType = 7
	SharedObjectClear         SharedObjectEventType = 8
	SharedObjectRemove        SharedObjectEventType = 9
	SharedObjectRequestRemove SharedObjectEventType = 10
	SharedObjectUseSuccess    SharedObjectEventType = 11
)

// SharedObjectEvent represents a single event within a shared object message.
type SharedObjectEvent struct {
	Type SharedObjectEventType
	Data []byte
}

// Amf0SharedObjectMessage represents an AMF0 Shared Object Message (type 19).
// A shared object is a collection of name-value pairs synchronized across
// multiple clients. Each message can contain multiple events.
// See rtmp_specification_1.0.pdf §7.1.3.
type Amf0SharedObjectMessage struct {
	MetadataFields
	Name           string
	CurrentVersion uint32
	Flags          [8]byte // Undocumented; semantics are not defined in the RTMP spec. Observed as all zeros in practice.
	Events         []SharedObjectEvent
}

func init() { RegisterType(new(Amf0SharedObjectMessage)) }

func (m Amf0SharedObjectMessage) Type() Type {
	return TypeAmf0SharedObjectMessage
}

func (m Amf0SharedObjectMessage) Marshal() ([]byte, error) {
	out := &bytes.Buffer{}

	// Shared Object Name (AMF0 string body: uint16 length + UTF-8)
	if err := amf0.String(m.Name).Write(out); err != nil {
		return nil, fmt.Errorf("could not write shared object name: %w", err)
	}

	// Current Version
	if err := binary.Write(out, binary.BigEndian, m.CurrentVersion); err != nil {
		return nil, fmt.Errorf("could not write shared object version: %w", err)
	}

	// Flags (8 bytes)
	out.Write(m.Flags[:])

	// Events
	for _, e := range m.Events {
		out.WriteByte(byte(e.Type))
		if err := binary.Write(out, binary.BigEndian, uint32(len(e.Data))); err != nil {
			return nil, fmt.Errorf("could not write shared object event length: %w", err)
		}
		out.Write(e.Data)
	}

	return out.Bytes(), nil
}

func (m *Amf0SharedObjectMessage) Unmarshal(data []byte) error {
	buf := bytes.NewBuffer(data)

	// Read shared object name (AMF0 string body: uint16 length + UTF-8)
	var name amf0.String
	if err := name.Read(buf); err != nil {
		return fmt.Errorf("shared object message name: %w", err)
	}
	m.Name = string(name)

	// Read current version (uint32)
	if err := binary.Read(buf, binary.BigEndian, &m.CurrentVersion); err != nil {
		return fmt.Errorf("shared object message version: %w", err)
	}

	// Read flags (8 bytes)
	if buf.Len() < 8 {
		return fmt.Errorf("shared object message flags: %w", ErrShortMessage)
	}
	copy(m.Flags[:], buf.Next(8))

	// Read events
	m.Events = nil
	for buf.Len() > 0 {
		if buf.Len() < 5 { // 1 byte type + 4 bytes length
			return fmt.Errorf("shared object event header: %w", ErrShortMessage)
		}
		eventType, _ := buf.ReadByte()
		var eventDataLen uint32
		binary.Read(buf, binary.BigEndian, &eventDataLen)

		if buf.Len() < int(eventDataLen) {
			return fmt.Errorf("shared object event data: %w", ErrShortMessage)
		}
		eventData := make([]byte, eventDataLen)
		copy(eventData, buf.Next(int(eventDataLen)))

		m.Events = append(m.Events, SharedObjectEvent{
			Type: SharedObjectEventType(eventType),
			Data: eventData,
		})
	}

	return nil
}

func (m Amf0SharedObjectMessage) String() string {
	return fmt.Sprintf("%v: %+v Name=%q Version=%d Flags=%x Events=%v",
		m.Type(), m.MetadataFields, m.Name, m.CurrentVersion, m.Flags, m.Events)
}

// Convenience: event String for debugging
func (e SharedObjectEvent) String() string {
	return fmt.Sprintf("{%v len=%d}", e.Type, len(e.Data))
}

// Helper methods to encode/decode event data for common event patterns.

// DecodeEvent decodes the event at index i, returning its type and the
// name-value pair (name as AMF0 string body, value as AMF0 value). Used by
// Change, RequestChange, Success events.
func (m Amf0SharedObjectMessage) DecodeEvent(i int) (eventType SharedObjectEventType, name string, value any, err error) {
	eventType = m.Events[i].Type
	buf := bytes.NewBuffer(m.Events[i].Data)

	var nameStr amf0.String
	if err := nameStr.Read(buf); err != nil {
		return eventType, "", nil, fmt.Errorf("could not read shared object event name: %w", err)
	}
	name = string(nameStr)

	if buf.Len() > 0 {
		value, err = amf0.Read(buf)
		if err != nil {
			return eventType, name, nil, fmt.Errorf("could not read shared object event value: %w", err)
		}
	}
	return eventType, name, value, nil
}

// AddEvent encodes a name-value pair as event data using AMF0 and appends it
// to the Events slice.
func (m *Amf0SharedObjectMessage) AddEvent(eventType SharedObjectEventType, name string, value any) error {
	out := &bytes.Buffer{}

	if err := amf0.String(name).Write(out); err != nil {
		return fmt.Errorf("could not write shared object event name: %w", err)
	}

	if value != nil {
		if err := amf0.Write(out, value); err != nil {
			return fmt.Errorf("could not write shared object event value: %w", err)
		}
	}
	m.Events = append(m.Events, SharedObjectEvent{Type: eventType, Data: out.Bytes()})
	return nil
}
