package message_test

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tachode/rtmp-go/amf0"
	"github.com/tachode/rtmp-go/amf3"
	"github.com/tachode/rtmp-go/message"
)

func TestAmf3DataMessage_MarshalUnmarshal(t *testing.T) {
	original := message.Amf3DataMessage{
		Handler:    "onMetaData",
		Parameters: []any{amf3.String("param1"), amf3.Integer(100)},
	}

	data, err := original.Marshal()
	assert.NoError(t, err)
	assert.NotNil(t, data)

	// First byte should be format selector (0x00)
	assert.Equal(t, byte(0x00), data[0])

	var unmarshaled message.Amf3DataMessage
	err = unmarshaled.Unmarshal(data)
	assert.NoError(t, err)

	assert.Equal(t, original.Handler, unmarshaled.Handler)
	assert.Equal(t, len(original.Parameters), len(unmarshaled.Parameters))
	assert.Equal(t, amf3.String("param1"), unmarshaled.Parameters[0])
	assert.Equal(t, amf3.Integer(100), unmarshaled.Parameters[1])
}

func TestAmf3DataMessage_NoParameters(t *testing.T) {
	original := message.Amf3DataMessage{
		Handler:    "onStatus",
		Parameters: nil,
	}

	data, err := original.Marshal()
	assert.NoError(t, err)

	var unmarshaled message.Amf3DataMessage
	err = unmarshaled.Unmarshal(data)
	assert.NoError(t, err)

	assert.Equal(t, "onStatus", unmarshaled.Handler)
	assert.Nil(t, unmarshaled.Parameters)
}

func TestAmf3DataMessage_UnmarshalShortData(t *testing.T) {
	var msg message.Amf3DataMessage
	err := msg.Unmarshal([]byte{})
	assert.Error(t, err)
}

func TestAmf3DataMessage_Type(t *testing.T) {
	msg := message.Amf3DataMessage{}
	assert.Equal(t, message.TypeAmf3DataMessage, msg.Type())
}

func TestAmf3DataMessage_String(t *testing.T) {
	msg := message.Amf3DataMessage{
		Handler:    "onMetaData",
		Parameters: []any{amf3.String("hello")},
	}
	s := msg.String()
	assert.Contains(t, s, "onMetaData")
}

func TestAmf3DataMessage_ViaRegistry(t *testing.T) {
	original := message.Amf3DataMessage{
		Handler:    "registryTest",
		Parameters: []any{amf3.String("val")},
	}

	payload, err := original.Marshal()
	assert.NoError(t, err)

	result, err := message.NewContext().Unmarshal(0, message.TypeAmf3DataMessage, 0, payload)
	assert.NoError(t, err)

	dm, ok := result.(*message.Amf3DataMessage)
	assert.True(t, ok)
	assert.Equal(t, "registryTest", dm.Handler)
	assert.Equal(t, 1, len(dm.Parameters))
}

func TestAmf3DataMessage_Amf0Parameters(t *testing.T) {
	// Build a payload where parameters are AMF0-encoded (no 0x11 prefix),
	// which is explicitly allowed by the spec.
	out := bytes.NewBuffer(nil)

	// Format selector
	out.WriteByte(0x00)

	// Handler name (AMF0 string)
	amf0.Write(out, amf0.String("onMetaData"))

	// First parameter as AMF0 number (no 0x11 prefix)
	amf0.Write(out, amf0.Number(42.0))

	// Second parameter as AMF3 string (with 0x11 prefix)
	out.WriteByte(0x11)
	w := amf3.NewWriter(out)
	w.WriteValue(amf3.String("amf3val"))

	var msg message.Amf3DataMessage
	err := msg.Unmarshal(out.Bytes())
	assert.NoError(t, err)

	assert.Equal(t, "onMetaData", msg.Handler)
	assert.Equal(t, 2, len(msg.Parameters))
	assert.Equal(t, amf0.Number(42.0), msg.Parameters[0])
	assert.Equal(t, amf3.String("amf3val"), msg.Parameters[1])
}

func TestAmf3DataMessage_Amf3Handler(t *testing.T) {
	// Build a payload where the handler name is AMF3-encoded.
	out := bytes.NewBuffer(nil)

	// Format selector
	out.WriteByte(0x00)

	// Handler name as AMF3 string (with 0x11 prefix)
	out.WriteByte(0x11)
	w := amf3.NewWriter(out)
	w.WriteValue(amf3.String("onMetaData"))

	// Parameter as AMF3 integer
	out.WriteByte(0x11)
	w.WriteValue(amf3.Integer(99))

	var msg message.Amf3DataMessage
	err := msg.Unmarshal(out.Bytes())
	assert.NoError(t, err)

	assert.Equal(t, "onMetaData", msg.Handler)
	assert.Equal(t, 1, len(msg.Parameters))
	assert.Equal(t, amf3.Integer(99), msg.Parameters[0])
}
