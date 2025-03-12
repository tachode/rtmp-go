package amf0_test

import (
	"bytes"
	"testing"
	"time"

	"github.com/tachode/rtmp-go/amf0"
)

func TestWrite(t *testing.T) {
	tests := []struct {
		name    string
		value   any
		wantErr bool
		expect  []byte
	}{
		{"nil value", nil, false, []byte{0x05}},
		{"float64 value", float64(3.14), false, []byte{0x00, 0x40, 0x09, 0x1E, 0xB8, 0x51, 0xEB, 0x85, 0x1F}},
		{"int value", 42, false, []byte{0x00, 0x40, 0x45, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}},
		{"bool value", true, false, []byte{0x01, 0x01}},
		{"short string value", "hello", false, []byte{0x02, 0x00, 0x05, 0x68, 0x65, 0x6C, 0x6C, 0x6F}},
		{"long string value", string(make([]byte, 0x10000)), false, nil},
		{"time value", time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC), false,
			[]byte{0x0B, 0x42, 0x72, 0x4E, 0x05, 0x33, 0x58, 0x00, 0x00, 0x00, 0x00}},
		{"unsupported type", struct{}{}, true, nil},
		{"Number value", amf0.Number(3.14), false, []byte{0x00, 0x40, 0x09, 0x1E, 0xB8, 0x51, 0xEB, 0x85, 0x1F}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			err := amf0.Write(&buf, tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("Write() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.expect != nil {
				if !bytes.Equal(buf.Bytes(), tt.expect) {
					t.Errorf("Write() = %v, want %v", buf.Bytes(), tt.expect)
				}
			}
		})
	}
}
