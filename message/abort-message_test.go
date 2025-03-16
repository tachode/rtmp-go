package message_test

import (
	"bytes"
	"testing"

	"github.com/tachode/rtmp-go/message"
)

func TestAbortMessage_Marshal(t *testing.T) {
	msg := message.AbortMessage{
		ChunkStreamId: 12345,
	}

	data, err := msg.Marshal()
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	expected := []byte{0x00, 0x00, 0x30, 0x39} // 12345 in big-endian
	if !bytes.Equal(data, expected) {
		t.Errorf("Marshal output mismatch. Got %v, want %v", data, expected)
	}
}

func TestAbortMessage_Unmarshal(t *testing.T) {
	data := []byte{0x00, 0x00, 0x30, 0x39} // 12345 in big-endian

	var msg message.AbortMessage
	err := msg.Unmarshal(data)
	if err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if msg.ChunkStreamId != 12345 {
		t.Errorf("Unmarshal ChunkStreamId mismatch. Got %d, want %d", msg.ChunkStreamId, 12345)
	}
}

func TestAbortMessage_Unmarshal_ShortData(t *testing.T) {
	data := []byte{0x00, 0x00} // insufficient data

	var msg message.AbortMessage
	err := msg.Unmarshal(data)
	if err == nil {
		t.Fatal("Expected error for short data, got nil")
	}

	if err != message.ErrShortMessage {
		t.Errorf("Unexpected error. Got %v, want %v", err, message.ErrShortMessage)
	}
}
