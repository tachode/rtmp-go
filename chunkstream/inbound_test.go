package chunkstream_test

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tachode/rtmp-go/chunkstream"
	"github.com/tachode/rtmp-go/message"
)

var mc = message.NewContext()

func TestNewInboundChunkStream(t *testing.T) {
	chunkStreamId := uint32(3)
	inbound := chunkstream.NewInboundChunkStream(chunkStreamId, mc)
	assert.NotNil(t, inbound, "NewInboundChunkStream should return a non-nil Inbound")
	assert.Equal(t, uint32(128), inbound.MaxChunkSize, "Default MaxChunkSize should be 128")
}

func TestInbound_Read_InvalidChunkStreamId(t *testing.T) {
	inbound := chunkstream.NewInboundChunkStream(3, mc)
	mockReader := bytes.NewReader([]uint8{0xa, 0x0, 0x3, 0xe8, 0x0, 0x0, 0xa, 0x9, 0x1, 0x0, 0x0, 0x0, 0x0, 0x74, 0x65, 0x73, 0x74, 0x20, 0x64, 0x61, 0x74, 0x61})
	_, _, err := inbound.Read(mockReader)
	assert.ErrorIs(t, err, chunkstream.ErrInvalidChunkStreamId, "Expected ErrInvalidChunkStreamId")
}

func TestInbound_Read_CompleteMessage(t *testing.T) {
	inbound := chunkstream.NewInboundChunkStream(10, mc)

	mockData := []uint8{0xa, 0x0, 0x3, 0xe8, 0x0, 0x0, 0xa, 0x9, 0x1, 0x0, 0x0, 0x0, 0x0, 0x74, 0x65, 0x73, 0x74, 0x20, 0x64, 0x61, 0x74, 0x61}
	expectedMessage := &message.VideoMessage{
		MetadataFields: message.MetadataFields{
			Timestamp: 1000,
			StreamId:  1,
			Length:    10,
		},
		PacketType: message.ERTMPVideoPacketTypeCodedFrames,
		Tracks: []message.VideoTrack{{
			Payload: []byte("test data"),
		}},
	}
	expectedMessage.SetContext(mc)

	mockReader := bytes.NewReader(mockData)

	n, msg, err := inbound.Read(mockReader)

	require.NoError(t, err, "Read should not return an error")
	assert.Equal(t, len(mockData), n, "Read bytes count should match input data length")
	assert.NotNil(t, msg, "Message should not be nil")
	assert.Equal(t, expectedMessage, msg, "Message should match expected message")
}

func TestInbound_Read_Multichunk_Message(t *testing.T) {
	inbound := chunkstream.NewInboundChunkStream(10, mc)
	inbound.MaxChunkSize = 4

	mockData := []byte{
		0xa, 0x0, 0x3, 0xe8, 0x0, 0x0, 0x22, 0x9, 0x1, 0x0, 0x0, 0x0, 0x0, 0x74, 0x65, 0x73,
		0xca, 0x74, 0x20, 0x64, 0x61,
		0xca, 0x74, 0x61, 0x20, 0x74,
		0xca, 0x68, 0x61, 0x74, 0x20,
		0xca, 0x65, 0x78, 0x63, 0x65,
		0xca, 0x65, 0x64, 0x73, 0x20,
		0xca, 0x63, 0x68, 0x75, 0x6e,
		0xca, 0x6b, 0x20, 0x73, 0x69,
		0xca, 0x7a, 0x65,
	}
	mockReader := bytes.NewReader(mockData)
	expectedMessage := &message.VideoMessage{
		MetadataFields: message.MetadataFields{
			Timestamp: 1000,
			StreamId:  1,
			Length:    34,
		},
		PacketType: message.ERTMPVideoPacketTypeCodedFrames,
		Tracks: []message.VideoTrack{{
			Payload: []byte("test data that exceeds chunk size"),
		}},
	}
	expectedMessage.SetContext(mc)

	var msg message.Message
	var n int
	var err error
	var chunkCount int
	for msg == nil {
		var m int
		m, msg, err = inbound.Read(mockReader)
		n += m
		require.NoError(t, err, "Read should not return an error")
		chunkCount++
	}
	assert.Equal(t, 9, chunkCount, "Expected 9 chunks to be read")
	assert.Equal(t, len(mockData), n, "Read bytes count should match input data length")
	assert.Equal(t, expectedMessage, msg, "Message should match expected message")
}

