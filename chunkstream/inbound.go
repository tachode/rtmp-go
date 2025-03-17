package chunkstream

import (
	"io"

	"github.com/tachode/rtmp-go/message"
)

type Inbound struct {
	MaxChunkSize  uint32
	chunkStreamId uint32
	data          []byte
	unreadData    []byte
	chunkHeader   ChunkHeader
	streamTime    uint32
}

func NewInboundChunkStream(chunkStreamId uint32) *Inbound {
	return &Inbound{
		// RTMP Specification §5.4.1: "The maximum chunk size defaults to 128 bytes
		MaxChunkSize:  128,
		chunkStreamId: chunkStreamId,
	}
}

func (i *Inbound) Read(r io.Reader) (n int, msg message.Message, err error) {
	m, err := i.chunkHeader.Read(r)
	if err != nil {
		return
	}
	n += m
	if i.chunkHeader.ChunkStreamId != i.chunkStreamId {
		err = ErrInvalidChunkStreamId
		return
	}

	// At the start of a new message, i.data will be nil
	if i.data == nil {
		i.data = make([]byte, i.chunkHeader.MessageLength)
		i.unreadData = i.data[:]
		if i.chunkHeader.TimestampIsDelta {
			i.streamTime += i.chunkHeader.Timestamp
		} else {
			i.streamTime = i.chunkHeader.Timestamp
		}

	}

	chunkSize := min(int(i.MaxChunkSize), len(i.unreadData))
	m, err = io.ReadFull(r, i.unreadData[:chunkSize])
	n += m
	if err != nil {
		return
	}
	i.unreadData = i.unreadData[chunkSize:]

	// Once we've read a complete message, parse it and return it
	if len(i.unreadData) == 0 {
		msg, err = message.Unmarshal(i.streamTime, i.chunkHeader.MessageType, i.chunkHeader.MessageStreamId, i.data)
		i.data = nil
		i.unreadData = nil
	}

	return
}
