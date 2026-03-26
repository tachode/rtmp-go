package examples

import (
	"log"
	"net"

	rtmp "github.com/tachode/rtmp-go"
	"github.com/tachode/rtmp-go/command"
	"github.com/tachode/rtmp-go/message"
)

// ConnectAndCreateStream dials the given address, performs the RTMP client
// handshake, sends a Connect command, and creates a stream. It returns the
// ready-to-use RTMP connection.
func ConnectAndCreateStream(addr string) rtmp.Conn {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}

	rtmpConn, err := rtmp.NewClientConn(conn, 3)
	if err != nil {
		conn.Close()
		log.Fatal("Handshake error:", err)
	}
	log.Println("Client handshake completed")

	err = rtmpConn.CreateOutboundChunkstream(3, 1)
	if err != nil {
		log.Fatal(err)
	}

	SendCommand(rtmpConn, 3, &command.Connect{
		Transaction:    1,
		App:            "live",
		FlashVer:       "FMLE/3.0",
		TcUrl:          "rtmp://" + addr + "/live",
		AudioCodecs:    command.SupportSndAAC,
		VideoCodecs:    command.SupportVidH264,
		ObjectEncoding: command.ObjectEncodingAMF0,
	})
	WaitForResult(rtmpConn, 1)

	SendCommand(rtmpConn, 3, &command.CreateStream{
		Transaction: 2,
	})
	WaitForResult(rtmpConn, 2)

	return rtmpConn
}

// SendCommand serializes a command and writes it to the connection.
func SendCommand(c rtmp.Conn, chunkStream int, cmd command.Command) {
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

// WaitForResult reads messages until a _result or _error with the given
// transaction ID is received.
func WaitForResult(c rtmp.Conn, transactionId float64) {
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

// WaitForStatus reads messages until an onStatus command with the given status
// code is received.
func WaitForStatus(c rtmp.Conn, code command.StatusCode) {
	for {
		msg, err := c.ReadMessage()
		if err != nil {
			log.Fatal("Connection closed while waiting for status:", err)
		}
		log.Printf("<<< %v", msg)

		if cmd, ok := msg.(message.Command); ok {
			if cmd.GetCommand() == "onStatus" {
				parsed, err := command.FromMessageCommand(cmd)
				if err == nil {
					if onStatus, ok := parsed.(*command.OnStatus); ok {
						log.Printf("Received onStatus: %s", onStatus.Code)
						if onStatus.Code == code {
							return
						}
					}
				}
			}
		}
	}
}
