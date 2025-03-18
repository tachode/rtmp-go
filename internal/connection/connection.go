package connection

import (
	"bufio"
	"bytes"
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"reflect"
	"sync"
	"sync/atomic"
	"time"

	"github.com/tachode/rtmp-go/chunkstream"
	"github.com/tachode/rtmp-go/message"
)

type Role uint8

//go:generate stringer -type=Role
const (
	Client Role = iota
	Server
	NoHandshake
)

// TODO -- work out how we set outbound chunk sizes

type connection struct {
	conn                net.Conn
	outboundQueue       *PriorityQueue[[]byte]
	inboundQueue        chan (message.Message)
	outboundChunkStream map[uint32]*outboundChunkStream

	// mu protects error
	mu  sync.RWMutex
	err error

	// statistics
	bytesWritten atomic.Uint64
	bytesRead    atomic.Uint64
}

type outboundChunkStream struct {
	mu       sync.Mutex
	cs       *chunkstream.Outbound
	priority int
}

func New(netConn net.Conn, priorityCount int, role Role) (*connection, error) {
	c := &connection{
		conn:                netConn,
		outboundQueue:       NewPriorityQueue[[]byte](priorityCount, 10),
		inboundQueue:        make(chan message.Message),
		outboundChunkStream: make(map[uint32]*outboundChunkStream),
	}
	c.CreateOutboundChunkstream(2, 0) // CS ID 2, at high priority, for command messages

	var err error
	switch role {
	case Client:
		err = c.clientHandshake()
	case Server:
		err = c.serverHandshake()
	}
	if err != nil {
		return nil, fmt.Errorf("handshake: %w", err)
	}

	go c.writeMessages()
	go c.readMessages()
	return c, nil
}

func (c *connection) ReadMessage() (msg message.Message, err error) {
	err = c.getError()
	if err != nil {
		return
	}
	msg = <-c.inboundQueue
	if msg == nil {
		err = c.getError()
		if err == nil {
			err = io.EOF
		}
	}
	return
}

func (c *connection) WriteMessage(msg message.Message, chunkStreamId int) (err error) {
	err = c.getError()
	if err != nil {
		return
	}
	cs, found := c.outboundChunkStream[uint32(chunkStreamId)]
	if !found {
		return ErrNoSuchChunkstream
	}
	cs.mu.Lock()
	defer cs.mu.Unlock()
	chunks, err := cs.cs.Marshal(msg)
	if err != nil {
		return
	}
	for _, chunk := range chunks {
		c.outboundQueue.Enqueue(chunk, cs.priority)
	}
	return
}

// This is not synchronized with WriteMessage -- callers must take care not
// to call them from different goroutines without proper locking.
func (c *connection) CreateOutboundChunkstream(chunkStreamId int, priority int) (err error) {
	err = c.getError()
	if err != nil {
		return
	}
	_, found := c.outboundChunkStream[uint32(chunkStreamId)]
	if found {
		return ErrChunkStreamAlreadyExists
	}
	c.outboundChunkStream[uint32(chunkStreamId)] = &outboundChunkStream{
		cs:       chunkstream.NewOutboundChunkStream(uint32(chunkStreamId)),
		priority: priority,
	}
	return
}

func (c *connection) Read(b []byte) (n int, err error) {
	// Read a message and serialize it (including reconstructed chunkheader)
	msg, err := c.ReadMessage()
	if err != nil {
		return
	}
	csHeader := chunkstream.ChunkHeader{
		Type:             chunkstream.HeaderTypeFull,
		ChunkStreamId:    0, // We don't populate this, but it's not important at this layer
		Timestamp:        msg.Metadata().Timestamp,
		MessageLength:    msg.Metadata().Length,
		MessageType:      msg.Type(),
		MessageStreamId:  msg.Metadata().StreamId,
		TimestampIsDelta: false,
	}
	w := bytes.Buffer{}
	_, err = csHeader.Write(&w)
	if err != nil {
		return
	}

	body, err := msg.Marshal()
	if err != nil {
		return
	}
	_, err = w.Write(body)
	if err != nil {
		return
	}

	n = copy(b, w.Bytes())

	return
}

func (c *connection) Write(b []byte) (n int, err error) {
	// Parse the header and message and call WriteMessage() -- this really only works correctly with full chunk headers.
	r := bytes.NewBuffer(b)
	csHeader := chunkstream.ChunkHeader{}
	_, err = csHeader.Read(r)
	if err != nil {
		return
	}
	msg, err := message.Unmarshal(csHeader.Timestamp, csHeader.MessageType, csHeader.MessageStreamId, r.Bytes())
	if err != nil {
		return
	}
	c.WriteMessage(msg, int(csHeader.ChunkStreamId))
	n = len(b)
	return
}

func (c *connection) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.err == nil {
		c.err = ErrConnectionClosed
	} else {
		return c.err
	}
	c.outboundQueue.Close()
	err := c.conn.Close()
	if err != nil {
		c.err = err
	}
	return c.err
}

func (c *connection) LocalAddr() net.Addr {
	return c.conn.LocalAddr()
}

func (c *connection) RemoteAddr() net.Addr {
	return c.conn.RemoteAddr()
}

func (c *connection) SetDeadline(t time.Time) error {
	return c.conn.SetDeadline(t)
}

func (c *connection) SetReadDeadline(t time.Time) error {
	return c.conn.SetReadDeadline(t)
}

func (c *connection) SetWriteDeadline(t time.Time) error {
	return c.conn.SetWriteDeadline(t)
}

///////////////////////////////////////////////////////////////////////////

