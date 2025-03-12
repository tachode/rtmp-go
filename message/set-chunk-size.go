package message

import (
	"encoding/binary"
	"fmt"
)

type SetChunkSize struct {
	MetadataFields
	ChunkSize uint32
}

func init() { RegisterType(new(SetChunkSize)) }

func (m SetChunkSize) Type() Type {
	return TypeSetChunkSize
}

func (m SetChunkSize) Marshal() ([]byte, error) {
	out := make([]byte, 4)
	binary.BigEndian.PutUint32(out, m.ChunkSize&0x7fffffff)
	return out, nil
}

func (m *SetChunkSize) Unmarshal(data []byte) error {
	if len(data) < 4 {
		return ErrShortMessage
	}
	m.ChunkSize = binary.BigEndian.Uint32(data) & 0x7fffffff
	return nil
}

func (m SetChunkSize) String() string {
	return fmt.Sprintf("%v: %+v ChunkSize=%d", m.Type(), m.MetadataFields, m.ChunkSize)
}
