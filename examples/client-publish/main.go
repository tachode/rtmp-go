package main

import (
	"log"

	"github.com/tachode/rtmp-go/command"
	"github.com/tachode/rtmp-go/data"
	"github.com/tachode/rtmp-go/examples"
	"github.com/tachode/rtmp-go/message"
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

	// Send stream metadata
	metadata := &data.OnMetaData{
		AudioCodecId:    message.AudioCodecIdAAC,
		AudioDataRate:   160,
		AudioSampleRate: 48000,
		AudioSampleSize: 16,
		Stereo:          true,
		VideoCodecId:    message.VideoCodecIdAvc,
		VideoDataRate:   6000,
		Width:           160,
		Height:          120,
		FrameRate:       24,
	}
	dataMsg, err := metadata.ToDataMessage()
	if err != nil {
		log.Fatal("Error creating onMetaData message:", err)
	}
	log.Printf(">>> %s", dataMsg.(message.Message))
	err = rtmpConn.WriteMessage(dataMsg.(message.Message), 3)
	if err != nil {
		log.Fatal("Error sending onMetaData:", err)
	}

	// Stream blank media until the connection is closed
	stop := make(chan struct{})
	defer close(stop)
	examples.StreamBlankMedia(rtmpConn, 3, 1, stop)
}