func Test_Second_Message_Headertype_Full(t *testing.T) {
	// Generate the bytestream
	outbound := chunkstream.NewOutboundChunkStream(10, mc)
	msg1 := &message.VideoMessage{
		MetadataFields: message.MetadataFields{
			Timestamp: 1000,
			StreamId:  1,
			Length:    14,
		},
		PacketType: message.ERTMPVideoPacketTypeCodedFrames,
		Tracks: []message.VideoTrack{{
			Payload: []byte("first message"),
		}},
	}
	msg2 := &message.VideoMessage{
		MetadataFields: message.MetadataFields{
			Timestamp: 2000,
			StreamId:  2,
			Length:    15,
		},
		PacketType: message.ERTMPVideoPacketTypeCodedFrames,
		Tracks: []message.VideoTrack{{
			Payload: []byte("second message"),
		}},
	}

	data := bytes.Buffer{}

	chunks, err := outbound.Marshal(msg1)
	require.NoError(t, err)
	for _, chunk := range chunks {
		data.Write(chunk)
	}

	chunks, err = outbound.Marshal(msg2)
	require.NoError(t, err)
	for _, chunk := range chunks {
		data.Write(chunk)
	}
	assert.Equal(t, chunkstream.HeaderTypeFull, chunkstream.HeaderType(chunks[0][0]>>6))

	// Read the bytestream
	inbound := chunkstream.NewInboundChunkStream(10, mc)
	var msg message.Message
	_, msg, err = inbound.Read(&data)
	require.NoError(t, err)
	assert.Equal(t, msg1, msg)
	_, msg, err = inbound.Read(&data)
	require.NoError(t, err)
	assert.Equal(t, msg2, msg)
}

func Test_Second_Message_Headertype_SameStream(t *testing.T) {
	// Generate the bytestream
	outbound := chunkstream.NewOutboundChunkStream(10, mc)
	msg1 := &message.VideoMessage{
		MetadataFields: message.MetadataFields{
			Timestamp: 1000,
			StreamId:  1,
			Length:    14,
		},
		PacketType: message.ERTMPVideoPacketTypeCodedFrames,
		Tracks: []message.VideoTrack{{
			Payload: []byte("first message"),
		}},
	}
	msg2 := &message.VideoMessage{
		MetadataFields: message.MetadataFields{
			Timestamp: 2000,
			StreamId:  1,
			Length:    15,
		},
		PacketType: message.ERTMPVideoPacketTypeCodedFrames,
		Tracks: []message.VideoTrack{{
			Payload: []byte("second message"),
		}},
	}

	data := bytes.Buffer{}

	chunks, err := outbound.Marshal(msg1)
	require.NoError(t, err)
	for _, chunk := range chunks {
		data.Write(chunk)
	}

	chunks, err = outbound.Marshal(msg2)
	require.NoError(t, err)
	for _, chunk := range chunks {
		data.Write(chunk)
	}
	assert.Equal(t, chunkstream.HeaderTypeSameStream, chunkstream.HeaderType(chunks[0][0]>>6))

	// Read the bytestream
	inbound := chunkstream.NewInboundChunkStream(10, mc)
	var msg message.Message
	_, msg, err = inbound.Read(&data)
	require.NoError(t, err)
	assert.Equal(t, msg1, msg)
	_, msg, err = inbound.Read(&data)
	require.NoError(t, err)
	assert.Equal(t, msg2, msg)
}

func Test_Second_Message_Headertype_SameStreamAndLength(t *testing.T) {
	// Generate the bytestream
	outbound := chunkstream.NewOutboundChunkStream(10, mc)
	msg1 := &message.VideoMessage{
		MetadataFields: message.MetadataFields{
			Timestamp: 1000,
			StreamId:  1,
			Length:    15,
		},
		PacketType: message.ERTMPVideoPacketTypeCodedFrames,
		Tracks: []message.VideoTrack{{
			Payload: []byte("first  message"),
		}},
	}
	msg2 := &message.VideoMessage{
		MetadataFields: message.MetadataFields{
			Timestamp: 1000,
			StreamId:  1,
			Length:    15,
		},
		PacketType: message.ERTMPVideoPacketTypeCodedFrames,
		Tracks: []message.VideoTrack{{
			Payload: []byte("second message"),
		}},
	}

	data := bytes.Buffer{}

	chunks, err := outbound.Marshal(msg1)
	require.NoError(t, err)
	for _, chunk := range chunks {
		data.Write(chunk)
	}

	chunks, err = outbound.Marshal(msg2)
	require.NoError(t, err)
	for _, chunk := range chunks {
		data.Write(chunk)
	}
	assert.Equal(t, chunkstream.HeaderTypeSameStreamAndLength, chunkstream.HeaderType(chunks[0][0]>>6))

	// Read the bytestream
	inbound := chunkstream.NewInboundChunkStream(10, mc)
	var msg message.Message
	_, msg, err = inbound.Read(&data)
	require.NoError(t, err)
	assert.Equal(t, msg1, msg)
	_, msg, err = inbound.Read(&data)
	require.NoError(t, err)
	assert.Equal(t, msg2, msg)
}

