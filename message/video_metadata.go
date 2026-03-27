package message

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/tachode/rtmp-go/amf0"
)

// VideoMetadata holds the metadata carried by a VideoPacketType.Metadata
// frame. The metadata frame contains a series of AMF-encoded [name, value]
// pairs. Currently, the only defined pair is "colorInfo".
type VideoMetadata struct {
	// ColorInfo holds structured color and HDR metadata, if present.
	ColorInfo *ColorInfo

	// Other preserves any name-value pairs whose name is not
	// recognized. Keys are pair names; values are the parsed AMF
	// representations. This allows round-tripping of metadata pairs
	// added in future spec revisions.
	Other map[string]any
}

// ColorInfo describes color space and HDR properties of a video stream,
// as defined in the E-RTMP v2 specification.
type ColorInfo struct {
	ColorConfig *ColorConfig `amf:"colorConfig"`
	HdrCll      *HdrCll      `amf:"hdrCll"`
	HdrMdcv     *HdrMdcv     `amf:"hdrMdcv"`
}

// ColorConfig describes the color encoding parameters.
type ColorConfig struct {
	BitDepth                float64 `amf:"bitDepth,omitempty"`
	ColorPrimaries          float64 `amf:"colorPrimaries,omitempty"`
	TransferCharacteristics float64 `amf:"transferCharacteristics,omitempty"`
	MatrixCoefficients      float64 `amf:"matrixCoefficients,omitempty"`
}

// HdrCll carries Content Light Level information (CEA-861.3).
type HdrCll struct {
	MaxFall float64 `amf:"maxFall,omitempty"`
	MaxCLL  float64 `amf:"maxCLL,omitempty"`
}

// HdrMdcv carries Mastering Display Color Volume metadata (SMPTE ST 2086).
type HdrMdcv struct {
	RedX         float64 `amf:"redX,omitempty"`
	RedY         float64 `amf:"redY,omitempty"`
	GreenX       float64 `amf:"greenX,omitempty"`
	GreenY       float64 `amf:"greenY,omitempty"`
	BlueX        float64 `amf:"blueX,omitempty"`
	BlueY        float64 `amf:"blueY,omitempty"`
	WhitePointX  float64 `amf:"whitePointX,omitempty"`
	WhitePointY  float64 `amf:"whitePointY,omitempty"`
	MaxLuminance float64 `amf:"maxLuminance,omitempty"`
	MinLuminance float64 `amf:"minLuminance,omitempty"`
}

func (m VideoMetadata) String() string {
	b, _ := json.Marshal(m)
	return string(b)
}

// UnmarshalAMF parses a series of AMF0 [name, value] pairs from the given
// bytes. If a pair's name is "colorInfo", its value is decoded into the
// ColorInfo struct. All other pairs are stored as raw AMF value bytes in
// UnknownMetadata.
func (m *VideoMetadata) UnmarshalAMF(data []byte) error {
	r := bytes.NewReader(data)

	for r.Len() > 0 {
		// Read name (AMF0 string)
		nameVal, err := amf0.Read(r)
		if err != nil {
			return fmt.Errorf("reading metadata pair name: %w", err)
		}
		name, ok := ToString(nameVal)
		if !ok {
			return fmt.Errorf("metadata pair name is not a string: %T", nameVal)
		}

		if name == "colorInfo" {
			// Read the value and decode into ColorInfo
			val, err := amf0.Read(r)
			if err != nil {
				return fmt.Errorf("reading colorInfo value: %w", err)
			}
			obj, ok := val.(Object)
			if !ok {
				return fmt.Errorf("colorInfo value is not an object: %T", val)
			}
			m.ColorInfo = &ColorInfo{}
			ReadFields(obj, m.ColorInfo)
		} else {
			// Parse and store the AMF value for unknown pairs
			val, err := amf0.Read(r)
			if err != nil {
				return fmt.Errorf("reading unknown metadata value %q: %w", name, err)
			}
			if m.Other == nil {
				m.Other = make(map[string]any)
			}
			m.Other[name] = val
		}
	}
	return nil
}

// MarshalAMF serializes the metadata as a series of AMF0 [name, value] pairs.
func (m *VideoMetadata) MarshalAMF() ([]byte, error) {
	var buf bytes.Buffer

	// Write colorInfo pair if present
	if m.ColorInfo != nil {
		if err := amf0.Write(&buf, "colorInfo"); err != nil {
			return nil, fmt.Errorf("writing colorInfo name: %w", err)
		}
		obj := WriteFields(m.ColorInfo)
		if err := amf0.Write(&buf, obj); err != nil {
			return nil, fmt.Errorf("writing colorInfo value: %w", err)
		}
	}

	// Write unknown metadata pairs
	for name, val := range m.Other {
		if err := amf0.Write(&buf, name); err != nil {
			return nil, fmt.Errorf("writing unknown metadata name %q: %w", name, err)
		}
		if err := amf0.Write(&buf, val); err != nil {
			return nil, fmt.Errorf("writing unknown metadata value %q: %w", name, err)
		}
	}

	return buf.Bytes(), nil
}
