package main

import (
	"log"
	"net"

	rtmp "github.com/tachode/rtmp-go"
	"github.com/tachode/rtmp-go/command"
	"github.com/tachode/rtmp-go/message"
	"github.com/tachode/rtmp-go/usercontrol"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:1935")
	if err != nil {
		log.Fatal(err)
	}

	rtmpConn, err := rtmp.NewClientConn(conn, 3)
	if err != nil {
		conn.Close()
		log.Fatal("Handshake error:", err)
	}
	defer rtmpConn.Close()
	log.Println("Client handshake completed")

	err = rtmpConn.CreateOutboundChunkstream(3, 1)
	if err != nil {
		log.Fatal(err)
	}

	// Send connect
	sendCommand(rtmpConn, 3, &command.Connect{
		Transaction:    1,
		App:            "live",
		FlashVer:       "FMLE/3.0",
		TcUrl:          "rtmp://localhost:1935/live",
		AudioCodecs:    command.SupportSndAAC,
		VideoCodecs:    command.SupportVidH264,
		ObjectEncoding: command.ObjectEncodingAMF0,
	})

	// Read until we get the connect response
	waitForResult(rtmpConn, 1)

	// Send createStream
	sendCommand(rtmpConn, 3, &command.CreateStream{
		Transaction: 2,
	})

	// Read until we get the createStream response
	waitForResult(rtmpConn, 2)

	// Send play
	sendCommand(rtmpConn, 3, &command.Play{
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

func sendCommand(c rtmp.Conn, chunkStream int, cmd command.Command) {
	msg, err := cmd.ToMessageCommand()
	if err != nil {
		log.Fatal("Error creating command message:", err)
	}
	log.Printf(">>> %s", msg)
	err = c.WriteMessage(msg, chunkStream)
	if err != nil {
		log.Fatal("Error sending command:", err)
	}
}

func waitForResult(c rtmp.Conn, transactionId float64) {
	for {
		msg, err := c.ReadMessage()
		if err != nil {
			log.Fatal("Connection closed while waiting for response:", err)
		}
		log.Printf("<<< %v", msg)

		if cmd, ok := msg.(message.Command); ok {
			name := cmd.GetCommand()
			if (name == "_result" || name == "_error") && cmd.GetTransactionId() == transactionId {
				if name == "_error" {
					log.Fatalf("Server returned error for transaction %v", transactionId)
				}
				return
			}
		}
	}
}
