package amf3_test

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tachode/rtmp-go/amf3"
)

func TestArray_ReadWrite(t *testing.T) {
	t.Run("dense only array", func(t *testing.T) {
		// array-marker + U29A-value(count=2, low bit=1) + empty-string + 2 values
		var data bytes.Buffer
		data.WriteByte(byte(amf3.ArrayMarker))
		data.Write(u29Encode((2 << 1) | 1)) // dense count = 2
		data.Write(utf8vrEmpty())           // empty associative portion
		// Dense value 1: integer 42
		data.WriteByte(byte(amf3.IntegerMarker))
		data.Write(u29Encode(42))
		// Dense value 2: string "hi"
		data.WriteByte(byte(amf3.StringMarker))
		data.Write(utf8vr("hi"))

		r := amf3.NewReader(&data)
		val, err := r.ReadValue()
		assert.NoError(t, err)

		arr, ok := val.(*amf3.Array)
		assert.True(t, ok)
		assert.Equal(t, 2, len(arr.Dense))
		assert.Equal(t, amf3.Integer(42), arr.Dense[0])
		assert.Equal(t, amf3.String("hi"), arr.Dense[1])
		assert.Equal(t, 0, len(arr.Associative))
	})

	t.Run("associative only array", func(t *testing.T) {
		var data bytes.Buffer
		data.WriteByte(byte(amf3.ArrayMarker))
		data.Write(u29Encode((0 << 1) | 1)) // dense count = 0
		// Associative: "key" -> integer 99
		data.Write(utf8vr("key"))
		data.WriteByte(byte(amf3.IntegerMarker))
		data.Write(u29Encode(99))
		// End associative
		data.Write(utf8vrEmpty())

		r := amf3.NewReader(&data)
		val, err := r.ReadValue()
		assert.NoError(t, err)

		arr, ok := val.(*amf3.Array)
		assert.True(t, ok)
		assert.Equal(t, 0, len(arr.Dense))
		assert.Equal(t, amf3.Integer(99), arr.Associative["key"])
	})

	t.Run("mixed array", func(t *testing.T) {
		var data bytes.Buffer
		data.WriteByte(byte(amf3.ArrayMarker))
		data.Write(u29Encode((1 << 1) | 1)) // dense count = 1
		// Associative: "name" -> string "test"
		data.Write(utf8vr("name"))
		data.WriteByte(byte(amf3.StringMarker))
		data.Write(utf8vr("test"))
		data.Write(utf8vrEmpty()) // end associative
		// Dense value: true
		data.WriteByte(byte(amf3.TrueMarker))

		r := amf3.NewReader(&data)
		val, err := r.ReadValue()
		assert.NoError(t, err)

		arr, ok := val.(*amf3.Array)
		assert.True(t, ok)
		assert.Equal(t, 1, len(arr.Dense))
		assert.Equal(t, amf3.Boolean(true), arr.Dense[0])
		assert.Equal(t, amf3.String("test"), arr.Associative["name"])
	})

	t.Run("write dense array", func(t *testing.T) {
		arr := &amf3.Array{
			Associative: map[string]any{},
			Dense:       []any{amf3.Integer(1), amf3.Integer(2)},
		}

		var buf bytes.Buffer
		w := amf3.NewWriter(&buf)
		err := w.WriteValue(arr)
		assert.NoError(t, err)

		// Read it back
		r := amf3.NewReader(&buf)
		val, err := r.ReadValue()
		assert.NoError(t, err)

		readArr, ok := val.(*amf3.Array)
		assert.True(t, ok)
		assert.Equal(t, 2, len(readArr.Dense))
		assert.Equal(t, amf3.Integer(1), readArr.Dense[0])
		assert.Equal(t, amf3.Integer(2), readArr.Dense[1])
	})

	t.Run("write empty array", func(t *testing.T) {
		arr := &amf3.Array{
			Associative: map[string]any{},
			Dense:       []any{},
		}

		var buf bytes.Buffer
		w := amf3.NewWriter(&buf)
		err := w.WriteValue(arr)
		assert.NoError(t, err)

		r := amf3.NewReader(&buf)
		val, err := r.ReadValue()
		assert.NoError(t, err)

		readArr, ok := val.(*amf3.Array)
		assert.True(t, ok)
		assert.Equal(t, 0, len(readArr.Dense))
		assert.Equal(t, 0, len(readArr.Associative))
	})
}
