package amf0_test

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tachode/rtmp-go/amf0"
)

func TestStrictArray_Read(t *testing.T) {
	data := []byte{
		0x00, 0x00, 0x00, 0x02, // length: 2
		0x02, 0x00, 0x03, 'f', 'o', 'o', // string: "foo"
		0x02, 0x00, 0x03, 'b', 'a', 'r', // string: "bar"
	}
	r := bytes.NewReader(data)
	var arr amf0.StrictArray
	err := arr.Read(r)
	assert.NoError(t, err)
	assert.EqualValues(t, amf0.StrictArray{amf0.String("foo"), amf0.String("bar")}, arr)
}

func TestStrictArray_Write(t *testing.T) {
	arr := amf0.StrictArray{"foo", "bar"}
	var buf bytes.Buffer
	err := arr.Write(&buf)
	assert.NoError(t, err)
	expected := []byte{
		0x00, 0x00, 0x00, 0x02, // length: 2
		0x02, 0x00, 0x03, 'f', 'o', 'o', // string: "foo"
		0x02, 0x00, 0x03, 'b', 'a', 'r', // string: "bar"
	}
	assert.Equal(t, expected, buf.Bytes())
}
