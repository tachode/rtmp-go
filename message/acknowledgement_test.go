package message_test

import (
	"bytes"
	"encoding/binary"
	"testing"

	"github.com/tachode/rtmp-go/message"
)

func TestAcknowledgement_Marshal(t *testing.T) {
	ack := message.Acknowledgement{
		SequenceNumber: 12345,
	}

	data, err := ack.Marshal()
	if err != nil {
		t.Fatalf("Marshal() failed: %v", err)
	}

	expected := make([]byte, 4)
	binary.BigEndian.PutUint32(expected, 12345)

	if !bytes.Equal(data, expected) {
		t.Errorf("Marshal() = %v, want %v", data, expected)
	}
}

func TestAcknowledgement_Unmarshal(t *testing.T) {
	data := make([]byte, 4)
	binary.BigEndian.PutUint32(data, 12345)

	var ack message.Acknowledgement
	err := ack.Unmarshal(data)
	if err != nil {
		t.Fatalf("Unmarshal() failed: %v", err)
	}

	if ack.SequenceNumber != 12345 {
		t.Errorf("Unmarshal() SequenceNumber = %d, want %d", ack.SequenceNumber, 12345)
	}
}

func TestAcknowledgement_Unmarshal_ShortData(t *testing.T) {
	data := []byte{1, 2, 3} // Less than 4 bytes

	var ack message.Acknowledgement
	err := ack.Unmarshal(data)
	if err == nil {
		t.Fatal("Unmarshal() expected error for short data, got nil")
	}
}
