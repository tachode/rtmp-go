package amf0_test

import (
	"bytes"
	"testing"

	"github.com/tachode/rtmp-go/amf0"
)

func TestBooleanRead(t *testing.T) {
	var b amf0.Boolean
	data := []byte{1} // true
	err := b.Read(bytes.NewReader(data))
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if b != true {
		t.Fatalf("expected true, got %v", b)
	}

	data = []byte{0} // false
	err = b.Read(bytes.NewReader(data))
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if b != false {
		t.Fatalf("expected false, got %v", b)
	}
}

func TestBooleanWrite(t *testing.T) {
	var buf bytes.Buffer
	b := amf0.Boolean(true)
	err := b.Write(&buf)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !bytes.Equal(buf.Bytes(), []byte{1}) {
		t.Fatalf("expected [1], got %v", buf.Bytes())
	}

	buf.Reset()
	b = amf0.Boolean(false)
	err = b.Write(&buf)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !bytes.Equal(buf.Bytes(), []byte{0}) {
		t.Fatalf("expected [0], got %v", buf.Bytes())
	}
}
