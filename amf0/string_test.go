package amf0_test

import (
	"bytes"
	"testing"

	"github.com/tachode/rtmp-go/amf0"
)

func TestString_Read(t *testing.T) {
	var ls amf0.String
	data := []byte{0x00, 0x05, 'H', 'e', 'l', 'l', 'o'}
	err := ls.Read(bytes.NewReader(data))
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	expected := amf0.String("Hello")
	if ls != expected {
		t.Fatalf("expected %v, got %v", expected, ls)
	}
}

func TestString_Write(t *testing.T) {
	ls := amf0.String("Hello")
	var buf bytes.Buffer
	err := ls.Write(&buf)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	expected := []byte{0x00, 0x05, 'H', 'e', 'l', 'l', 'o'}
	if !bytes.Equal(buf.Bytes(), expected) {
		t.Fatalf("expected %v, got %v", expected, buf.Bytes())
	}
}
