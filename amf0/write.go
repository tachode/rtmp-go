package amf0

import (
	"encoding/binary"
	"fmt"
	"io"
	"reflect"
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
	default:
		// Handle named types whose underlying type is string or numeric
		// (e.g., type Level string, type EnumName int)
		// Only apply if the value doesn't already implement Value.
		if _, ok := value.(Value); !ok {
			rv := reflect.ValueOf(value)
			switch rv.Kind() {
			case reflect.String:
				s := rv.String()
				if len(s) < 0x10000 {
					value = String(s)
				} else {
					value = LongString(s)
				}
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				value = Number(rv.Int())
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				value = Number(rv.Uint())
			case reflect.Float32, reflect.Float64:
				value = Number(rv.Float())
			}
		}
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
