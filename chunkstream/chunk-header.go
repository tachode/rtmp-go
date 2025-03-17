package chunkstream

import (
	"io"

	"github.com/tachode/rtmp-go/message"
)

const ExtendedTimestampMarker = 0xFFFFFF

type HeaderType uint8

//go:generate stringer -type=HeaderType -trimprefix=HeaderType
const (
	HeaderTypeFull HeaderType = iota
	HeaderTypeSameStream
	HeaderTypeSameStreamAndLength
	HeaderTypeContinuation
)

type ChunkHeader struct {
	Type             HeaderType
	ChunkStreamId    uint32
	Timestamp        uint32
	MessageLength    uint32
	MessageType      message.Type
	MessageStreamId  uint32
	TimestampIsDelta bool
}

func (h ChunkHeader) Write(w io.Writer) (n int, err error) {
	var m int
	switch {
	case h.ChunkStreamId < 64:
		// Chunk basic header 1: 1 byte
		m, err = w.Write([]byte{byte(h.Type<<6) | byte(h.ChunkStreamId)})
		n += m
		if err != nil {
			return
		}
	case h.ChunkStreamId < 320:
		// Chunk basic header 2: 2 bytes
		m, err = w.Write([]byte{byte(h.Type << 6), byte(h.ChunkStreamId - 64)})
		n += m
		if err != nil {
			return
		}
	case h.ChunkStreamId <= 65599:
		// Chunk basic header 3: 3 bytes
		m, err = w.Write([]byte{byte(h.Type<<6) | 1, byte((h.ChunkStreamId - 64) >> 8), byte(h.ChunkStreamId - 64)})
		n += m
		if err != nil {
			return
		}
	default:
		return n, ErrInvalidChunkStreamId
	}

	extendedTimestamp := h.Timestamp >= ExtendedTimestampMarker
	timestampField := h.Timestamp
	if extendedTimestamp {
		timestampField = ExtendedTimestampMarker
	}

	switch h.Type {
	case HeaderTypeFull:
		// Chunk message header 1: 11 bytes
		if h.TimestampIsDelta {
			return n, ErrDeltaTimePassedToFullHeader
		}
		m, err = w.Write([]byte{
			byte(timestampField >> 16),
			byte(timestampField >> 8),
			byte(timestampField),
			byte(h.MessageLength >> 16),
			byte(h.MessageLength >> 8),
			byte(h.MessageLength),
			byte(h.MessageType),
			// Stream ID is little-endian
			byte(h.MessageStreamId),
			byte(h.MessageStreamId >> 8),
			byte(h.MessageStreamId >> 16),
			byte(h.MessageStreamId >> 24),
		})
		n += m
		if err != nil {
			return
		}
	case HeaderTypeSameStream:
		// Chunk message header 2: 7 bytes
		if !h.TimestampIsDelta {
			return n, ErrNonDeltaTimestampPassedToShortHeader
		}
		m, err = w.Write([]byte{
			byte(timestampField >> 16),
			byte(timestampField >> 8),
			byte(timestampField),
			byte(h.MessageLength >> 16),
			byte(h.MessageLength >> 8),
			byte(h.MessageLength),
			byte(h.MessageType),
		})
		n += m
		if err != nil {
			return
		}
	case HeaderTypeSameStreamAndLength:
		// Chunk message header 3: 3 bytes
		if !h.TimestampIsDelta {
			return n, ErrNonDeltaTimestampPassedToShortHeader
		}
		m, err = w.Write([]byte{
			byte(timestampField >> 16),
			byte(timestampField >> 8),
			byte(timestampField),
		})
		n += m
		if err != nil {
			return
		}
	case HeaderTypeContinuation:
		// Chunk message header 4: 0 bytes
	}

	if extendedTimestamp {
		m, err = w.Write([]byte{
			byte(h.Timestamp >> 24),
			byte(h.Timestamp >> 16),
			byte(h.Timestamp >> 8),
			byte(h.Timestamp),
		})
		n += m
		if err != nil {
			return
		}
	}

	return
}

func (h *ChunkHeader) Read(r io.Reader) (n int, err error) {
	var m int
	var basicHeader [3]byte
	if m, err = r.Read(basicHeader[:1]); err != nil {
		return
	}
	n += m
	h.Type = HeaderType(basicHeader[0] >> 6)
	h.ChunkStreamId = uint32(basicHeader[0] & 0x3F)
	if h.ChunkStreamId == 0 {
		if m, err = r.Read(basicHeader[1:2]); err != nil {
			return
		}
		n += m
		h.ChunkStreamId = 64 + uint32(basicHeader[1])
	} else if h.ChunkStreamId == 1 {
		if m, err = r.Read(basicHeader[1:3]); err != nil {
			return
		}
		n += m
		h.ChunkStreamId = 64 + uint32(basicHeader[1])<<8 + uint32(basicHeader[2])
	}

	switch h.Type {
	case HeaderTypeFull:
		h.TimestampIsDelta = false
		var fullHeader [11]byte
		if m, err = r.Read(fullHeader[:]); err != nil {
			return
		}
		n += m
		h.Timestamp = uint32(fullHeader[0])<<16 | uint32(fullHeader[1])<<8 | uint32(fullHeader[2])
		h.MessageLength = uint32(fullHeader[3])<<16 | uint32(fullHeader[4])<<8 | uint32(fullHeader[5])
		h.MessageType = message.Type(fullHeader[6])
		// Little-endian
		h.MessageStreamId = uint32(fullHeader[7]) | uint32(fullHeader[8])<<8 | uint32(fullHeader[9])<<16 | uint32(fullHeader[10])<<24
	case HeaderTypeSameStream:
		h.TimestampIsDelta = true
		var sameStreamHeader [7]byte
		if m, err = r.Read(sameStreamHeader[:]); err != nil {
			return
		}
		n += m
		h.Timestamp = uint32(sameStreamHeader[0])<<16 | uint32(sameStreamHeader[1])<<8 | uint32(sameStreamHeader[2])
		h.MessageLength = uint32(sameStreamHeader[3])<<16 | uint32(sameStreamHeader[4])<<8 | uint32(sameStreamHeader[5])
		h.MessageType = message.Type(sameStreamHeader[6])
	case HeaderTypeSameStreamAndLength:
		h.TimestampIsDelta = true
		var sameLengthAndStreamHeader [3]byte
		if m, err = r.Read(sameLengthAndStreamHeader[:]); err != nil {
			return
		}
		n += m
		h.Timestamp = uint32(sameLengthAndStreamHeader[0])<<16 | uint32(sameLengthAndStreamHeader[1])<<8 | uint32(sameLengthAndStreamHeader[2])
	case HeaderTypeContinuation:
		h.TimestampIsDelta = true
	}

	if h.Timestamp == ExtendedTimestampMarker {
		var extendedTimestamp [4]byte
		if m, err = r.Read(extendedTimestamp[:]); err != nil {
			return
		}
		n += m
		h.Timestamp = uint32(extendedTimestamp[0])<<24 | uint32(extendedTimestamp[1])<<16 | uint32(extendedTimestamp[2])<<8 | uint32(extendedTimestamp[3])
	}

	return
}
