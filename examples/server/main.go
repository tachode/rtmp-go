package main

import (
	"fmt"
	"log"
	"net"

	rtmp "github.com/tachode/rtmp-go"
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

	for {
		msg, err := rtmpConn.ReadMessage()
		if err != nil {
			log.Printf("Connection %s closed: %v", conn.RemoteAddr(), err)
			return
		}
		fmt.Println(msg)
	}
}
