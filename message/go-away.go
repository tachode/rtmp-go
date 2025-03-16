package message

import (
	"fmt"
)

type GoAway struct {
	MetadataFields
}

func init() { RegisterType(new(GoAway)) }

func (m GoAway) Type() Type {
	return TypeGoAway
}

func (m GoAway) Marshal() ([]byte, error) {
	return []byte{}, nil
}

func (m *GoAway) Unmarshal(data []byte) error {
	return nil
}

func (m GoAway) String() string {
	return fmt.Sprintf("%v: %+v", m.Type(), m.MetadataFields)
}
