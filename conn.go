package rtmp

import (
	"net"

	"github.com/tachode/rtmp-go/connection"
	"github.com/tachode/rtmp-go/message"
)

type Conn interface {
	net.Conn

	ReadMessage() (msg message.Message, err error)
	WriteMessage(msg message.Message, chunkStreamId int) error
	AddChunkStreamId(chunkStreamId int, priority int) error
}

func NewConn(conn net.Conn, priorityCount int) (Conn, error) {
	return connection.New(conn, priorityCount)
}
