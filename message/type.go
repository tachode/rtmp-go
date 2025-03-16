package message

type Type uint8

//go:generate stringer -type=Type -trimprefix=Type
const (
	TypeSetChunkSize              Type = 1
	TypeAbortMessage              Type = 2
	TypeAcknowledgement           Type = 3
	TypeUserControlMessage        Type = 4
	TypeWindowAcknowledgementSize Type = 5
	TypeSetPeerBandwidth          Type = 6
	TypeVirtualControl            Type = 7 // Described as "used between edge server and origin server" but not defined
	TypeAudioMessage              Type = 8
	TypeVideoMessage              Type = 9
	TypeAmf3DataMessage           Type = 15 // Not implemented
	TypeAmf3SharedObjectMessage   Type = 16 // Reserved, not defined in RTMP spec
	TypeAmf3CommandMessage        Type = 17 // Not implemented
	TypeAmf0DataMessage           Type = 18
	TypeAmf0SharedObjectMessage   Type = 19 // Reserved, not defined in RTMP spec
	TypeAmf0CommandMessage        Type = 20
	TypeAggregateMessage          Type = 22 // TODO
	TypeGoAway                    Type = 32 // Defined in https://github.com/facebookarchive/rtmp-go-away

	// Aliases based on https://en.wikipedia.org/wiki/Real-Time_Messaging_Protocol
	TypeSetPacketSize     Type = TypeSetChunkSize
	TypeControlMessage    Type = TypeUserControlMessage
	TypeServerBandwidth   Type = TypeWindowAcknowledgementSize
	TypeClientBandwidth   Type = TypeSetPeerBandwidth
	TypeDataExtended      Type = TypeAmf3DataMessage
	TypeContainerExtended Type = TypeAmf3SharedObjectMessage
	TypeCommandExtended   Type = TypeAmf3CommandMessage
	TypeData              Type = TypeAmf0DataMessage
	TypeContainer         Type = TypeAmf0SharedObjectMessage
	TypeCommand           Type = TypeAmf0CommandMessage
	TypeUdp               Type = 0x15 // Not mentioned in the RTMP spec
	TypePresent           Type = 0x17 // Not mentioned in the RTMP spec
)
