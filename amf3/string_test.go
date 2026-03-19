package amf3_test

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tachode/rtmp-go/amf3"
)

func TestString_ReadWrite(t *testing.T) {
	t.Run("simple string", func(t *testing.T) {
		// string-marker + UTF-8-vr("Hello")
		data := []byte{byte(amf3.StringMarker)}
		data = append(data, utf8vr("Hello")...)

		r := amf3.NewReader(bytes.NewReader(data))
		val, err := r.ReadValue()
		assert.NoError(t, err)
		assert.Equal(t, amf3.String("Hello"), val)
	})

	t.Run("empty string", func(t *testing.T) {
		data := []byte{byte(amf3.StringMarker)}
		data = append(data, utf8vrEmpty()...)

		r := amf3.NewReader(bytes.NewReader(data))
		val, err := r.ReadValue()
		assert.NoError(t, err)
		assert.Equal(t, amf3.String(""), val)
	})

	t.Run("write simple string", func(t *testing.T) {
		var buf bytes.Buffer
		w := amf3.NewWriter(&buf)
		err := w.WriteValue(amf3.String("Hello"))
		assert.NoError(t, err)

		expected := []byte{byte(amf3.StringMarker)}
		expected = append(expected, utf8vr("Hello")...)
		assert.Equal(t, expected, buf.Bytes())
	})

	t.Run("write empty string", func(t *testing.T) {
		var buf bytes.Buffer
		w := amf3.NewWriter(&buf)
		err := w.WriteValue(amf3.String(""))
		assert.NoError(t, err)

		expected := []byte{byte(amf3.StringMarker)}
		expected = append(expected, utf8vrEmpty()...)
		assert.Equal(t, expected, buf.Bytes())
	})

	t.Run("write Go string via convenience", func(t *testing.T) {
		var buf bytes.Buffer
		w := amf3.NewWriter(&buf)
		err := w.WriteValue("Hello")
		assert.NoError(t, err)

		expected := []byte{byte(amf3.StringMarker)}
		expected = append(expected, utf8vr("Hello")...)
		assert.Equal(t, expected, buf.Bytes())
	})
}

func TestString_References(t *testing.T) {
	t.Run("read string reference", func(t *testing.T) {
		// Write two strings: first as literal, second as reference to first
		var data []byte
		data = append(data, byte(amf3.StringMarker))
		data = append(data, utf8vr("Hello")...)
		data = append(data, byte(amf3.StringMarker))
		data = append(data, utf8vrRef(0)...) // reference to index 0

		r := amf3.NewReader(bytes.NewReader(data))

		val1, err := r.ReadValue()
		assert.NoError(t, err)
		assert.Equal(t, amf3.String("Hello"), val1)

		val2, err := r.ReadValue()
		assert.NoError(t, err)
		assert.Equal(t, amf3.String("Hello"), val2)
	})

	t.Run("write deduplicates strings", func(t *testing.T) {
		var buf bytes.Buffer
		w := amf3.NewWriter(&buf)

		err := w.WriteValue(amf3.String("Hello"))
		assert.NoError(t, err)

		err = w.WriteValue(amf3.String("Hello"))
		assert.NoError(t, err)

		// First: marker + literal "Hello"
		// Second: marker + reference to index 0
		expected := []byte{byte(amf3.StringMarker)}
		expected = append(expected, utf8vr("Hello")...)
		expected = append(expected, byte(amf3.StringMarker))
		expected = append(expected, utf8vrRef(0)...)
		assert.Equal(t, expected, buf.Bytes())
	})

	t.Run("empty string never referenced", func(t *testing.T) {
		var buf bytes.Buffer
		w := amf3.NewWriter(&buf)

		err := w.WriteValue(amf3.String(""))
		assert.NoError(t, err)
		err = w.WriteValue(amf3.String(""))
		assert.NoError(t, err)

		// Both should be literal empty string (0x01), not references
		expected := []byte{
			byte(amf3.StringMarker), 0x01,
			byte(amf3.StringMarker), 0x01,
		}
		assert.Equal(t, expected, buf.Bytes())
	})
}
