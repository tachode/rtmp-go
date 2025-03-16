package message_test

import (
	"bytes"
	"testing"

	"github.com/tachode/rtmp-go/message"
)

func TestSetPeerBandwidth_Marshal(t *testing.T) {
	tests := []struct {
		name      string
		input     message.SetPeerBandwidth
		wantBytes []byte
		wantErr   bool
	}{
		{
			name: "Valid Marshal",
			input: message.SetPeerBandwidth{
				WindowSize: 123456,
				LimitType:  message.BandwidthLimitSoft,
			},
			wantBytes: []byte{0x00, 0x01, 0xe2, 0x40, 0x01},
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotBytes, err := tt.input.Marshal()
			if (err != nil) != tt.wantErr {
				t.Errorf("Marshal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !bytes.Equal(gotBytes, tt.wantBytes) {
				t.Errorf("Marshal() gotBytes = %v, want %v", gotBytes, tt.wantBytes)
			}
		})
	}
}

func TestSetPeerBandwidth_Unmarshal(t *testing.T) {
	tests := []struct {
		name    string
		input   []byte
		want    message.SetPeerBandwidth
		wantErr bool
	}{
		{
			name:  "Valid Unmarshal",
			input: []byte{0x00, 0x01, 0xe2, 0x40, 0x00},
			want: message.SetPeerBandwidth{
				WindowSize: 123456,
				LimitType:  message.BandwidthLimitHard,
			},
			wantErr: false,
		},
		{
			name:    "Invalid Unmarshal - Short Data",
			input:   []byte{0x00, 0x01, 0xe2},
			want:    message.SetPeerBandwidth{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got message.SetPeerBandwidth
			err := got.Unmarshal(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Unmarshal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Unmarshal() got = %v, want %v", got, tt.want)
			}
		})
	}
}
