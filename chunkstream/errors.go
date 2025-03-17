package chunkstream

import "errors"

var (
	ErrInvalidChunkSize                     = errors.New("invalid chunk size")
	ErrInvalidChunkStreamId                 = errors.New("invalid chunk stream id")
	ErrDeltaTimePassedToFullHeader          = errors.New("delta time passed to full header")
	ErrNonDeltaTimestampPassedToShortHeader = errors.New("non-delta timestamp passed to short header")
)
