package message

import (
	"bytes"
	"errors"
	"fmt"
	"io"

	"github.com/tachode/rtmp-go/amf0"
)

type Amf0DataMessage struct {
	MetadataFields
	Handler    string
	Parameters []any
}

func init() { RegisterType(new(Amf0DataMessage)) }

func (m Amf0DataMessage) Type() Type {
	return TypeAmf0DataMessage
}

func (m Amf0DataMessage) Marshal() ([]byte, error) {
	out := bytes.NewBuffer(nil)
	amf0.Write(out, amf0.String(m.Handler))
	if m.Parameters != nil {
		for _, param := range m.Parameters {
			amf0.Write(out, param)
		}
	}
	return out.Bytes(), nil
}

func (m *Amf0DataMessage) Unmarshal(data []byte) error {
	buf := bytes.NewBuffer(data)
	var err error
	if m.Handler, err = amf0.ReadString(buf); err != nil {
		return fmt.Errorf("could not read handler string: %w", err)
	}
	m.Parameters = nil
	for buf.Len() > 0 {
		param, err := amf0.Read(buf)
		// some implementations send extra bytes at the end of the message
		if errors.Is(err, io.ErrUnexpectedEOF) {
			break
		}
		if err != nil {
			return fmt.Errorf("could not read parameter: %w", err)
		}
		m.Parameters = append(m.Parameters, param)
	}

	return nil
}

func (m Amf0DataMessage) String() string {
	return fmt.Sprintf("%v: %+v Handler=%v(%+v)", m.Type(),
		m.MetadataFields, m.Handler, m.Parameters)
}
