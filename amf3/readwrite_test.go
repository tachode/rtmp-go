package amf3_test

import (
	"bytes"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/tachode/rtmp-go/amf3"
)

func TestWriteValue_GoTypeConversions(t *testing.T) {
	tests := []struct {
		name       string
		value      any
		expectType amf3.Type
	}{
		{"float32", float32(1.5), amf3.DoubleMarker},
		{"float64", float64(1.5), amf3.DoubleMarker},
		{"int small", 42, amf3.IntegerMarker},
		{"int large", 0x20000000, amf3.DoubleMarker}, // >= 2^29, overflows U29
		{"int8", int8(1), amf3.IntegerMarker},
		{"int16", int16(1), amf3.IntegerMarker},
		{"int32", int32(1), amf3.IntegerMarker},
		{"int64", int64(1), amf3.IntegerMarker},
		{"int64 large", int64(0x10000000), amf3.DoubleMarker}, // >= 2^28 signed
		{"uint small", uint(42), amf3.IntegerMarker},
		{"uint8", uint8(1), amf3.IntegerMarker},
		{"uint16", uint16(1), amf3.IntegerMarker},
		{"uint32 small", uint32(1), amf3.IntegerMarker},
		{"uint32 large", uint32(0x20000000), amf3.DoubleMarker},
		{"uint64", uint64(1), amf3.IntegerMarker},
		{"bool true", true, amf3.TrueMarker},
		{"bool false", false, amf3.FalseMarker},
		{"string", "hello", amf3.StringMarker},
		{"time.Time", time.Now(), amf3.DateMarker},
		{"[]byte", []byte{1, 2, 3}, amf3.ByteArrayMarker},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			w := amf3.NewWriter(&buf)
			err := w.WriteValue(tt.value)
			assert.NoError(t, err)
			assert.Equal(t, byte(tt.expectType), buf.Bytes()[0],
				"expected marker 0x%02x, got 0x%02x", byte(tt.expectType), buf.Bytes()[0])
		})
	}
}

func TestWriteValue_Nil(t *testing.T) {
	var buf bytes.Buffer
	w := amf3.NewWriter(&buf)
	err := w.WriteValue(nil)
	assert.NoError(t, err)
	assert.Equal(t, []byte{byte(amf3.NullMarker)}, buf.Bytes())
}

func TestWriteValue_UnsupportedType(t *testing.T) {
	var buf bytes.Buffer
	w := amf3.NewWriter(&buf)
	err := w.WriteValue(struct{}{})
	assert.Error(t, err)
}

// Named types for testing that WriteValue handles types with underlying
// string/numeric kinds (e.g., type Level string, type EnumName int).
type namedString string
type namedInt int
type namedInt64 int64
type namedUint uint32
type namedFloat float64

func TestWriteValue_NamedTypes(t *testing.T) {
	tests := []struct {
		name       string
		value      any
		expectType amf3.Type
	}{
		{"named string", namedString("hello"), amf3.StringMarker},
		{"named int small", namedInt(42), amf3.IntegerMarker},
		{"named int large", namedInt(0x20000000), amf3.DoubleMarker},
		{"named int64 small", namedInt64(1), amf3.IntegerMarker},
		{"named int64 large", namedInt64(0x10000000), amf3.DoubleMarker},
		{"named uint small", namedUint(42), amf3.IntegerMarker},
		{"named uint large", namedUint(0x20000000), amf3.DoubleMarker},
		{"named float", namedFloat(3.14), amf3.DoubleMarker},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			w := amf3.NewWriter(&buf)
			err := w.WriteValue(tt.value)
			assert.NoError(t, err)
			assert.Equal(t, byte(tt.expectType), buf.Bytes()[0],
				"expected marker 0x%02x, got 0x%02x", byte(tt.expectType), buf.Bytes()[0])
		})
	}
}

func TestReadValue_UnknownMarker(t *testing.T) {
	data := []byte{0xFF}
	r := amf3.NewReader(bytes.NewReader(data))
	_, err := r.ReadValue()
	assert.Error(t, err)
}
