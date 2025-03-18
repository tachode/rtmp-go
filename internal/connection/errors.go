package connection

import "errors"

var (
	ErrNoSuchChunkstream        error = errors.New("no such chunk stream")
	ErrChunkStreamAlreadyExists error = errors.New("chunk stream already exists")
	ErrConnectionClosed         error = errors.New("connection closed")
	ErrInvalidVersion           error = errors.New("invalid peer version")
	ErrHandshakeMismatch        error = errors.New("response mismatch in handshake")
)
