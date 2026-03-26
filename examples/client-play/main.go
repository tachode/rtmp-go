package main

import (
	"log"

	"github.com/tachode/rtmp-go/command"
	"github.com/tachode/rtmp-go/examples"
	"github.com/tachode/rtmp-go/message"
	"github.com/tachode/rtmp-go/usercontrol"
)

func main() {
	rtmpConn := examples.ConnectAndCreateStream("localhost:1935")
	defer rtmpConn.Close()

	// Send play
	examples.SendCommand(rtmpConn, 3, &command.Play{
		StreamId:      1,
		Transaction:   0,
		StreamKey:     "test",
		StartPosition: -2,
		Duration:      -1,
		Reset:         true,
	})

	// Read and log incoming messages
	for {
		msg, err := rtmpConn.ReadMessage()
		if err != nil {
			log.Fatal("Connection closed:", err)
		}

		switch m := msg.(type) {
		case message.Command:
			cmd, err := command.FromMessageCommand(m)
			if err != nil {
				log.Printf("<<< %v (command: %s)", msg, m.GetCommand())
			} else {
				log.Printf("<<< %T: %+v", cmd, cmd)
			}
		case *message.AudioMessage:
			log.Printf("<<< Audio: timestamp=%d packetType=%v len=%d",
				m.Metadata().Timestamp, m.PacketType, len(m.Tracks[0].Payload))
		case *message.VideoMessage:
			log.Printf("<<< Video: timestamp=%d frameType=%v packetType=%v len=%d",
				m.Metadata().Timestamp, m.FrameType, m.PacketType, len(m.Tracks[0].Payload))
		case *message.UserControlMessage:
			event, err := usercontrol.FromMessage(m)
			if err != nil {
				log.Printf("<<< UserControl: %v", err)
			} else {
				log.Printf("<<< %T: %+v", event, event)
			}
		default:
			log.Printf("<<< %v", msg)
		}
	}
}
