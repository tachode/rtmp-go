package message_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tachode/rtmp-go/amf0"
	"github.com/tachode/rtmp-go/message"
)

func TestAmf0CommandMessage_MarshalUnmarshal(t *testing.T) {
	original := message.Amf0CommandMessage{
		Command:       "testCommand",
		TransactionId: 1.0,
		Object:        amf0.Object{"key": amf0.String("value")},
		Parameters:    []any{amf0.String("param1"), amf0.Number(42), amf0.Boolean(true)},
	}

	// Marshal the message
	data, err := original.Marshal()
	assert.NoError(t, err)
	assert.NotNil(t, data)

	// Unmarshal the message
	var unmarshaled message.Amf0CommandMessage
	err = unmarshaled.Unmarshal(data)
	assert.NoError(t, err)

	// Assert equality
	assert.Equal(t, original.Command, unmarshaled.Command)
	assert.Equal(t, original.TransactionId, unmarshaled.TransactionId)
	assert.Equal(t, original.Object, unmarshaled.Object)
	assert.Equal(t, original.Parameters, unmarshaled.Parameters)
}

func TestAmf0CommandMessage_UnmarshalInvalidData(t *testing.T) {
	invalidData := []byte{0x01, 0x02, 0x03}

	var msg message.Amf0CommandMessage
	err := msg.Unmarshal(invalidData)
	assert.Error(t, err)
}

func TestAmf0CommandMessage_MarshalEmptyObject(t *testing.T) {
	msg := message.Amf0CommandMessage{
		Command:       "emptyObjectTest",
		TransactionId: 2.0,
		Object:        nil,
		Parameters:    nil,
	}

	data, err := msg.Marshal()
	assert.NoError(t, err)
	assert.NotNil(t, data)

	var unmarshaled message.Amf0CommandMessage
	err = unmarshaled.Unmarshal(data)
	assert.NoError(t, err)

	assert.Equal(t, msg.Command, unmarshaled.Command)
	assert.Equal(t, msg.TransactionId, unmarshaled.TransactionId)
	assert.Nil(t, unmarshaled.Object)
	assert.Nil(t, unmarshaled.Parameters)
}

func TestAmf0CommandMessage_String(t *testing.T) {
	msg := message.Amf0CommandMessage{
		Command:       "stringTest",
		TransactionId: 3.0,
		Object:        amf0.Object{"key": "value"},
		Parameters:    []any{"param1", 42},
	}

	str := msg.String()
	assert.Contains(t, str, "stringTest")
	assert.Contains(t, str, "tid=3")
	assert.Contains(t, str, "key:value")
	assert.Contains(t, str, "param1")
	assert.Contains(t, str, "42")
}
