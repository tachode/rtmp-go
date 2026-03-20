package message_test

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tachode/rtmp-go/amf0"
	"github.com/tachode/rtmp-go/amf3"
	"github.com/tachode/rtmp-go/message"
)

func TestAmf3CommandMessage_MarshalUnmarshal(t *testing.T) {
	original := message.Amf3CommandMessage{
		Command:       "testCommand",
		TransactionId: 1.0,
		Object:        nil,
		Parameters:    []any{amf3.String("param1"), amf3.Integer(42), amf3.Boolean(true)},
	}

	data, err := original.Marshal()
	assert.NoError(t, err)
	assert.NotNil(t, data)

	// First byte should be the format selector (0x00)
	assert.Equal(t, byte(0x00), data[0])

	var unmarshaled message.Amf3CommandMessage
	err = unmarshaled.Unmarshal(data)
	assert.NoError(t, err)

	assert.Equal(t, original.Command, unmarshaled.Command)
	assert.Equal(t, original.TransactionId, unmarshaled.TransactionId)
	assert.Equal(t, amf3.Null{}, unmarshaled.Object)
	assert.Equal(t, len(original.Parameters), len(unmarshaled.Parameters))
	assert.Equal(t, amf3.String("param1"), unmarshaled.Parameters[0])
	assert.Equal(t, amf3.Integer(42), unmarshaled.Parameters[1])
	assert.Equal(t, amf3.Boolean(true), unmarshaled.Parameters[2])
}

func TestAmf3CommandMessage_WithObject(t *testing.T) {
	obj := &amf3.Object{
		Traits: &amf3.TraitInfo{
			IsDynamic: true,
		},
		DynamicMembers: map[string]any{
			"key": amf3.String("value"),
		},
	}
	original := message.Amf3CommandMessage{
		Command:       "connect",
		TransactionId: 1.0,
		Object:        obj,
		Parameters:    nil,
	}

	data, err := original.Marshal()
	assert.NoError(t, err)

	var unmarshaled message.Amf3CommandMessage
	err = unmarshaled.Unmarshal(data)
	assert.NoError(t, err)

	assert.Equal(t, "connect", unmarshaled.Command)
	assert.Equal(t, 1.0, unmarshaled.TransactionId)
	assert.NotNil(t, unmarshaled.Object)
	result, ok := unmarshaled.Object.(*amf3.Object)
	assert.True(t, ok)
	assert.Equal(t, amf3.String("value"), result.DynamicMembers["key"])
}

func TestAmf3CommandMessage_EmptyParameters(t *testing.T) {
	original := message.Amf3CommandMessage{
		Command:       "ping",
		TransactionId: 0,
		Object:        nil,
		Parameters:    nil,
	}

	data, err := original.Marshal()
	assert.NoError(t, err)

	var unmarshaled message.Amf3CommandMessage
	err = unmarshaled.Unmarshal(data)
	assert.NoError(t, err)

	assert.Equal(t, "ping", unmarshaled.Command)
	assert.Equal(t, 0.0, unmarshaled.TransactionId)
	assert.Equal(t, amf3.Null{}, unmarshaled.Object)
	assert.Nil(t, unmarshaled.Parameters)
}

func TestAmf3CommandMessage_UnmarshalShortData(t *testing.T) {
	var msg message.Amf3CommandMessage
	err := msg.Unmarshal([]byte{})
	assert.Error(t, err)
}

func TestAmf3CommandMessage_Type(t *testing.T) {
	msg := message.Amf3CommandMessage{}
	assert.Equal(t, message.TypeAmf3CommandMessage, msg.Type())
}

func TestAmf3CommandMessage_String(t *testing.T) {
	msg := message.Amf3CommandMessage{
		Command:       "testCmd",
		TransactionId: 5.0,
		Parameters:    []any{amf3.String("hello")},
	}
	s := msg.String()
	assert.Contains(t, s, "testCmd")
	assert.Contains(t, s, "tid=5")
}

func TestAmf3CommandMessage_ViaRegistry(t *testing.T) {
	original := message.Amf3CommandMessage{
		Command:       "registryTest",
		TransactionId: 7.0,
		Object:        nil,
		Parameters:    []any{amf3.String("val")},
	}

	payload, err := original.Marshal()
	assert.NoError(t, err)

	result, err := message.NewContext().Unmarshal(0, message.TypeAmf3CommandMessage, 0, payload)
	assert.NoError(t, err)

	cmd, ok := result.(*message.Amf3CommandMessage)
	assert.True(t, ok)
	assert.Equal(t, "registryTest", cmd.Command)
	assert.Equal(t, 7.0, cmd.TransactionId)
	assert.Equal(t, 1, len(cmd.Parameters))
}

func TestAmf3CommandMessage_Amf0ObjectAndParams(t *testing.T) {
	// Build a payload where the command object and parameters are AMF0-encoded
	// (no 0x11 prefix), which is explicitly allowed by the spec.
	out := bytes.NewBuffer(nil)

	// Format selector
	out.WriteByte(0x00)

	// Command name (AMF0 string)
	amf0.Write(out, amf0.String("mixed"))
	// Transaction ID (AMF0 number)
	amf0.Write(out, amf0.Number(3.0))

	// Command object as AMF0 null (no 0x11 prefix)
	amf0.Write(out, amf0.Null{})

	// First parameter as AMF0 string (no 0x11 prefix)
	amf0.Write(out, amf0.String("amf0param"))

	// Second parameter as AMF3 string (with 0x11 prefix)
	out.WriteByte(0x11)
	w := amf3.NewWriter(out)
	w.WriteValue(amf3.String("amf3param"))

	var msg message.Amf3CommandMessage
	err := msg.Unmarshal(out.Bytes())
	assert.NoError(t, err)

	assert.Equal(t, "mixed", msg.Command)
	assert.Equal(t, 3.0, msg.TransactionId)
	assert.Equal(t, amf0.Null{}, msg.Object)
	assert.Equal(t, 2, len(msg.Parameters))
	assert.Equal(t, amf0.String("amf0param"), msg.Parameters[0])
	assert.Equal(t, amf3.String("amf3param"), msg.Parameters[1])
}

func TestAmf3CommandMessage_Amf3CommandAndTxId(t *testing.T) {
	// Build a payload where command name and transaction ID are AMF3-encoded.
	out := bytes.NewBuffer(nil)

	// Format selector
	out.WriteByte(0x00)

	// Command name as AMF3 string (with 0x11 prefix)
	out.WriteByte(0x11)
	w := amf3.NewWriter(out)
	w.WriteValue(amf3.String("amf3cmd"))

	// Transaction ID as AMF3 integer (with 0x11 prefix)
	out.WriteByte(0x11)
	w.WriteValue(amf3.Integer(5))

	// Command object as AMF3 null (with 0x11 prefix)
	out.WriteByte(0x11)
	w.WriteValue(amf3.Null{})

	var msg message.Amf3CommandMessage
	err := msg.Unmarshal(out.Bytes())
	assert.NoError(t, err)

	assert.Equal(t, "amf3cmd", msg.Command)
	assert.Equal(t, 5.0, msg.TransactionId)
	assert.Equal(t, amf3.Null{}, msg.Object)
	assert.Nil(t, msg.Parameters)
}