func (c *connection) OutboundQueueLength() int {
	return c.outboundQueue.Length()
}

func (c *connection) BytesWritten() uint64 {
	return c.bytesWritten.Load()
}

func (c *connection) BytesRead() uint64 {
	return c.bytesRead.Load()
}

func (c *connection) writeMessages() {
	for {
		chunk := c.outboundQueue.Dequeue()
		if chunk == nil {
			return
		}
		n, err := c.conn.Write(chunk)
		c.bytesWritten.Add(uint64(n))
		if err != nil {
			c.setError("write", err)
		}
	}
}

func (c *connection) readMessages() {
	defer func() { close(c.inboundQueue) }()
	r := bufio.NewReader(c.conn)
	chunkStream := make(map[uint32]*chunkstream.Inbound)
	bytesRead := 0
	lastAck := 0
	remoteWindowSize := 2_500_000 // Everyone seems to use this value as the default
	for {
		chunkStreamId := uint32(0)
		basicHeader, err := r.Peek(3)
		if err != nil {
			c.setError("peek", err)
		}
		basicHeader[0] &= 0x3F
		switch basicHeader[0] {
		case 0:
			chunkStreamId = uint32(basicHeader[1] + 64)
		case 1:
			chunkStreamId = (uint32(basicHeader[1]) << 8) + uint32(basicHeader[2]) + 64
		default:
			chunkStreamId = uint32(basicHeader[0])
		}
		cs := chunkStream[chunkStreamId]
		if cs == nil {
			cs = chunkstream.NewInboundChunkStream(chunkStreamId)
			chunkStream[chunkStreamId] = cs
		}

		n, msg, err := cs.Read(r)
		bytesRead += n
		c.bytesRead.Add(uint64(n))
		if err != nil {
			c.setError("read", err)
		}
		if bytesRead-lastAck > (remoteWindowSize / 2) {
			c.WriteMessage(&message.Acknowledgement{SequenceNumber: uint32(bytesRead)}, 2)
			lastAck = bytesRead
		}

		if msg != nil {
			switch m := msg.(type) {
			case *message.SetChunkSize:
				for _, cs := range chunkStream {
					cs.MaxChunkSize = m.ChunkSize
				}
			case *message.UserControlMessage:
				if m.Event == message.UserControlPingRequest {
					pong := *m
					pong.Event = message.UserControlPingResponse
					c.WriteMessage(&pong, 2)
				}
				// TODO -- we may need to act on other message types here
			}
			c.inboundQueue <- msg
		}

	}
}

func (c *connection) getError() error {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.err
}

func (c *connection) setError(where string, err error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.err = fmt.Errorf("%s: %w", where, err)
	c.outboundQueue.Close()
	c.conn.Close()
}

type handshakeMessage struct {
	Time   uint32
	Time2  uint32
	Random [1528]byte
}

func (c *connection) clientHandshake() error {
	start := time.Now()

	// Write C0 (protocol version)
	_, err := c.conn.Write([]byte{3})
	if err != nil {
		return err
	}

	// Write C1
	var c1 handshakeMessage
	rand.Read(c1.Random[:])
	err = binary.Write(c.conn, binary.BigEndian, &c1)
	if err != nil {
		return err
	}

	// Wait for S0 / S1
	var serverVersion uint8
	err = binary.Read(c.conn, binary.BigEndian, &serverVersion)
	if err != nil {
		return err
	}
	if serverVersion != 3 {
		return ErrInvalidVersion
	}
	var s1 handshakeMessage
	err = binary.Read(c.conn, binary.BigEndian, &s1)
	if err != nil {
		return err
	}

	// Send C2
	c2 := handshakeMessage{
		Time:   s1.Time,
		Time2:  uint32(time.Since(start) / time.Millisecond),
		Random: s1.Random,
	}
	err = binary.Write(c.conn, binary.BigEndian, &c2)
	if err != nil {
		return err
	}

	// Wait for S2
	var s2 handshakeMessage
	err = binary.Read(c.conn, binary.BigEndian, &s2)
	if err != nil {
		return err
	}
	if s2.Time != c1.Time || !reflect.DeepEqual(s2.Random, c1.Random) {
		return ErrHandshakeMismatch
	}

	return nil
}

func (c *connection) serverHandshake() error {
	start := time.Now()

	// Read client version (C0)
	var clientVersion uint8
	err := binary.Read(c.conn, binary.BigEndian, &clientVersion)
	if err != nil {
		return err
	}
	if clientVersion != 3 {
		return ErrInvalidVersion
	}

	// Send server version (S0)
	_, err = c.conn.Write([]byte{3})
	if err != nil {
		return err
	}

	// Send S1
	var s1 handshakeMessage
	rand.Read(s1.Random[:])
	err = binary.Write(c.conn, binary.BigEndian, &s1)
	if err != nil {
		return err
	}

	// Read C1
	var c1 handshakeMessage
	err = binary.Read(c.conn, binary.BigEndian, &c1)
	if err != nil {
		return err
	}

	// Send S2
	s2 := handshakeMessage{
		Time:   c1.Time,
		Time2:  uint32(time.Since(start) / time.Millisecond),
		Random: c1.Random,
	}
	err = binary.Write(c.conn, binary.BigEndian, &s2)
	if err != nil {
		return err
	}

	// Read C2
	var c2 handshakeMessage
	err = binary.Read(c.conn, binary.BigEndian, &c2)
	if err != nil {
		return err
	}
	if c2.Time != s1.Time || !reflect.DeepEqual(c2.Random, s1.Random) {
		return ErrHandshakeMismatch
	}

	return nil
}
