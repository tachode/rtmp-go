package message

// MarshalText implements encoding.TextMarshaler so that JSON encoding
// uses the human-readable String() representation.
func (i AudioCodecId) MarshalText() ([]byte, error) {
	return []byte(i.String()), nil
}

// MarshalText implements encoding.TextMarshaler so that JSON encoding
// uses the human-readable String() representation.
func (i VideoCodecId) MarshalText() ([]byte, error) {
	return []byte(i.String()), nil
}
