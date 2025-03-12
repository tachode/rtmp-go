package amf0

import (
	"encoding/binary"
	"fmt"
	"io"
	"time"
)

func Write(w io.Writer, value any) error {
	if value == nil {
		value = new(Null)
	}

	// Convenience handling to convert common go types into Values
	switch v := value.(type) {
	case float32:
		value = Number(v)
	case float64:
		value = Number(v)
	case int:
		value = Number(v)
	case int8:
		value = Number(v)
	case int16:
		value = Number(v)
	case int32:
		value = Number(v)
	case int64:
		value = Number(v)
	case uint:
		value = Number(v)
	case uint8:
		value = Number(v)
	case uint16:
		value = Number(v)
	case uint32:
		value = Number(v)
	case uint64:
		value = Number(v)
	case bool:
		value = Boolean(v)
	case string:
		if len(v) < 0x10000 {
			value = String(v)
		} else {
			value = LongString(v)
		}
	case time.Time:
		value = Date(v)
	}

	if v, ok := value.(Value); ok {
		err := binary.Write(w, binary.BigEndian, v.Type())
		if err != nil {
			return err
		}
		return v.Write(w)
	} else {
		return fmt.Errorf("cannot convert type %T into AMF 0", value)
	}
}
