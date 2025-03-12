package amf0_test

import (
	"bytes"
	"testing"

	"github.com/tachode/rtmp-go/amf0"
)

func TestReference_Read(t *testing.T) {
	var ref amf0.Reference
	data := []byte{0x00, 0x01} // Example data
	buf := bytes.NewReader(data)

	err := ref.Read(buf)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if ref != 1 {
		t.Fatalf("Expected reference to be 1, got %v", ref)
	}
}

func TestReference_Write(t *testing.T) {
	ref := amf0.Reference(1)
	buf := new(bytes.Buffer)

	err := ref.Write(buf)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	expected := []byte{0x00, 0x01}
	if !bytes.Equal(buf.Bytes(), expected) {
		t.Fatalf("Expected %v, got %v", expected, buf.Bytes())
	}
}
