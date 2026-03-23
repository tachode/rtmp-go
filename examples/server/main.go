package main

import (
	"log"
	"net"

	rtmp "github.com/tachode/rtmp-go"
	"github.com/tachode/rtmp-go/command"
	"github.com/tachode/rtmp-go/message"
)

func main() {
	listener, err := net.Listen("tcp", ":1935")
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()
	log.Println("RTMP server listening on :1935")

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Accept error:", err)
			continue
		}
		go handleConn(conn)
	}
}

func handleConn(conn net.Conn) {
	log.Printf("New connection from %s", conn.RemoteAddr())

	rtmpConn, err := rtmp.NewServerConn(conn, 3)
	if err != nil {
		conn.Close()
		log.Printf("Handshake error from %s: %v", conn.RemoteAddr(), err)
		return
	}
	defer rtmpConn.Close()
	log.Printf("Server handshake completed")

	err = rtmpConn.CreateOutboundChunkstream(3, 1)
	if err != nil {
		log.Fatal(err)
	}

	for {
		msg, err := rtmpConn.ReadMessage()
		if err != nil {
			log.Printf("Connection %s closed: %v", conn.RemoteAddr(), err)
			return
		}
		log.Printf("<<< %v", msg)

		// Handle messages that require responses
		commandMessage, isCommand := msg.(message.Command)
		if isCommand {
			cmd, err := command.FromMessageCommand(commandMessage)
			if err != nil {
				log.Println(err)
			} else {
				log.Printf("Received command message %T: %+v", cmd, cmd)
				switch c := cmd.(type) {
				case *command.Connect:
					send(rtmpConn, 2, &message.WindowAcknowledgementSize{AcknowledgementWindowSize: 2_500_000})
					send(rtmpConn, 2, &message.SetPeerBandwidth{WindowSize: 2_500_000, LimitType: message.BandwidthLimitDynamic})
					send(rtmpConn, 2, &message.UserControlMessage{Event: message.UserControlStreamBegin, Parameters: []uint32{0}})
					send(rtmpConn, 3, c.MakeResponse(command.NewStatus(command.NetConnectionConnectSuccess), 0))
				case *command.ReleaseStream:
					send(rtmpConn, 3, c.MakeResponse(command.NewStatus(command.NetStreamUnpublishSuccess)))
				case *command.FCPublish:
					send(rtmpConn, 3, c.MakeResponse(command.NewStatus(command.NetStreamPublishStart, "FCPublish Ignored")))
				case *command.CreateStream:
					send(rtmpConn, 3, c.MakeResponse(1))
					send(rtmpConn, 2, &message.UserControlMessage{Event: message.UserControlStreamBegin, Parameters: []uint32{1}})
				case *command.Publish:
					send(rtmpConn, 3, c.MakeResponse(command.NewStatus(command.NetStreamPublishStart), 1))
				}

			}
		}

	}
}

func send(c rtmp.Conn, chunkStream int, m message.Message) {
	log.Printf(">>> %s\n", m)
	err := c.WriteMessage(m, chunkStream)
	if err != nil {
		log.Fatal(err)
	}
}
