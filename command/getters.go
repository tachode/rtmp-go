package command

import "github.com/tachode/rtmp-go/message"

// GetString is shorthand for message.GetString.
var GetString = message.GetString

// GetFloat64 is shorthand for message.GetFloat64.
var GetFloat64 = message.GetFloat64

// GetBool is shorthand for message.GetBool.
var GetBool = message.GetBool

// GetStringSlice is shorthand for message.GetStringSlice.
var GetStringSlice = message.GetStringSlice

// GetStringMap is shorthand for message.GetStringMap.
var GetStringMap = message.GetStringMap

func GetFourCcInfoMap(obj message.Object, key string) FourCcInfoMap {
	val, found := obj.Get(key)
	if !found {
		return nil
	}
	m, ok := message.ToStringMap(val)
	if !ok {
		return nil
	}
	result := make(FourCcInfoMap, len(m))
	for k, v := range m {
		if f, ok := message.ToFloat64(v); ok {
			result[k] = FourCcInfoMask(f)
		}
	}
	if len(result) == 0 {
		return nil
	}
	return result
}
