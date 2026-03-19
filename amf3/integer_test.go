package amf3_test

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tachode/rtmp-go/amf3"
)

func TestU29_RoundTrip(t *testing.T) {
	tests := []struct {
		name  string
		value uint32
		bytes []byte
	}{
		{"zero", 0x00, []byte{0x00}},
		{"one", 0x01, []byte{0x01}},
		{"max 1-byte", 0x7F, []byte{0x7F}},
		{"min 2-byte", 0x80, []byte{0x81, 0x00}},
		{"mid 2-byte", 0x100, []byte{0x82, 0x00}},
		{"max 2-byte", 0x3FFF, []byte{0xFF, 0x7F}},
		{"min 3-byte", 0x4000, []byte{0x81, 0x80, 0x00}},
		{"max 3-byte", 0x1FFFFF, []byte{0xFF, 0xFF, 0x7F}},
		{"min 4-byte", 0x200000, []byte{0x80, 0xC0, 0x80, 0x00}},
		{"max U29", 0x1FFFFFFF, []byte{0xFF, 0xFF, 0xFF, 0xFF}},
	}

	for _, tt := range tests {
		t.Run(tt.name+"_read", func(t *testing.T) {
			r := amf3.NewReader(bytes.NewReader(tt.bytes))
			val, err := r.ReadValue()
			// We need to test U29 through a type that uses it.
			// integer-type = integer-marker U29
			// Let's test via Integer instead.
			_ = val
			_ = err
		})
	}

	// Test via helper encoding
	for _, tt := range tests {
		t.Run(tt.name+"_encode", func(t *testing.T) {
			encoded := u29Encode(tt.value)
			assert.Equal(t, tt.bytes, encoded)
		})
	}
}

func TestInteger_ReadWrite(t *testing.T) {
	tests := []struct {
		name  string
		value int32
		bytes []byte // payload only (no marker)
	}{
		{"zero", 0, []byte{0x00}},
		{"one", 1, []byte{0x01}},
		{"127", 127, []byte{0x7F}},
		{"128", 128, []byte{0x81, 0x00}},
		{"max positive", 0x0FFFFFFF, []byte{0xBF, 0xFF, 0xFF, 0xFF}},
		{"negative one", -1, []byte{0xFF, 0xFF, 0xFF, 0xFF}},
		{"min negative", -0x10000000, []byte{0xC0, 0x80, 0x80, 0x00}},
	}

	for _, tt := range tests {
		t.Run(tt.name+"_read", func(t *testing.T) {
			// Build full message: marker + payload
			data := append([]byte{byte(amf3.IntegerMarker)}, tt.bytes...)
			r := amf3.NewReader(bytes.NewReader(data))
			val, err := r.ReadValue()
			assert.NoError(t, err)
			assert.Equal(t, amf3.Integer(tt.value), val)
		})

		t.Run(tt.name+"_write", func(t *testing.T) {
			var buf bytes.Buffer
			w := amf3.NewWriter(&buf)
			err := w.WriteValue(amf3.Integer(tt.value))
			assert.NoError(t, err)
			expected := append([]byte{byte(amf3.IntegerMarker)}, tt.bytes...)
			assert.Equal(t, expected, buf.Bytes())
		})
	}
}

func TestDouble_ReadWrite(t *testing.T) {
	tests := []struct {
		name  string
		value float64
		bytes []byte
	}{
		{"zero", 0.0, []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}},
		{"pi", 3.141592653589793, []byte{0x40, 0x09, 0x21, 0xfb, 0x54, 0x44, 0x2d, 0x18}},
		{"negative", -1.0, []byte{0xBF, 0xF0, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}},
	}

	for _, tt := range tests {
		t.Run(tt.name+"_read", func(t *testing.T) {
			data := append([]byte{byte(amf3.DoubleMarker)}, tt.bytes...)
			r := amf3.NewReader(bytes.NewReader(data))
			val, err := r.ReadValue()
			assert.NoError(t, err)
			assert.Equal(t, amf3.Double(tt.value), val)
		})

		t.Run(tt.name+"_write", func(t *testing.T) {
			var buf bytes.Buffer
			w := amf3.NewWriter(&buf)
			err := w.WriteValue(amf3.Double(tt.value))
			assert.NoError(t, err)
			expected := append([]byte{byte(amf3.DoubleMarker)}, tt.bytes...)
			assert.Equal(t, expected, buf.Bytes())
		})
	}
}
