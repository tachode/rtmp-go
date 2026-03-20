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
	// Name as AMF0 string (uint16 length prefix + UTF-8)
	nameBytes := []byte(m.Name)
	size := 2 + len(nameBytes) + 4 + 8 // name-length + name + version + flags
	for _, e := range m.Events {
		size += 1 + 4 + len(e.Data) // type + data-length + data
	}

	out := make([]byte, 0, size)

	// Shared Object Name (AMF0-style uint16-length-prefixed string)
	out = binary.BigEndian.AppendUint16(out, uint16(len(nameBytes)))
	out = append(out, nameBytes...)

	// Current Version
	out = binary.BigEndian.AppendUint32(out, m.CurrentVersion)

	// Flags (8 bytes)
	out = append(out, m.Flags[:]...)

	// Events
	for _, e := range m.Events {
		out = append(out, byte(e.Type))
		out = binary.BigEndian.AppendUint32(out, uint32(len(e.Data)))
		out = append(out, e.Data...)
	}

	return out, nil
}

func (m *Amf0SharedObjectMessage) Unmarshal(data []byte) error {
	if len(data) < 2 {
		return fmt.Errorf("shared object message: %w", ErrShortMessage)
	}

	// Read shared object name (uint16-length-prefixed string)
	nameLen := int(binary.BigEndian.Uint16(data[0:2]))
	data = data[2:]
	if len(data) < nameLen {
		return fmt.Errorf("shared object message name: %w", ErrShortMessage)
	}
	m.Name = string(data[:nameLen])
	data = data[nameLen:]

	// Read current version (uint32)
	if len(data) < 4 {
		return fmt.Errorf("shared object message version: %w", ErrShortMessage)
	}
	m.CurrentVersion = binary.BigEndian.Uint32(data[0:4])
	data = data[4:]

	// Read flags (8 bytes)
	if len(data) < 8 {
		return fmt.Errorf("shared object message flags: %w", ErrShortMessage)
	}
	copy(m.Flags[:], data[0:8])
	data = data[8:]

	// Read events
	m.Events = nil
	for len(data) > 0 {
		if len(data) < 5 { // 1 byte type + 4 bytes length
			return fmt.Errorf("shared object event header: %w", ErrShortMessage)
		}
		eventType := SharedObjectEventType(data[0])
		eventDataLen := int(binary.BigEndian.Uint32(data[1:5]))
		data = data[5:]

		if len(data) < eventDataLen {
			return fmt.Errorf("shared object event data: %w", ErrShortMessage)
		}
		eventData := make([]byte, eventDataLen)
		copy(eventData, data[:eventDataLen])
		data = data[eventDataLen:]

		m.Events = append(m.Events, SharedObjectEvent{
			Type: eventType,
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

// DecodeEventValue decodes event data as a name-value pair (name as AMF0 string
// prefix, value as AMF0 value). Used by Change, RequestChange, Success events.
func (e SharedObjectEvent) DecodeEventValue() (name string, value any, err error) {
	if len(e.Data) < 2 {
		return "", nil, fmt.Errorf("shared object event data too short for name-value pair")
	}

	nameLen := int(binary.BigEndian.Uint16(e.Data[0:2]))
	if len(e.Data) < 2+nameLen {
		return "", nil, fmt.Errorf("shared object event data too short for name")
	}
	name = string(e.Data[2 : 2+nameLen])

	remaining := e.Data[2+nameLen:]
	if len(remaining) > 0 {
		value, err = amf0.Read(bytes.NewReader(remaining))
		if err != nil {
			return name, nil, fmt.Errorf("could not read shared object event value: %w", err)
		}
	}
	return name, value, nil
}

// EncodeEventValue encodes a name-value pair as event data.
func EncodeEventValue(name string, value any) ([]byte, error) {
	nameBytes := []byte(name)
	out := make([]byte, 0, 2+len(nameBytes)+16)
	out = binary.BigEndian.AppendUint16(out, uint16(len(nameBytes)))
	out = append(out, nameBytes...)

	if value != nil {
		var buf bytes.Buffer
		err := amf0.Write(&buf, value)
		if err != nil {
			return nil, fmt.Errorf("could not write shared object event value: %w", err)
		}
		out = append(out, buf.Bytes()...)
	}
	return out, nil
}
