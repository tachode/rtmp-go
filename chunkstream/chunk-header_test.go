package chunkstream_test

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tachode/rtmp-go/chunkstream"
	"github.com/tachode/rtmp-go/message"
)

func TestChunkMessageHeader_WriteAndRead_FullHeader(t *testing.T) {
	header := chunkstream.ChunkHeader{
		Type:            chunkstream.HeaderTypeFull,
		ChunkStreamId:   3,
		Timestamp:       123456,
		MessageLength:   789,
		MessageType:     message.Type(8),
		MessageStreamId: 1,
	}

	var buf bytes.Buffer
	_, err := header.Write(&buf)
	require.NoError(t, err, "Write should not fail")

	var readHeader chunkstream.ChunkHeader
	_, err = readHeader.Read(&buf)
	require.NoError(t, err, "Read should not fail")

	assert.Equal(t, header, readHeader, "Headers should be equal")
}

func TestChunkMessageHeader_WriteAndRead_SameStreamHeader(t *testing.T) {
	header := chunkstream.ChunkHeader{
		Type:             chunkstream.HeaderTypeSameStream,
		ChunkStreamId:    5,
		Timestamp:        654321,
		MessageLength:    456,
		MessageType:      message.Type(9),
		TimestampIsDelta: true,
	}

	var buf bytes.Buffer
	_, err := header.Write(&buf)
	require.NoError(t, err, "Write should not fail")

	var readHeader chunkstream.ChunkHeader
	_, err = readHeader.Read(&buf)
	require.NoError(t, err, "Read should not fail")

	assert.Equal(t, header, readHeader, "Headers should be equal")
}

func TestChunkMessageHeader_WriteAndRead_SameLengthAndStreamHeader(t *testing.T) {
	header := chunkstream.ChunkHeader{
		Type:             chunkstream.HeaderTypeSameStreamAndLength,
		ChunkStreamId:    10,
		Timestamp:        98765,
		TimestampIsDelta: true,
	}

	var buf bytes.Buffer
	_, err := header.Write(&buf)
	require.NoError(t, err, "Write should not fail")

	var readHeader chunkstream.ChunkHeader
	_, err = readHeader.Read(&buf)
	require.NoError(t, err, "Read should not fail")

	assert.Equal(t, header, readHeader, "Headers should be equal")
}

func TestChunkMessageHeader_WriteAndRead_ContinuationHeader(t *testing.T) {
	header := chunkstream.ChunkHeader{
		Type:             chunkstream.HeaderTypeContinuation,
		ChunkStreamId:    15,
		TimestampIsDelta: true,
	}

	var buf bytes.Buffer
	_, err := header.Write(&buf)
	require.NoError(t, err, "Write should not fail")

	var readHeader chunkstream.ChunkHeader
	_, err = readHeader.Read(&buf)
	require.NoError(t, err, "Read should not fail")

	assert.Equal(t, header, readHeader, "Headers should be equal")
}

func TestChunkMessageHeader_Write_InvalidChunkStreamId(t *testing.T) {
	header := chunkstream.ChunkHeader{
		Type:          chunkstream.HeaderTypeFull,
		ChunkStreamId: 70000, // Invalid ChunkStreamId
	}

	var buf bytes.Buffer
	_, err := header.Write(&buf)
	assert.Error(t, err, "Expected error for invalid ChunkStreamId")
}

func TestChunkMessageHeader_WriteAndRead_TimestampIsDeltaMismatch_FullHeader(t *testing.T) {
	header := chunkstream.ChunkHeader{
		Type:             chunkstream.HeaderTypeFull,
		ChunkStreamId:    3,
		Timestamp:        123456,
		MessageLength:    789,
		MessageType:      message.Type(8),
		MessageStreamId:  1,
		TimestampIsDelta: true, // Invalid for FullHeader
	}

	var buf bytes.Buffer
	_, err := header.Write(&buf)
	assert.Error(t, err, "Expected error for TimestampIsDelta mismatch in FullHeader")
	assert.Equal(t, chunkstream.ErrDeltaTimePassedToFullHeader, err, "Error should match expected")
}

func TestChunkMessageHeader_WriteAndRead_TimestampIsDeltaMismatch_SameStreamHeader(t *testing.T) {
	header := chunkstream.ChunkHeader{
		Type:             chunkstream.HeaderTypeSameStream,
		ChunkStreamId:    5,
		Timestamp:        654321,
		MessageLength:    456,
		MessageType:      message.Type(9),
		TimestampIsDelta: false, // Invalid for SameStreamHeader
	}

	var buf bytes.Buffer
	_, err := header.Write(&buf)
	assert.Error(t, err, "Expected error for TimestampIsDelta mismatch in SameStreamHeader")
	assert.Equal(t, chunkstream.ErrNonDeltaTimestampPassedToShortHeader, err, "Error should match expected")
}

func TestChunkMessageHeader_WriteAndRead_TimestampIsDeltaMismatch_SameLengthAndStreamHeader(t *testing.T) {
	header := chunkstream.ChunkHeader{
		Type:             chunkstream.HeaderTypeSameStreamAndLength,
		ChunkStreamId:    10,
		Timestamp:        98765,
		TimestampIsDelta: false, // Invalid for SameLengthAndStreamHeader
	}

	var buf bytes.Buffer
	_, err := header.Write(&buf)
	assert.Error(t, err, "Expected error for TimestampIsDelta mismatch in SameLengthAndStreamHeader")
	assert.Equal(t, chunkstream.ErrNonDeltaTimestampPassedToShortHeader, err, "Error should match expected")
}
