package amf0

import (
	"encoding/binary"
	"io"
)

type StrictArray []any

func init() { RegisterType(new(StrictArray)) }

func (v StrictArray) Type() Type { return StrictArrayMarker }

func (v *StrictArray) Read(r io.Reader) error {
	var length uint32
	err := binary.Read(r, binary.BigEndian, &length)
	if err != nil {
		return err
	}
	*v = make([]any, length)
	for i := 0; i < int(length); i++ {
		(*v)[i], err = Read(r)
		if err != nil {
			return err
		}
	}
	return nil
}

func (v StrictArray) Write(w io.Writer) error {
	err := binary.Write(w, binary.BigEndian, uint32(len(v)))
	if err != nil {
		return err
	}
	for _, value := range v {
		err = Write(w, value)
		if err != nil {
			return err
		}
	}
	return nil
}