func Test_Second_Message_Headertype_Continuation(t *testing.T) {
	// Generate the bytestream
	outbound := chunkstream.NewOutboundChunkStream(10, mc)
	msg1 := &message.VideoMessage{
		MetadataFields: message.MetadataFields{
			Timestamp: 1000,
			StreamId:  1,
			Length:    15,
		},
		PacketType: message.ERTMPVideoPacketTypeCodedFrames,
		Tracks: []message.VideoTrack{{
			Payload: []byte("first  message"),
		}},
	}
	msg2 := &message.VideoMessage{
		MetadataFields: message.MetadataFields{
			Timestamp: 2000,
			StreamId:  1,
			Length:    15,
		},
		PacketType: message.ERTMPVideoPacketTypeCodedFrames,
		Tracks: []message.VideoTrack{{
			Payload: []byte("second message"),
		}},
	}

	data := bytes.Buffer{}

	chunks, err := outbound.Marshal(msg1)
	require.NoError(t, err)
	for _, chunk := range chunks {
		data.Write(chunk)
	}

	chunks, err = outbound.Marshal(msg2)
	require.NoError(t, err)
	for _, chunk := range chunks {
		data.Write(chunk)
	}
	assert.Equal(t, chunkstream.HeaderTypeContinuation, chunkstream.HeaderType(chunks[0][0]>>6))

	// Read the bytestream
	inbound := chunkstream.NewInboundChunkStream(10, mc)
	var msg message.Message
	_, msg, err = inbound.Read(&data)
	require.NoError(t, err)
	assert.Equal(t, msg1, msg)
	_, msg, err = inbound.Read(&data)
	require.NoError(t, err)
	assert.Equal(t, msg2, msg)
}

func Test_ExtendedTimestamp_MultichunkRoundTrip(t *testing.T) {
	// A message whose timestamp exceeds 0xFFFFFF triggers extended timestamps.
	// When the message is larger than the chunk size, continuation (type 3) chunks
	// must also carry the 4-byte extended timestamp. This test verifies the
	// full outbound→inbound round-trip for that case.
	outbound := chunkstream.NewOutboundChunkStream(10, mc)
	outbound.MaxChunkSize = 16

	payload := []byte("extended timestamp data that exceeds chunk size")
	msg1 := &message.VideoMessage{
		MetadataFields: message.MetadataFields{
			Timestamp: 0x1234567, // > 0xFFFFFF, triggers extended timestamp
			StreamId:  1,
		},
		PacketType: message.ERTMPVideoPacketTypeCodedFrames,
		Tracks: []message.VideoTrack{{
			Payload: payload,
		}},
	}

	chunks, err := outbound.Marshal(msg1)
	require.NoError(t, err)
	require.Greater(t, len(chunks), 1, "message should span multiple chunks")

	// Verify the continuation chunks are type 3 with extended timestamp bytes
	for i, chunk := range chunks[1:] {
		assert.Equal(t, chunkstream.HeaderTypeContinuation, chunkstream.HeaderType(chunk[0]>>6),
			"chunk %d should be type 3 continuation", i+1)
	}

	// Round-trip through inbound
	var data bytes.Buffer
	for _, chunk := range chunks {
		data.Write(chunk)
	}

	inbound := chunkstream.NewInboundChunkStream(10, mc)
	inbound.MaxChunkSize = 16

	var result message.Message
	for result == nil {
		var m int
		m, result, err = inbound.Read(&data)
		require.NoError(t, err, "Read should not return an error")
		require.Greater(t, m, 0, "should read at least one byte per chunk")
	}

	msg1.MetadataFields.Length = uint32(len(payload) + 1) // +1 for the video header byte
	assert.Equal(t, msg1, result, "round-tripped message should match original")
	assert.Equal(t, 0, data.Len(), "all bytes should be consumed")
}
