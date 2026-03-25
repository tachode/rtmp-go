package command

import "github.com/tachode/rtmp-go/message"

// NetStream.seek() command
// Seeks to the offset (in milliseconds) within a media file or playlist.

func init() { RegisterCommand(new(Seek)) }

type Seek struct {
	StreamId     int
	Transaction  int
	MilliSeconds float64 // Number of milliseconds to seek into the playlist.
}

func (s Seek) CommandName() string { return "seek" }

func (s *Seek) FromMessageCommand(cmd message.Command) error {
	s.StreamId = int(cmd.Metadata().StreamId)
	s.Transaction = int(cmd.GetTransactionId())
	params := cmd.GetParameters()
	if len(params) > 0 {
		if n, ok := message.ToFloat64(params[0]); ok {
			s.MilliSeconds = n
		}
	}
	return nil
}

func (s *Seek) ToMessageCommand() (message.Command, error) {
	cmd := &message.Amf0CommandMessage{
		MetadataFields: message.MetadataFields{
			StreamId: uint32(s.StreamId),
		},
		Command:       s.CommandName(),
		TransactionId: float64(s.Transaction),
		Parameters:    []any{s.MilliSeconds},
	}
	return cmd, nil
}

func (s *Seek) MakeStatus(status Status) message.Command {
	command := "onStatus"
	if status.Level == LevelError {
		command = "_error"
	}

	cmd := &message.Amf0CommandMessage{
		MetadataFields: message.MetadataFields{
			StreamId: uint32(s.StreamId),
		},
		Command:       command,
		TransactionId: float64(s.Transaction),
		Object:        nil,
		Parameters:    []any{status.ToObject()},
	}
	return cmd
}
