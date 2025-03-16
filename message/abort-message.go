package message

import (
	"encoding/binary"
	"fmt"
)

type AbortMessage struct {
	MetadataFields
	ChunkStreamId uint32
}

func init() { RegisterType(new(AbortMessage)) }

func (m AbortMessage) Type() Type {
	return TypeAbortMessage
}

func (m AbortMessage) Marshal() ([]byte, error) {
	out := make([]byte, 4)
	binary.BigEndian.PutUint32(out, m.ChunkStreamId)
	return out, nil
}

func (m *AbortMessage) Unmarshal(data []byte) error {
	if len(data) < 4 {
		return ErrShortMessage
	}
	m.ChunkStreamId = binary.BigEndian.Uint32(data)
	return nil
}

func (m AbortMessage) String() string {
	return fmt.Sprintf("%v: %+v ChunkStreamId=%d", m.Type(), m.MetadataFields, m.ChunkStreamId)
}
