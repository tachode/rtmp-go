package message

import (
	"github.com/tachode/rtmp-go/amf0"
)

// ToStringMap extracts a map[string]any from a value that implements
// message.Object or is an amf0.EcmaArray. It returns the map and true
// if the conversion succeeded, or nil and false otherwise.
func ToStringMap(v any) (map[string]any, bool) {
	switch m := v.(type) {
	case amf0.Object:
		return map[string]any(m), true
	case amf0.EcmaArray:
		return map[string]any(m), true
	case map[string]any:
		return m, true
	default:
		return nil, false
	}
}
