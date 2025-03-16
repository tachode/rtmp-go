package message

import (
	"encoding/binary"
	"fmt"
)

type UserControlMessageEvent uint16

//go:generate stringer -type=UserControlMessageEvent -trimprefix=UserControl
const (
	UserControlStreamBegin UserControlMessageEvent = iota
	UserControlStreamEOF
	UserControlStreamDry
	UserControlSetBufferLength
	UserControlStreamIsRecorded
	UserControlReserved
	UserControlPingRequest
	UserControlPingResponse

	// These are not described in the RTMP spec, but mentioned at https://en.wikipedia.org/wiki/Real-Time_Messaging_Protocol
	UserControlUDPRequest
	UserControlUDPResponse
	UserControlBandwidthLimit
	UserControlBandwidth
	UserControlThrottleBandwidth
	UserControlStreamCreated
	UserControlStreamDeleted
	UserControlSetReadAccess
	UserControlSetWriteAccess
	UserControlStreamMetaRequest
	UserControlStreamMetaResponse
	UserControlGetSegmentBoundary
	UserControlSetSegmentBoundary
	UserControlOnDisconnect
	UserControlSetCriticalLink
	UserControlDisconnect
	UserControlHashUpdate
	UserControlHashTimeout
	UserControlHashRequest
	UserControlHashResponse
	UserControlCheckBandwidth
	UserControlSetAudioSampleAccess
	UserControlSetVideoSampleAccess
	UserControlThrottleBegin
	UserControlThrottleEnd
	UserControlDRMNotify
	UserControlRTMFPSync
	UserControlQueryIHello
	UserControlForwardIHello
	UserControlRedirectIHello
	UserControlNotifyEOF
	UserControlProxyContinue
	UserControlProxyRemoveUpstream
	UserControlRTMFPSetKeepalives
	UserControlSegmentNotFound UserControlMessageEvent = 46
)

type UserControlMessage struct {
	MetadataFields
	Event      UserControlMessageEvent
	Parameters []uint32
}

func init() { RegisterType(new(UserControlMessage)) }

func (m UserControlMessage) Type() Type {
	return TypeUserControlMessage
}

func (m UserControlMessage) Marshal() ([]byte, error) {
	paramCount := 0
	if m.Parameters != nil {
		paramCount = len(m.Parameters)
	}
	out := make([]byte, 2+4*paramCount)
	binary.BigEndian.PutUint16(out, uint16(m.Event))
	for i, param := range m.Parameters {
		binary.BigEndian.PutUint32(out[2+4*i:], param)
	}
	return out, nil
}

func (m *UserControlMessage) Unmarshal(data []byte) error {
	if len(data) < 2 {
		return ErrShortMessage
	}
	m.Event = UserControlMessageEvent(binary.BigEndian.Uint16(data))
	data = data[2:]
	for len(data) >= 4 {
		param := binary.BigEndian.Uint32(data)
		m.Parameters = append(m.Parameters, param)
		data = data[4:]
	}
	return nil
}

func (m UserControlMessage) String() string {
	return fmt.Sprintf("%v: %+v Event=%d(%v)", m.Type(), m.MetadataFields, m.Event, m.Parameters)
}
