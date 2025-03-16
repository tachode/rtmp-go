package message_test

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tachode/rtmp-go/amf0"
	"github.com/tachode/rtmp-go/message"
)

func TestAmf0DataMessage_Marshal(t *testing.T) {
	msg := message.Amf0DataMessage{
		Handler: "testHandler",
		Parameters: []any{
			amf0.String("param1"),
			amf0.Number(42),
			amf0.Boolean(true),
		},
	}

	data, err := msg.Marshal()
	assert.NoError(t, err)

	expected := bytes.NewBuffer(nil)
	amf0.Write(expected, amf0.String("testHandler"))
	amf0.Write(expected, "param1")
	amf0.Write(expected, 42)
	amf0.Write(expected, true)

	assert.Equal(t, expected.Bytes(), data)
}

func TestAmf0DataMessage_Unmarshal(t *testing.T) {
	data := bytes.NewBuffer(nil)
	amf0.Write(data, amf0.String("testHandler"))
	amf0.Write(data, "param1")
	amf0.Write(data, 42)
	amf0.Write(data, true)

	var msg message.Amf0DataMessage
	err := msg.Unmarshal(data.Bytes())
	assert.NoError(t, err)

	assert.Equal(t, "testHandler", msg.Handler)
	assert.Equal(t, []any{amf0.String("param1"), amf0.Number(42), amf0.Boolean(true)}, msg.Parameters)
}

func TestAmf0DataMessage_UnmarshalWithExtraBytes(t *testing.T) {
	data := bytes.NewBuffer(nil)
	amf0.Write(data, amf0.String("testHandler"))
	amf0.Write(data, "param1")
	amf0.Write(data, 42)
	amf0.Write(data, true)
	data.Write([]byte{0x00, 0x01, 0x02}) // Extra bytes

	var msg message.Amf0DataMessage
	err := msg.Unmarshal(data.Bytes())
	assert.NoError(t, err)

	assert.Equal(t, "testHandler", msg.Handler)
	assert.Equal(t, []any{amf0.String("param1"), amf0.Number(42), amf0.Boolean(true)}, msg.Parameters)
}
