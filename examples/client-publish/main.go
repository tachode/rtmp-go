package main

import (
	"log"
	"net"

	rtmp "github.com/tachode/rtmp-go"
	"github.com/tachode/rtmp-go/command"
	"github.com/tachode/rtmp-go/examples"
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

	// Send publish
	sendCommand(rtmpConn, 3, &command.Publish{
		StreamId:     1,
		Transaction:  0,
		StreamKey:    "live",
		HowToPublish: command.HowToPublishLive,
	})

	// Wait for onStatus indicating publish started
	waitForPublishStart(rtmpConn)

	// Stream blank media until the connection is closed
	stop := make(chan struct{})
	defer close(stop)
	examples.StreamBlankMedia(rtmpConn, 3, 1, stop)
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

func waitForPublishStart(c rtmp.Conn) {
	for {
		msg, err := c.ReadMessage()
		if err != nil {
			log.Fatal("Connection closed while waiting for publish status:", err)
		}
		log.Printf("<<< %v", msg)

		if cmd, ok := msg.(message.Command); ok {
			if cmd.GetCommand() == "onStatus" {
				parsed, err := command.FromMessageCommand(cmd)
				if err == nil {
					if onStatus, ok := parsed.(*command.OnStatus); ok {
						log.Printf("Received onStatus: %s", onStatus.Code)
						if onStatus.Code == command.NetStreamPublishStart {
							return
						}
					}
				}
			}
		}

		if ucm, ok := msg.(*message.UserControlMessage); ok {
			event, err := usercontrol.FromMessage(ucm)
			if err == nil {
				log.Printf("<<< %T: %+v", event, event)
			}
		}
	}
}
