package amf0_test

import (
	"bytes"
	"io"
	"testing"

	"github.com/tachode/rtmp-go/amf0"
)

type mockValue struct {
	data []byte
}

func (m mockValue) Type() amf0.Type {
	return 0x01 // Example type marker
}

func (m *mockValue) Read(r io.Reader) error {
	_, err := r.Read(m.data)
	return err
}

func (m mockValue) Write(w io.Writer) error {
	_, err := w.Write(m.data)
	return err
}

func TestRead(t *testing.T) {
	// Register mock value type
	amf0.RegisterType(new(mockValue))

	tests := []struct {
		name    string
		input   []byte
		wantErr bool
	}{
		{
			name:    "Valid type marker",
			input:   []byte{0x01, 0x02, 0x03},
			wantErr: false,
		},
		{
			name:    "Unknown type marker",
			input:   []byte{0xFF},
			wantErr: true,
		},
		{
			name:    "Incomplete data",
			input:   []byte{0x01},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := bytes.NewReader(tt.input)
			_, err := amf0.Read(r)
			if (err != nil) != tt.wantErr {
				t.Errorf("Read() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
