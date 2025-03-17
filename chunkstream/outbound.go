package chunkstream

import (
	"bytes"
	"slices"

	"github.com/tachode/rtmp-go/message"
)

type Outbound struct {
	MaxChunkSize    uint32
	chunkStreamId   uint32
	lastTimestamp   uint32
	lastChunkHeader *ChunkHeader
}

func NewOutboundChunkStream(chunkStreamId uint32) *Outbound {
	return &Outbound{
		MaxChunkSize:  128,
		chunkStreamId: chunkStreamId,
	}
}

// Marshal serializes a message into a series of chunks.
func (o *Outbound) Marshal(msg message.Message) ([][]byte, error) {
	if o.chunkStreamId < 2 || o.chunkStreamId > 65599 {
		return nil, ErrInvalidChunkStreamId
	}
	data, err := msg.Marshal()
	if err != nil {
		return nil, err
	}
	msg.Metadata().Length = uint32(len(data))

	delta := msg.Metadata().Timestamp - o.lastTimestamp
	o.lastTimestamp = msg.Metadata().Timestamp

	chunks := make([][]byte, 0, len(data)/int(o.MaxChunkSize)+1)

	for len(data) > 0 {
		header := &ChunkHeader{
			Type:             HeaderTypeFull,
			ChunkStreamId:    o.chunkStreamId,
			Timestamp:        msg.Metadata().Timestamp,
			MessageLength:    msg.Metadata().Length,
			MessageType:      msg.Type(),
			MessageStreamId:  msg.Metadata().StreamId,
			TimestampIsDelta: false,
		}

		// Figure out the optimal chunk header format
		if o.lastChunkHeader != nil && header.MessageStreamId == o.lastChunkHeader.MessageStreamId {
			header.Timestamp = delta
			header.TimestampIsDelta = true
			header.Type = HeaderTypeSameStream
			if header.MessageLength == o.lastChunkHeader.MessageLength {
				header.Type = HeaderTypeSameStreamAndLength
				if delta == o.lastChunkHeader.Timestamp {
					header.Type = HeaderTypeContinuation
				}
			}
		}

		chunk := bytes.NewBuffer(nil)
		header.Write(chunk)
		chunkSize := min(o.MaxChunkSize, uint32(len(data)))
		chunk.Write(data[:chunkSize])
		data = data[chunkSize:]
		chunks = append(chunks, slices.Clone(chunk.Bytes()))
		o.lastChunkHeader = header
	}

	return chunks, nil
}
