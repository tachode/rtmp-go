package message

import (
	"encoding/binary"
	"fmt"
)

type BandwidthLimitType uint8

//go:generate stringer -type=BandwidthLimitType -trimprefix=BandwidthLimit
const (
	BandwidthLimitHard BandwidthLimitType = iota
	BandwidthLimitSoft
	BandwidthLimitDynamic
)

type SetPeerBandwidth struct {
	MetadataFields
	WindowSize uint32
	LimitType  BandwidthLimitType
}

func init() { RegisterType(new(SetPeerBandwidth)) }

func (m SetPeerBandwidth) Type() Type {
	return TypeSetPeerBandwidth
}

func (m SetPeerBandwidth) Marshal() ([]byte, error) {
	out := make([]byte, 5)
	binary.BigEndian.PutUint32(out, m.WindowSize)
	out[4] = byte(m.LimitType)
	return out, nil
}

func (m *SetPeerBandwidth) Unmarshal(data []byte) error {
	if len(data) < 5 {
		return ErrShortMessage
	}
	m.WindowSize = binary.BigEndian.Uint32(data)
	m.LimitType = BandwidthLimitType(data[4])
	return nil
}

func (m SetPeerBandwidth) String() string {
	return fmt.Sprintf("%v: %+v WindowSize=%d", m.Type(), m.MetadataFields, m.WindowSize)
}
