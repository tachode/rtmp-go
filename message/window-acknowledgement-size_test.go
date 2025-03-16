package message_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tachode/rtmp-go/message"
)

func TestWindowAcknowledgementSize_Marshal(t *testing.T) {
	msg := message.WindowAcknowledgementSize{
		AcknowledgementWindowSize: 12345,
	}

	data, err := msg.Marshal()
	assert.NoError(t, err)
	assert.Equal(t, []byte{0x00, 0x00, 0x30, 0x39}, data)
}

func TestWindowAcknowledgementSize_Unmarshal(t *testing.T) {
	data := []byte{0x00, 0x00, 0x30, 0x39}
	var msg message.WindowAcknowledgementSize

	err := msg.Unmarshal(data)
	assert.NoError(t, err)
	assert.Equal(t, uint32(12345), msg.AcknowledgementWindowSize)
}

func TestWindowAcknowledgementSize_Unmarshal_ShortData(t *testing.T) {
	data := []byte{0x00, 0x00, 0x30}
	var msg message.WindowAcknowledgementSize

	err := msg.Unmarshal(data)
	assert.Error(t, err)
	assert.Equal(t, message.ErrShortMessage, err)
}
