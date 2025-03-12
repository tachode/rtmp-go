package message

type Type uint8

//go:generate stringer -type=Type -trimprefix=Type
const (
	TypeSetChunkSize               Type = 1
	TypeAbortMessage               Type = 2
	TypeAcknowledgement            Type = 3
	TypeUserControlMessage         Type = 4
	TypeWindowAcknowledgementSize  Type = 5
	TypeSetPeerBandwidth           Type = 6
	TypeEdgeAndOriginServerCommand Type = 7
	TypeAudioMessage               Type = 8
	TypeVideoMessage               Type = 9
	TypeAmf3DataMessage            Type = 15
	TypeAmf3SharedObjectMessage    Type = 16
	TypeAmf3CommandMessage         Type = 17
	TypeAmf0DataMessage            Type = 18
	TypeAmf0SharedObjectMessage    Type = 19
	TypeAmf0CommandMessage         Type = 20
	TypeAggregateMessage           Type = 22
	TypeGoAway                     Type = 32 // Defined in https://github.com/facebookarchive/rtmp-go-away
)
