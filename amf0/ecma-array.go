package amf0

import (
	"encoding/binary"
	"io"
)

type EcmaArray map[string]any

func init() { RegisterType(new(EcmaArray)) }

func (v EcmaArray) Type() Type { return EcmaArrayMarker }

func (v *EcmaArray) Read(r io.Reader) error {
	var length uint32
	err := binary.Read(r, binary.BigEndian, &length)
	if err != nil {
		return err
	}

	*v = make(EcmaArray)

	for range length {
		var key String
		err = key.Read(r)
		if err != nil {
			return err
		}
		value, err := Read(r)
		if err != nil {
			return err
		}
		(*v)[string(key)] = value
	}
	return nil
}

func (v EcmaArray) Write(w io.Writer) error {
	err := binary.Write(w, binary.BigEndian, uint32(len(v)))
	if err != nil {
		return err
	}
	for key, val := range v {
		err = writeString[uint16](w, key)
		if err != nil {
			return err
		}
		err = Write(w, val)
		if err != nil {
			return err
		}

	}
	return nil
}
