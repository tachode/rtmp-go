package message_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tachode/rtmp-go/message"
)

func TestSetChunkSize_Type(t *testing.T) {
	msg := message.SetChunkSize{}
	assert.Equal(t, message.TypeSetChunkSize, msg.Type())
}

func TestSetChunkSize_Marshal(t *testing.T) {
	msg := message.SetChunkSize{ChunkSize: 4096}
	data, err := msg.Marshal()
	assert.NoError(t, err)
	assert.Equal(t, []byte{0x00, 0x00, 0x10, 0x00}, data)
}

func TestSetChunkSize_Unmarshal(t *testing.T) {
	data := []byte{0x00, 0x00, 0x10, 0x00}
	var msg message.SetChunkSize
	err := msg.Unmarshal(data)
	assert.NoError(t, err)
	assert.Equal(t, uint32(4096), msg.ChunkSize)
}

func TestSetChunkSize_Unmarshal_ShortMessage(t *testing.T) {
	data := []byte{0x00, 0x00, 0x10}
	var msg message.SetChunkSize
	err := msg.Unmarshal(data)
	assert.Error(t, err)
	assert.ErrorIs(t, err, message.ErrShortMessage)
}
