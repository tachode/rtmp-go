package message

import (
	"github.com/tachode/rtmp-go/amf0"
	"github.com/tachode/rtmp-go/amf3"
)

// ToFloat64 extracts a Go float64 from a value that may be a plain numeric
// type, amf0.Number, amf3.Integer, or amf3.Double. It returns the value and
// true if the conversion succeeded, or 0 and false otherwise.
func ToFloat64(v any) (float64, bool) {
	switch n := v.(type) {
	case float64:
		return n, true
	case float32:
		return float64(n), true
	case int:
		return float64(n), true
	case int32:
		return float64(n), true
	case int64:
		return float64(n), true
	case amf0.Number:
		return float64(n), true
	case amf3.Integer:
		return float64(n), true
	case amf3.Double:
		return float64(n), true
	default:
		return 0, false
	}
}
