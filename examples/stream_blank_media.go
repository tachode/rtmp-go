package examples

import (
	"encoding/binary"
	"log"
	"time"

	rtmp "github.com/tachode/rtmp-go"
	"github.com/tachode/rtmp-go/message"
)

// StreamBlankMedia sends a continuous stream of silent AAC audio at 48kHz and
// black H.264 keyframes at 24fps on the given RTMP connection and chunk stream.
// It blocks until the stop channel is closed.
func StreamBlankMedia(rtmpConn rtmp.Conn, chunkStream int, streamId uint32, stop <-chan struct{}) {
	// AAC AudioSpecificConfig: AAC-LC (2), 48kHz (3), Stereo (2)
	aacConfig := []byte{0x11, 0x90}

	// H.264 SPS for 160x120 baseline profile, level 1.2
	sps := []byte{
		0x67, 0x42, 0xc0, 0x0c, 0xdc, 0x28, 0x47, 0xe5,
		0xc0, 0x44, 0x00, 0x00, 0x03, 0x00, 0x04, 0x00,
		0x00, 0x03, 0x00, 0x08, 0x3c, 0x50, 0xae,
	}

	// H.264 PPS
	pps := []byte{0x68, 0xce, 0x0f, 0x13, 0x20}

	// AVCDecoderConfigurationRecord
	avcConfig := []byte{
		0x01,   // configurationVersion
		sps[1], // AVCProfileIndication
		sps[2], // profile_compatibility
		sps[3], // AVCLevelIndication
		0xff,   // lengthSizeMinusOne = 3
		0xe1,   // numOfSequenceParameterSets = 1
		0x00, byte(len(sps)),
	}
	avcConfig = append(avcConfig, sps...)
	avcConfig = append(avcConfig, 0x01, 0x00, byte(len(pps)))
	avcConfig = append(avcConfig, pps...)

	// H.264 IDR frame (black, 160x120)
	idrNal := []byte{
		0x65, 0x88, 0x84, 0x05, 0x73, 0x9f, 0xff, 0xff,
		0x0f, 0x45, 0x00, 0x01, 0x42, 0xdf, 0x27, 0x27,
		0x27, 0x27, 0x27, 0x27, 0x27, 0x27, 0x27, 0x5d,
		0x75, 0xd7, 0x5d, 0x75, 0xd7, 0x5d, 0x75, 0xd7,
		0x5d, 0x75, 0xd7, 0x5d, 0x75, 0xd7, 0x5d, 0x75,
		0xd7, 0x5d, 0x75, 0xd7, 0x5d, 0x75, 0xd7, 0x5d,
		0x75, 0xd7, 0x5d, 0x75, 0xd7, 0x5d, 0x75, 0xd7,
		0x5d, 0x75, 0xd7, 0x5d, 0x75, 0xd7, 0x5d, 0x75,
		0xd7, 0x5d, 0x75, 0xd7, 0x5d, 0x75, 0xd7, 0x5d,
		0x75, 0xd7, 0x5d, 0x78,
	}

	// Wrap IDR NAL in AVCC format (4-byte length prefix)
	idrAvcc := make([]byte, 4+len(idrNal))
	binary.BigEndian.PutUint32(idrAvcc, uint32(len(idrNal)))
	copy(idrAvcc[4:], idrNal)

	// AAC silent frame (1024 samples of silence at 48kHz)
	silentAac := []byte{0x21, 0x00, 0x49, 0x90, 0x02, 0x19, 0x00, 0x23, 0x80}

	// Send audio sequence header
	err := rtmpConn.WriteMessage(&message.AudioMessage{
		PacketType: message.ERTMPAudioPacketTypeSequenceStart,
		Rate:       message.AudioRate44kHz, // Legacy field; true rate is in AudioSpecificConfig
		SampleSize: message.AudioSize16Bit,
		Stereo:     true,
		Tracks:     []message.AudioTrack{{CodecId: message.AudioCodecIdAAC, Payload: aacConfig}},
	}, chunkStream)
	if err != nil {
		log.Printf("Error sending audio sequence header: %v", err)
		return
	}

	// Send video sequence header
	err = rtmpConn.WriteMessage(&message.VideoMessage{
		FrameType:  message.VideoFrameTypeKeyframe,
		PacketType: message.ERTMPVideoPacketTypeSequenceStart,
		Tracks:     []message.VideoTrack{{CodecId: message.VideoCodecIdAvc, Payload: avcConfig}},
	}, chunkStream)
	if err != nil {
		log.Printf("Error sending video sequence header: %v", err)
		return
	}

	log.Printf("Started streaming blank media on stream %d", streamId)

	// Audio: 48000 Hz / 1024 samples per frame ~ 46.875 fps
	// Video: 24 fps
	audioInterval := time.Second * 1024 / 48000
	videoInterval := time.Second / 24

	audioTicker := time.NewTicker(audioInterval)
	videoTicker := time.NewTicker(videoInterval)
	defer audioTicker.Stop()
	defer videoTicker.Stop()

	var audioSamples, videoFrames uint64

	for {
		select {
		case <-stop:
			log.Printf("Stopping stream on stream %d", streamId)
			return
		case <-audioTicker.C:
			audioTs := uint32(audioSamples * 1000 / 48000)
			err := rtmpConn.WriteMessage(&message.AudioMessage{
				MetadataFields: message.MetadataFields{StreamId: streamId, Timestamp: audioTs},
				PacketType:     message.ERTMPAudioPacketTypeCodedFrames,
				Rate:           message.AudioRate44kHz,
				SampleSize:     message.AudioSize16Bit,
				Stereo:         true,
				Tracks:         []message.AudioTrack{{CodecId: message.AudioCodecIdAAC, Payload: silentAac}},
			}, chunkStream)
			if err != nil {
				log.Printf("Error sending audio: %v", err)
				return
			}
			audioSamples += 1024
		case <-videoTicker.C:
			videoTs := uint32(videoFrames * 1000 / 24)
			err := rtmpConn.WriteMessage(&message.VideoMessage{
				MetadataFields: message.MetadataFields{StreamId: streamId, Timestamp: videoTs},
				FrameType:      message.VideoFrameTypeKeyframe,
				PacketType:     message.ERTMPVideoPacketTypeCodedFrames,
				Tracks:         []message.VideoTrack{{CodecId: message.VideoCodecIdAvc, Payload: idrAvcc}},
			}, chunkStream)
			if err != nil {
				log.Printf("Error sending video: %v", err)
				return
			}
			videoFrames++
		}
	}
}
