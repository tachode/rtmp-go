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

func GetStringSlice(obj message.Object, key string) []string {
	val, found := obj.Get(key)
	if !found {
		return nil
	}
	items, ok := message.ToSlice(val)
	if !ok {
		return nil
	}
	result := make([]string, 0, len(items))
	for _, item := range items {
		if s, ok := message.ToString(item); ok {
			result = append(result, s)
		}
	}
	if len(result) == 0 {
		return nil
	}
	return result
}

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
