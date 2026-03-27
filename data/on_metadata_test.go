package data

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tachode/rtmp-go/amf0"
	"github.com/tachode/rtmp-go/message"
)

func boolPtr(b bool) *bool { return &b }

func TestOnMetaData_FromDataMessage_Direct(t *testing.T) {
	msg := &message.Amf0DataMessage{
		Handler: "onMetaData",
		Parameters: []any{
			amf0.EcmaArray{
				"audiocodecid":    float64(10),
				"audiodatarate":   float64(160),
				"audiosamplerate": float64(48000),
				"audiosamplesize": float64(16),
				"stereo":          true,
				"videocodecid":    float64(1.752589105e+09),
				"videodatarate":   float64(6000),
				"width":           float64(1280),
				"height":          float64(720),
				"framerate":       float64(60),
				"duration":        float64(0),
				"fileSize":        float64(0),
				"audiochannels":   float64(2),
				"encoder":         "obs-output module (libobs version 32.1.0)",
				"2.1":             false,
				"3.1":             false,
				"4.0":             false,
				"4.1":             false,
				"5.1":             false,
				"7.1":             false,
			},
		},
	}

	handler, err := FromDataMessage(msg)
	require.NoError(t, err)
	md, ok := handler.(*OnMetaData)
	require.True(t, ok)

	assert.Equal(t, message.AudioCodecId(10), md.AudioCodecId)
	assert.Equal(t, float64(160), md.AudioDataRate)
	assert.Equal(t, float64(48000), md.AudioSampleRate)
	assert.Equal(t, float64(16), md.AudioSampleSize)
	assert.True(t, md.Stereo)
	assert.Equal(t, message.VideoCodecId(1752589105), md.VideoCodecId)
	assert.Equal(t, float64(6000), md.VideoDataRate)
	assert.Equal(t, float64(1280), md.Width)
	assert.Equal(t, float64(720), md.Height)
	assert.Equal(t, float64(60), md.FrameRate)
	assert.Equal(t, float64(0), md.Duration)
	assert.Equal(t, float64(0), md.FileSize)
	assert.Equal(t, float64(2), md.AudioChannels)
	assert.Equal(t, "obs-output module (libobs version 32.1.0)", md.Encoder)
	require.NotNil(t, md.Surround2Point1)
	assert.False(t, *md.Surround2Point1)
	require.NotNil(t, md.Surround3Point1)
	assert.False(t, *md.Surround3Point1)
	require.NotNil(t, md.Surround4Point0)
	assert.False(t, *md.Surround4Point0)
	require.NotNil(t, md.Surround4Point1)
	assert.False(t, *md.Surround4Point1)
	require.NotNil(t, md.Surround5Point1)
	assert.False(t, *md.Surround5Point1)
	require.NotNil(t, md.Surround7Point1)
	assert.False(t, *md.Surround7Point1)
}

func TestOnMetaData_FromDataMessage_SetDataFrame(t *testing.T) {
	msg := &message.Amf0DataMessage{
		Handler: "@setDataFrame",
		Parameters: []any{
			"onMetaData",
			amf0.EcmaArray{
				"width":  float64(1920),
				"height": float64(1080),
			},
		},
	}

	handler, err := FromDataMessage(msg)
	require.NoError(t, err)
	md, ok := handler.(*OnMetaData)
	require.True(t, ok)

	assert.Equal(t, float64(1920), md.Width)
	assert.Equal(t, float64(1080), md.Height)
}

func TestOnMetaData_FromDataMessage_SetDataFrameAmf0String(t *testing.T) {
	msg := &message.Amf0DataMessage{
		Handler: "@setDataFrame",
		Parameters: []any{
			amf0.String("onMetaData"),
			amf0.EcmaArray{
				"width": float64(640),
			},
		},
	}

	handler, err := FromDataMessage(msg)
	require.NoError(t, err)
	md, ok := handler.(*OnMetaData)
	require.True(t, ok)

	assert.Equal(t, float64(640), md.Width)
}

func TestOnMetaData_RoundTrip(t *testing.T) {
	original := OnMetaData{
		AudioCodecId:    message.AudioCodecIdAAC,
		AudioDataRate:   160,
		AudioDelay:      0.5,
		AudioSampleRate: 48000,
		AudioSampleSize: 16,
		Stereo:          true,
		VideoCodecId:    message.VideoCodecIdHevc_ERTMP,
		VideoDataRate:   6000,
		Width:           1280,
		Height:          720,
		FrameRate:       60,
		Duration:        120.5,
		FileSize:        1024000,
		CanSeekToEnd:    true,
		CreationDate:    "2026-03-27",
		AudioChannels:   2,
		Encoder:         "obs-output module",
		Surround5Point1: boolPtr(true),
		Surround7Point1: boolPtr(true),
	}

	msg, err := original.ToDataMessage()
	require.NoError(t, err)

	var parsed OnMetaData
	err = parsed.FromDataMessage(msg)
	require.NoError(t, err)

	assert.Equal(t, original, parsed)
}

