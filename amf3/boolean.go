package amf3

// Boolean represents the AMF 3 true/false types (§3.4, §3.5).
// In AMF 3, the boolean value is encoded directly in the type marker:
// false-marker (0x02) or true-marker (0x03). No additional payload.
//
// Boolean is not registered in the type registry because its value is
// determined by which marker is read. Instead, Reader.ReadValue() handles
// FalseMarker and TrueMarker as special cases.
type Boolean bool

// Type returns TrueMarker or FalseMarker depending on the boolean value.
func (v Boolean) Type() Type {
	if v {
		return TrueMarker
	}
	return FalseMarker
}

func (v Boolean) Write(w *Writer) error {
	return nil
}
