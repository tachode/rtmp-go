package main

import (
	"github.com/tachode/rtmp-go/command"
	"github.com/tachode/rtmp-go/examples"
)

func main() {
	rtmpConn := examples.ConnectAndCreateStream("localhost:1935")
	defer rtmpConn.Close()

	// Send publish
	examples.SendCommand(rtmpConn, 3, &command.Publish{
		StreamId:     1,
		Transaction:  0,
		StreamKey:    "live",
		HowToPublish: command.HowToPublishLive,
	})

	// Wait for onStatus indicating publish started
	examples.WaitForStatus(rtmpConn, command.NetStreamPublishStart)

	// Stream blank media until the connection is closed
	stop := make(chan struct{})
	defer close(stop)
	examples.StreamBlankMedia(rtmpConn, 3, 1, stop)
}