func TestOnMetaData_RoundTrip_TrackIdInfoMaps(t *testing.T) {
	original := OnMetaData{
		VideoCodecId:  message.VideoCodecIdAV1_ERTMP,
		Width:         1280,
		Height:        720,
		VideoDataRate: 6000,
		VideoTrackIdInfoMap: map[int]VideoTrackInfo{
			1: {
				Width:         1024,
				Height:        768,
				VideoDataRate: 2000,
			},
			2: {
				Width:         3840,
				Height:        2160,
				VideoDataRate: 30000,
				VideoCodecId:  message.VideoCodecIdAV1_ERTMP,
			},
		},
		AudioTrackIdInfoMap: map[int]AudioTrackInfo{
			1: {
				AudioDataRate: 256,
				Channels:      2,
				SampleRate:    44100,
				AudioCodecId:  message.AudioCodecIdAAC_ERTMP,
			},
		},
	}

	msg, err := original.ToDataMessage()
	require.NoError(t, err)

	var parsed OnMetaData
	err = parsed.FromDataMessage(msg)
	require.NoError(t, err)

	assert.Equal(t, original.Width, parsed.Width)
	assert.Equal(t, original.Height, parsed.Height)
	require.NotNil(t, parsed.VideoTrackIdInfoMap)
	assert.Equal(t, float64(1024), parsed.VideoTrackIdInfoMap[1].Width)
	assert.Equal(t, float64(768), parsed.VideoTrackIdInfoMap[1].Height)
	assert.Equal(t, float64(2000), parsed.VideoTrackIdInfoMap[1].VideoDataRate)
	assert.Equal(t, float64(3840), parsed.VideoTrackIdInfoMap[2].Width)
	assert.Equal(t, float64(30000), parsed.VideoTrackIdInfoMap[2].VideoDataRate)
	assert.Equal(t, message.VideoCodecIdAV1_ERTMP, parsed.VideoTrackIdInfoMap[2].VideoCodecId)
	require.NotNil(t, parsed.AudioTrackIdInfoMap)
	assert.Equal(t, float64(256), parsed.AudioTrackIdInfoMap[1].AudioDataRate)
	assert.Equal(t, float64(2), parsed.AudioTrackIdInfoMap[1].Channels)
	assert.Equal(t, float64(44100), parsed.AudioTrackIdInfoMap[1].SampleRate)
	assert.Equal(t, message.AudioCodecIdAAC_ERTMP, parsed.AudioTrackIdInfoMap[1].AudioCodecId)
}

func TestOnMetaData_ToDataMessage_OmitsZeroFields(t *testing.T) {
	md := OnMetaData{
		Width:  1280,
		Height: 720,
	}

	msg, err := md.ToDataMessage()
	require.NoError(t, err)

	dataMsg, ok := msg.(*message.Amf0DataMessage)
	require.True(t, ok)
	assert.Equal(t, "onMetaData", dataMsg.Handler)

	props := dataMsg.Parameters[0].(amf0.EcmaArray)
	assert.Equal(t, float64(1280), props["width"])
	assert.Equal(t, float64(720), props["height"])
	_, hasAudioCodec := props["audiocodecid"]
	assert.False(t, hasAudioCodec)
	_, hasStereo := props["stereo"]
	assert.False(t, hasStereo)
	_, hasEncoder := props["encoder"]
	assert.False(t, hasEncoder)
}

func TestOnMetaData_EmptyParameters(t *testing.T) {
	msg := &message.Amf0DataMessage{
		Handler:    "onMetaData",
		Parameters: nil,
	}

	handler, err := FromDataMessage(msg)
	require.NoError(t, err)
	md, ok := handler.(*OnMetaData)
	require.True(t, ok)
	assert.Equal(t, float64(0), md.Width)
}

func TestFromDataMessage_UnknownHandler(t *testing.T) {
	msg := &message.Amf0DataMessage{
		Handler: "unknownHandler",
	}

	_, err := FromDataMessage(msg)
	assert.ErrorIs(t, err, UnknownHandlerError)
}

func TestFromDataMessage_SetDataFrameUnknownInner(t *testing.T) {
	msg := &message.Amf0DataMessage{
		Handler:    "@setDataFrame",
		Parameters: []any{"unknownHandler"},
	}

	_, err := FromDataMessage(msg)
	assert.ErrorIs(t, err, UnknownHandlerError)
}

func TestOnMetaData_TagAlias_ParsesAlternateNames(t *testing.T) {
	// ffmpeg sends "filesize" (lowercase) instead of "fileSize" (camelCase).
	// The tag `amf:"fileSize,filesize"` should accept both.
	msg := &message.Amf0DataMessage{
		Handler: "onMetaData",
		Parameters: []any{
			amf0.EcmaArray{
				"filesize": float64(999),
			},
		},
	}

	handler, err := FromDataMessage(msg)
	require.NoError(t, err)
	md := handler.(*OnMetaData)
	assert.Equal(t, float64(999), md.FileSize)

	// Round-trip should serialize as "fileSize" (the first/canonical name).
	out, err := md.ToDataMessage()
	require.NoError(t, err)
	props := out.(*message.Amf0DataMessage).Parameters[0].(amf0.EcmaArray)
	_, hasCanonical := props["fileSize"]
	assert.True(t, hasCanonical)
	_, hasAlias := props["filesize"]
	assert.False(t, hasAlias)
}

func TestOnMetaData_TagAlias_PrefersFirstName(t *testing.T) {
	// When both names are present, the first (canonical) one wins.
	msg := &message.Amf0DataMessage{
		Handler: "onMetaData",
		Parameters: []any{
			amf0.EcmaArray{
				"fileSize": float64(100),
				"filesize": float64(200),
			},
		},
	}

	handler, err := FromDataMessage(msg)
	require.NoError(t, err)
	md := handler.(*OnMetaData)
	assert.Equal(t, float64(100), md.FileSize)
}
