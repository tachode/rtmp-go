package amf0

import (
	"encoding/binary"
	"io"
)

func readString[LengthType uint16 | uint32](r io.Reader) (string, error) {
	var length LengthType
	err := binary.Read(r, binary.BigEndian, &length)
	if err != nil {
		return "", err
	}
	value := make([]byte, int(length))
	_, err = io.ReadFull(r, value)
	if err != nil {
		return "", err
	}
	return string(value), nil
}

func writeString[LengthType uint16 | uint32](w io.Writer, value string) error {
	err := binary.Write(w, binary.BigEndian, LengthType(len(value)))
	if err != nil {
		return err
	}
	_, err = w.Write([]byte(value))
	return err
}
