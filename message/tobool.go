package message

import (
	"github.com/tachode/rtmp-go/amf0"
	"github.com/tachode/rtmp-go/amf3"
)

// ToBool extracts a Go bool from a value that may be a plain bool,
// amf0.Boolean, or amf3.Boolean. It returns the value and true if the
// conversion succeeded, or false and false otherwise.
func ToBool(v any) (bool, bool) {
	switch b := v.(type) {
	case bool:
		return b, true
	case amf0.Boolean:
		return bool(b), true
	case amf3.Boolean:
		return bool(b), true
	default:
		return false, false
	}
}
