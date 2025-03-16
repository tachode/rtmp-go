package message

import (
	"encoding/binary"
	"fmt"
)

type WindowAcknowledgementSize struct {
	MetadataFields
	AcknowledgementWindowSize uint32
}

func init() { RegisterType(new(WindowAcknowledgementSize)) }

func (m WindowAcknowledgementSize) Type() Type {
	return TypeWindowAcknowledgementSize
}

func (m WindowAcknowledgementSize) Marshal() ([]byte, error) {
	out := make([]byte, 4)
	binary.BigEndian.PutUint32(out, m.AcknowledgementWindowSize)
	return out, nil
}

func (m *WindowAcknowledgementSize) Unmarshal(data []byte) error {
	if len(data) < 4 {
		return ErrShortMessage
	}
	m.AcknowledgementWindowSize = binary.BigEndian.Uint32(data)
	return nil
}

func (m WindowAcknowledgementSize) String() string {
	return fmt.Sprintf("%v: %+v AcknowledgementWindowSize=%d", m.Type(), m.MetadataFields, m.AcknowledgementWindowSize)
}
