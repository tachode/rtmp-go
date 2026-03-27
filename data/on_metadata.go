package data

import (
	"encoding/json"

	"github.com/tachode/rtmp-go/amf0"
	"github.com/tachode/rtmp-go/message"
)

func init() { RegisterHandler(new(OnMetaData)) }

// OnMetaData carries stream metadata, typically sent at the start of a stream
// via an AMF0 data message with handler "onMetaData" (or wrapped in
// "@setDataFrame").
type OnMetaData struct {
	// Audio properties (E-RTMP spec)

	// AudioCodecId is the audio codec ID. For legacy codecs this is a small
	// integer matching AudioTagHeader CodecID values. When FourCC signaling
	// is used, this is a big-endian FourCC value (e.g. mp4a == 0x6d703461).
	AudioCodecId message.AudioCodecId `amf:"audiocodecid"`

	// AudioDataRate is the audio bitrate in kilobits per second.
	AudioDataRate float64 `amf:"audiodatarate"`

	// AudioDelay is the delay introduced by the audio codec, in seconds.
	AudioDelay float64 `amf:"audiodelay"`

	// AudioSampleRate is the frequency at which the audio stream is replayed.
	AudioSampleRate float64 `amf:"audiosamplerate"`

	// AudioSampleSize is the number of bits used to represent each audio sample.
	AudioSampleSize float64 `amf:"audiosamplesize"`

	// Stereo indicates stereo audio.
	Stereo bool `amf:"stereo"`

	// Video properties (E-RTMP spec)

	// VideoCodecId is the video codec ID. For legacy codecs this is a small
	// integer matching VideoTagHeader CodecID values. When FourCC signaling
	// is used, this is a big-endian FourCC value (e.g. avc1 == 0x61766331).
	VideoCodecId message.VideoCodecId `amf:"videocodecid"`

	// VideoDataRate is the video bitrate in kilobits per second.
	VideoDataRate float64 `amf:"videodatarate"`

	// Width is the width of the video in pixels.
	Width float64 `amf:"width"`

	// Height is the height of the video in pixels.
	Height float64 `amf:"height"`

	// FrameRate is the number of frames per second.
	FrameRate float64 `amf:"framerate"`

	// General properties (E-RTMP spec)

	// Duration is the total duration of the file in seconds.
	Duration float64 `amf:"duration"`

	// FileSize is the total size of the file in bytes.
	// ffmpeg sends this as lowercase "filesize"; both forms are accepted.
	FileSize float64 `amf:"fileSize,filesize"`

	// CanSeekToEnd indicates that the last video frame is a key frame.
	CanSeekToEnd bool `amf:"canSeekToEnd"`

	// CreationDate is the creation date and time as a string.
	CreationDate string `amf:"creationdate"`

	// E-RTMP multitrack properties

	// AudioTrackIdInfoMap provides per-track metadata for additional audio
	// tracks beyond the default track. Keys are track IDs starting at 1;
	// the default track (ID 0) is described by the top-level fields.
	AudioTrackIdInfoMap map[int]AudioTrackInfo `amf:"audioTrackIdInfoMap" json:",omitempty"`

	// VideoTrackIdInfoMap provides per-track metadata for additional video
	// tracks beyond the default track. Keys are track IDs starting at 1;
	// the default track (ID 0) is described by the top-level fields.
	VideoTrackIdInfoMap map[int]VideoTrackInfo `amf:"videoTrackIdInfoMap" json:",omitempty"`

	// Properties observed in OBS metadata (not in the E-RTMP spec)

	// AudioChannels is the number of audio channels. Observed in OBS metadata.
	AudioChannels float64 `amf:"audiochannels"`

	// Encoder is a string identifying the software that produced the stream.
	// Observed in OBS and ffmpeg/libav metadata.
	Encoder string `amf:"encoder"`

	// Surround sound channel layout flags. Each indicates whether the
	// corresponding speaker configuration is available. Observed in OBS metadata.
	// A nil value means the field was not present in the metadata.
	Surround2Point1 *bool `amf:"2.1" json:",omitempty"` // stereo + LFE
	Surround3Point1 *bool `amf:"3.1" json:",omitempty"` // 3 channels + LFE
	Surround4Point0 *bool `amf:"4.0" json:",omitempty"` // quadraphonic
	Surround4Point1 *bool `amf:"4.1" json:",omitempty"` // quad + LFE
	Surround5Point1 *bool `amf:"5.1" json:",omitempty"` // 5 channels + LFE
	Surround7Point1 *bool `amf:"7.1" json:",omitempty"` // 7 channels + LFE

	// Properties observed in ffmpeg metadata (not in the E-RTMP spec)

	// MajorBrand is the ISO base media file format major brand (e.g. "mp42").
	// Observed in ffmpeg/Lavf metadata.
	MajorBrand string `amf:"major_brand"`

	// MinorVersion is the ISO base media file format minor version.
	// Observed in ffmpeg/Lavf metadata.
	MinorVersion float64 `amf:"minor_version"`

	// CompatibleBrands lists the ISO base media file format compatible brands
	// (e.g. "isommp41mp42"). Observed in ffmpeg/Lavf metadata.
	CompatibleBrands string `amf:"compatible_brands"`
}

func (m OnMetaData) HandlerName() string { return "onMetaData" }

func (m OnMetaData) String() string {
	b, _ := json.Marshal(m)
	return string(b)
}

func (m *OnMetaData) FromDataMessage(msg message.Data) error {
	params := msg.GetParameters()

	// Unwrap @setDataFrame: skip the leading "onMetaData" string parameter.
	if msg.GetHandler() == "@setDataFrame" && len(params) > 0 {
		if _, ok := params[0].(string); ok {
			params = params[1:]
		} else if s, ok := params[0].(amf0.String); ok && string(s) == "onMetaData" {
			params = params[1:]
		}
	}

	if len(params) == 0 {
		return nil
	}

	obj, ok := params[0].(message.Object)
	if !ok {
		return nil
	}

	message.ReadFields(obj, m)
	return nil
}

func (m *OnMetaData) ToDataMessage() (message.Data, error) {
	return &message.Amf0DataMessage{
		Handler:    "onMetaData",
		Parameters: []any{message.WriteFields(m)},
	}, nil
}

// VideoTrackInfo describes the characteristics of a single video track
// in the videoTrackIdInfoMap. Fields may differ from or repeat the
// top-level OnMetaData video fields.
type VideoTrackInfo struct {
	// Width is the width of this track in pixels.
	Width float64 `amf:"width"`
	// Height is the height of this track in pixels.
	Height float64 `amf:"height"`
	// VideoDataRate is the video bitrate in kilobits per second.
	VideoDataRate float64 `amf:"videodatarate"`
	// VideoCodecId is the video codec ID (legacy integer or FourCC value).
	VideoCodecId message.VideoCodecId `amf:"videocodecid"`
}

// AudioTrackInfo describes the characteristics of a single audio track
// in the audioTrackIdInfoMap. Fields may differ from or repeat the
// top-level OnMetaData audio fields.
type AudioTrackInfo struct {
	// AudioDataRate is the audio bitrate in kilobits per second.
	AudioDataRate float64 `amf:"audiodatarate"`
	// Channels is the number of audio channels.
	Channels float64 `amf:"channels"`
	// SampleRate is the audio sample rate in Hz.
	SampleRate float64 `amf:"samplerate"`
	// AudioCodecId is the audio codec ID (legacy integer or FourCC value).
	AudioCodecId message.AudioCodecId `amf:"audiocodecid"`
}
