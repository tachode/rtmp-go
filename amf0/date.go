package amf0

import (
	"encoding/binary"
	"io"
	"time"
)

type Date time.Time

func init() { RegisterType(new(Date)) }

func (v Date) Type() Type { return DateMarker }

func (v *Date) Read(r io.Reader) error {
	var value float64
	err := binary.Read(r, binary.BigEndian, &value)
	if err != nil {
		return err
	}

	// AMF 0 spec §2.13: "While the design of this type reserves room for time
	// zone offset information, it should not be filled in, nor used"
	var timezone int16
	err = binary.Read(r, binary.BigEndian, &timezone)
	if err != nil {
		return err
	}

	*v = Date(time.UnixMilli(int64(value)))
	return nil
}

func (v Date) Write(w io.Writer) error {
	err := binary.Write(w, binary.BigEndian, float64(time.Time(v).UnixMilli()))
	if err != nil {
		return err
	}

	// AMF 0 spec §2.13: "While the design of this type reserves room for time
	// zone offset information, it should not be filled in, nor used"
	var timezone int16
	return binary.Write(w, binary.BigEndian, timezone)
}
