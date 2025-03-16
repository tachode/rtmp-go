package message

import (
	"fmt"
)

type AudioCodecId uint8

//go:generate stringer -type=AudioCodecId -trimprefix=AudioCodecId
const (
	AudioCodecIdLinearPCMPlatformEndian AudioCodecId = iota
	AudioCodecIdADPCM
	AudioCodecIdMP3
	AudioCodecIdLinearPCMLittleEndian
	AudioCodecIdNellymoser16kHzMono
	AudioCodecIdNellymoser8kHzMono
	AudioCodecIdNellymoser
	AudioCodecIdG711ALaw
	AudioCodecIdG711MuLaw
	AudioCodecIdExtendedHeader // Defined by Enhanced RTMP
	AudioCodecIdAAC
	AudioCodecIdSpeex
	AudioCodecIdReserved
	AudioCodecIdOpus // Nonstandard, used by ZLMediaKit
	AudioCodecIdMP38kHz
	AudioCodecIdDeviceSpecific
)

type AudioRate uint8

//go:generate stringer -type=AudioRate -trimprefix=AudioRate -linecomment
const (
	AudioRate5_5kHz AudioRate = iota // 5.5kHz
	AudioRate11kHz
	AudioRate22kHz
	AudioRate44kHz
)

type AudioSize uint8

//go:generate stringer -type=AudioSize -trimprefix=AudioSize
const (
	AudioSize8Bit AudioSize = iota
	AudioSize16Bit
)

type AacPacketType uint8

//go:generate stringer -type=AacPacketType -trimprefix=AacPacketType
const (
	AacPacketTypeSequenceHeader AacPacketType = iota
	AacPacketTypeRaw
)

type AudioMessage struct {
	MetadataFields
	CodecId    AudioCodecId
	Rate       AudioRate
	SampleSize AudioSize
	Stereo     bool
	AacType    AacPacketType
	Payload    []byte
}

func init() { RegisterType(new(AudioMessage)) }

func (m AudioMessage) Type() Type {
	return TypeAudioMessage
}

func (m AudioMessage) Marshal() ([]byte, error) {
	headerByte := byte(m.CodecId)<<4 | byte(m.Rate)<<2 | byte(m.SampleSize)<<1
	if m.Stereo {
		headerByte |= 0x01
	}
	out := []byte{headerByte}
	if m.CodecId == AudioCodecIdAAC {
		out = append(out, byte(m.AacType))
	}
	out = append(out, m.Payload...)
	return out, nil
}

func (m *AudioMessage) Unmarshal(data []byte) error {
	if len(data) < 1 {
		return ErrShortMessage
	}
	headerByte := data[0]
	m.CodecId = AudioCodecId(headerByte >> 4)
	m.Rate = AudioRate(headerByte >> 2 & 0x03)
	m.SampleSize = AudioSize(headerByte >> 1 & 0x01)
	m.Stereo = headerByte&0x01 != 0
	data = data[1:]
	if m.CodecId == AudioCodecIdAAC {
		if len(data) < 1 {
			return ErrShortMessage
		}
		m.AacType = AacPacketType(data[0])
		data = data[1:]
	}
	m.Payload = data
	return nil
}

func (m AudioMessage) String() string {
	return fmt.Sprintf("%v: %+v fmt:%v rate:%v sample:%v stereo:%v, aac:%v payload:%v bytes", m.Type(), m.MetadataFields,
		m.CodecId, m.Rate, m.SampleSize, m.Stereo, m.AacType, len(m.Payload))
}
