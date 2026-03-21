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
	VideoFrameTypeVideoCommand
)

type VideoCommand uint8

//go:generate stringer -type=VideoCommand -trimprefix=VideoCommand
const (
	VideoCommmandStartSeek VideoCommand = iota
	VideoCommandEndSeek
)

type VideoCodecId uint32

//go:generate stringer -type=VideoCodecId -trimprefix=VideoCodecId
const (
	// Legacy Codec IDs
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

	// E-RTMP video codecs
	VideoCodecIdVP8_ERTMP  = VideoCodecId('v'<<24 | 'p'<<16 | '0'<<8 | '8')
	VideoCodecIdVP9_ERTMP  = VideoCodecId('v'<<24 | 'p'<<16 | '0'<<8 | '9')
	VideoCodecIdAV1_ERTMP  = VideoCodecId('a'<<24 | 'v'<<16 | '0'<<8 | '1')
	VideoCodecIdAvc_ERTMP  = VideoCodecId('a'<<24 | 'v'<<16 | 'c'<<8 | '1')
	VideoCodecIdHevc_ERTMP = VideoCodecId('h'<<24 | 'v'<<16 | 'c'<<8 | '1')
	VideoCodecIdVVC_ERTMP  = VideoCodecId('v'<<24 | 'v'<<16 | 'c'<<8 | '1')
)

type AvcPacketType uint8

//go:generate stringer -type=AvcPacketType -trimprefix=AvcPacketType
const (
	AvcPacketTypeSequenceHeader AvcPacketType = iota
	AvcPacketTypeNALU
	AvcPacketTypeEndOfSequence
)

type ERTMPVideoPacketType uint8

//go:generate stringer -type=ERTMPVideoPacketType -trimprefix=ERTMPVideoPacketType
const (
	ERTMPVideoPacketTypeSequenceStart        ERTMPVideoPacketType = 0
	ERTMPVideoPacketTypeCodedFrames          ERTMPVideoPacketType = 1
	ERTMPVideoPacketTypeSequenceEnd          ERTMPVideoPacketType = 2
	ERTMPVideoPacketTypeCodedFramesX         ERTMPVideoPacketType = 3
	ERTMPVideoPacketTypeMetadata             ERTMPVideoPacketType = 4
	ERTMPVideoPacketTypeMPEG2TSSequenceStart ERTMPVideoPacketType = 5
	ERTMPVideoPacketTypeMultitrack           ERTMPVideoPacketType = 6
	ERTMPVideoPacketTypeModEx                ERTMPVideoPacketType = 7
)

type ERTMPVideoPacketModExType uint8

//go:generate stringer -type=ERTMPVideoPacketModExType -trimprefix=ERTMPVideoPacketModExType
const (
	ERTMPVideoPacketModExTypeTimestampOffsetNano ERTMPVideoPacketModExType = 0
)

type VideoTrack struct {
	TrackId         uint8
	CodecId         VideoCodecId
	CompositionTime int32
	Payload         []byte
}

type VideoMessage struct {
	MetadataFields
	FrameType  VideoFrameType
	PacketType ERTMPVideoPacketType
	Tracks     []VideoTrack

	// E-RTMP multitrack
	MultitrackType ERTMPAvMultitrackType

	// ModEx
	TimestampNanoOffset uint32

	// Video command (when FrameType == VideoFrameTypeVideoCommand)
	Command VideoCommand

	// VideoMetadata holds AMF-encoded metadata (when PacketType == ERTMPVideoPacketTypeMetadata)
	VideoMetadata []byte
}

func init() { RegisterType(new(VideoMessage)) }

func (m VideoMessage) Type() Type {
	return TypeVideoMessage
}

func (m VideoMessage) IsERTMP() bool {
	if len(m.Tracks) > 0 {
		return m.Tracks[0].CodecId > 15
	}
	// Metadata-only messages with no tracks are E-RTMP
	if m.PacketType == ERTMPVideoPacketTypeMetadata && len(m.VideoMetadata) > 0 {
		return true
	}
	return false
}

func (m *VideoMessage) Unmarshal(data []byte) error {
	if len(data) < 1 {
		return ErrShortMessage
	}
	isExVideoHeader := data[0]&0x80 != 0
	m.FrameType = VideoFrameType((data[0] >> 4) & 0x07)

	if !isExVideoHeader {
		return m.unmarshalLegacy(data)
	}
	return m.unmarshalERTMP(data)
}

