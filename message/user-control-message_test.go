package message_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tachode/rtmp-go/message"
)

func TestUserControlMessage_MarshalUnmarshal(t *testing.T) {
	tests := []struct {
		name       string
		event      message.UserControlMessageEvent
		parameters []uint32
	}{
		{
			name:       "StreamBegin with no parameters",
			event:      message.UserControlStreamBegin,
			parameters: nil,
		},
		{
			name:       "PingRequest with parameters",
			event:      message.UserControlPingRequest,
			parameters: []uint32{12345, 67890},
		},
		{
			name:       "StreamEOF with one parameter",
			event:      message.UserControlStreamEOF,
			parameters: []uint32{42},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := message.UserControlMessage{
				Event:      tt.event,
				Parameters: tt.parameters,
			}

			// Marshal the message
			data, err := msg.Marshal()
			assert.NoError(t, err)

			// Unmarshal the message
			var unmarshaledMsg message.UserControlMessage
			err = unmarshaledMsg.Unmarshal(data)
			assert.NoError(t, err)

			// Assert the unmarshaled message matches the original
			assert.Equal(t, msg.Event, unmarshaledMsg.Event)
			assert.Equal(t, msg.Parameters, unmarshaledMsg.Parameters)
		})
	}
}

func TestUserControlMessage_Unmarshal_InvalidData(t *testing.T) {
	tests := []struct {
		name string
		data []byte
	}{
		{
			name: "Empty data",
			data: []byte{},
		},
		{
			name: "Data too short",
			data: []byte{0x00},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var msg message.UserControlMessage
			err := msg.Unmarshal(tt.data)
			assert.Error(t, err)
		})
	}
}
