package amf0_test

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tachode/rtmp-go/amf0"
)

func TestEcmaArray_Read(t *testing.T) {
	data := []byte{
		0x00, 0x00, 0x00, 0x02, // length
		0x00, 0x03, 'k', 'e', 'y', // key "key"
		0x02, 0x00, 0x03, 'v', 'a', 'l', // value "val"
		0x00, 0x03, 'k', 'e', 'y', // key "key"
		0x02, 0x00, 0x03, 'v', 'a', 'l', // value "val"
	}
	r := bytes.NewReader(data)
	var arr amf0.EcmaArray
	err := arr.Read(r)
	assert.NoError(t, err)
	assert.Equal(t, amf0.EcmaArray{"key": amf0.String("val")}, arr)
}

func TestEcmaArray_Write(t *testing.T) {
	arr := amf0.EcmaArray{"key": "val"}
	var buf bytes.Buffer
	err := arr.Write(&buf)
	assert.NoError(t, err)
	expected := []byte{
		0x00, 0x00, 0x00, 0x01, // length
		0x00, 0x03, 'k', 'e', 'y', // key "key"
		0x02, 0x00, 0x03, 'v', 'a', 'l', // value "val"
	}
	assert.Equal(t, expected, buf.Bytes())
}
