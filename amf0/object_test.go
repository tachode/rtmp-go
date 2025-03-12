package amf0_test

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tachode/rtmp-go/amf0"
)

func TestObject_Read(t *testing.T) {
	data := []byte{
		0x00, 0x03, 'f', 'o', 'o', // key "foo"
		0x02, 0x00, 0x03, 'b', 'a', 'r', // value "bar"
		0x00, 0x00, 0x09, // ObjectEnd marker
	}
	r := bytes.NewReader(data)
	var obj amf0.Object
	err := obj.Read(r)
	assert.NoError(t, err)
	assert.Equal(t, amf0.Object{"foo": amf0.String("bar")}, obj)
}

func TestObject_Write(t *testing.T) {
	obj := amf0.Object{"foo": "bar"}
	var buf bytes.Buffer
	err := obj.Write(&buf)
	assert.NoError(t, err)
	expected := []byte{
		0x00, 0x03, 'f', 'o', 'o', // key "foo"
		0x02, 0x00, 0x03, 'b', 'a', 'r', // value "bar"
		0x00, 0x00, 0x09, // ObjectEnd marker
	}
	assert.Equal(t, expected, buf.Bytes())
}

func TestObject_Read_Empty(t *testing.T) {
	data := []byte{
		0x00, 0x00, 0x09, // ObjectEnd marker
	}
	r := bytes.NewReader(data)
	var obj amf0.Object
	err := obj.Read(r)
	assert.NoError(t, err)
	assert.Equal(t, amf0.Object{}, obj)
}

func TestObject_Write_Empty(t *testing.T) {
	obj := amf0.Object{}
	var buf bytes.Buffer
	err := obj.Write(&buf)
	assert.NoError(t, err)
	expected := []byte{
		0x00, 0x00, 0x09, // ObjectEnd marker
	}
	assert.Equal(t, expected, buf.Bytes())
}
