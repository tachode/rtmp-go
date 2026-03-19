package amf3

import (
	"fmt"
	"io"
	"reflect"
)

// Reader deserializes AMF 3 data, maintaining reference tables for strings,
// objects, and traits as required by the spec (§2.2).
type Reader struct {
	r       io.Reader
	strings []string
	objects []any
	traits  []*TraitInfo
}

// NewReader creates a new AMF 3 reader wrapping the given io.Reader.
// Reference tables are initialized empty and populated as values are read.
func NewReader(r io.Reader) *Reader {
	return &Reader{r: r}
}

// ReadValue reads a single AMF 3 typed value (marker + payload).
func (r *Reader) ReadValue() (any, error) {
	var marker [1]byte
	_, err := io.ReadFull(r.r, marker[:])
	if err != nil {
		return nil, err
	}

	// Boolean values are encoded entirely in the marker (§3.4, §3.5).
	switch Type(marker[0]) {
	case FalseMarker:
		return Boolean(false), nil
	case TrueMarker:
		return Boolean(true), nil
	}

	prototype, ok := typeRegistry[Type(marker[0])]
	if !ok {
		return nil, fmt.Errorf("unknown AMF 3 type marker 0x%02x", marker[0])
	}

	instance := reflect.New(reflect.Indirect(reflect.ValueOf(prototype)).Type()).Interface()
	readableOut, ok := instance.(MutableValue)
	if !ok {
		return nil, fmt.Errorf("invalid registered type %T does not implement MutableValue interface", prototype)
	}

	err = readableOut.Read(r)
	if err != nil {
		return nil, err
	}

	// Types that participate in the object reference table are returned as
	// pointers so that the reference table and all holders of the value
	// see the same object (important for circular references).
	if _, ok := readableOut.(objectRefType); ok {
		return readableOut, nil
	}

	// Simple value types (Integer, Double, String, Null, Undefined) are
	// dereferenced and returned as values, matching the AMF0 convention.
	val, ok := reflect.Indirect(reflect.ValueOf(readableOut)).Interface().(Value)
	if !ok {
		return nil, fmt.Errorf("invalid registered type %T does not implement Value interface", prototype)
	}
	return val, nil
}

// objectRefType is implemented by types that participate in the AMF 3
// object reference table (§2.2). ReadValue returns these types as pointers.
type objectRefType interface {
	isObjectRefType()
}

// readString reads a UTF-8-vr encoded string (§1.3.2, §3.8).
// The low bit of the U29 header is a flag:
//   - 0: string reference (remaining bits = index into string table)
//   - 1: string literal (remaining bits = byte-length of UTF-8 string)
//
// The empty string is never sent by reference.
func (r *Reader) readString() (string, error) {
	header, err := readU29(r.r)
	if err != nil {
		return "", err
	}

	// Low bit = 0: string reference
	if header&1 == 0 {
		index := header >> 1
		if int(index) >= len(r.strings) {
			return "", fmt.Errorf("AMF3 string reference index %d out of range (table size %d)", index, len(r.strings))
		}
		return r.strings[index], nil
	}

	// Low bit = 1: string literal
	length := header >> 1
	if length == 0 {
		return "", nil // empty string, not added to reference table
	}

	buf := make([]byte, length)
	_, err = io.ReadFull(r.r, buf)
	if err != nil {
		return "", err
	}

	s := string(buf)
	r.strings = append(r.strings, s)
	return s, nil
}

// readObjectRef attempts to read an object reference. If the low bit of the
// U29 header is 0, it returns the referenced object and true. Otherwise it
// returns the remaining header bits (shifted right by 1) and false.
func (r *Reader) readObjectRef() (any, uint32, bool, error) {
	header, err := readU29(r.r)
	if err != nil {
		return nil, 0, false, err
	}

	if header&1 == 0 {
		index := header >> 1
		if int(index) >= len(r.objects) {
			return nil, 0, false, fmt.Errorf("AMF3 object reference index %d out of range (table size %d)", index, len(r.objects))
		}
		return r.objects[index], 0, true, nil
	}

	return nil, header >> 1, false, nil
}

// addObjectRef adds an object to the object reference table.
func (r *Reader) addObjectRef(obj any) {
	r.objects = append(r.objects, obj)
}

// addTraitRef adds a trait to the traits reference table and returns its index.
func (r *Reader) addTraitRef(trait *TraitInfo) {
	r.traits = append(r.traits, trait)
}

// getTraitRef returns the trait at the given index.
func (r *Reader) getTraitRef(index uint32) (*TraitInfo, error) {
	if int(index) >= len(r.traits) {
		return nil, fmt.Errorf("AMF3 trait reference index %d out of range (table size %d)", index, len(r.traits))
	}
	return r.traits[index], nil
}
