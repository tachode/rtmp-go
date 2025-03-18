package rtmp

import (
	"github.com/tachode/rtmp-go/chunkstream"
	"github.com/tachode/rtmp-go/internal/connection"
	"github.com/tachode/rtmp-go/message"
)

// These are aliases of errors from the sub-modules, to make it easier for applications to refer to them
var (
	ErrNoSuchChunkstream                    = connection.ErrNoSuchChunkstream
	ErrChunkStreamAlreadyExists             = connection.ErrChunkStreamAlreadyExists
	ErrConnectionClosed                     = connection.ErrConnectionClosed
	ErrInvalidVersion                       = connection.ErrInvalidVersion
	ErrHandshakeMismatch                    = connection.ErrHandshakeMismatch
	ErrShortMessage                         = message.ErrShortMessage
	ErrInvalidChunkSize                     = chunkstream.ErrInvalidChunkSize
	ErrInvalidChunkStreamId                 = chunkstream.ErrInvalidChunkStreamId
	ErrDeltaTimePassedToFullHeader          = chunkstream.ErrDeltaTimePassedToFullHeader
	ErrNonDeltaTimestampPassedToShortHeader = chunkstream.ErrNonDeltaTimestampPassedToShortHeader
)
