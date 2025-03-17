package chunkstream_test

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tachode/rtmp-go/chunkstream"
	"github.com/tachode/rtmp-go/message"
)

func TestNewInboundChunkStream(t *testing.T) {
	chunkStreamId := uint32(3)
	inbound := chunkstream.NewInboundChunkStream(chunkStreamId)
	assert.NotNil(t, inbound, "NewInboundChunkStream should return a non-nil Inbound")
	assert.Equal(t, uint32(128), inbound.MaxChunkSize, "Default MaxChunkSize should be 128")
}

func TestInbound_Read_InvalidChunkStreamId(t *testing.T) {
	inbound := chunkstream.NewInboundChunkStream(3)
	mockReader := bytes.NewReader([]uint8{0xa, 0x0, 0x3, 0xe8, 0x0, 0x0, 0xa, 0x9, 0x1, 0x0, 0x0, 0x0, 0x0, 0x74, 0x65, 0x73, 0x74, 0x20, 0x64, 0x61, 0x74, 0x61})
	_, _, err := inbound.Read(mockReader)
	assert.ErrorIs(t, err, chunkstream.ErrInvalidChunkStreamId, "Expected ErrInvalidChunkStreamId")
}

func TestInbound_Read_CompleteMessage(t *testing.T) {
	inbound := chunkstream.NewInboundChunkStream(10)

	mockData := []uint8{0xa, 0x0, 0x3, 0xe8, 0x0, 0x0, 0xa, 0x9, 0x1, 0x0, 0x0, 0x0, 0x0, 0x74, 0x65, 0x73, 0x74, 0x20, 0x64, 0x61, 0x74, 0x61}
	expectedMessage := &message.VideoMessage{
		MetadataFields: message.MetadataFields{
			Timestamp: 1000,
			StreamId:  1,
			Length:    10,
		},
		Payload: []byte("test data"),
	}

	mockReader := bytes.NewReader(mockData)

	n, msg, err := inbound.Read(mockReader)

	require.NoError(t, err, "Read should not return an error")
	assert.Equal(t, len(mockData), n, "Read bytes count should match input data length")
	assert.NotNil(t, msg, "Message should not be nil")
	assert.Equal(t, expectedMessage, msg, "Message should match expected message")
}

func TestInbound_Read_Multichunk_Message(t *testing.T) {
	inbound := chunkstream.NewInboundChunkStream(10)
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
		Payload: []byte("test data that exceeds chunk size"),
	}

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
	outbound := chunkstream.NewOutboundChunkStream(10)
	msg1 := &message.VideoMessage{
		MetadataFields: message.MetadataFields{
			Timestamp: 1000,
			StreamId:  1,
			Length:    14,
		},
		Payload: []byte("first message"),
	}
	msg2 := &message.VideoMessage{
		MetadataFields: message.MetadataFields{
			Timestamp: 2000,
			StreamId:  2,
			Length:    15,
		},
		Payload: []byte("second message"),
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
	inbound := chunkstream.NewInboundChunkStream(10)
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
	outbound := chunkstream.NewOutboundChunkStream(10)
	msg1 := &message.VideoMessage{
		MetadataFields: message.MetadataFields{
			Timestamp: 1000,
			StreamId:  1,
			Length:    14,
		},
		Payload: []byte("first message"),
	}
	msg2 := &message.VideoMessage{
		MetadataFields: message.MetadataFields{
			Timestamp: 2000,
			StreamId:  1,
			Length:    15,
		},
		Payload: []byte("second message"),
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
	inbound := chunkstream.NewInboundChunkStream(10)
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
	outbound := chunkstream.NewOutboundChunkStream(10)
	msg1 := &message.VideoMessage{
		MetadataFields: message.MetadataFields{
			Timestamp: 1000,
			StreamId:  1,
			Length:    15,
		},
		Payload: []byte("first  message"),
	}
	msg2 := &message.VideoMessage{
		MetadataFields: message.MetadataFields{
			Timestamp: 1000,
			StreamId:  1,
			Length:    15,
		},
		Payload: []byte("second message"),
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
	inbound := chunkstream.NewInboundChunkStream(10)
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
	outbound := chunkstream.NewOutboundChunkStream(10)
	msg1 := &message.VideoMessage{
		MetadataFields: message.MetadataFields{
			Timestamp: 1000,
			StreamId:  1,
			Length:    15,
		},
		Payload: []byte("first  message"),
	}
	msg2 := &message.VideoMessage{
		MetadataFields: message.MetadataFields{
			Timestamp: 2000,
			StreamId:  1,
			Length:    15,
		},
		Payload: []byte("second message"),
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
	inbound := chunkstream.NewInboundChunkStream(10)
	var msg message.Message
	_, msg, err = inbound.Read(&data)
	require.NoError(t, err)
	assert.Equal(t, msg1, msg)
	_, msg, err = inbound.Read(&data)
	require.NoError(t, err)
	assert.Equal(t, msg2, msg)
}
