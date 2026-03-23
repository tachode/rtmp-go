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
