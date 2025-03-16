package message

import (
	"encoding/binary"
	"fmt"
)

type Acknowledgement struct {
	MetadataFields
	SequenceNumber uint32
}

func init() { RegisterType(new(Acknowledgement)) }

func (m Acknowledgement) Type() Type {
	return TypeAcknowledgement
}

func (m Acknowledgement) Marshal() ([]byte, error) {
	out := make([]byte, 4)
	binary.BigEndian.PutUint32(out, m.SequenceNumber)
	return out, nil
}

func (m *Acknowledgement) Unmarshal(data []byte) error {
	if len(data) < 4 {
		return ErrShortMessage
	}
	m.SequenceNumber = binary.BigEndian.Uint32(data)
	return nil
}

func (m Acknowledgement) String() string {
	return fmt.Sprintf("%v: %+v SequenceNumber=%d", m.Type(), m.MetadataFields, m.SequenceNumber)
}
