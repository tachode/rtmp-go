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

		switch m := msg.(type) {
		case message.Command:
			cmd, err := command.FromMessageCommand(m)
			if err != nil {
				log.Println(err)
			} else {
				log.Printf("Received command message %T: %+v", cmd, cmd)
				switch c := cmd.(type) {
				case *command.Connect:
					send(rtmpConn, 2, &message.WindowAcknowledgementSize{AcknowledgementWindowSize: 2_500_000})
					send(rtmpConn, 2, &message.SetPeerBandwidth{WindowSize: 2_500_000, LimitType: message.BandwidthLimitDynamic})
					sendUserControl(rtmpConn, &usercontrol.StreamBegin{StreamID: 0})
					send(rtmpConn, 3, c.MakeResponse(command.NewStatus(command.NetConnectionConnectSuccess), 0))
				case *command.ReleaseStream:
					send(rtmpConn, 3, c.MakeResponse(command.NewStatus(command.NetStreamUnpublishSuccess)))
				case *command.FCPublish:
					send(rtmpConn, 3, c.MakeResponse(command.NewStatus(command.NetStreamPublishStart, "FCPublish Ignored")))
				case *command.CreateStream:
					send(rtmpConn, 3, c.MakeResponse(1))
					sendUserControl(rtmpConn, &usercontrol.StreamBegin{StreamID: 1})
				case *command.Publish:
					send(rtmpConn, 3, c.MakeStatus(command.NewStatus(command.NetStreamPublishStart), 1))
				case *command.Play:
					send(rtmpConn, 3, c.MakeStatus(command.NewStatus(command.NetStreamPlayStart)))
				case *command.GetStreamLength:
					send(rtmpConn, 3, c.MakeResponse(0)) // 0 == live
				case *command.FCUnpublish:
					send(rtmpConn, 3, c.MakeResponse(command.NewStatus(command.NetStreamPublishStart, "FCUnpublish Ignored")))
				case *command.DeleteStream:
					sendUserControl(rtmpConn, &usercontrol.StreamEOF{StreamID: 1})
					// Note: use NetStream.Play.UnpublishNotify if stream was outbound instead of inbound
					send(rtmpConn, 3, c.MakeResponse(command.NewStatus(command.NetStreamUnpublishSuccess)))
				}
			}

		case *message.UserControlMessage:
			event, err := usercontrol.FromMessage(m)
			if err != nil {
				log.Println(err)
			} else {
				log.Printf("Received user control event %T: %+v", event, event)
				switch e := event.(type) {
				case *usercontrol.SetBufferLength:
					log.Printf("Client set buffer length: stream=%d, length=%dms", e.StreamID, e.BufferLength)
				case *usercontrol.PingRequest:
					sendUserControl(rtmpConn, &usercontrol.PingResponse{Timestamp: e.Timestamp})
				case *usercontrol.PingResponse:
					log.Printf("Received ping response: timestamp=%d", e.Timestamp)
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

func sendUserControl(c rtmp.Conn, event usercontrol.Event) {
	msg, err := event.ToMessage()
	if err != nil {
		log.Fatal(err)
	}
	send(c, 2, msg)
}
