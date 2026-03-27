package message

import (
	"reflect"
	"strconv"
	"strings"

	"github.com/tachode/rtmp-go/amf0"
)

// parseTag splits an amf struct tag into its components.
// The tag format is: "name[|alias...][,omitempty]"
// Names are separated by pipes; the optional ",omitempty" flag is at the end.
func parseTag(tag string) (names []string, omitempty bool) {
	if i := strings.Index(tag, ",omitempty"); i >= 0 {
		omitempty = true
		tag = tag[:i]
	}
	names = strings.Split(tag, "|")
	return
}

// ReadFields populates the struct pointed to by target from the given
// Object, using `amf` struct tags as property names. Tags may
// contain pipe-separated aliases (e.g. `amf:"fileSize|filesize"`); the
// first matching name wins.
//
// Supported field types: float64, uint16, uint32, int, bool, string, *bool,
// *struct (with amf tags), []string, and map[int]T (where T is a struct with amf tags).
func ReadFields(obj Object, target any) {
	v := reflect.ValueOf(target).Elem()
	t := v.Type()
	for i := range t.NumField() {
		tag := t.Field(i).Tag.Get("amf")
		if tag == "" {
			continue
		}
		names, _ := parseTag(tag)
		fv := v.Field(i)
		switch fv.Kind() {
		case reflect.Float64:
			for _, name := range names {
				if val := GetFloat64(obj, name); val != 0 {
					fv.SetFloat(val)
					break
				}
			}
		case reflect.Uint16:
			for _, name := range names {
				if val := GetFloat64(obj, name); val != 0 {
					fv.SetUint(uint64(val))
					break
				}
			}
		case reflect.Uint32:
			for _, name := range names {
				if val := GetFloat64(obj, name); val != 0 {
					fv.SetUint(uint64(val))
					break
				}
			}
		case reflect.Int:
			for _, name := range names {
				if val := GetFloat64(obj, name); val != 0 {
					fv.SetInt(int64(val))
					break
				}
			}
		case reflect.Bool:
			for _, name := range names {
				if val := GetBool(obj, name); val {
					fv.SetBool(val)
					break
				}
			}
		case reflect.String:
			for _, name := range names {
				if val := GetString(obj, name); val != "" {
					fv.SetString(val)
					break
				}
			}
		case reflect.Pointer:
			elemKind := fv.Type().Elem().Kind()
			if elemKind == reflect.Bool {
				for _, name := range names {
					if bp := GetBoolPtr(obj, name); bp != nil {
						fv.Set(reflect.ValueOf(bp))
						break
					}
				}
			} else if elemKind == reflect.Struct {
				for _, name := range names {
					if sub := getObject(obj, name); sub != nil {
						elem := reflect.New(fv.Type().Elem())
						ReadFields(sub, elem.Interface())
						fv.Set(elem)
						break
					}
				}
			}
		case reflect.Slice:
			if fv.Type().Elem().Kind() == reflect.String {
				for _, name := range names {
					if val := GetStringSlice(obj, name); val != nil {
						fv.Set(reflect.ValueOf(val))
						break
					}
				}
			}
		case reflect.Map:
			if fv.Type().Key().Kind() == reflect.Int {
				for _, name := range names {
					readTrackIdInfoMap(obj, name, fv)
					if !fv.IsNil() {
						break
					}
				}
			}
		}
	}
}

// readTrackIdInfoMap reads a map[int]T field from an Object property,
// where T is a struct with amf tags on its fields.
func readTrackIdInfoMap(obj Object, key string, fv reflect.Value) {
	m := GetStringMap(obj, key)
	if m == nil {
		return
	}
	elemType := fv.Type().Elem()
	mapVal := reflect.MakeMap(fv.Type())
	for k, v := range m {
		id, err := strconv.Atoi(k)
		if err != nil {
			continue
		}
		trackObj, ok := v.(Object)
		if !ok {
			continue
		}
		elem := reflect.New(elemType)
		ReadFields(trackObj, elem.Interface())
		mapVal.SetMapIndex(reflect.ValueOf(id), elem.Elem())
	}
	if mapVal.Len() > 0 {
		fv.Set(mapVal)
	}
}

// WriteFields serializes the struct into an amf0.EcmaArray, using `amf`
// struct tags as property names. When a tag contains pipe-separated aliases,
// only the first name is used for serialization. Fields with zero values
// are included unless the tag contains ",omitempty".
//
// Supported field types: float64, uint16, uint32, int, bool, string, *bool,
// *struct (with amf tags), []string, and map[int]T (where T is a struct with amf tags).
func WriteFields(source any) amf0.EcmaArray {
	props := amf0.EcmaArray{}
	v := reflect.ValueOf(source)
	if v.Kind() == reflect.Pointer {
		v = v.Elem()
	}
	t := v.Type()
	for i := range t.NumField() {
		tag := t.Field(i).Tag.Get("amf")
		if tag == "" {
			continue
		}
		names, omitempty := parseTag(tag)
		name := names[0]
		fv := v.Field(i)
		switch fv.Kind() {
		case reflect.Float64:
			if !omitempty || fv.Float() != 0 {
				props[name] = fv.Float()
			}
		case reflect.Uint16:
			if !omitempty || fv.Uint() != 0 {
				props[name] = float64(fv.Uint())
			}
		case reflect.Uint32:
			if !omitempty || fv.Uint() != 0 {
				props[name] = float64(fv.Uint())
			}
		case reflect.Int:
			if !omitempty || fv.Int() != 0 {
				props[name] = float64(fv.Int())
			}
		case reflect.Bool:
			if !omitempty || fv.Bool() {
				props[name] = fv.Bool()
			}
		case reflect.String:
			if !omitempty || fv.String() != "" {
				props[name] = fv.String()
			}
		case reflect.Pointer:
			if !fv.IsNil() {
				elemKind := fv.Type().Elem().Kind()
				switch elemKind {
				case reflect.Bool:
					props[name] = fv.Elem().Bool()
				case reflect.Struct:
					props[name] = WriteFields(fv.Interface())
				}
			}
		case reflect.Slice:
			if fv.Type().Elem().Kind() == reflect.String {
				if !omitempty || fv.Len() > 0 {
					arr := make(amf0.StrictArray, fv.Len())
					for j := range fv.Len() {
						arr[j] = fv.Index(j).String()
					}
					props[name] = arr
				}
			}
		case reflect.Map:
			if fv.Len() > 0 && fv.Type().Key().Kind() == reflect.Int {
				innerMap := make(amf0.EcmaArray, fv.Len())
				for _, key := range fv.MapKeys() {
					innerMap[strconv.Itoa(int(key.Int()))] = WriteFields(fv.MapIndex(key).Interface())
				}
				props[name] = innerMap
			}
		}
	}
	return props
}

// getObject extracts an Object value from an Object by key.
func getObject(obj Object, key string) Object {
	val, found := obj.Get(key)
	if !found {
		return nil
	}
	// Try direct Object interface
	if o, ok := val.(Object); ok {
		return o
	}
	// Try map[string]any (from ToStringMap conversions)
	if m, ok := ToStringMap(val); ok {
		return mapObject(m)
	}
	return nil
}

// mapObject wraps a map[string]any to implement the Object interface.
type mapObject map[string]any

func (m mapObject) Get(key string) (any, bool) {
	v, ok := m[key]
	return v, ok
}
