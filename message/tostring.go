package message

import (
	"github.com/tachode/rtmp-go/amf0"
	"github.com/tachode/rtmp-go/amf3"
)

// ToString extracts a Go string from a value that may be a plain string,
// amf0.String, or amf3.String. It returns the string and true if the
// conversion succeeded, or "" and false otherwise.
func ToString(v any) (string, bool) {
	switch s := v.(type) {
	case string:
		return s, true
	case amf0.String:
		return string(s), true
	case amf3.String:
		return string(s), true
	default:
		return "", false
	}
}
