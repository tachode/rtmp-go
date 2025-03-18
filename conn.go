package rtmp

import (
	"net"

	"github.com/tachode/rtmp-go/internal/connection"
	"github.com/tachode/rtmp-go/message"
)

// Convenience (readability) constants for a common High/Medium/Low priority scheme
const (
	HighPriority   = 0
	MediumPriority = 1
	LowPriority    = 2
)

type Conn interface {
	net.Conn
	ReadMessage() (msg message.Message, err error)
	WriteMessage(msg message.Message, chunkStreamId int) error
	CreateOutboundChunkstream(chunkStreamId int, priority int) error
}

func NewClientConn(conn net.Conn, priorityCount int) (Conn, error) {
	return connection.New(conn, priorityCount, connection.Client)
}

func NewServerConn(conn net.Conn, priorityCount int) (Conn, error) {
	return connection.New(conn, priorityCount, connection.Server)
}
