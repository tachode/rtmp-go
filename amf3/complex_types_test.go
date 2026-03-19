package amf3_test

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tachode/rtmp-go/amf3"
)

func TestXmlDocument_ReadWrite(t *testing.T) {
	t.Run("read xml document", func(t *testing.T) {
		xml := "<root><child/></root>"
		var data bytes.Buffer
		data.WriteByte(byte(amf3.XmlDocumentMarker))
		// U29X-value: (byte-length << 1) | 1
		data.Write(u29Encode((uint32(len(xml)) << 1) | 1))
		data.WriteString(xml)

		r := amf3.NewReader(&data)
		val, err := r.ReadValue()
		assert.NoError(t, err)
		x, ok := val.(*amf3.XmlDocument)
		assert.True(t, ok)
		assert.Equal(t, amf3.XmlDocument(xml), *x)
	})

	t.Run("write xml document", func(t *testing.T) {
		xml := amf3.XmlDocument("<root/>")
		var buf bytes.Buffer
		w := amf3.NewWriter(&buf)
		err := w.WriteValue(xml)
		assert.NoError(t, err)

		r := amf3.NewReader(&buf)
		val, err := r.ReadValue()
		assert.NoError(t, err)
		x2, ok := val.(*amf3.XmlDocument)
		assert.True(t, ok)
		assert.Equal(t, xml, *x2)
	})
}

func TestXml_ReadWrite(t *testing.T) {
	t.Run("read xml", func(t *testing.T) {
		xml := "<data attr=\"val\"/>"
		var data bytes.Buffer
		data.WriteByte(byte(amf3.XmlMarker))
		data.Write(u29Encode((uint32(len(xml)) << 1) | 1))
		data.WriteString(xml)

		r := amf3.NewReader(&data)
		val, err := r.ReadValue()
		assert.NoError(t, err)
		xv, ok := val.(*amf3.Xml)
		assert.True(t, ok)
		assert.Equal(t, amf3.Xml(xml), *xv)
	})

	t.Run("write xml", func(t *testing.T) {
		xml := amf3.Xml("<data/>")
		var buf bytes.Buffer
		w := amf3.NewWriter(&buf)
		err := w.WriteValue(xml)
		assert.NoError(t, err)

		r := amf3.NewReader(&buf)
		val, err := r.ReadValue()
		assert.NoError(t, err)
		xv2, ok := val.(*amf3.Xml)
		assert.True(t, ok)
		assert.Equal(t, xml, *xv2)
	})
}

func TestByteArray_ReadWrite(t *testing.T) {
	t.Run("read byte array", func(t *testing.T) {
		payload := []byte{0xDE, 0xAD, 0xBE, 0xEF}
		var data bytes.Buffer
		data.WriteByte(byte(amf3.ByteArrayMarker))
		// U29B-value: (byte-length << 1) | 1
		data.Write(u29Encode((uint32(len(payload)) << 1) | 1))
		data.Write(payload)

		r := amf3.NewReader(&data)
		val, err := r.ReadValue()
		assert.NoError(t, err)
		ba, ok := val.(*amf3.ByteArray)
		assert.True(t, ok)
		assert.Equal(t, amf3.ByteArray(payload), *ba)
	})

	t.Run("write byte array", func(t *testing.T) {
		ba := amf3.ByteArray([]byte{0x01, 0x02, 0x03})
		var buf bytes.Buffer
		w := amf3.NewWriter(&buf)
		err := w.WriteValue(ba)
		assert.NoError(t, err)

		r := amf3.NewReader(&buf)
		val, err := r.ReadValue()
		assert.NoError(t, err)
		ba2, ok := val.(*amf3.ByteArray)
		assert.True(t, ok)
		assert.Equal(t, ba, *ba2)
	})

	t.Run("write Go []byte via convenience", func(t *testing.T) {
		payload := []byte{0xAA, 0xBB}
		var buf bytes.Buffer
		w := amf3.NewWriter(&buf)
		err := w.WriteValue(payload)
		assert.NoError(t, err)
		assert.Equal(t, byte(amf3.ByteArrayMarker), buf.Bytes()[0])
	})

	t.Run("empty byte array", func(t *testing.T) {
		ba := amf3.ByteArray([]byte{})
		var buf bytes.Buffer
		w := amf3.NewWriter(&buf)
		err := w.WriteValue(ba)
		assert.NoError(t, err)

		r := amf3.NewReader(&buf)
		val, err := r.ReadValue()
		assert.NoError(t, err)
		ba3, ok := val.(*amf3.ByteArray)
		assert.True(t, ok)
		assert.Equal(t, ba, *ba3)
	})
}
