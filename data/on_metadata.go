package data

import (
	"encoding/json"
	"reflect"
	"strconv"
	"strings"

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

	readFields(obj, m)
	return nil
}

func (m *OnMetaData) ToDataMessage() (message.Data, error) {
	return &message.Amf0DataMessage{
		Handler:    "onMetaData",
		Parameters: []any{writeFields(m)},
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

// readFields populates the struct pointed to by target from the given
// message.Object, using `amf` struct tags as property names. Tags may
// contain comma-separated aliases (e.g. `amf:"fileSize,filesize"`); the
// first matching name wins.
func readFields(obj message.Object, target any) {
	v := reflect.ValueOf(target).Elem()
	t := v.Type()
	for i := range t.NumField() {
		tag := t.Field(i).Tag.Get("amf")
		if tag == "" {
			continue
		}
		names := strings.Split(tag, ",")
		fv := v.Field(i)
		switch fv.Kind() {
		case reflect.Float64:
			for _, name := range names {
				if val := message.GetFloat64(obj, name); val != 0 {
					fv.SetFloat(val)
					break
				}
			}
		case reflect.Uint32:
			for _, name := range names {
				if val := message.GetFloat64(obj, name); val != 0 {
					fv.SetUint(uint64(val))
					break
				}
			}
		case reflect.Bool:
			for _, name := range names {
				if val := message.GetBool(obj, name); val {
					fv.SetBool(val)
					break
				}
			}
		case reflect.String:
			for _, name := range names {
				if val := message.GetString(obj, name); val != "" {
					fv.SetString(val)
					break
				}
			}
		case reflect.Pointer:
			if fv.Type().Elem().Kind() == reflect.Bool {
				for _, name := range names {
					if bp := message.GetBoolPtr(obj, name); bp != nil {
						fv.Set(reflect.ValueOf(bp))
						break
					}
				}
			}
		case reflect.Map:
			if fv.Type().Key().Kind() == reflect.Int {
				for _, name := range names {
					readTrackIdInfoMap(obj, name, fv)
					if !fv.IsNil() {
						break
					}
				}
			}
		}
	}
}

// readTrackIdInfoMap reads a map[int]T field from a message.Object property,
// where T is a struct with amf tags on its fields.
func readTrackIdInfoMap(obj message.Object, key string, fv reflect.Value) {
	m := message.GetStringMap(obj, key)
	if m == nil {
		return
	}
	elemType := fv.Type().Elem()
	mapVal := reflect.MakeMap(fv.Type())
	for k, v := range m {
		id, err := strconv.Atoi(k)
		if err != nil {
			continue
		}
		trackObj, ok := v.(message.Object)
		if !ok {
			continue
		}
		elem := reflect.New(elemType)
		readFields(trackObj, elem.Interface())
		mapVal.SetMapIndex(reflect.ValueOf(id), elem.Elem())
	}
	if mapVal.Len() > 0 {
		fv.Set(mapVal)
	}
}

// writeFields serializes the struct into an amf0.EcmaArray, using `amf`
// struct tags as property names. When a tag contains comma-separated aliases,
// only the first name is used for serialization. Zero-valued fields are omitted.
func writeFields(source any) amf0.EcmaArray {
	props := amf0.EcmaArray{}
	v := reflect.ValueOf(source)
	if v.Kind() == reflect.Pointer {
		v = v.Elem()
	}
	t := v.Type()
	for i := range t.NumField() {
		tag := t.Field(i).Tag.Get("amf")
		if tag == "" {
			continue
		}
		name, _, _ := strings.Cut(tag, ",")
		fv := v.Field(i)
		switch fv.Kind() {
		case reflect.Float64:
			if fv.Float() != 0 {
				props[name] = fv.Float()
			}
		case reflect.Uint32:
			if fv.Uint() != 0 {
				props[name] = float64(fv.Uint())
			}
		case reflect.Bool:
			if fv.Bool() {
				props[name] = true
			}
		case reflect.String:
			if fv.String() != "" {
				props[name] = fv.String()
			}
		case reflect.Pointer:
			if !fv.IsNil() {
				props[name] = fv.Elem().Bool()
			}
		case reflect.Map:
			if fv.Len() > 0 && fv.Type().Key().Kind() == reflect.Int {
				innerMap := make(amf0.EcmaArray, fv.Len())
				for _, key := range fv.MapKeys() {
					innerMap[strconv.Itoa(int(key.Int()))] = writeFields(fv.MapIndex(key).Interface())
				}
				props[name] = innerMap
			}
		}
	}
	return props
}
