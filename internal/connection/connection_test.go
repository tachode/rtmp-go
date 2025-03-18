package connection_test

import (
	"net"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/tachode/rtmp-go/chunkstream"
	"github.com/tachode/rtmp-go/internal/connection"
	"github.com/tachode/rtmp-go/message"
)

type mockConn struct {
	Out             chan ([]byte)
	In              chan ([]byte)
	WriteTranscript []byte
	readBuffer      []byte
}

func (m *mockConn) Read(b []byte) (n int, err error) {
	if m.readBuffer == nil {
		m.readBuffer = <-m.In
	}
	copy(b, m.readBuffer)
	if len(m.readBuffer) <= len(b) {
		n = len(m.readBuffer)
		m.readBuffer = nil
	} else {
		n = len(b)
		m.readBuffer = m.readBuffer[n:]
	}
	return
}

func (m *mockConn) Write(b []byte) (n int, err error) {
	bCopy := make([]byte, len(b))
	copy(bCopy, b)
	m.Out <- bCopy
	m.WriteTranscript = append(m.WriteTranscript, b...)
	return len(b), nil
}

func (m *mockConn) Close() error {
	return nil
}

func (m *mockConn) LocalAddr() net.Addr {
	return &mockAddr{}
}

func (m *mockConn) RemoteAddr() net.Addr {
	return &mockAddr{}
}

func (m *mockConn) SetDeadline(t time.Time) error {
	return nil
}

func (m *mockConn) SetReadDeadline(t time.Time) error {
	return nil
}

func (m *mockConn) SetWriteDeadline(t time.Time) error {
	return nil
}

type mockAddr struct{}

func (m *mockAddr) Network() string {
	return "mock"
}

func (m *mockAddr) String() string {
	return "mock.address"
}

func TestNewConnection(t *testing.T) {
	c, err := connection.New(&mockConn{}, 3, connection.NoHandshake)
	assert.NoError(t, err)
	assert.NotNil(t, c)
}

func TestHandshake(t *testing.T) {
	clientToServer := make(chan []byte, 1024)
	serverToClient := make(chan []byte, 1024)
	clientConn := &mockConn{In: serverToClient, Out: clientToServer}
	serverConn := &mockConn{In: clientToServer, Out: serverToClient}
	serverHandshakeDone := make(chan struct{})
	clientHandshakeDone := make(chan struct{})

	go func() {
		client, err := connection.New(clientConn, 3, connection.Client)
		assert.NoError(t, err)
		assert.NotNil(t, client)
		close(clientHandshakeDone)
	}()

	go func() {
		server, err := connection.New(serverConn, 3, connection.Server)
		assert.NoError(t, err)
		assert.NotNil(t, server)
		close(serverHandshakeDone)
	}()

	<-clientHandshakeDone
	<-serverHandshakeDone

	assert.Equal(t, 1+(4+4+1528)*2, len(clientConn.WriteTranscript))
	assert.Equal(t, 1+(4+4+1528)*2, len(serverConn.WriteTranscript))
}

func TestReadMessage(t *testing.T) {
	conn := &mockConn{In: make(chan []byte, 1024)}
	c, err := connection.New(conn, 3, connection.NoHandshake)
	assert.NoError(t, err)
	conn.In <- []byte{
		(byte(chunkstream.HeaderTypeFull) << 6) | 0x5, // full header, stream 5
		0x0, 0x0, 0x0, // message timestamp
		0x0, 0x0, 0x4, // message length
		byte(message.TypeAbortMessage), // message type
		0x0, 0x0, 0x0, 0x0,             // message stream ID
		0x0, 0x0, 0x4, 0xd2, // payload: ChunkStreamid
	}
	msg, err := c.ReadMessage()
	assert.NoError(t, err)
	assert.Equal(t, &message.AbortMessage{ChunkStreamId: 1234, MetadataFields: message.MetadataFields{Length: 4}}, msg)
}

func TestWriteMessageToUnsetStreamId(t *testing.T) {
	conn := &mockConn{Out: make(chan []byte, 1024)}
	c, err := connection.New(conn, 3, connection.NoHandshake)
	assert.NoError(t, err)
	err = c.WriteMessage(&message.AbortMessage{ChunkStreamId: 1234}, 5)
	assert.ErrorIs(t, err, connection.ErrNoSuchChunkstream)
}

func TestWriteMessage(t *testing.T) {
	conn := &mockConn{In: nil, Out: make(chan []byte, 1024)}
	c, err := connection.New(conn, 3, connection.NoHandshake)
	assert.NoError(t, err)
	err = c.CreateOutboundChunkstream(5, 2)
	assert.NoError(t, err)

	expected := []byte{
		(byte(chunkstream.HeaderTypeFull) << 6) | 0x5, // full header, stream 5
		0x0, 0x0, 0x0, // message timestamp
		0x0, 0x0, 0x4, // message length
		byte(message.TypeAbortMessage), // message type
		0x0, 0x0, 0x0, 0x0,             // message stream ID
		0x0, 0x0, 0x4, 0xd2, // payload: ChunkStreamid
	}
	err = c.WriteMessage(&message.AbortMessage{ChunkStreamId: 1234}, 5)
	assert.NoError(t, err)
	// quick spin loop to wait for the messages to be written
	for c.BytesWritten() == 0 || c.OutboundQueueLength() > 0 {
		time.Sleep(1 * time.Millisecond)
	}
	assert.Equal(t, expected, conn.WriteTranscript)
}

func TestRemoteChunkSizeChange(t *testing.T) {
	pc, _, _, _ := runtime.Caller(0)
	t.Skip("TODO:", strings.Split(runtime.FuncForPC(pc).Name(), ".")[1])

}

func TestReadMessageError(t *testing.T) {
	pc, _, _, _ := runtime.Caller(0)
	t.Skip("TODO:", strings.Split(runtime.FuncForPC(pc).Name(), ".")[1])
}

func TestWriteMessageError(t *testing.T) {
	pc, _, _, _ := runtime.Caller(0)
	t.Skip("TODO:", strings.Split(runtime.FuncForPC(pc).Name(), ".")[1])
}

func TestReadMessageTimeout(t *testing.T) {
	// os.ErrDeadlineExceeded should be recoverable
	pc, _, _, _ := runtime.Caller(0)
	t.Skip("TODO:", strings.Split(runtime.FuncForPC(pc).Name(), ".")[1])
}

func TestWriteMessageTimeout(t *testing.T) {
	// os.ErrDeadlineExceeded should be recoverable
	pc, _, _, _ := runtime.Caller(0)
	t.Skip("TODO:", strings.Split(runtime.FuncForPC(pc).Name(), ".")[1])
}

func TestClose(t *testing.T) {
	pc, _, _, _ := runtime.Caller(0)
	t.Skip("TODO:", strings.Split(runtime.FuncForPC(pc).Name(), ".")[1])
}

func TestSeveralMessages(t *testing.T) {
	// Set up a client and server and send several messages between them in both directions
	pc, _, _, _ := runtime.Caller(0)
	t.Skip("TODO:", strings.Split(runtime.FuncForPC(pc).Name(), ".")[1])
}
