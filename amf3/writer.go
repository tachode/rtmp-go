package amf3

import (
	"encoding/binary"
	"fmt"
	"io"
	"time"
)

// Writer serializes AMF 3 data, maintaining reference tables for strings
// and traits as required by the spec (§2.2).
//
// Note: write-side object deduplication is not currently implemented.
// The spec allows complex objects (Array, Object, Date, XML, XMLDocument,
// ByteArray) to be sent by reference, but this Writer always sends them
// inline. Strings and traits are still deduplicated normally.
type Writer struct {
	w       io.Writer
	strings map[string]int
	traits  map[string]int // keyed by a canonical trait description
}

// NewWriter creates a new AMF 3 writer wrapping the given io.Writer.
// Reference tables are initialized empty and populated as values are written.
func NewWriter(w io.Writer) *Writer {
	return &Writer{
		w:       w,
		strings: make(map[string]int),
		traits:  make(map[string]int),
	}
}

// SetWriter changes the underlying writer without affecting reference tables.
// This is useful when the same logical AMF 3 stream is split across multiple
// buffers for transport (e.g., chunk streams). Pass nil to release the
// current buffer for garbage collection.
func (w *Writer) SetWriter(newWriter io.Writer) {
	w.w = newWriter
}

// WriteValue writes a single AMF 3 typed value (marker + payload).
func (w *Writer) WriteValue(value any) error {
	if w.w == nil {
		return ErrNilWriter
	}
	if value == nil {
		return w.writeMarker(NullMarker)
	}

	// Convenience handling to convert common Go types into AMF 3 Values
	switch v := value.(type) {
	case float32:
		value = Double(v)
	case float64:
		value = Double(v)
	case int:
		value = w.intOrDouble(int64(v))
	case int8:
		value = w.intOrDouble(int64(v))
	case int16:
		value = w.intOrDouble(int64(v))
	case int32:
		value = w.intOrDouble(int64(v))
	case int64:
		value = w.intOrDouble(int64(v))
	case uint:
		value = w.uintOrDouble(uint64(v))
	case uint8:
		value = w.uintOrDouble(uint64(v))
	case uint16:
		value = w.uintOrDouble(uint64(v))
	case uint32:
		value = w.uintOrDouble(uint64(v))
	case uint64:
		value = w.uintOrDouble(uint64(v))
	case bool:
		value = Boolean(v)
	case string:
		value = String(v)
	case time.Time:
		value = Date(v)
	case []byte:
		value = ByteArray(v)
	}

	if v, ok := value.(Value); ok {
		err := w.writeMarker(v.Type())
		if err != nil {
			return err
		}
		return v.Write(w)
	}
	return fmt.Errorf("cannot convert type %T into AMF 3", value)
}

// intOrDouble returns an Integer if the signed value fits in 29 bits,
// otherwise a Double. Per spec §3.6: signed int >= 2^28 becomes double.
func (w *Writer) intOrDouble(v int64) Value {
	if v >= -0x10000000 && v <= 0x0FFFFFFF {
		return Integer(v)
	}
	return Double(v)
}

// uintOrDouble returns an Integer if the unsigned value fits in 29 bits,
// otherwise a Double. Per spec §3.6: unsigned uint >= 2^29 becomes double.
func (w *Writer) uintOrDouble(v uint64) Value {
	if v <= 0x1FFFFFFF {
		return Integer(v)
	}
	return Double(v)
}

func (w *Writer) writeMarker(t Type) error {
	return binary.Write(w.w, binary.BigEndian, t)
}

// writeString writes a UTF-8-vr encoded string (§1.3.2, §3.8).
// If the string was previously written, a reference is sent instead.
// The empty string is always sent as a literal (0x01) and never referenced.
func (w *Writer) writeString(s string) error {
	if s == "" {
		return writeU29(w.w, 0x01) // UTF-8-empty
	}

	if idx, ok := w.strings[s]; ok {
		return writeU29(w.w, uint32(idx<<1)) // reference: low bit = 0
	}

	w.strings[s] = len(w.strings)

	length := uint32(len(s))
	err := writeU29(w.w, (length<<1)|1) // literal: low bit = 1
	if err != nil {
		return err
	}
	_, err = w.Write([]byte(s))
	return err
}

// writeTraitRef checks if the trait has been written before. If so, it
// returns the index and true. Otherwise it returns false.
func (w *Writer) writeTraitRef(key string) (int, bool) {
	idx, ok := w.traits[key]
	return idx, ok
}

// addTraitRef adds a trait to the traits reference table.
func (w *Writer) addTraitRef(key string) {
	w.traits[key] = len(w.traits)
}

// Write passes through to the underlying writer.
func (w *Writer) Write(p []byte) (int, error) {
	return w.w.Write(p)
}
