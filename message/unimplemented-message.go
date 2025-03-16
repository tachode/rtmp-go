package message

import (
	"fmt"
)

type UnimplementedMessage struct {
	MetadataFields
	messageType Type
	Payload     []byte
}

func (m UnimplementedMessage) Type() Type {
	return m.messageType
}

func (m UnimplementedMessage) Marshal() ([]byte, error) {
	return m.Payload, nil
}

func (m *UnimplementedMessage) Unmarshal(data []byte) error {
	m.Payload = data
	return nil
}

func (m UnimplementedMessage) String() string {
	return fmt.Sprintf("%v (unimplemented): %+v len=%d", m.Type(), m.MetadataFields, len(m.Payload))
}
