package message_test

import (
	"encoding/binary"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tachode/rtmp-go/message"
)

func TestVideoMessage_Marshal(t *testing.T) {
	tests := []struct {
		name     string
		msg      message.VideoMessage
		expected []byte
	}{
		{
			name: "Legacy Keyframe AVC NALU",
			msg: message.VideoMessage{
				FrameType:  message.VideoFrameTypeKeyframe,
				PacketType: message.ERTMPVideoPacketTypeCodedFrames,
				Tracks: []message.VideoTrack{{
					CodecId:         message.VideoCodecIdAvc,
					CompositionTime: 0,
					Payload:         []byte{0x01, 0x02, 0x03},
				}},
			},
			// 0x17 = keyframe(1)<<4 | AVC(7), then [packetType=1][CT 00 00 00][payload]
			expected: []byte{0x17, 0x01, 0x00, 0x00, 0x00, 0x01, 0x02, 0x03},
		},
		{
			name: "Legacy Keyframe AVC negative CT",
			msg: message.VideoMessage{
				FrameType:  message.VideoFrameTypeKeyframe,
				PacketType: message.ERTMPVideoPacketTypeCodedFrames,
				Tracks: []message.VideoTrack{{
					CodecId:         message.VideoCodecIdAvc,
					CompositionTime: -1,
					Payload:         []byte{0x01, 0x02, 0x03},
				}},
			},
			expected: []byte{0x17, 0x01, 0xff, 0xff, 0xff, 0x01, 0x02, 0x03},
		},
		{
			name: "Legacy Interframe non-AVC",
			msg: message.VideoMessage{
				FrameType:  message.VideoFrameTypeInterframe,
				PacketType: message.ERTMPVideoPacketTypeCodedFrames,
				Tracks: []message.VideoTrack{{
					CodecId: message.VideoCodecIdSorensonH263,
					Payload: []byte{0x04, 0x05},
				}},
			},
			expected: []byte{0x22, 0x04, 0x05},
		},
		{
			name: "E-RTMP AV1 SequenceStart",
			msg: message.VideoMessage{
				FrameType:  message.VideoFrameTypeKeyframe,
				PacketType: message.ERTMPVideoPacketTypeSequenceStart,
				Tracks: []message.VideoTrack{{
					CodecId: message.VideoCodecIdAV1_ERTMP,
					Payload: []byte{0xAA, 0xBB},
				}},
			},
			// [1:1][001:3][0000:4] = 0x90, then FOURCC "av01", then payload
			expected: []byte{0x90, 0x61, 0x76, 0x30, 0x31, 0xAA, 0xBB},
		},
		{
			name: "E-RTMP HEVC CodedFrames with CT",
			msg: message.VideoMessage{
				FrameType:  message.VideoFrameTypeInterframe,
				PacketType: message.ERTMPVideoPacketTypeCodedFrames,
				Tracks: []message.VideoTrack{{
					CodecId:         message.VideoCodecIdHevc_ERTMP,
					CompositionTime: 40,
					Payload:         []byte{0x01},
				}},
			},
			// [1:1][010:3][0001:4] = 0xA1, FOURCC "hvc1", SI24(40), payload
			expected: []byte{0xA1, 0x68, 0x76, 0x63, 0x31, 0x00, 0x00, 0x28, 0x01},
		},
		{
			name: "E-RTMP VVC CodedFrames with CT",
			msg: message.VideoMessage{
				FrameType:  message.VideoFrameTypeInterframe,
				PacketType: message.ERTMPVideoPacketTypeCodedFrames,
				Tracks: []message.VideoTrack{{
					CodecId:         message.VideoCodecIdVVC_ERTMP,
					CompositionTime: 40,
					Payload:         []byte{0x01},
				}},
			},
			// [1:1][010:3][0001:4] = 0xA1, FOURCC "vvc1", SI24(40), payload
			expected: []byte{0xA1, 0x76, 0x76, 0x63, 0x31, 0x00, 0x00, 0x28, 0x01},
		},
		{
			name: "E-RTMP VP9 CodedFramesX (CT implicit 0)",
			msg: message.VideoMessage{
				FrameType:  message.VideoFrameTypeKeyframe,
				PacketType: message.ERTMPVideoPacketTypeCodedFramesX,
				Tracks: []message.VideoTrack{{
					CodecId: message.VideoCodecIdVP9_ERTMP,
					Payload: []byte{0xDE, 0xAD},
				}},
			},
			// [1:1][001:3][0011:4] = 0x93, FOURCC "vp09", payload (no CT)
			expected: []byte{0x93, 0x76, 0x70, 0x30, 0x39, 0xDE, 0xAD},
		},
		{
			name: "E-RTMP SequenceEnd",
			msg: message.VideoMessage{
				FrameType:  message.VideoFrameTypeKeyframe,
				PacketType: message.ERTMPVideoPacketTypeSequenceEnd,
				Tracks: []message.VideoTrack{{
					CodecId: message.VideoCodecIdAV1_ERTMP,
				}},
			},
			// [1:1][001:3][0010:4] = 0x92, FOURCC "av01"
			expected: []byte{0x92, 0x61, 0x76, 0x30, 0x31},
		},
		{
			name: "E-RTMP Metadata",
			msg: message.VideoMessage{
				FrameType:     message.VideoFrameTypeKeyframe,
				PacketType:    message.ERTMPVideoPacketTypeMetadata,
				VideoMetadata: []byte{0x02, 0x00, 0x09}, // AMF0 string prefix
			},
			// [1:1][001:3][0100:4] = 0x94, then metadata (no FOURCC)
			expected: []byte{0x94, 0x02, 0x00, 0x09},
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
			name: "Legacy Keyframe AVC",
			data: []byte{0x17, 0x01, 0x00, 0x00, 0x00, 0x01, 0x02, 0x03},
			expected: message.VideoMessage{
				FrameType:  message.VideoFrameTypeKeyframe,
				PacketType: message.ERTMPVideoPacketTypeCodedFrames,
				Tracks: []message.VideoTrack{{
					CodecId:         message.VideoCodecIdAvc,
					CompositionTime: 0,
					Payload:         []byte{0x01, 0x02, 0x03},
				}},
			},
		},
		{
			name: "Legacy Interframe non-AVC",
			data: []byte{0x22, 0x04, 0x05},
			expected: message.VideoMessage{
				FrameType:  message.VideoFrameTypeInterframe,
				PacketType: message.ERTMPVideoPacketTypeCodedFrames,
				Tracks: []message.VideoTrack{{
					CodecId: message.VideoCodecIdSorensonH263,
					Payload: []byte{0x04, 0x05},
				}},
			},
		},
		{
			name: "E-RTMP AV1 SequenceStart",
			data: []byte{0x90, 0x61, 0x76, 0x30, 0x31, 0xAA, 0xBB},
			expected: message.VideoMessage{
				FrameType:  message.VideoFrameTypeKeyframe,
				PacketType: message.ERTMPVideoPacketTypeSequenceStart,
				Tracks: []message.VideoTrack{{
					CodecId: message.VideoCodecIdAV1_ERTMP,
					Payload: []byte{0xAA, 0xBB},
				}},
			},
		},
		{
			name: "E-RTMP SequenceEnd",
			data: []byte{0x92, 0x61, 0x76, 0x30, 0x31},
			expected: message.VideoMessage{
				FrameType:  message.VideoFrameTypeKeyframe,
				PacketType: message.ERTMPVideoPacketTypeSequenceEnd,
				Tracks: []message.VideoTrack{{
					CodecId: message.VideoCodecIdAV1_ERTMP,
				}},
			},
		},
		{
			name: "E-RTMP Metadata",
			data: []byte{0x94, 0x02, 0x00, 0x09},
			expected: message.VideoMessage{
				FrameType:     message.VideoFrameTypeKeyframe,
				PacketType:    message.ERTMPVideoPacketTypeMetadata,
				VideoMetadata: []byte{0x02, 0x00, 0x09},
			},
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

func TestVideoMessage_ERTMP_CodedFrames_WithCompositionTime(t *testing.T) {
	// HEVC and AVC include SI24 compositionTimeOffset in CodedFrames
	t.Run("HEVC positive CT", func(t *testing.T) {
		msg := message.VideoMessage{
			FrameType:  message.VideoFrameTypeInterframe,
			PacketType: message.ERTMPVideoPacketTypeCodedFrames,
			Tracks: []message.VideoTrack{{
				CodecId:         message.VideoCodecIdHevc_ERTMP,
				CompositionTime: 40,
				Payload:         []byte{0x01, 0x02},
			}},
		}
		data, err := msg.Marshal()
		require.NoError(t, err)

		var parsed message.VideoMessage
		err = parsed.Unmarshal(data)
		require.NoError(t, err)

		require.Len(t, parsed.Tracks, 1)
		assert.Equal(t, int32(40), parsed.Tracks[0].CompositionTime)
		assert.Equal(t, []byte{0x01, 0x02}, parsed.Tracks[0].Payload)
	})

	t.Run("AVC negative CT", func(t *testing.T) {
		msg := message.VideoMessage{
			FrameType:  message.VideoFrameTypeKeyframe,
			PacketType: message.ERTMPVideoPacketTypeCodedFrames,
			Tracks: []message.VideoTrack{{
				CodecId:         message.VideoCodecIdAvc_ERTMP,
				CompositionTime: -1,
				Payload:         []byte{0x03},
			}},
		}
		data, err := msg.Marshal()
		require.NoError(t, err)

		var parsed message.VideoMessage
		err = parsed.Unmarshal(data)
		require.NoError(t, err)

		require.Len(t, parsed.Tracks, 1)
		assert.Equal(t, int32(-1), parsed.Tracks[0].CompositionTime)
		assert.Equal(t, []byte{0x03}, parsed.Tracks[0].Payload)
	})

	t.Run("VVC positive CT", func(t *testing.T) {
		msg := message.VideoMessage{
			FrameType:  message.VideoFrameTypeInterframe,
			PacketType: message.ERTMPVideoPacketTypeCodedFrames,
			Tracks: []message.VideoTrack{{
				CodecId:         message.VideoCodecIdVVC_ERTMP,
				CompositionTime: 80,
				Payload:         []byte{0x05, 0x06},
			}},
		}
		data, err := msg.Marshal()
		require.NoError(t, err)

		var parsed message.VideoMessage
		err = parsed.Unmarshal(data)
		require.NoError(t, err)

		require.Len(t, parsed.Tracks, 1)
		assert.Equal(t, int32(80), parsed.Tracks[0].CompositionTime)
		assert.Equal(t, []byte{0x05, 0x06}, parsed.Tracks[0].Payload)
	})

	t.Run("VVC negative CT", func(t *testing.T) {
		msg := message.VideoMessage{
			FrameType:  message.VideoFrameTypeKeyframe,
			PacketType: message.ERTMPVideoPacketTypeCodedFrames,
			Tracks: []message.VideoTrack{{
				CodecId:         message.VideoCodecIdVVC_ERTMP,
				CompositionTime: -10,
				Payload:         []byte{0x07},
			}},
		}
		data, err := msg.Marshal()
		require.NoError(t, err)

		var parsed message.VideoMessage
		err = parsed.Unmarshal(data)
		require.NoError(t, err)

		require.Len(t, parsed.Tracks, 1)
		assert.Equal(t, int32(-10), parsed.Tracks[0].CompositionTime)
		assert.Equal(t, []byte{0x07}, parsed.Tracks[0].Payload)
	})

	t.Run("VP9 no CT", func(t *testing.T) {
		// VP9 CodedFrames should NOT have a compositionTime field
		msg := message.VideoMessage{
			FrameType:  message.VideoFrameTypeKeyframe,
			PacketType: message.ERTMPVideoPacketTypeCodedFrames,
			Tracks: []message.VideoTrack{{
				CodecId: message.VideoCodecIdVP9_ERTMP,
				Payload: []byte{0x01, 0x02},
			}},
		}
		data, err := msg.Marshal()
		require.NoError(t, err)

		var parsed message.VideoMessage
		err = parsed.Unmarshal(data)
		require.NoError(t, err)

		require.Len(t, parsed.Tracks, 1)
		assert.Equal(t, int32(0), parsed.Tracks[0].CompositionTime)
		assert.Equal(t, []byte{0x01, 0x02}, parsed.Tracks[0].Payload)
	})
}

func TestVideoMessage_ERTMP_CodedFramesX(t *testing.T) {
	msg := message.VideoMessage{
		FrameType:  message.VideoFrameTypeKeyframe,
		PacketType: message.ERTMPVideoPacketTypeCodedFramesX,
		Tracks: []message.VideoTrack{{
			CodecId: message.VideoCodecIdAvc_ERTMP,
			Payload: []byte{0xAA, 0xBB},
		}},
	}
	data, err := msg.Marshal()
	require.NoError(t, err)

	var parsed message.VideoMessage
	err = parsed.Unmarshal(data)
	require.NoError(t, err)

	assert.Equal(t, message.ERTMPVideoPacketTypeCodedFramesX, parsed.PacketType)
	require.Len(t, parsed.Tracks, 1)
	assert.Equal(t, int32(0), parsed.Tracks[0].CompositionTime)
	assert.Equal(t, []byte{0xAA, 0xBB}, parsed.Tracks[0].Payload)
}

func TestVideoMessage_ERTMP_Multitrack(t *testing.T) {
	t.Run("OneTrack", func(t *testing.T) {
		msg := message.VideoMessage{
			FrameType:      message.VideoFrameTypeKeyframe,
			PacketType:     message.ERTMPVideoPacketTypeCodedFramesX,
			MultitrackType: message.ERTMPAvMultitrackTypeOneTrack,
			Tracks: []message.VideoTrack{{
				TrackId: 0,
				CodecId: message.VideoCodecIdAV1_ERTMP,
				Payload: []byte{0x01, 0x02},
			}},
		}
		data, err := msg.Marshal()
		require.NoError(t, err)

		var parsed message.VideoMessage
		err = parsed.Unmarshal(data)
		require.NoError(t, err)

		assert.Equal(t, message.ERTMPVideoPacketTypeCodedFramesX, parsed.PacketType)
		assert.Equal(t, message.ERTMPAvMultitrackTypeOneTrack, parsed.MultitrackType)
		require.Len(t, parsed.Tracks, 1)
		assert.Equal(t, uint8(0), parsed.Tracks[0].TrackId)
		assert.Equal(t, message.VideoCodecIdAV1_ERTMP, parsed.Tracks[0].CodecId)
		assert.Equal(t, []byte{0x01, 0x02}, parsed.Tracks[0].Payload)
	})

	t.Run("ManyTracks", func(t *testing.T) {
		msg := message.VideoMessage{
			FrameType:      message.VideoFrameTypeKeyframe,
			PacketType:     message.ERTMPVideoPacketTypeCodedFramesX,
			MultitrackType: message.ERTMPAvMultitrackTypeManyTracks,
			Tracks: []message.VideoTrack{
				{TrackId: 0, CodecId: message.VideoCodecIdAV1_ERTMP, Payload: []byte{0xAA}},
				{TrackId: 1, CodecId: message.VideoCodecIdAV1_ERTMP, Payload: []byte{0xBB, 0xCC}},
			},
		}
		data, err := msg.Marshal()
		require.NoError(t, err)

		var parsed message.VideoMessage
		err = parsed.Unmarshal(data)
		require.NoError(t, err)

		assert.Equal(t, message.ERTMPVideoPacketTypeCodedFramesX, parsed.PacketType)
		assert.Equal(t, message.ERTMPAvMultitrackTypeManyTracks, parsed.MultitrackType)
		require.Len(t, parsed.Tracks, 2)
		assert.Equal(t, uint8(0), parsed.Tracks[0].TrackId)
		assert.Equal(t, []byte{0xAA}, parsed.Tracks[0].Payload)
		assert.Equal(t, uint8(1), parsed.Tracks[1].TrackId)
		assert.Equal(t, []byte{0xBB, 0xCC}, parsed.Tracks[1].Payload)
	})

	t.Run("ManyTracksManyCodecs", func(t *testing.T) {
		msg := message.VideoMessage{
			FrameType:      message.VideoFrameTypeKeyframe,
			PacketType:     message.ERTMPVideoPacketTypeCodedFramesX,
			MultitrackType: message.ERTMPAvMultitrackTypeManyTracksManyCodecs,
			Tracks: []message.VideoTrack{
				{TrackId: 0, CodecId: message.VideoCodecIdAV1_ERTMP, Payload: []byte{0xAA}},
				{TrackId: 1, CodecId: message.VideoCodecIdHevc_ERTMP, Payload: []byte{0xBB}},
			},
		}
		data, err := msg.Marshal()
		require.NoError(t, err)

		var parsed message.VideoMessage
		err = parsed.Unmarshal(data)
		require.NoError(t, err)

		assert.Equal(t, message.ERTMPVideoPacketTypeCodedFramesX, parsed.PacketType)
		assert.Equal(t, message.ERTMPAvMultitrackTypeManyTracksManyCodecs, parsed.MultitrackType)
		require.Len(t, parsed.Tracks, 2)
		assert.Equal(t, message.VideoCodecIdAV1_ERTMP, parsed.Tracks[0].CodecId)
		assert.Equal(t, message.VideoCodecIdHevc_ERTMP, parsed.Tracks[1].CodecId)
	})
}

func TestVideoMessage_ERTMP_ModEx(t *testing.T) {
	msg := message.VideoMessage{
		FrameType:           message.VideoFrameTypeKeyframe,
		PacketType:          message.ERTMPVideoPacketTypeCodedFramesX,
		TimestampNanoOffset: 500000,
		Tracks: []message.VideoTrack{{
			CodecId: message.VideoCodecIdAV1_ERTMP,
			Payload: []byte{0x01},
		}},
	}
	data, err := msg.Marshal()
	require.NoError(t, err)

	// First byte should have ModEx in low nibble: [1:1][001:3][0111:4] = 0x97
	assert.Equal(t, byte(0x97), data[0])

	var parsed message.VideoMessage
	err = parsed.Unmarshal(data)
	require.NoError(t, err)

	assert.Equal(t, message.ERTMPVideoPacketTypeCodedFramesX, parsed.PacketType)
	assert.Equal(t, uint32(500000), parsed.TimestampNanoOffset)
	require.Len(t, parsed.Tracks, 1)
	assert.Equal(t, message.VideoCodecIdAV1_ERTMP, parsed.Tracks[0].CodecId)
	assert.Equal(t, []byte{0x01}, parsed.Tracks[0].Payload)
}

func TestVideoMessage_ERTMP_Metadata(t *testing.T) {
	// Metadata message has no FOURCC, just AMF-encoded body
	metadata := []byte{0x02, 0x00, 0x09, 0x63, 0x6F, 0x6C, 0x6F, 0x72, 0x49, 0x6E, 0x66, 0x6F}

	msg := message.VideoMessage{
		FrameType:     message.VideoFrameTypeKeyframe,
		PacketType:    message.ERTMPVideoPacketTypeMetadata,
		VideoMetadata: metadata,
	}
	data, err := msg.Marshal()
	require.NoError(t, err)

	// First byte: [1:1][001:3][0100:4] = 0x94
	assert.Equal(t, byte(0x94), data[0])
	assert.Equal(t, metadata, data[1:])

	var parsed message.VideoMessage
	err = parsed.Unmarshal(data)
	require.NoError(t, err)

	assert.Equal(t, message.ERTMPVideoPacketTypeMetadata, parsed.PacketType)
	assert.Equal(t, metadata, parsed.VideoMetadata)
}

func TestVideoMessage_ERTMP_MPEG2TSSequenceStart(t *testing.T) {
	msg := message.VideoMessage{
		FrameType:  message.VideoFrameTypeKeyframe,
		PacketType: message.ERTMPVideoPacketTypeMPEG2TSSequenceStart,
		Tracks: []message.VideoTrack{{
			CodecId: message.VideoCodecIdAV1_ERTMP,
			Payload: []byte{0x01, 0x02, 0x03},
		}},
	}
	data, err := msg.Marshal()
	require.NoError(t, err)

	var parsed message.VideoMessage
	err = parsed.Unmarshal(data)
	require.NoError(t, err)

	assert.Equal(t, message.ERTMPVideoPacketTypeMPEG2TSSequenceStart, parsed.PacketType)
	require.Len(t, parsed.Tracks, 1)
	assert.Equal(t, message.VideoCodecIdAV1_ERTMP, parsed.Tracks[0].CodecId)
	assert.Equal(t, []byte{0x01, 0x02, 0x03}, parsed.Tracks[0].Payload)
}

func TestVideoMessage_ERTMP_RoundTrip(t *testing.T) {
	messages := []message.VideoMessage{
		{
			FrameType:  message.VideoFrameTypeKeyframe,
			PacketType: message.ERTMPVideoPacketTypeSequenceStart,
			Tracks: []message.VideoTrack{{
				CodecId: message.VideoCodecIdAV1_ERTMP,
				Payload: []byte{0x01, 0x02, 0x03, 0x04},
			}},
		},
		{
			FrameType:  message.VideoFrameTypeInterframe,
			PacketType: message.ERTMPVideoPacketTypeCodedFramesX,
			Tracks: []message.VideoTrack{{
				CodecId: message.VideoCodecIdVP8_ERTMP,
				Payload: []byte{0xAA, 0xBB},
			}},
		},
		{
			FrameType:  message.VideoFrameTypeKeyframe,
			PacketType: message.ERTMPVideoPacketTypeSequenceEnd,
			Tracks:     []message.VideoTrack{{CodecId: message.VideoCodecIdHevc_ERTMP}},
		},
	}

	for _, msg := range messages {
		data, err := msg.Marshal()
		require.NoError(t, err)

		var parsed message.VideoMessage
		err = parsed.Unmarshal(data)
		require.NoError(t, err)

		assert.Equal(t, msg.FrameType, parsed.FrameType)
		assert.Equal(t, msg.PacketType, parsed.PacketType)
		require.Len(t, parsed.Tracks, len(msg.Tracks))
		for i, track := range msg.Tracks {
			assert.Equal(t, track.CodecId, parsed.Tracks[i].CodecId)
			assert.Equal(t, track.Payload, parsed.Tracks[i].Payload)
		}
	}
}

func TestVideoMessage_IsERTMP(t *testing.T) {
	legacy := message.VideoMessage{
		Tracks: []message.VideoTrack{{CodecId: message.VideoCodecIdAvc}},
	}
	assert.False(t, legacy.IsERTMP())

	ertmp := message.VideoMessage{
		Tracks: []message.VideoTrack{{CodecId: message.VideoCodecIdAV1_ERTMP}},
	}
	assert.True(t, ertmp.IsERTMP())
}

func TestVideoMessage_ERTMP_WireFormat(t *testing.T) {
	// Verify exact wire format for E-RTMP single-track CodedFramesX
	// Format: [1:1][frameType:3][packetType:4] [FOURCC:32] [payload...]
	input := message.VideoMessage{
		FrameType:  message.VideoFrameTypeKeyframe,
		PacketType: message.ERTMPVideoPacketTypeCodedFramesX,
		Tracks: []message.VideoTrack{{
			CodecId: message.VideoCodecIdAV1_ERTMP,
			Payload: []byte{0xDE, 0xAD},
		}},
	}

	data, err := input.Marshal()
	require.NoError(t, err)
	require.Len(t, data, 7) // 1 header + 4 fourcc + 2 payload

	// Header byte: [1][001][0011] = 0x93
	assert.Equal(t, byte(0x93), data[0])

	// FOURCC "av01" = 0x61763031
	fourcc := binary.BigEndian.Uint32(data[1:5])
	assert.Equal(t, uint32(message.VideoCodecIdAV1_ERTMP), fourcc)

	// Payload
	assert.Equal(t, []byte{0xDE, 0xAD}, data[5:])
}
