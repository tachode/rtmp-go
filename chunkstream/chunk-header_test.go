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

func TestChunkBasicHeader_ThreeByteFormat(t *testing.T) {
	// Spec §5.3.1.1: for 3-byte basic headers (cs id 64–65599), the wire
	// format is: [fmt:2 | 0b000001:6] [cs_id-64 low byte] [cs_id-64 high byte]
	// i.e. the two ID bytes are little-endian.
	for _, csid := range []uint32{320, 1000, 12345, 65599} {
		header := chunkstream.ChunkHeader{
			Type:            chunkstream.HeaderTypeFull,
			ChunkStreamId:   csid,
			Timestamp:       100,
			MessageLength:   50,
			MessageType:     message.Type(8),
			MessageStreamId: 1,
		}

		var buf bytes.Buffer
		_, err := header.Write(&buf)
		require.NoError(t, err, "Write should not fail for csid %d", csid)

		// Verify the wire bytes: first byte low 6 bits must be 0x01,
		// second byte is low byte of (csid-64), third byte is high byte.
		wire := buf.Bytes()
		assert.Equal(t, byte(0x01), wire[0]&0x3F, "low 6 bits should be 1 for 3-byte header (csid %d)", csid)
		assert.Equal(t, byte(csid-64), wire[1], "second byte should be low byte of csid-64 (csid %d)", csid)
		assert.Equal(t, byte((csid-64)>>8), wire[2], "third byte should be high byte of csid-64 (csid %d)", csid)

		// Round-trip
		var readHeader chunkstream.ChunkHeader
		readBuf := bytes.NewBuffer(buf.Bytes())
		_, err = readHeader.Read(readBuf)
		require.NoError(t, err, "Read should not fail for csid %d", csid)
		assert.Equal(t, header, readHeader, "Headers should round-trip for csid %d", csid)
	}
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

func TestChunkMessageHeader_ExtendedTimestamp_Type3Continuation(t *testing.T) {
	// When a message uses extended timestamps (>= 0xFFFFFF), every type 3
	// continuation chunk also carries the 4-byte extended timestamp field.
	// Verify that Read correctly consumes those bytes on type 3 chunks.
	extendedTS := uint32(0x1234567)

	// Write a type 0 header with an extended timestamp.
	type0 := chunkstream.ChunkHeader{
		Type:            chunkstream.HeaderTypeFull,
		ChunkStreamId:   3,
		Timestamp:       extendedTS,
		MessageLength:   500,
		MessageType:     message.Type(8),
		MessageStreamId: 1,
	}
	var buf bytes.Buffer
	_, err := type0.Write(&buf)
	require.NoError(t, err)

	// Write a type 3 continuation with the same extended timestamp.
	type3 := chunkstream.ChunkHeader{
		Type:             chunkstream.HeaderTypeContinuation,
		ChunkStreamId:    3,
		Timestamp:        extendedTS,
		TimestampIsDelta: true,
	}
	_, err = type3.Write(&buf)
	require.NoError(t, err)

	// Read the type 0 header.
	var h chunkstream.ChunkHeader
	_, err = h.Read(&buf)
	require.NoError(t, err)
	assert.Equal(t, chunkstream.HeaderTypeFull, h.Type)
	assert.Equal(t, extendedTS, h.Timestamp)

	// Read the type 3 header on the same ChunkHeader (simulates Inbound reuse).
	// Before the fix, this would fail to consume the 4-byte extended timestamp,
	// leaving stale bytes in the buffer.
	_, err = h.Read(&buf)
	require.NoError(t, err)
	assert.Equal(t, chunkstream.HeaderTypeContinuation, h.Type)
	assert.Equal(t, extendedTS, h.Timestamp)

	// Buffer should be fully consumed.
	assert.Equal(t, 0, buf.Len(), "all bytes should be consumed")
}
