package message_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tachode/rtmp-go/message"
)

func TestVideoMetadata_RoundTrip_AllSubObjects(t *testing.T) {
	vm := &message.VideoMetadata{
		ColorInfo: &message.ColorInfo{
			ColorConfig: &message.ColorConfig{
				BitDepth:                10,
				ColorPrimaries:          9,
				TransferCharacteristics: 16,
				MatrixCoefficients:      9,
			},
			HdrCll: &message.HdrCll{
				MaxFall: 100,
				MaxCLL:  1000,
			},
			HdrMdcv: &message.HdrMdcv{
				RedX:         0.708,
				RedY:         0.292,
				GreenX:       0.170,
				GreenY:       0.797,
				BlueX:        0.131,
				BlueY:        0.046,
				WhitePointX:  0.3127,
				WhitePointY:  0.3290,
				MaxLuminance: 1000,
				MinLuminance: 0.0001,
			},
		},
	}

	data, err := vm.MarshalAMF()
	require.NoError(t, err)
	require.NotEmpty(t, data)

	var parsed message.VideoMetadata
	err = parsed.UnmarshalAMF(data)
	require.NoError(t, err)

	require.NotNil(t, parsed.ColorInfo)
	require.NotNil(t, parsed.ColorInfo.ColorConfig)
	assert.Equal(t, float64(10), parsed.ColorInfo.ColorConfig.BitDepth)
	assert.Equal(t, float64(9), parsed.ColorInfo.ColorConfig.ColorPrimaries)
	assert.Equal(t, float64(16), parsed.ColorInfo.ColorConfig.TransferCharacteristics)
	assert.Equal(t, float64(9), parsed.ColorInfo.ColorConfig.MatrixCoefficients)

	require.NotNil(t, parsed.ColorInfo.HdrCll)
	assert.Equal(t, float64(100), parsed.ColorInfo.HdrCll.MaxFall)
	assert.Equal(t, float64(1000), parsed.ColorInfo.HdrCll.MaxCLL)

	require.NotNil(t, parsed.ColorInfo.HdrMdcv)
	assert.Equal(t, 0.708, parsed.ColorInfo.HdrMdcv.RedX)
	assert.Equal(t, 0.292, parsed.ColorInfo.HdrMdcv.RedY)
	assert.Equal(t, 0.170, parsed.ColorInfo.HdrMdcv.GreenX)
	assert.Equal(t, 0.797, parsed.ColorInfo.HdrMdcv.GreenY)
	assert.Equal(t, 0.131, parsed.ColorInfo.HdrMdcv.BlueX)
	assert.Equal(t, 0.046, parsed.ColorInfo.HdrMdcv.BlueY)
	assert.Equal(t, 0.3127, parsed.ColorInfo.HdrMdcv.WhitePointX)
	assert.Equal(t, 0.3290, parsed.ColorInfo.HdrMdcv.WhitePointY)
	assert.Equal(t, float64(1000), parsed.ColorInfo.HdrMdcv.MaxLuminance)
	assert.Equal(t, 0.0001, parsed.ColorInfo.HdrMdcv.MinLuminance)
}

func TestVideoMetadata_RoundTrip_PartialSubObjects(t *testing.T) {
	vm := &message.VideoMetadata{
		ColorInfo: &message.ColorInfo{
			ColorConfig: &message.ColorConfig{
				BitDepth:       10,
				ColorPrimaries: 1,
			},
		},
	}

	data, err := vm.MarshalAMF()
	require.NoError(t, err)

	var parsed message.VideoMetadata
	err = parsed.UnmarshalAMF(data)
	require.NoError(t, err)

	require.NotNil(t, parsed.ColorInfo)
	require.NotNil(t, parsed.ColorInfo.ColorConfig)
	assert.Equal(t, float64(10), parsed.ColorInfo.ColorConfig.BitDepth)
	assert.Equal(t, float64(1), parsed.ColorInfo.ColorConfig.ColorPrimaries)
	assert.Nil(t, parsed.ColorInfo.HdrCll)
	assert.Nil(t, parsed.ColorInfo.HdrMdcv)
}

func TestVideoMetadata_RoundTrip_NilColorInfo(t *testing.T) {
	vm := &message.VideoMetadata{}

	data, err := vm.MarshalAMF()
	require.NoError(t, err)
	assert.Empty(t, data)

	var parsed message.VideoMetadata
	err = parsed.UnmarshalAMF(data)
	require.NoError(t, err)
	assert.Nil(t, parsed.ColorInfo)
	assert.Nil(t, parsed.Other)
}

func TestVideoMetadata_UnknownMetadata_Preserved(t *testing.T) {
	original := &message.VideoMetadata{
		ColorInfo: &message.ColorInfo{
			ColorConfig: &message.ColorConfig{
				BitDepth: 8,
			},
		},
	}

	colorInfoBytes, err := original.MarshalAMF()
	require.NoError(t, err)

	unknown := &message.VideoMetadata{
		Other: map[string]any{
			"futureInfo": nil,
		},
	}
	unknownBytes, err := unknown.MarshalAMF()
	require.NoError(t, err)

	combined := append(colorInfoBytes, unknownBytes...)

	var parsed message.VideoMetadata
	err = parsed.UnmarshalAMF(combined)
	require.NoError(t, err)

	require.NotNil(t, parsed.ColorInfo)
	require.NotNil(t, parsed.ColorInfo.ColorConfig)
	assert.Equal(t, float64(8), parsed.ColorInfo.ColorConfig.BitDepth)

	require.Contains(t, parsed.Other, "futureInfo")

	// Round-trip
	roundTripped, err := parsed.MarshalAMF()
	require.NoError(t, err)

	var reparsed message.VideoMetadata
	err = reparsed.UnmarshalAMF(roundTripped)
	require.NoError(t, err)

	require.NotNil(t, reparsed.ColorInfo)
	assert.Equal(t, float64(8), reparsed.ColorInfo.ColorConfig.BitDepth)
	require.Contains(t, reparsed.Other, "futureInfo")
}

func TestVideoMetadata_String(t *testing.T) {
	vm := &message.VideoMetadata{
		ColorInfo: &message.ColorInfo{
			ColorConfig: &message.ColorConfig{
				BitDepth:       10,
				ColorPrimaries: 9,
			},
		},
	}
	s := vm.String()
	assert.Contains(t, s, "BitDepth")
	assert.Contains(t, s, "ColorPrimaries")
}
