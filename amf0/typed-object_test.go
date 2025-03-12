package amf0_test

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tachode/rtmp-go/amf0"
)

func TestTypedObject_Read(t *testing.T) {
	data := []byte{0x00, 0x05, 'H', 'e', 'l', 'l', 'o', 0x00, 0x03, 'k', 'e', 'y', 0x02, 0x00, 0x03, 'v', 'a', 'l', 0x00, 0x00, 0x09}
	r := bytes.NewReader(data)
	obj := &amf0.TypedObject{}

	err := obj.Read(r)
	assert.NoError(t, err)
	assert.Equal(t, "Hello", obj.ClassName)
	assert.Equal(t, amf0.String("val"), obj.Object["key"])
}

func TestTypedObject_Write(t *testing.T) {
	obj := amf0.TypedObject{
		ClassName: "Hello",
		Object:    amf0.Object{"key": "val"},
	}
	var buf bytes.Buffer

	err := obj.Write(&buf)
	assert.NoError(t, err)
	expected := []byte{0x00, 0x05, 'H', 'e', 'l', 'l', 'o', 0x00, 0x03, 'k', 'e', 'y', 0x02, 0x00, 0x03, 'v', 'a', 'l', 0x00, 0x00, 0x09}
	assert.Equal(t, expected, buf.Bytes())
}
