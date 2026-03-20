package message_test

import (
	"encoding/binary"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tachode/rtmp-go/amf0"
	"github.com/tachode/rtmp-go/message"
)

func TestAmf0SharedObjectMessage_RoundTrip(t *testing.T) {
	original := message.Amf0SharedObjectMessage{
		Name:           "mySharedObj",
		CurrentVersion: 5,
		Flags:          [8]byte{},
		Events: []message.SharedObjectEvent{
			{Type: message.SharedObjectUse, Data: []byte{}},
			{Type: message.SharedObjectChange, Data: []byte{0x01, 0x02, 0x03}},
		},
	}

	data, err := original.Marshal()
	assert.NoError(t, err)

	var unmarshaled message.Amf0SharedObjectMessage
	err = unmarshaled.Unmarshal(data)
	assert.NoError(t, err)

	assert.Equal(t, original.Name, unmarshaled.Name)
	assert.Equal(t, original.CurrentVersion, unmarshaled.CurrentVersion)
	assert.Equal(t, original.Flags, unmarshaled.Flags)
	assert.Equal(t, len(original.Events), len(unmarshaled.Events))
	for i, e := range original.Events {
		assert.Equal(t, e.Type, unmarshaled.Events[i].Type)
		assert.Equal(t, e.Data, unmarshaled.Events[i].Data)
	}
}

func TestAmf0SharedObjectMessage_MarshalFormat(t *testing.T) {
	msg := message.Amf0SharedObjectMessage{
		Name:           "test",
		CurrentVersion: 1,
		Flags:          [8]byte{},
		Events: []message.SharedObjectEvent{
			{Type: message.SharedObjectClear, Data: []byte{}},
		},
	}

	data, err := msg.Marshal()
	assert.NoError(t, err)

	// Verify wire format manually:
	// uint16(4) + "test" + uint32(1) + [8]byte{} + uint8(8) + uint32(0)
	expected := []byte{
		0x00, 0x04, 't', 'e', 's', 't', // name
		0x00, 0x00, 0x00, 0x01, // version
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // flags (8 bytes)
		0x08,                   // event type: Clear
		0x00, 0x00, 0x00, 0x00, // event data length: 0
	}
	assert.Equal(t, expected, data)
}

func TestAmf0SharedObjectMessage_UnmarshalShortData(t *testing.T) {
	tests := []struct {
		name string
		data []byte
	}{
		{"empty", []byte{}},
		{"short name length", []byte{0x00}},
		{"short name", []byte{0x00, 0x05, 't'}},
		{"short version", []byte{0x00, 0x01, 't', 0x00, 0x00}},
		{"short flags", []byte{0x00, 0x01, 't', 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00}},
		{"short event header", []byte{0x00, 0x01, 't', 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01}},
		{"short event data", []byte{0x00, 0x01, 't', 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x05, 0xAA}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var msg message.Amf0SharedObjectMessage
			err := msg.Unmarshal(tt.data)
			assert.Error(t, err)
		})
	}
}

func TestAmf0SharedObjectMessage_NoEvents(t *testing.T) {
	msg := message.Amf0SharedObjectMessage{
		Name:           "empty",
		CurrentVersion: 0,
		Flags:          [8]byte{},
		Events:         nil,
	}

	data, err := msg.Marshal()
	assert.NoError(t, err)

	var unmarshaled message.Amf0SharedObjectMessage
	err = unmarshaled.Unmarshal(data)
	assert.NoError(t, err)

	assert.Equal(t, "empty", unmarshaled.Name)
	assert.Equal(t, uint32(0), unmarshaled.CurrentVersion)
	assert.Nil(t, unmarshaled.Events)
}

func TestAmf0SharedObjectMessage_MultipleEvents(t *testing.T) {
	msg := message.Amf0SharedObjectMessage{
		Name:           "multi",
		CurrentVersion: 3,
		Flags:          [8]byte{0x01},
		Events: []message.SharedObjectEvent{
			{Type: message.SharedObjectUse, Data: []byte{}},
			{Type: message.SharedObjectRequestChange, Data: []byte{0x01, 0x02}},
			{Type: message.SharedObjectClear, Data: []byte{}},
		},
	}

	data, err := msg.Marshal()
	assert.NoError(t, err)

	var unmarshaled message.Amf0SharedObjectMessage
	err = unmarshaled.Unmarshal(data)
	assert.NoError(t, err)

	assert.Equal(t, 3, len(unmarshaled.Events))
	assert.Equal(t, message.SharedObjectUse, unmarshaled.Events[0].Type)
	assert.Equal(t, message.SharedObjectRequestChange, unmarshaled.Events[1].Type)
	assert.Equal(t, []byte{0x01, 0x02}, unmarshaled.Events[1].Data)
	assert.Equal(t, message.SharedObjectClear, unmarshaled.Events[2].Type)
}

func TestAmf0SharedObjectMessage_Type(t *testing.T) {
	msg := message.Amf0SharedObjectMessage{}
	assert.Equal(t, message.TypeAmf0SharedObjectMessage, msg.Type())
}

func TestAmf0SharedObjectMessage_String(t *testing.T) {
	msg := message.Amf0SharedObjectMessage{
		Name:           "test",
		CurrentVersion: 1,
		Events: []message.SharedObjectEvent{
			{Type: message.SharedObjectUse, Data: []byte{}},
		},
	}
	s := msg.String()
	assert.Contains(t, s, "test")
	assert.Contains(t, s, "Version=1")
}

func TestAmf0SharedObjectMessage_DecodeEvent(t *testing.T) {
	msg := message.Amf0SharedObjectMessage{}
	err := msg.AddEvent(message.SharedObjectChange, "key", amf0.String("val"))
	assert.NoError(t, err)

	eventType, name, value, err := msg.DecodeEvent(0)
	assert.NoError(t, err)
	assert.Equal(t, message.SharedObjectChange, eventType)
	assert.Equal(t, "key", name)
	assert.Equal(t, amf0.String("val"), value)
}

func TestAmf0SharedObjectMessage_AddEvent_NameOnly(t *testing.T) {
	msg := message.Amf0SharedObjectMessage{}
	err := msg.AddEvent(message.SharedObjectRemove, "slot", nil)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(msg.Events))
	assert.Equal(t, message.SharedObjectRemove, msg.Events[0].Type)
	// Data should be: uint16(4) + "slot"
	assert.Equal(t, 2+4, len(msg.Events[0].Data))
	assert.Equal(t, uint16(4), binary.BigEndian.Uint16(msg.Events[0].Data[0:2]))
	assert.Equal(t, "slot", string(msg.Events[0].Data[2:]))
}

func TestAmf0SharedObjectMessage_Unmarshal_ViaRegistry(t *testing.T) {
	// Build raw payload for a shared object message
	msg := message.Amf0SharedObjectMessage{
		Name:           "so1",
		CurrentVersion: 2,
		Flags:          [8]byte{},
		Events: []message.SharedObjectEvent{
			{Type: message.SharedObjectUse, Data: []byte{}},
		},
	}

	payload, err := msg.Marshal()
	assert.NoError(t, err)

	// Unmarshal through the registry
	result, err := message.NewContext().Unmarshal(0, message.TypeAmf0SharedObjectMessage, 0, payload)
	assert.NoError(t, err)

	so, ok := result.(*message.Amf0SharedObjectMessage)
	assert.True(t, ok)
	assert.Equal(t, "so1", so.Name)
	assert.Equal(t, uint32(2), so.CurrentVersion)
	assert.Equal(t, 1, len(so.Events))
	assert.Equal(t, message.SharedObjectUse, so.Events[0].Type)
}
