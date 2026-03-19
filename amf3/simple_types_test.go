package amf3_test

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tachode/rtmp-go/amf3"
)

func TestUndefined_ReadWrite(t *testing.T) {
	// Read
	data := []byte{byte(amf3.UndefinedMarker)}
	r := amf3.NewReader(bytes.NewReader(data))
	val, err := r.ReadValue()
	assert.NoError(t, err)
	assert.Equal(t, amf3.Undefined{}, val)

	// Write
	var buf bytes.Buffer
	w := amf3.NewWriter(&buf)
	err = w.WriteValue(amf3.Undefined{})
	assert.NoError(t, err)
	assert.Equal(t, []byte{byte(amf3.UndefinedMarker)}, buf.Bytes())
}

func TestNull_ReadWrite(t *testing.T) {
	// Read
	data := []byte{byte(amf3.NullMarker)}
	r := amf3.NewReader(bytes.NewReader(data))
	val, err := r.ReadValue()
	assert.NoError(t, err)
	assert.Equal(t, amf3.Null{}, val)

	// Write nil
	var buf bytes.Buffer
	w := amf3.NewWriter(&buf)
	err = w.WriteValue(nil)
	assert.NoError(t, err)
	assert.Equal(t, []byte{byte(amf3.NullMarker)}, buf.Bytes())
}

func TestBoolean_ReadWrite(t *testing.T) {
	// Read false
	data := []byte{byte(amf3.FalseMarker)}
	r := amf3.NewReader(bytes.NewReader(data))
	val, err := r.ReadValue()
	assert.NoError(t, err)
	assert.Equal(t, amf3.Boolean(false), val)

	// Read true
	data = []byte{byte(amf3.TrueMarker)}
	r = amf3.NewReader(bytes.NewReader(data))
	val, err = r.ReadValue()
	assert.NoError(t, err)
	assert.Equal(t, amf3.Boolean(true), val)

	// Write false
	var buf bytes.Buffer
	w := amf3.NewWriter(&buf)
	err = w.WriteValue(amf3.Boolean(false))
	assert.NoError(t, err)
	assert.Equal(t, []byte{byte(amf3.FalseMarker)}, buf.Bytes())

	// Write true
	buf.Reset()
	w = amf3.NewWriter(&buf)
	err = w.WriteValue(amf3.Boolean(true))
	assert.NoError(t, err)
	assert.Equal(t, []byte{byte(amf3.TrueMarker)}, buf.Bytes())

	// Write Go bool via convenience
	buf.Reset()
	w = amf3.NewWriter(&buf)
	err = w.WriteValue(true)
	assert.NoError(t, err)
	assert.Equal(t, []byte{byte(amf3.TrueMarker)}, buf.Bytes())

	buf.Reset()
	w = amf3.NewWriter(&buf)
	err = w.WriteValue(false)
	assert.NoError(t, err)
	assert.Equal(t, []byte{byte(amf3.FalseMarker)}, buf.Bytes())
}
