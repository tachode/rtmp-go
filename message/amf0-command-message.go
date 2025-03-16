package message

import (
	"bytes"
	"errors"
	"fmt"
	"io"

	"github.com/tachode/rtmp-go/amf0"
)

type Amf0CommandMessage struct {
	MetadataFields
	Command       string
	TransactionId float64
	Object        amf0.Object
	Parameters    []any
}

func init() { RegisterType(new(Amf0CommandMessage)) }

func (m Amf0CommandMessage) Type() Type {
	return TypeAmf0CommandMessage
}

func (m Amf0CommandMessage) Marshal() ([]byte, error) {
	out := bytes.NewBuffer(nil)
	amf0.Write(out, amf0.String(m.Command))
	amf0.Write(out, amf0.Number(m.TransactionId))
	if len(m.Object) > 0 {
		amf0.Write(out, m.Object)
	} else {
		amf0.Write(out, nil)
	}
	if m.Parameters != nil {
		for _, param := range m.Parameters {
			amf0.Write(out, param)
		}
	}
	return out.Bytes(), nil
}

func (m *Amf0CommandMessage) Unmarshal(data []byte) error {
	buf := bytes.NewBuffer(data)
	var err error
	if m.Command, err = amf0.ReadString(buf); err != nil {
		return fmt.Errorf("could not read command string: %w", err)
	}
	if m.TransactionId, err = amf0.ReadNumber(buf); err != nil {
		return fmt.Errorf("could not read transaction id: %w", err)
	}
	if m.Object, err = amf0.ReadObject(buf); err != nil {
		return fmt.Errorf("could not read object: %w", err)
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

func (m Amf0CommandMessage) String() string {
	return fmt.Sprintf("%v: %+v Command=%v(tid=%v, obj=%+v, %+v)", m.Type(),
		m.MetadataFields, m.Command, m.TransactionId, m.Object, m.Parameters)
}
