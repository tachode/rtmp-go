package command

import "github.com/tachode/rtmp-go/message"

// NetStream.seek() command
// Seeks to the offset (in milliseconds) within a media file or playlist.

func init() { RegisterCommand(new(Seek)) }

type Seek struct {
	StreamId     int
	Transaction  int
	MilliSeconds float64 `amfParameter:"0"` // Number of milliseconds to seek into the playlist.
}

func (s Seek) CommandName() string { return "seek" }

func (s *Seek) FromMessageCommand(cmd message.Command) error {
	message.ReadFromCommand(cmd, s)
	return nil
}

func (s *Seek) ToMessageCommand() (message.Command, error) {
	return message.BuildCommand(s.CommandName(), s), nil
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