func (m *VideoMessage) unmarshalLegacy(data []byte) error {
	codecId := VideoCodecId(data[0] & 0x0F)
	data = data[1:]

	if m.FrameType == VideoFrameTypeVideoCommand {
		if len(data) < 1 {
			return ErrShortMessage
		}
		m.Command = VideoCommand(data[0])
		return nil
	}

	packetType := ERTMPVideoPacketTypeCodedFrames
	var compositionTime int32

	if codecId == VideoCodecIdAvc {
		if len(data) < 4 {
			return ErrShortMessage
		}
		packetType = ERTMPVideoPacketType(data[0])
		compositionTime = int32(binary.BigEndian.Uint32(data)) & 0xFF_FF_FF
		if compositionTime > 0x7F_FF_FF {
			compositionTime -= 0x1_00_00_00
		}
		data = data[4:]
	}

	m.PacketType = packetType
	m.Tracks = []VideoTrack{{
		CodecId:         codecId,
		CompositionTime: compositionTime,
		Payload:         data,
	}}
	return nil
}

func (m *VideoMessage) unmarshalERTMP(data []byte) error {
	// First byte: [isExVideoHeader:1][videoFrameType:3][videoPacketType:4]
	videoPacketType := ERTMPVideoPacketType(data[0] & 0x0F)
	pos := 1

	// Process ModEx loop
	for videoPacketType == ERTMPVideoPacketTypeModEx {
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

		videoPacketModExType := ERTMPVideoPacketModExType(data[pos] >> 4)
		videoPacketType = ERTMPVideoPacketType(data[pos] & 0x0F)
		pos++

		if videoPacketModExType == ERTMPVideoPacketModExTypeTimestampOffsetNano && len(modExData) >= 3 {
			m.TimestampNanoOffset = uint32(modExData[0])<<16 | uint32(modExData[1])<<8 | uint32(modExData[2])
		}
	}

	m.PacketType = videoPacketType

	// Handle VideoCommand (with ExVideoTagHeader)
	if videoPacketType != ERTMPVideoPacketTypeMetadata && m.FrameType == VideoFrameTypeVideoCommand {
		if pos >= len(data) {
			return ErrShortMessage
		}
		m.Command = VideoCommand(data[pos])
		return nil
	}

	isMultitrack := false
	var sharedCodecId VideoCodecId

	switch videoPacketType {
	case ERTMPVideoPacketTypeMultitrack:
		isMultitrack = true
		if pos >= len(data) {
			return ErrShortMessage
		}
		m.MultitrackType = ERTMPAvMultitrackType(data[pos] >> 4)
		m.PacketType = ERTMPVideoPacketType(data[pos] & 0x0F)
		pos++

		if m.MultitrackType != ERTMPAvMultitrackTypeManyTracksManyCodecs {
			if pos+4 > len(data) {
				return ErrShortMessage
			}
			sharedCodecId = VideoCodecId(binary.BigEndian.Uint32(data[pos:]))
			pos += 4
		}
	case ERTMPVideoPacketTypeMetadata:
		// Metadata has no FOURCC, body is AMF-encoded metadata
		m.VideoMetadata = data[pos:]
		return nil
	default:
		if pos+4 > len(data) {
			return ErrShortMessage
		}
		sharedCodecId = VideoCodecId(binary.BigEndian.Uint32(data[pos:]))
		pos += 4
	}

	// Body loop
	m.Tracks = nil
	for {
		track := VideoTrack{}

		if isMultitrack && m.MultitrackType == ERTMPAvMultitrackTypeManyTracksManyCodecs {
			if pos+4 > len(data) {
				return ErrShortMessage
			}
			track.CodecId = VideoCodecId(binary.BigEndian.Uint32(data[pos:]))
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
		case ERTMPVideoPacketTypeSequenceEnd:
			// No payload

		case ERTMPVideoPacketTypeMetadata:
			track.Payload = data[pos:trackEnd]
			pos = trackEnd

		case ERTMPVideoPacketTypeCodedFrames:
			// CodedFrames for AVC/HEVC includes SI24 compositionTimeOffset
			if track.CodecId == VideoCodecIdAvc_ERTMP || track.CodecId == VideoCodecIdHevc_ERTMP {
				if pos+3 > trackEnd {
					return ErrShortMessage
				}
				ct := int32(data[pos])<<16 | int32(data[pos+1])<<8 | int32(data[pos+2])
				if ct > 0x7F_FF_FF {
					ct -= 0x1_00_00_00
				}
				track.CompositionTime = ct
				pos += 3
			}
			track.Payload = data[pos:trackEnd]
			pos = trackEnd

		case ERTMPVideoPacketTypeCodedFramesX:
			// CompositionTime implicitly 0
			track.Payload = data[pos:trackEnd]
			pos = trackEnd

		case ERTMPVideoPacketTypeSequenceStart, ERTMPVideoPacketTypeMPEG2TSSequenceStart:
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

func (m VideoMessage) Marshal() ([]byte, error) {
	if !m.IsERTMP() {
		return m.marshalLegacy()
	}
	return m.marshalERTMP()
}

func (m VideoMessage) marshalLegacy() ([]byte, error) {
	if m.FrameType == VideoFrameTypeVideoCommand {
		out := []byte{byte(m.FrameType) << 4, byte(m.Command)}
		return out, nil
	}

	if len(m.Tracks) == 0 {
		return nil, fmt.Errorf("video message has no tracks")
	}
	track := m.Tracks[0]

	out := make([]byte, 1)
	out[0] = byte(m.FrameType)<<4 | byte(track.CodecId&0x0F)
	if track.CodecId == VideoCodecIdAvc {
		out = binary.BigEndian.AppendUint32(out, uint32(track.CompositionTime))
		out[1] = byte(m.PacketType)
	}
	out = append(out, track.Payload...)
	return out, nil
}

func (m VideoMessage) marshalERTMP() ([]byte, error) {
	out := make([]byte, 0, 256)

	isMultitrack := len(m.Tracks) > 1 || m.MultitrackType != 0

	packetType := m.PacketType
	if isMultitrack {
		packetType = ERTMPVideoPacketTypeMultitrack
	}

	// ModEx: TimestampNanoOffset
	if m.TimestampNanoOffset > 0 {
		headerByte := byte(0x80) | byte(m.FrameType&0x07)<<4 | byte(ERTMPVideoPacketTypeModEx)
		out = append(out, headerByte)

		modExData := []byte{
			byte(m.TimestampNanoOffset >> 16),
			byte(m.TimestampNanoOffset >> 8),
			byte(m.TimestampNanoOffset),
		}
		out = append(out, byte(len(modExData)-1)) // UI8: size - 1
		out = append(out, modExData...)

		// [VideoPacketModExType:4][packetType:4]
		out = append(out, byte(ERTMPVideoPacketModExTypeTimestampOffsetNano)<<4|byte(packetType))
	} else {
		headerByte := byte(0x80) | byte(m.FrameType&0x07)<<4 | byte(packetType)
		out = append(out, headerByte)
	}

	// Handle command frame
	if packetType != ERTMPVideoPacketTypeMetadata && m.FrameType == VideoFrameTypeVideoCommand {
		out = append(out, byte(m.Command))
		return out, nil
	}

	// Metadata: no FOURCC, just AMF-encoded body
	if m.PacketType == ERTMPVideoPacketTypeMetadata && !isMultitrack {
		out = append(out, m.VideoMetadata...)
		return out, nil
	}

	if len(m.Tracks) == 0 {
		return nil, fmt.Errorf("video message has no tracks")
	}

	if isMultitrack {
		// [AvMultitrackType:4][VideoPacketType:4]
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

func (m VideoMessage) marshalTrackBody(out []byte, track VideoTrack) []byte {
	switch m.PacketType {
	case ERTMPVideoPacketTypeSequenceEnd:
		// No payload

	case ERTMPVideoPacketTypeMetadata:
		out = append(out, track.Payload...)

	case ERTMPVideoPacketTypeCodedFrames:
		// AVC/HEVC include SI24 compositionTimeOffset
		if track.CodecId == VideoCodecIdAvc_ERTMP || track.CodecId == VideoCodecIdHevc_ERTMP {
			ct := uint32(track.CompositionTime) & 0xFF_FF_FF
			out = append(out, byte(ct>>16), byte(ct>>8), byte(ct))
		}
		out = append(out, track.Payload...)

	case ERTMPVideoPacketTypeCodedFramesX:
		// CompositionTime implicitly 0
		out = append(out, track.Payload...)

	case ERTMPVideoPacketTypeSequenceStart, ERTMPVideoPacketTypeMPEG2TSSequenceStart:
		out = append(out, track.Payload...)
	}
	return out
}

func (m VideoMessage) String() string {
	if len(m.Tracks) == 0 {
		if m.FrameType == VideoFrameTypeVideoCommand {
			return fmt.Sprintf("%v: %+v frame:%v cmd:%v", m.Type(), m.MetadataFields,
				m.FrameType, m.Command)
		}
		if m.PacketType == ERTMPVideoPacketTypeMetadata {
			return fmt.Sprintf("%v: %+v frame:%v pkt:Metadata metadata:%d bytes",
				m.Type(), m.MetadataFields, m.FrameType, len(m.VideoMetadata))
		}
		return fmt.Sprintf("%v: %+v (no tracks)", m.Type(), m.MetadataFields)
	}
	if !m.IsERTMP() {
		track := m.Tracks[0]
		return fmt.Sprintf("%v: %+v frame:%v codec:%v pkt:%v ct:%v payload:%v bytes",
			m.Type(), m.MetadataFields, m.FrameType, track.CodecId,
			m.PacketType, track.CompositionTime, len(track.Payload))
	}
	return fmt.Sprintf("%v: %+v frame:%v pkt:%v tracks:%d multitrack:%v payload:%v bytes",
		m.Type(), m.MetadataFields, m.FrameType, m.PacketType,
		len(m.Tracks), m.MultitrackType, len(m.Tracks[0].Payload))
}
