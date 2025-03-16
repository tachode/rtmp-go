package message

import (
	"encoding/binary"
	"fmt"
)

type VideoFrameType uint8

//go:generate stringer -type=VideoFrameType -trimprefix=VideoFrameType
const (
	VideoFrameTypeInvalid VideoFrameType = iota
	VideoFrameTypeKeyframe
	VideoFrameTypeInterframe
	VideoFrameTypeDisposableInterframe
	VideoFrameTypeGeneratedKeyframe
	VideoFrameTypeVideoInfo
)

type VideoCodecId uint8

//go:generate stringer -type=VideoCodecId -trimprefix=VideoCodecId
const (
	VideoCodecIdReserved0 VideoCodecId = iota
	VideoCodecIdReserved1
	VideoCodecIdSorensonH263
	VideoCodecIdScreen1
	VideoCodecIdOn2VP6
	VideoCodecIdOn2VP6Alpha
	VideoCodecIdScreen2
	VideoCodecIdAvc
	VideoCodecIdRealH263
	VideoCodecIdMpeg4

	VideoCodecIdH264 = VideoCodecIdAvc
)

type AvcPacketType uint8

//go:generate stringer -type=AvcPacketType -trimprefix=AvcPacketType
const (
	AvcPacketTypeSequenceHeader AvcPacketType = iota
	AvcPacketTypeNALU
	AvcPacketTypeEndOfSequence
)

type VideoMessage struct {
	MetadataFields
	FrameType       VideoFrameType
	CodecId         VideoCodecId
	PacketType      AvcPacketType
	CompositionTime int32
	Payload         []byte
}

func init() { RegisterType(new(VideoMessage)) }

func (m VideoMessage) Type() Type {
	return TypeVideoMessage
}

func (m VideoMessage) Marshal() ([]byte, error) {
	out := make([]byte, 1)
	out[0] = byte(m.FrameType)<<4 | byte(m.CodecId&0x0f)
	if m.CodecId == VideoCodecIdAvc {
		out = binary.BigEndian.AppendUint32(out, uint32(m.CompositionTime))
		out[1] = byte(m.PacketType)
	}
	out = append(out, m.Payload...)
	return out, nil
}

func (m *VideoMessage) Unmarshal(data []byte) error {
	if len(data) < 1 {
		return ErrShortMessage
	}
	headerByte := data[0]
	m.FrameType = VideoFrameType(headerByte >> 4)
	m.CodecId = VideoCodecId(headerByte & 0x0f)
	data = data[1:]
	if m.CodecId == VideoCodecIdAvc {
		if len(data) < 4 {
			return ErrShortMessage
		}
		m.PacketType = AvcPacketType(data[0])
		m.CompositionTime = int32(binary.BigEndian.Uint32(data)) & 0xFF_FF_FF
		// Sign extend the 24-bit value
		if m.CompositionTime > 0x7F_FF_FF {
			m.CompositionTime -= 0x1_00_00_00
		}
		data = data[4:]
	}
	m.Payload = data
	return nil
}

func (m VideoMessage) String() string {
	return fmt.Sprintf("%v: %+v frame:%v codec:%v packet:%v ct:%v payload:%v bytes", m.Type(), m.MetadataFields,
		m.FrameType, m.CodecId, m.PacketType, m.CompositionTime, len(m.Payload))
}
