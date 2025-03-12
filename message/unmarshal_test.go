package message_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tachode/rtmp-go/message"
)

type MockMessage struct {
	message.MetadataFields
	payload []byte
}

func (m *MockMessage) Type() message.Type {
	return 1
}

func (m *MockMessage) Unmarshal(payload []byte) error {
	m.payload = payload
	return nil
}

func (m *MockMessage) Marshal() ([]byte, error) {
	return m.payload, nil
}

func TestUnmarshal(t *testing.T) {
	message.RegisterType(new(MockMessage))

	t.Run("unknown message type", func(t *testing.T) {
		_, err := message.Unmarshal(123, 99, 456, []byte{1, 2, 3})
		assert.Error(t, err)
		assert.Equal(t, "unknown RTMP message Type(99)", err.Error())
	})

	t.Run("successful unmarshal", func(t *testing.T) {
		msg, err := message.Unmarshal(123, 1, 456, []byte{1, 2, 3})
		assert.NoError(t, err)
		assert.NotNil(t, msg)
		assert.Equal(t, message.Type(1), msg.Type())
		assert.Equal(t, uint32(3), msg.Metadata().Length)
		assert.Equal(t, uint32(123), msg.Metadata().Timestamp)
		assert.Equal(t, uint32(456), msg.Metadata().StreamId)
		assert.Equal(t, []byte{1, 2, 3}, msg.(*MockMessage).payload)
	})
}
