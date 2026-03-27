package message

// GetString extracts a string value from an Object by key.
func GetString(obj Object, key string) string {
	val, found := obj.Get(key)
	if found {
		if s, ok := ToString(val); ok {
			return s
		}
	}
	return ""
}

// GetFloat64 extracts a float64 value from an Object by key.
func GetFloat64(obj Object, key string) float64 {
	val, found := obj.Get(key)
	if found {
		if f, ok := ToFloat64(val); ok {
			return f
		}
	}
	return 0
}

// GetBool extracts a bool value from an Object by key.
func GetBool(obj Object, key string) bool {
	val, found := obj.Get(key)
	if found {
		if b, ok := ToBool(val); ok {
			return b
		}
	}
	return false
}

// GetBoolPtr extracts a *bool from an Object by key.
// Returns nil if the key is not present or cannot be converted to bool.
func GetBoolPtr(obj Object, key string) *bool {
	val, found := obj.Get(key)
	if found {
		if b, ok := ToBool(val); ok {
			return &b
		}
	}
	return nil
}

// GetStringSlice extracts a []string from an Object by key,
// converting each element via ToString.
func GetStringSlice(obj Object, key string) []string {
	val, found := obj.Get(key)
	if !found {
		return nil
	}
	items, ok := ToSlice(val)
	if !ok {
		return nil
	}
	result := make([]string, 0, len(items))
	for _, item := range items {
		if s, ok := ToString(item); ok {
			result = append(result, s)
		}
	}
	if len(result) == 0 {
		return nil
	}
	return result
}

// GetStringMap extracts a map[string]any from an Object by key.
func GetStringMap(obj Object, key string) map[string]any {
	val, found := obj.Get(key)
	if !found {
		return nil
	}
	m, ok := ToStringMap(val)
	if !ok {
		return nil
	}
	return m
}
