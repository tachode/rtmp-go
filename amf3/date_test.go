package amf3_test

import (
	"bytes"
	"encoding/binary"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/tachode/rtmp-go/amf3"
)

func TestDate_ReadWrite(t *testing.T) {
	t.Run("read date", func(t *testing.T) {
		expectedTime := time.Date(2021, time.January, 1, 0, 0, 0, 0, time.UTC)
		var data bytes.Buffer
		data.WriteByte(byte(amf3.DateMarker))
		// U29D-value: low bit = 1, remaining bits unused → 0x01
		data.WriteByte(0x01)
		binary.Write(&data, binary.BigEndian, float64(expectedTime.UnixMilli()))

		r := amf3.NewReader(&data)
		val, err := r.ReadValue()
		assert.NoError(t, err)
		d, ok := val.(*amf3.Date)
		assert.True(t, ok)
		assert.True(t, time.Time(*d).Equal(expectedTime))
	})

	t.Run("write date", func(t *testing.T) {
		expectedTime := time.Date(2021, time.January, 1, 0, 0, 0, 0, time.UTC)
		var buf bytes.Buffer
		w := amf3.NewWriter(&buf)
		err := w.WriteValue(amf3.Date(expectedTime))
		assert.NoError(t, err)

		// Verify: marker + U29D-value(0x01) + 8-byte double
		assert.Equal(t, byte(amf3.DateMarker), buf.Bytes()[0])
		assert.Equal(t, byte(0x01), buf.Bytes()[1])

		var millis float64
		binary.Read(bytes.NewReader(buf.Bytes()[2:]), binary.BigEndian, &millis)
		assert.Equal(t, float64(expectedTime.UnixMilli()), millis)
	})

	t.Run("write Go time.Time via convenience", func(t *testing.T) {
		expectedTime := time.Date(2021, time.January, 1, 0, 0, 0, 0, time.UTC)
		var buf bytes.Buffer
		w := amf3.NewWriter(&buf)
		err := w.WriteValue(expectedTime)
		assert.NoError(t, err)
		assert.Equal(t, byte(amf3.DateMarker), buf.Bytes()[0])
	})
}
