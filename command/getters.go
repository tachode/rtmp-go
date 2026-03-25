package command

import "github.com/tachode/rtmp-go/message"

func GetString(obj message.Object, key string) string {
	val, found := obj.Get(key)
	if found {
		if s, ok := message.ToString(val); ok {
			return s
		}
	}
	return ""
}

func GetFloat64(obj message.Object, key string) float64 {
	val, found := obj.Get(key)
	if found {
		if f, ok := message.ToFloat64(val); ok {
			return f
		}
	}
	return 0
}

func GetBool(obj message.Object, key string) bool {
	val, found := obj.Get(key)
	if found {
		if b, ok := message.ToBool(val); ok {
			return b
		}
	}
	return false
}
