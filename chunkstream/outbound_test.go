package chunkstream_test

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tachode/rtmp-go/chunkstream"
	"github.com/tachode/rtmp-go/message"
)

func TestNewOutboundChunkStream(t *testing.T) {
	outbound := chunkstream.NewOutboundChunkStream(10)
	assert.NotNil(t, outbound)
	assert.Equal(t, uint32(128), outbound.MaxChunkSize)
}

func TestOutboundMarshal_SingleChunk(t *testing.T) {
	outbound := chunkstream.NewOutboundChunkStream(10)
	msg := &message.VideoMessage{
		MetadataFields: message.MetadataFields{
			Timestamp: 1000,
			StreamId:  1,
		},
		Payload: []byte("test data"),
	}
	msgBytes, err := msg.Marshal()
	require.NoError(t, err)

	chunks, err := outbound.Marshal(msg)
	require.NoError(t, err)
	require.Len(t, chunks, 1)

	csHeader := chunkstream.ChunkHeader{}
	_, err = csHeader.Read(bytes.NewReader(chunks[0]))
	require.NoError(t, err)
	assert.Equal(t, chunkstream.HeaderTypeFull, csHeader.Type)
	assert.Equal(t, uint32(10), csHeader.ChunkStreamId)
	assert.Equal(t, uint32(1000), csHeader.Timestamp)
	assert.Equal(t, uint32(len(msgBytes)), csHeader.MessageLength)

	require.Equal(t, msgBytes, chunks[0][12:]) // Skip header (1 byte for CSID, 11 bytes for header)

	assert.Greater(t, len(chunks[0]), 0)
}

func TestOutboundMarshal_MultipleChunks(t *testing.T) {
	outbound := chunkstream.NewOutboundChunkStream(10)
	outbound.MaxChunkSize = 4
	msg := &message.VideoMessage{
		MetadataFields: message.MetadataFields{
			Timestamp: 1000,
			StreamId:  1,
		},
		Payload: []byte("test data that exceeds chunk size"),
	}
	msgBytes, err := msg.Marshal()
	require.NoError(t, err)

	chunks, err := outbound.Marshal(msg)
	require.NoError(t, err)
	require.Greater(t, len(chunks), 1)

	totalData := []byte{}
	for i, chunk := range chunks {
		headerSize := 1
		if i == 0 {
			headerSize = 12
		}
		totalData = append(totalData, chunk[headerSize:]...) // Skip header
	}
	assert.Equal(t, msgBytes, totalData)
}

func TestOutboundMarshal_HeaderOptimization_SameStream(t *testing.T) {
	outbound := chunkstream.NewOutboundChunkStream(10)
	msg1 := &message.VideoMessage{
		MetadataFields: message.MetadataFields{
			Timestamp: 1000,
			StreamId:  1,
		},
		Payload: []byte("first message"),
	}
	msg2 := &message.VideoMessage{
		MetadataFields: message.MetadataFields{
			Timestamp: 2000,
			StreamId:  1,
		},
		Payload: []byte("second message"),
	}

	_, err := outbound.Marshal(msg1)
	require.NoError(t, err)

	chunks, err := outbound.Marshal(msg2)
	require.NoError(t, err)
	require.Len(t, chunks, 1)

	assert.Equal(t, chunkstream.HeaderTypeSameStream, chunkstream.HeaderType(chunks[0][0]>>6))
}

func TestOutboundMarshal_HeaderOptimization_SameStreamAndLength(t *testing.T) {
	outbound := chunkstream.NewOutboundChunkStream(10)
	msg1 := &message.VideoMessage{
		MetadataFields: message.MetadataFields{
			Timestamp: 1000,
			StreamId:  1,
		},
		Payload: []byte("first--message"),
	}
	msg2 := &message.VideoMessage{
		MetadataFields: message.MetadataFields{
			Timestamp: 1000,
			StreamId:  1,
		},
		Payload: []byte("second message"),
	}

	_, err := outbound.Marshal(msg1)
	require.NoError(t, err)

	chunks, err := outbound.Marshal(msg2)
	require.NoError(t, err)
	require.Len(t, chunks, 1)

	assert.Equal(t, chunkstream.HeaderTypeSameStreamAndLength, chunkstream.HeaderType(chunks[0][0]>>6))
}

func TestOutboundMarshal_HeaderOptimization_Continuation(t *testing.T) {
	outbound := chunkstream.NewOutboundChunkStream(10)
	msg1 := &message.VideoMessage{
		MetadataFields: message.MetadataFields{
			Timestamp: 1000,
			StreamId:  1,
		},
		Payload: []byte("first--message"),
	}
	msg2 := &message.VideoMessage{
		MetadataFields: message.MetadataFields{
			Timestamp: 2000,
			StreamId:  1,
		},
		Payload: []byte("second message"),
	}

	_, err := outbound.Marshal(msg1)
	require.NoError(t, err)

	chunks, err := outbound.Marshal(msg2)
	require.NoError(t, err)
	require.Len(t, chunks, 1)

	assert.Equal(t, chunkstream.HeaderTypeContinuation, chunkstream.HeaderType(chunks[0][0]>>6))
}
