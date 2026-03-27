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

func TestSetChunkSize_Unmarshal_Zero(t *testing.T) {
	data := []byte{0x00, 0x00, 0x00, 0x00}
	var msg message.SetChunkSize
	err := msg.Unmarshal(data)
	assert.ErrorIs(t, err, message.ErrInvalidChunkSize)
}

func TestSetChunkSize_Unmarshal_TooLarge(t *testing.T) {
	// 0x01000000 = 16777216, exceeds max of 0xFFFFFF
	data := []byte{0x01, 0x00, 0x00, 0x00}
	var msg message.SetChunkSize
	err := msg.Unmarshal(data)
	assert.ErrorIs(t, err, message.ErrInvalidChunkSize)
}

func TestSetChunkSize_Unmarshal_MaxValid(t *testing.T) {
	// 0x00FFFFFF = 16777215, the maximum valid chunk size
	data := []byte{0x00, 0xFF, 0xFF, 0xFF}
	var msg message.SetChunkSize
	err := msg.Unmarshal(data)
	assert.NoError(t, err)
	assert.Equal(t, uint32(0xFFFFFF), msg.ChunkSize)
}

func TestSetChunkSize_Unmarshal_MinValid(t *testing.T) {
	data := []byte{0x00, 0x00, 0x00, 0x01}
	var msg message.SetChunkSize
	err := msg.Unmarshal(data)
	assert.NoError(t, err)
	assert.Equal(t, uint32(1), msg.ChunkSize)
}
