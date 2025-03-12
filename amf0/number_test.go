package amf0_test

import (
	"bytes"
	"testing"

	"github.com/tachode/rtmp-go/amf0"
)

func TestNumberRead(t *testing.T) {
	var num amf0.Number
	data := []byte{0x40, 0x09, 0x21, 0xfb, 0x54, 0x44, 0x2d, 0x18} // 3.141592653589793
	buf := bytes.NewReader(data)

	err := num.Read(buf)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	expected := amf0.Number(3.141592653589793)
	if num != expected {
		t.Fatalf("Expected %v, got %v", expected, num)
	}
}

func TestNumberWrite(t *testing.T) {
	num := amf0.Number(3.141592653589793)
	buf := new(bytes.Buffer)

	err := num.Write(buf)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	expected := []byte{0x40, 0x09, 0x21, 0xfb, 0x54, 0x44, 0x2d, 0x18}
	if !bytes.Equal(buf.Bytes(), expected) {
		t.Fatalf("Expected %v, got %v", expected, buf.Bytes())
	}
}
