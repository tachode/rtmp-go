package message

import (
	"encoding/binary"
	"fmt"
)

type AudioCodecId uint32

//go:generate stringer -type=AudioCodecId -trimprefix=AudioCodecId
const (
	// Legacy codec IDs
	AudioCodecIdLinearPCMPlatformEndian AudioCodecId = iota
	AudioCodecIdADPCM
	AudioCodecIdMP3
	AudioCodecIdLinearPCMLittleEndian
	AudioCodecIdNellymoser16kHzMono
	AudioCodecIdNellymoser8kHzMono
	AudioCodecIdNellymoser
	AudioCodecIdG711ALaw
	AudioCodecIdG711MuLaw
	audioCodecIdExtendedHeader // Defined by E-RTMP, private
	AudioCodecIdAAC
	AudioCodecIdSpeex
	AudioCodecIdReserved
	AudioCodecIdOpus // Nonstandard, used by ZLMediaKit
	AudioCodecIdMP3_8kHz
	AudioCodecIdDeviceSpecific

	// E-RTMP codec IDs
	AudioCodecIdAC3_ERTMP  = AudioCodecId('a'<<24 | 'c'<<16 | '-'<<8 | '3')
	AudioCodecIdEAC3_ERTMP = AudioCodecId('e'<<24 | 'c'<<16 | '-'<<8 | '3')
	AudioCodecIdOpus_ERTMP = AudioCodecId('O'<<24 | 'p'<<16 | 'u'<<8 | 's')
	AudioCodecIdMP3_ERTMP  = AudioCodecId('.'<<24 | 'm'<<16 | 'p'<<8 | '3')
	AudioCodecIdFlac_ERTMP = AudioCodecId('f'<<24 | 'L'<<16 | 'a'<<8 | 'C')
	AudioCodecIdAAC_ERTMP  = AudioCodecId('m'<<24 | 'p'<<16 | '4'<<8 | 'a')
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

type AudioTrack struct {
	TrackId uint8
	CodecId AudioCodecId
	Payload []byte

	// Multichannel config (when PacketType == ERTMPAudioPacketTypeMultichannelConfig)
	ChannelOrder   ERTMPAudioChannelOrder
	ChannelCount   uint8
	ChannelMapping []ERTMPAudioChannel
	ChannelFlags   ERTMPAudioChannelMask
}

type AudioMessage struct {
	MetadataFields
	PacketType ERTMPAudioPacketType
	Tracks     []AudioTrack

	// Legacy fields (only relevant when codec ID <= 15)
	Rate       AudioRate
	SampleSize AudioSize
	Stereo     bool

	// E-RTMP multitrack
	MultitrackType ERTMPAvMultitrackType

	// ModEx
	TimestampNanoOffset uint32
}

type ERTMPAudioPacketType uint8

//go:generate stringer -type=ERTMPAudioPacketType -trimprefix=ERTMPAudioPacketType
const (
	ERTMPAudioPacketTypeSequenceStart      ERTMPAudioPacketType = 0
	ERTMPAudioPacketTypeCodedFrames        ERTMPAudioPacketType = 1
	ERTMPAudioPacketTypeSequenceEnd        ERTMPAudioPacketType = 2
	ERTMPAudioPacketTypeMultichannelConfig ERTMPAudioPacketType = 4
	ERTMPAudioPacketTypeMultitrack         ERTMPAudioPacketType = 5
	ERTMPAudioPacketTypeModEx              ERTMPAudioPacketType = 7
)

type ERTMPAudioPacketModExType uint8

//go:generate stringer -type=ERTMPAudioPacketModExType -trimprefix=ERTMPAudioPacketModExType
const (
	ERTMPAudioPacketModExTypeTimestampOffsetNano ERTMPAudioPacketModExType = 0
)

type ERTMPAvMultitrackType uint8

//go:generate stringer -type=ERTMPAvMultitrackType -trimprefix=ERTMPAvMultitrackType
const (
	ERTMPAvMultitrackTypeOneTrack             ERTMPAvMultitrackType = 0
	ERTMPAvMultitrackTypeManyTracks           ERTMPAvMultitrackType = 1
	ERTMPAvMultitrackTypeManyTracksManyCodecs ERTMPAvMultitrackType = 2
)

type ERTMPAudioChannelOrder uint8

//go:generate stringer -type=ERTMPAudioChannelOrder -trimprefix=ERTMPAudioChannelOrder
const (
	ERTMPAudioChannelOrderUnspecified ERTMPAudioChannelOrder = 0
	ERTMPAudioChannelOrderNative      ERTMPAudioChannelOrder = 1
	ERTMPAudioChannelOrderCustom      ERTMPAudioChannelOrder = 2
)

type ERTMPAudioChannel uint8

//go:generate stringer -type=ERTMPAudioChannel -trimprefix=ERTMPAudioChannel
const (
	ERTMPAudioChannelFrontLeft         ERTMPAudioChannel = iota // 0
	ERTMPAudioChannelFrontRight                                 // 1
	ERTMPAudioChannelFrontCenter                                // 2
	ERTMPAudioChannelLowFrequency1                              // 3
	ERTMPAudioChannelBackLeft                                   // 4
	ERTMPAudioChannelBackRight                                  // 5
	ERTMPAudioChannelFrontLeftCenter                            // 6
	ERTMPAudioChannelFrontRightCenter                           // 7
	ERTMPAudioChannelBackCenter                                 // 8
	ERTMPAudioChannelSideLeft                                   // 9
	ERTMPAudioChannelSideRight                                  // 10
	ERTMPAudioChannelTopCenter                                  // 11
	ERTMPAudioChannelTopFrontLeft                               // 12
	ERTMPAudioChannelTopFrontCenter                             // 13
	ERTMPAudioChannelTopFrontRight                              // 14
	ERTMPAudioChannelTopBackLeft                                // 15
	ERTMPAudioChannelTopBackCenter                              // 16
	ERTMPAudioChannelTopBackRight                               // 17
	ERTMPAudioChannelLowFrequency2                              // 18
	ERTMPAudioChannelTopSideLeft                                // 19
	ERTMPAudioChannelTopSideRight                               // 20
	ERTMPAudioChannelBottomFrontCenter                          // 21
	ERTMPAudioChannelBottomFrontLeft                            // 22
	ERTMPAudioChannelBottomFrontRight                           // 23

	ERTMPAudioChannelUnused  ERTMPAudioChannel = 0xfe
	ERTMPAudioChannelUnknown ERTMPAudioChannel = 0xff
)

type ERTMPAudioChannelMask uint32

//go:generate stringer -type=ERTMPAudioChannelMask -trimprefix=ERTMPAudioChannelMask
const (
	ERTMPAudioChannelMaskFrontLeft         ERTMPAudioChannelMask = 0x000001
	ERTMPAudioChannelMaskFrontRight        ERTMPAudioChannelMask = 0x000002
	ERTMPAudioChannelMaskFrontCenter       ERTMPAudioChannelMask = 0x000004
	ERTMPAudioChannelMaskLowFrequency1     ERTMPAudioChannelMask = 0x000008
	ERTMPAudioChannelMaskBackLeft          ERTMPAudioChannelMask = 0x000010
	ERTMPAudioChannelMaskBackRight         ERTMPAudioChannelMask = 0x000020
	ERTMPAudioChannelMaskFrontLeftCenter   ERTMPAudioChannelMask = 0x000040
	ERTMPAudioChannelMaskFrontRightCenter  ERTMPAudioChannelMask = 0x000080
	ERTMPAudioChannelMaskBackCenter        ERTMPAudioChannelMask = 0x000100
	ERTMPAudioChannelMaskSideLeft          ERTMPAudioChannelMask = 0x000200
	ERTMPAudioChannelMaskSideRight         ERTMPAudioChannelMask = 0x000400
	ERTMPAudioChannelMaskTopCenter         ERTMPAudioChannelMask = 0x000800
	ERTMPAudioChannelMaskTopFrontLeft      ERTMPAudioChannelMask = 0x001000
	ERTMPAudioChannelMaskTopFrontCenter    ERTMPAudioChannelMask = 0x002000
	ERTMPAudioChannelMaskTopFrontRight     ERTMPAudioChannelMask = 0x004000
	ERTMPAudioChannelMaskTopBackLeft       ERTMPAudioChannelMask = 0x008000
	ERTMPAudioChannelMaskTopBackCenter     ERTMPAudioChannelMask = 0x010000
	ERTMPAudioChannelMaskTopBackRight      ERTMPAudioChannelMask = 0x020000
	ERTMPAudioChannelMaskLowFrequency2     ERTMPAudioChannelMask = 0x040000
	ERTMPAudioChannelMaskTopSideLeft       ERTMPAudioChannelMask = 0x080000
	ERTMPAudioChannelMaskTopSideRight      ERTMPAudioChannelMask = 0x100000
	ERTMPAudioChannelMaskBottomFrontCenter ERTMPAudioChannelMask = 0x200000
	ERTMPAudioChannelMaskBottomFrontLeft   ERTMPAudioChannelMask = 0x400000
	ERTMPAudioChannelMaskBottomFrontRight  ERTMPAudioChannelMask = 0x800000
)

func init() { RegisterType(new(AudioMessage)) }

func (m AudioMessage) Type() Type {
	return TypeAudioMessage
}

func (m AudioMessage) IsERTMP() bool {
	if len(m.Tracks) > 0 {
		return m.Tracks[0].CodecId > 15
	}
	return false
}

func (m *AudioMessage) Unmarshal(data []byte) error {
	if len(data) < 1 {
		return ErrShortMessage
	}
	soundFormat := AudioCodecId(data[0] >> 4)
	if soundFormat != audioCodecIdExtendedHeader {
		return m.unmarshalLegacy(data)
	}
	return m.unmarshalERTMP(data)
}

func (m *AudioMessage) unmarshalLegacy(data []byte) error {
	headerByte := data[0]
	codecId := AudioCodecId(headerByte >> 4)
	m.Rate = AudioRate(headerByte >> 2 & 0x03)
	m.SampleSize = AudioSize(headerByte >> 1 & 0x01)
	m.Stereo = headerByte&0x01 != 0
	data = data[1:]

	packetType := ERTMPAudioPacketTypeCodedFrames
	if codecId == AudioCodecIdAAC {
		if len(data) < 1 {
			return ErrShortMessage
		}
		packetType = ERTMPAudioPacketType(data[0])
		data = data[1:]
	}

	m.PacketType = packetType
	m.Tracks = []AudioTrack{{
		CodecId: codecId,
		Payload: data,
	}}
	return nil
}

func (m *AudioMessage) unmarshalERTMP(data []byte) error {
	// First byte: [soundFormat=9:4][audioPacketType:4]
	audioPacketType := ERTMPAudioPacketType(data[0] & 0x0F)
	pos := 1

	// Process ModEx loop
	for audioPacketType == ERTMPAudioPacketTypeModEx {
		if pos >= len(data) {
			return ErrShortMessage
		}
		modExDataSize := int(data[pos]) + 1
		pos++

		if modExDataSize == 256 {
			if pos+2 > len(data) {
				return ErrShortMessage
			}
			modExDataSize = int(binary.BigEndian.Uint16(data[pos:])) + 1
			pos += 2
		}

		if pos+modExDataSize >= len(data) {
			return ErrShortMessage
		}
		modExData := data[pos : pos+modExDataSize]
		pos += modExDataSize

		audioPacketModExType := ERTMPAudioPacketModExType(data[pos] >> 4)
		audioPacketType = ERTMPAudioPacketType(data[pos] & 0x0F)
		pos++

		if audioPacketModExType == ERTMPAudioPacketModExTypeTimestampOffsetNano && len(modExData) >= 3 {
			m.TimestampNanoOffset = uint32(modExData[0])<<16 | uint32(modExData[1])<<8 | uint32(modExData[2])
		}
	}

	m.PacketType = audioPacketType

	isMultitrack := false
	var sharedCodecId AudioCodecId

	if audioPacketType == ERTMPAudioPacketTypeMultitrack {
		isMultitrack = true
		if pos >= len(data) {
			return ErrShortMessage
		}
		m.MultitrackType = ERTMPAvMultitrackType(data[pos] >> 4)
		m.PacketType = ERTMPAudioPacketType(data[pos] & 0x0F)
		pos++

		if m.MultitrackType != ERTMPAvMultitrackTypeManyTracksManyCodecs {
			if pos+4 > len(data) {
				return ErrShortMessage
			}
			sharedCodecId = AudioCodecId(binary.BigEndian.Uint32(data[pos:]))
			pos += 4
		}
	} else {
		if pos+4 > len(data) {
			return ErrShortMessage
		}
		sharedCodecId = AudioCodecId(binary.BigEndian.Uint32(data[pos:]))
		pos += 4
	}

	// Body loop
	m.Tracks = nil
	for {
		track := AudioTrack{}

		if isMultitrack && m.MultitrackType == ERTMPAvMultitrackTypeManyTracksManyCodecs {
			if pos+4 > len(data) {
				return ErrShortMessage
			}
			track.CodecId = AudioCodecId(binary.BigEndian.Uint32(data[pos:]))
			pos += 4
		} else {
			track.CodecId = sharedCodecId
		}

		trackEnd := len(data)
		if isMultitrack {
			if pos >= len(data) {
				return ErrShortMessage
			}
			track.TrackId = data[pos]
			pos++

			if m.MultitrackType != ERTMPAvMultitrackTypeOneTrack {
				if pos+3 > len(data) {
					return ErrShortMessage
				}
				trackSize := int(data[pos])<<16 | int(data[pos+1])<<8 | int(data[pos+2])
				pos += 3
				trackEnd = pos + trackSize
				if trackEnd > len(data) {
					return ErrShortMessage
				}
			}
		}

		switch m.PacketType {
		case ERTMPAudioPacketTypeMultichannelConfig:
			if pos+2 > trackEnd {
				return ErrShortMessage
			}
			track.ChannelOrder = ERTMPAudioChannelOrder(data[pos])
			pos++
			track.ChannelCount = data[pos]
			pos++

			if track.ChannelOrder == ERTMPAudioChannelOrderCustom {
				if pos+int(track.ChannelCount) > trackEnd {
					return ErrShortMessage
				}
				track.ChannelMapping = make([]ERTMPAudioChannel, track.ChannelCount)
				for i := range track.ChannelMapping {
					track.ChannelMapping[i] = ERTMPAudioChannel(data[pos])
					pos++
				}
			}

			if track.ChannelOrder == ERTMPAudioChannelOrderNative {
				if pos+4 > trackEnd {
					return ErrShortMessage
				}
				track.ChannelFlags = ERTMPAudioChannelMask(binary.BigEndian.Uint32(data[pos:]))
				pos += 4
			}

		case ERTMPAudioPacketTypeSequenceEnd:
			// No payload

		case ERTMPAudioPacketTypeSequenceStart, ERTMPAudioPacketTypeCodedFrames:
			track.Payload = data[pos:trackEnd]
			pos = trackEnd
		}

		m.Tracks = append(m.Tracks, track)

		if isMultitrack && m.MultitrackType != ERTMPAvMultitrackTypeOneTrack && pos < len(data) {
			continue
		}
		break
	}

	return nil
}

func (m AudioMessage) Marshal() ([]byte, error) {
	if len(m.Tracks) == 0 {
		return nil, fmt.Errorf("audio message has no tracks")
	}
	if !m.IsERTMP() {
		return m.marshalLegacy()
	}
	return m.marshalERTMP()
}

func (m AudioMessage) marshalLegacy() ([]byte, error) {
	track := m.Tracks[0]
	headerByte := byte(track.CodecId)<<4 | byte(m.Rate)<<2 | byte(m.SampleSize)<<1
	if m.Stereo {
		headerByte |= 0x01
	}
	out := []byte{headerByte}
	if track.CodecId == AudioCodecIdAAC {
		out = append(out, byte(m.PacketType))
	}
	out = append(out, track.Payload...)
	return out, nil
}

func (m AudioMessage) marshalERTMP() ([]byte, error) {
	out := make([]byte, 0, 256)

	isMultitrack := len(m.Tracks) > 1 || m.MultitrackType != 0

	packetType := m.PacketType
	if isMultitrack {
		packetType = ERTMPAudioPacketTypeMultitrack
	}

	// ModEx: TimestampNanoOffset
	if m.TimestampNanoOffset > 0 {
		out = append(out, byte(audioCodecIdExtendedHeader)<<4|byte(ERTMPAudioPacketTypeModEx))

		modExData := []byte{
			byte(m.TimestampNanoOffset >> 16),
			byte(m.TimestampNanoOffset >> 8),
			byte(m.TimestampNanoOffset),
		}
		out = append(out, byte(len(modExData)-1)) // UI8: size - 1
		out = append(out, modExData...)

		// [AudioPacketModExType:4][packetType:4]
		out = append(out, byte(ERTMPAudioPacketModExTypeTimestampOffsetNano)<<4|byte(packetType))
	} else {
		out = append(out, byte(audioCodecIdExtendedHeader)<<4|byte(packetType))
	}

	if isMultitrack {
		// [AvMultitrackType:4][AudioPacketType:4]
		out = append(out, byte(m.MultitrackType)<<4|byte(m.PacketType))

		if m.MultitrackType != ERTMPAvMultitrackTypeManyTracksManyCodecs {
			fourcc := make([]byte, 4)
			binary.BigEndian.PutUint32(fourcc, uint32(m.Tracks[0].CodecId))
			out = append(out, fourcc...)
		}
	} else {
		fourcc := make([]byte, 4)
		binary.BigEndian.PutUint32(fourcc, uint32(m.Tracks[0].CodecId))
		out = append(out, fourcc...)
	}

	// Body: per-track data
	for _, track := range m.Tracks {
		if isMultitrack && m.MultitrackType == ERTMPAvMultitrackTypeManyTracksManyCodecs {
			fourcc := make([]byte, 4)
			binary.BigEndian.PutUint32(fourcc, uint32(track.CodecId))
			out = append(out, fourcc...)
		}

		if isMultitrack {
			out = append(out, track.TrackId)

			if m.MultitrackType != ERTMPAvMultitrackTypeOneTrack {
				// Reserve 3 bytes for track size, fill in after writing body
				sizePos := len(out)
				out = append(out, 0, 0, 0)
				bodyStart := len(out)
				out = m.marshalTrackBody(out, track)
				trackSize := len(out) - bodyStart
				out[sizePos] = byte(trackSize >> 16)
				out[sizePos+1] = byte(trackSize >> 8)
				out[sizePos+2] = byte(trackSize)
				continue
			}
		}

		out = m.marshalTrackBody(out, track)
	}

	return out, nil
}

func (m AudioMessage) marshalTrackBody(out []byte, track AudioTrack) []byte {
	switch m.PacketType {
	case ERTMPAudioPacketTypeMultichannelConfig:
		out = append(out, byte(track.ChannelOrder))
		out = append(out, track.ChannelCount)
		if track.ChannelOrder == ERTMPAudioChannelOrderCustom {
			for _, ch := range track.ChannelMapping {
				out = append(out, byte(ch))
			}
		}
		if track.ChannelOrder == ERTMPAudioChannelOrderNative {
			flags := make([]byte, 4)
			binary.BigEndian.PutUint32(flags, uint32(track.ChannelFlags))
			out = append(out, flags...)
		}
	case ERTMPAudioPacketTypeSequenceEnd:
		// No payload
	case ERTMPAudioPacketTypeSequenceStart, ERTMPAudioPacketTypeCodedFrames:
		out = append(out, track.Payload...)
	}
	return out
}

func (m AudioMessage) String() string {
	if len(m.Tracks) == 0 {
		return fmt.Sprintf("%v: %+v (no tracks)", m.Type(), m.MetadataFields)
	}
	if !m.IsERTMP() {
		track := m.Tracks[0]
		return fmt.Sprintf("%v: %+v codec:%v rate:%v sample:%v stereo:%v pkt:%v payload:%v bytes",
			m.Type(), m.MetadataFields, track.CodecId, m.Rate, m.SampleSize, m.Stereo,
			m.PacketType, len(track.Payload))
	}
	return fmt.Sprintf("%v: %+v pkt:%v tracks:%d multitrack:%v payload:%v bytes",
		m.Type(), m.MetadataFields, m.PacketType, len(m.Tracks),
		m.MultitrackType, len(m.Tracks[0].Payload))
}
