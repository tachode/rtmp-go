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
	return message.ReadFromCommand(cmd, s)
}

func (s *Seek) ToMessageCommand() (message.Command, error) {
	return message.BuildCommand(s.CommandName(), s), nil
}

func (s *Seek) MakeStatus(status Status) message.Command {
	return streamStatusResponse(s.StreamId, s.Transaction, status)
}
