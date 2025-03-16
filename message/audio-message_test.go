package message_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tachode/rtmp-go/message"
)

func TestAudioMessage_Marshal(t *testing.T) {
	tests := []struct {
		name     string
		input    message.AudioMessage
		expected []byte
	}{
		{
			name: "AAC Stereo",
			input: message.AudioMessage{
				CodecId:    message.AudioCodecIdAAC,
				Rate:       message.AudioRate44kHz,
				SampleSize: message.AudioSize16Bit,
				Stereo:     true,
				AacType:    message.AacPacketTypeSequenceHeader,
				Payload:    []byte{0x01, 0x02, 0x03},
			},
			expected: []byte{0xAF, 0x00, 0x01, 0x02, 0x03},
		},
		{
			name: "MP3 Mono",
			input: message.AudioMessage{
				CodecId:    message.AudioCodecIdMP3,
				Rate:       message.AudioRate22kHz,
				SampleSize: message.AudioSize8Bit,
				Stereo:     false,
				Payload:    []byte{0x04, 0x05},
			},
			expected: []byte{0x28, 0x04, 0x05},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.input.Marshal()
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestAudioMessage_Unmarshal(t *testing.T) {
	tests := []struct {
		name     string
		input    []byte
		expected message.AudioMessage
		err      error
	}{
		{
			name:  "AAC Stereo",
			input: []byte{0xAF, 0x00, 0x01, 0x02, 0x03},
			expected: message.AudioMessage{
				CodecId:    message.AudioCodecIdAAC,
				Rate:       message.AudioRate44kHz,
				SampleSize: message.AudioSize16Bit,
				Stereo:     true,
				AacType:    message.AacPacketTypeSequenceHeader,
				Payload:    []byte{0x01, 0x02, 0x03},
			},
			err: nil,
		},
		{
			name:  "MP3 Mono",
			input: []byte{0x28, 0x04, 0x05},
			expected: message.AudioMessage{
				CodecId:    message.AudioCodecIdMP3,
				Rate:       message.AudioRate22kHz,
				SampleSize: message.AudioSize8Bit,
				Stereo:     false,
				Payload:    []byte{0x04, 0x05},
			},
			err: nil,
		},
		{
			name:  "Short Message",
			input: []byte{},
			err:   message.ErrShortMessage,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var result message.AudioMessage
			err := result.Unmarshal(tt.input)
			if tt.err != nil {
				assert.ErrorIs(t, err, tt.err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}
