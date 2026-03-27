package message

import (
	"github.com/tachode/rtmp-go/amf0"
	"github.com/tachode/rtmp-go/amf3"
)

// ToSlice extracts a Go []any from a value that may be an amf0.StrictArray
// or amf3.Array. It returns the slice and true if the conversion succeeded,
// or nil and false otherwise.
func ToSlice(v any) ([]any, bool) {
	switch s := v.(type) {
	case amf0.StrictArray:
		return []any(s), true
	case *amf3.Array:
		return s.Dense, true
	case []any:
		return s, true
	default:
		return nil, false
	}
}
