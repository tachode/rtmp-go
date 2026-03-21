package message_test

import (
	"encoding/binary"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tachode/rtmp-go/message"
)

func TestAudioMessage_Marshal(t *testing.T) {
	tests := []struct {
		name     string
		input    message.AudioMessage
		expected []byte
	}{
		{
			name: "Legacy AAC Stereo SequenceStart",
			input: message.AudioMessage{
				PacketType: message.ERTMPAudioPacketTypeSequenceStart,
				Rate:       message.AudioRate44kHz,
				SampleSize: message.AudioSize16Bit,
				Stereo:     true,
				Tracks: []message.AudioTrack{{
					CodecId: message.AudioCodecIdAAC,
					Payload: []byte{0x01, 0x02, 0x03},
				}},
			},
			expected: []byte{0xAF, 0x00, 0x01, 0x02, 0x03},
		},
		{
			name: "Legacy AAC CodedFrames",
			input: message.AudioMessage{
				PacketType: message.ERTMPAudioPacketTypeCodedFrames,
				Rate:       message.AudioRate44kHz,
				SampleSize: message.AudioSize16Bit,
				Stereo:     true,
				Tracks: []message.AudioTrack{{
					CodecId: message.AudioCodecIdAAC,
					Payload: []byte{0x04, 0x05},
				}},
			},
			expected: []byte{0xAF, 0x01, 0x04, 0x05},
		},
		{
			name: "Legacy MP3 Mono",
			input: message.AudioMessage{
				PacketType: message.ERTMPAudioPacketTypeCodedFrames,
				Rate:       message.AudioRate22kHz,
				SampleSize: message.AudioSize8Bit,
				Stereo:     false,
				Tracks: []message.AudioTrack{{
					CodecId: message.AudioCodecIdMP3,
					Payload: []byte{0x04, 0x05},
				}},
			},
			expected: []byte{0x28, 0x04, 0x05},
		},
		{
			name: "E-RTMP Opus SequenceStart",
			input: message.AudioMessage{
				PacketType: message.ERTMPAudioPacketTypeSequenceStart,
				Tracks: []message.AudioTrack{{
					CodecId: message.AudioCodecIdOpus_ERTMP,
					Payload: []byte{0xAA, 0xBB},
				}},
			},
			// Byte 0: [9:4][0:4] = 0x90
			// Bytes 1-4: FOURCC "Opus" = 0x4F707573
			// Bytes 5-6: payload
			expected: []byte{0x90, 0x4F, 0x70, 0x75, 0x73, 0xAA, 0xBB},
		},
		{
			name: "E-RTMP AAC CodedFrames",
			input: message.AudioMessage{
				PacketType: message.ERTMPAudioPacketTypeCodedFrames,
				Tracks: []message.AudioTrack{{
					CodecId: message.AudioCodecIdAAC_ERTMP,
					Payload: []byte{0x01, 0x02},
				}},
			},
			// Byte 0: [9:4][1:4] = 0x91
			// Bytes 1-4: FOURCC "mp4a" = 0x6D703461
			// Bytes 5-6: payload
			expected: []byte{0x91, 0x6D, 0x70, 0x34, 0x61, 0x01, 0x02},
		},
		{
			name: "E-RTMP SequenceEnd",
			input: message.AudioMessage{
				PacketType: message.ERTMPAudioPacketTypeSequenceEnd,
				Tracks: []message.AudioTrack{{
					CodecId: message.AudioCodecIdOpus_ERTMP,
				}},
			},
			// Byte 0: [9:4][2:4] = 0x92
			// Bytes 1-4: FOURCC "Opus"
			expected: []byte{0x92, 0x4F, 0x70, 0x75, 0x73},
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
			name:  "Legacy AAC Stereo SequenceStart",
			input: []byte{0xAF, 0x00, 0x01, 0x02, 0x03},
			expected: message.AudioMessage{
				PacketType: message.ERTMPAudioPacketTypeSequenceStart,
				Rate:       message.AudioRate44kHz,
				SampleSize: message.AudioSize16Bit,
				Stereo:     true,
				Tracks: []message.AudioTrack{{
					CodecId: message.AudioCodecIdAAC,
					Payload: []byte{0x01, 0x02, 0x03},
				}},
			},
		},
		{
			name:  "Legacy MP3 Mono",
			input: []byte{0x28, 0x04, 0x05},
			expected: message.AudioMessage{
				PacketType: message.ERTMPAudioPacketTypeCodedFrames,
				Rate:       message.AudioRate22kHz,
				SampleSize: message.AudioSize8Bit,
				Stereo:     false,
				Tracks: []message.AudioTrack{{
					CodecId: message.AudioCodecIdMP3,
					Payload: []byte{0x04, 0x05},
				}},
			},
		},
		{
			name:  "E-RTMP Opus SequenceStart",
			input: []byte{0x90, 0x4F, 0x70, 0x75, 0x73, 0xAA, 0xBB},
			expected: message.AudioMessage{
				PacketType: message.ERTMPAudioPacketTypeSequenceStart,
				Tracks: []message.AudioTrack{{
					CodecId: message.AudioCodecIdOpus_ERTMP,
					Payload: []byte{0xAA, 0xBB},
				}},
			},
		},
		{
			name:  "E-RTMP AAC CodedFrames",
			input: []byte{0x91, 0x6D, 0x70, 0x34, 0x61, 0x01, 0x02},
			expected: message.AudioMessage{
				PacketType: message.ERTMPAudioPacketTypeCodedFrames,
				Tracks: []message.AudioTrack{{
					CodecId: message.AudioCodecIdAAC_ERTMP,
					Payload: []byte{0x01, 0x02},
				}},
			},
		},
		{
			name:  "E-RTMP SequenceEnd",
			input: []byte{0x92, 0x4F, 0x70, 0x75, 0x73},
			expected: message.AudioMessage{
				PacketType: message.ERTMPAudioPacketTypeSequenceEnd,
				Tracks: []message.AudioTrack{{
					CodecId: message.AudioCodecIdOpus_ERTMP,
				}},
			},
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

func TestAudioMessage_ERTMP_Multitrack(t *testing.T) {
	t.Run("OneTrack", func(t *testing.T) {
		msg := message.AudioMessage{
			PacketType:     message.ERTMPAudioPacketTypeCodedFrames,
			MultitrackType: message.ERTMPAvMultitrackTypeOneTrack,
			Tracks: []message.AudioTrack{{
				TrackId: 0,
				CodecId: message.AudioCodecIdOpus_ERTMP,
				Payload: []byte{0x01, 0x02},
			}},
		}
		data, err := msg.Marshal()
		require.NoError(t, err)

		var parsed message.AudioMessage
		err = parsed.Unmarshal(data)
		require.NoError(t, err)

		assert.Equal(t, message.ERTMPAudioPacketTypeCodedFrames, parsed.PacketType)
		assert.Equal(t, message.ERTMPAvMultitrackTypeOneTrack, parsed.MultitrackType)
		require.Len(t, parsed.Tracks, 1)
		assert.Equal(t, uint8(0), parsed.Tracks[0].TrackId)
		assert.Equal(t, message.AudioCodecIdOpus_ERTMP, parsed.Tracks[0].CodecId)
		assert.Equal(t, []byte{0x01, 0x02}, parsed.Tracks[0].Payload)
	})

	t.Run("ManyTracks", func(t *testing.T) {
		msg := message.AudioMessage{
			PacketType:     message.ERTMPAudioPacketTypeCodedFrames,
			MultitrackType: message.ERTMPAvMultitrackTypeManyTracks,
			Tracks: []message.AudioTrack{
				{TrackId: 0, CodecId: message.AudioCodecIdOpus_ERTMP, Payload: []byte{0xAA}},
				{TrackId: 1, CodecId: message.AudioCodecIdOpus_ERTMP, Payload: []byte{0xBB, 0xCC}},
			},
		}
		data, err := msg.Marshal()
		require.NoError(t, err)

		var parsed message.AudioMessage
		err = parsed.Unmarshal(data)
		require.NoError(t, err)

		assert.Equal(t, message.ERTMPAudioPacketTypeCodedFrames, parsed.PacketType)
		assert.Equal(t, message.ERTMPAvMultitrackTypeManyTracks, parsed.MultitrackType)
		require.Len(t, parsed.Tracks, 2)
		assert.Equal(t, uint8(0), parsed.Tracks[0].TrackId)
		assert.Equal(t, []byte{0xAA}, parsed.Tracks[0].Payload)
		assert.Equal(t, uint8(1), parsed.Tracks[1].TrackId)
		assert.Equal(t, []byte{0xBB, 0xCC}, parsed.Tracks[1].Payload)
	})

	t.Run("ManyTracksManyCodecs", func(t *testing.T) {
		msg := message.AudioMessage{
			PacketType:     message.ERTMPAudioPacketTypeCodedFrames,
			MultitrackType: message.ERTMPAvMultitrackTypeManyTracksManyCodecs,
			Tracks: []message.AudioTrack{
				{TrackId: 0, CodecId: message.AudioCodecIdOpus_ERTMP, Payload: []byte{0xAA}},
				{TrackId: 1, CodecId: message.AudioCodecIdAAC_ERTMP, Payload: []byte{0xBB}},
			},
		}
		data, err := msg.Marshal()
		require.NoError(t, err)

		var parsed message.AudioMessage
		err = parsed.Unmarshal(data)
		require.NoError(t, err)

		assert.Equal(t, message.ERTMPAudioPacketTypeCodedFrames, parsed.PacketType)
		assert.Equal(t, message.ERTMPAvMultitrackTypeManyTracksManyCodecs, parsed.MultitrackType)
		require.Len(t, parsed.Tracks, 2)
		assert.Equal(t, message.AudioCodecIdOpus_ERTMP, parsed.Tracks[0].CodecId)
		assert.Equal(t, message.AudioCodecIdAAC_ERTMP, parsed.Tracks[1].CodecId)
	})
}

func TestAudioMessage_ERTMP_ModEx(t *testing.T) {
	// Build a message with TimestampNanoOffset
	msg := message.AudioMessage{
		PacketType:          message.ERTMPAudioPacketTypeCodedFrames,
		TimestampNanoOffset: 500000, // 0.5ms
		Tracks: []message.AudioTrack{{
			CodecId: message.AudioCodecIdOpus_ERTMP,
			Payload: []byte{0x01},
		}},
	}
	data, err := msg.Marshal()
	require.NoError(t, err)

	// First byte should be [9:4][ModEx=7:4] = 0x97
	assert.Equal(t, byte(0x97), data[0])

	var parsed message.AudioMessage
	err = parsed.Unmarshal(data)
	require.NoError(t, err)

	assert.Equal(t, message.ERTMPAudioPacketTypeCodedFrames, parsed.PacketType)
	assert.Equal(t, uint32(500000), parsed.TimestampNanoOffset)
	require.Len(t, parsed.Tracks, 1)
	assert.Equal(t, message.AudioCodecIdOpus_ERTMP, parsed.Tracks[0].CodecId)
	assert.Equal(t, []byte{0x01}, parsed.Tracks[0].Payload)
}

func TestAudioMessage_ERTMP_MultichannelConfig(t *testing.T) {
	t.Run("Native", func(t *testing.T) {
		msg := message.AudioMessage{
			PacketType: message.ERTMPAudioPacketTypeMultichannelConfig,
			Tracks: []message.AudioTrack{{
				CodecId:      message.AudioCodecIdOpus_ERTMP,
				ChannelOrder: message.ERTMPAudioChannelOrderNative,
				ChannelCount: 6,
				ChannelFlags: message.ERTMPAudioChannelMaskFrontLeft |
					message.ERTMPAudioChannelMaskFrontRight |
					message.ERTMPAudioChannelMaskFrontCenter |
					message.ERTMPAudioChannelMaskLowFrequency1 |
					message.ERTMPAudioChannelMaskBackLeft |
					message.ERTMPAudioChannelMaskBackRight,
			}},
		}
		data, err := msg.Marshal()
		require.NoError(t, err)

		var parsed message.AudioMessage
		err = parsed.Unmarshal(data)
		require.NoError(t, err)

		assert.Equal(t, message.ERTMPAudioPacketTypeMultichannelConfig, parsed.PacketType)
		require.Len(t, parsed.Tracks, 1)
		assert.Equal(t, message.ERTMPAudioChannelOrderNative, parsed.Tracks[0].ChannelOrder)
		assert.Equal(t, uint8(6), parsed.Tracks[0].ChannelCount)
		assert.Equal(t, message.ERTMPAudioChannelMask(0x3F), parsed.Tracks[0].ChannelFlags)
	})

	t.Run("Custom", func(t *testing.T) {
		msg := message.AudioMessage{
			PacketType: message.ERTMPAudioPacketTypeMultichannelConfig,
			Tracks: []message.AudioTrack{{
				CodecId:      message.AudioCodecIdOpus_ERTMP,
				ChannelOrder: message.ERTMPAudioChannelOrderCustom,
				ChannelCount: 2,
				ChannelMapping: []message.ERTMPAudioChannel{
					message.ERTMPAudioChannelFrontLeft,
					message.ERTMPAudioChannelFrontRight,
				},
			}},
		}
		data, err := msg.Marshal()
		require.NoError(t, err)

		var parsed message.AudioMessage
		err = parsed.Unmarshal(data)
		require.NoError(t, err)

		assert.Equal(t, message.ERTMPAudioChannelOrderCustom, parsed.Tracks[0].ChannelOrder)
		assert.Equal(t, uint8(2), parsed.Tracks[0].ChannelCount)
		assert.Equal(t, []message.ERTMPAudioChannel{
			message.ERTMPAudioChannelFrontLeft,
			message.ERTMPAudioChannelFrontRight,
		}, parsed.Tracks[0].ChannelMapping)
	})
}

func TestAudioMessage_ERTMP_RoundTrip(t *testing.T) {
	// Verify that marshal → unmarshal produces identical results for various E-RTMP messages
	messages := []message.AudioMessage{
		{
			PacketType: message.ERTMPAudioPacketTypeSequenceStart,
			Tracks: []message.AudioTrack{{
				CodecId: message.AudioCodecIdFlac_ERTMP,
				Payload: []byte{0x66, 0x4C, 0x61, 0x43, 0x00, 0x00, 0x00, 0x22},
			}},
		},
		{
			PacketType: message.ERTMPAudioPacketTypeCodedFrames,
			Tracks: []message.AudioTrack{{
				CodecId: message.AudioCodecIdAC3_ERTMP,
				Payload: []byte{0x0B, 0x77, 0x01, 0x02, 0x03},
			}},
		},
		{
			PacketType: message.ERTMPAudioPacketTypeSequenceEnd,
			Tracks:     []message.AudioTrack{{CodecId: message.AudioCodecIdMP3_ERTMP}},
		},
	}

	for _, msg := range messages {
		data, err := msg.Marshal()
		require.NoError(t, err)

		var parsed message.AudioMessage
		err = parsed.Unmarshal(data)
		require.NoError(t, err)

		assert.Equal(t, msg.PacketType, parsed.PacketType)
		require.Len(t, parsed.Tracks, len(msg.Tracks))
		for i, track := range msg.Tracks {
			assert.Equal(t, track.CodecId, parsed.Tracks[i].CodecId)
			assert.Equal(t, track.Payload, parsed.Tracks[i].Payload)
		}
	}
}

func TestAudioMessage_IsERTMP(t *testing.T) {
	legacy := message.AudioMessage{
		Tracks: []message.AudioTrack{{CodecId: message.AudioCodecIdAAC}},
	}
	assert.False(t, legacy.IsERTMP())

	ertmp := message.AudioMessage{
		Tracks: []message.AudioTrack{{CodecId: message.AudioCodecIdOpus_ERTMP}},
	}
	assert.True(t, ertmp.IsERTMP())
}

func TestAudioMessage_ERTMP_WireFormat(t *testing.T) {
	// Verify the exact wire format for an E-RTMP single-track CodedFrames message
	// Format: [9:4][PacketType:4] [FOURCC:32] [payload...]
	input := message.AudioMessage{
		PacketType: message.ERTMPAudioPacketTypeCodedFrames,
		Tracks: []message.AudioTrack{{
			CodecId: message.AudioCodecIdOpus_ERTMP,
			Payload: []byte{0xDE, 0xAD},
		}},
	}

	data, err := input.Marshal()
	require.NoError(t, err)
	require.Len(t, data, 7) // 1 header + 4 fourcc + 2 payload

	// Header byte
	assert.Equal(t, byte(0x91), data[0]) // soundFormat=9, packetType=1

	// FOURCC "Opus" = 0x4F707573
	fourcc := binary.BigEndian.Uint32(data[1:5])
	assert.Equal(t, uint32(message.AudioCodecIdOpus_ERTMP), fourcc)

	// Payload
	assert.Equal(t, []byte{0xDE, 0xAD}, data[5:])
}
