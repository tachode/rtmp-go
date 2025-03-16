package message_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tachode/rtmp-go/message"
)

func TestVideoMessage_Marshal(t *testing.T) {
	tests := []struct {
		name     string
		msg      message.VideoMessage
		expected []byte
	}{
		{
			name: "Keyframe with AVC codec",
			msg: message.VideoMessage{
				FrameType:       message.VideoFrameTypeKeyframe,
				CodecId:         message.VideoCodecIdAvc,
				PacketType:      message.AvcPacketTypeNALU,
				CompositionTime: 0,
				Payload:         []byte{0x01, 0x02, 0x03},
			},
			expected: []byte{0x17, 0x01, 0x00, 0x00, 0x00, 0x01, 0x02, 0x03},
		},
		{
			name: "Keyframe with AVC codec, negative composition time offset",
			msg: message.VideoMessage{
				FrameType:       message.VideoFrameTypeKeyframe,
				CodecId:         message.VideoCodecIdAvc,
				PacketType:      message.AvcPacketTypeNALU,
				CompositionTime: -1,
				Payload:         []byte{0x01, 0x02, 0x03},
			},
			expected: []byte{0x17, 0x01, 0xff, 0xff, 0xff, 0x01, 0x02, 0x03},
		},
		{
			name: "Interframe with non-AVC codec",
			msg: message.VideoMessage{
				FrameType: message.VideoFrameTypeInterframe,
				CodecId:   message.VideoCodecIdSorensonH263,
				Payload:   []byte{0x04, 0x05},
			},
			expected: []byte{0x22, 0x04, 0x05},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.msg.Marshal()
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestVideoMessage_Unmarshal(t *testing.T) {
	tests := []struct {
		name        string
		data        []byte
		expected    message.VideoMessage
		expectError bool
	}{
		{
			name: "Keyframe with AVC codec",
			data: []byte{0x17, 0x01, 0x00, 0x00, 0x00, 0x01, 0x02, 0x03},
			expected: message.VideoMessage{
				FrameType:       message.VideoFrameTypeKeyframe,
				CodecId:         message.VideoCodecIdAvc,
				PacketType:      message.AvcPacketTypeNALU,
				CompositionTime: 0,
				Payload:         []byte{0x01, 0x02, 0x03},
			},
			expectError: false,
		},
		{
			name: "Interframe with non-AVC codec",
			data: []byte{0x22, 0x04, 0x05},
			expected: message.VideoMessage{
				FrameType: message.VideoFrameTypeInterframe,
				CodecId:   message.VideoCodecIdSorensonH263,
				Payload:   []byte{0x04, 0x05},
			},
			expectError: false,
		},
		{
			name:        "Invalid data length",
			data:        []byte{},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var msg message.VideoMessage
			err := msg.Unmarshal(tt.data)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, msg)
			}
		})
	}
}
