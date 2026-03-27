package message

import "errors"

var (
	ErrShortMessage     = errors.New("message too short")
	ErrInvalidChunkSize = errors.New("invalid chunk size: must be between 1 and 16777215")
)
