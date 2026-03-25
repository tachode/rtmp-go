package command

import (
	"errors"

	"github.com/tachode/rtmp-go/amf0"
	"github.com/tachode/rtmp-go/message"
)

// NetConnection.connect() command

func init() { RegisterCommand(new(Connect)) }

// AudioCodecFlag is a bitmask indicating supported audio codecs.
// Each bit corresponds to a message.AudioCodecId value.
type AudioCodecFlag uint16

const (
	SupportSndNone    AudioCodecFlag = 1 << message.AudioCodecIdLinearPCMPlatformEndian // 0x0001
	SupportSndADPCM   AudioCodecFlag = 1 << message.AudioCodecIdADPCM                   // 0x0002
	SupportSndMP3     AudioCodecFlag = 1 << message.AudioCodecIdMP3                     // 0x0004
	SupportSndIntel   AudioCodecFlag = 1 << message.AudioCodecIdLinearPCMLittleEndian   // 0x0008
	SupportSndUnused  AudioCodecFlag = 0x0010
	SupportSndNelly8  AudioCodecFlag = 1 << message.AudioCodecIdNellymoser8kHzMono // 0x0020
	SupportSndNelly   AudioCodecFlag = 1 << message.AudioCodecIdNellymoser         // 0x0040
	SupportSndG711A   AudioCodecFlag = 1 << message.AudioCodecIdG711ALaw           // 0x0080
	SupportSndG711U   AudioCodecFlag = 1 << message.AudioCodecIdG711MuLaw          // 0x0100
	SupportSndNelly16 AudioCodecFlag = 0x0200
	SupportSndAAC     AudioCodecFlag = 1 << message.AudioCodecIdAAC   // 0x0400
	SupportSndSpeex   AudioCodecFlag = 1 << message.AudioCodecIdSpeex // 0x0800
	SupportSndAll     AudioCodecFlag = 0x0FFF
)

// VideoCodecFlag is a bitmask indicating supported video codecs.
// Each bit corresponds to a message.VideoCodecId value.
type VideoCodecFlag uint16

const (
	SupportVidUnused    VideoCodecFlag = 1 << message.VideoCodecIdReserved0    // 0x0001
	SupportVidJPEG      VideoCodecFlag = 1 << message.VideoCodecIdReserved1    // 0x0002
	SupportVidSorenson  VideoCodecFlag = 1 << message.VideoCodecIdSorensonH263 // 0x0004
	SupportVidHomebrew  VideoCodecFlag = 1 << message.VideoCodecIdScreen1      // 0x0008
	SupportVidVP6       VideoCodecFlag = 1 << message.VideoCodecIdOn2VP6       // 0x0010
	SupportVidVP6Alpha  VideoCodecFlag = 1 << message.VideoCodecIdOn2VP6Alpha  // 0x0020
	SupportVidHomebrewV VideoCodecFlag = 1 << message.VideoCodecIdScreen2      // 0x0040
	SupportVidH264      VideoCodecFlag = 1 << message.VideoCodecIdAvc          // 0x0080
	SupportVidAll       VideoCodecFlag = 0x00FF
)

// VideoFunction is a bitmask indicating supported video functions.
type VideoFunction uint16

const (
	SupportVidClientSeek VideoFunction = 1
)

// ObjectEncoding indicates the AMF encoding method.
type ObjectEncoding int

const (
	ObjectEncodingAMF0 ObjectEncoding = 0
	ObjectEncodingAMF3 ObjectEncoding = 3
)

type Connect struct {
	Transaction    int            // Always set to 1 for connect commands.
	App            string         // The server application name the client is connected to (e.g. "testapp").
	FlashVer       string         // Flash Player version string, as returned by getversion().
	SwfUrl         string         // URL of the source SWF file making the connection.
	TcUrl          string         // URL of the server: protocol://servername:port/appName/appInstance.
	Fpad           bool           // True if a proxy is being used.
	AudioCodecs    AudioCodecFlag // Bitmask indicating which audio codecs the client supports.
	VideoCodecs    VideoCodecFlag // Bitmask indicating which video codecs the client supports.
	VideoFunction  VideoFunction  // Bitmask indicating which special video functions are supported.
	PageUrl        string         // URL of the web page from where the SWF file was loaded.
	ObjectEncoding ObjectEncoding // AMF encoding method (AMF0 or AMF3).
}

func (c Connect) CommandName() string { return "connect" }

func (c *Connect) FromMessageCommand(cmd message.Command) error {
	c.Transaction = int(cmd.GetTransactionId())
	obj := cmd.GetObject()
	if obj == nil {
		return errors.New("connect command contains no command object")
	}
	c.App = GetString(obj, "app")
	c.FlashVer = GetString(obj, "flashver")
	c.SwfUrl = GetString(obj, "swfUrl")
	c.TcUrl = GetString(obj, "tcUrl")
	c.Fpad = GetBool(obj, "fpad")
	c.AudioCodecs = AudioCodecFlag(GetFloat64(obj, "audioCodecs"))
	c.VideoCodecs = VideoCodecFlag(GetFloat64(obj, "videoCodecs"))
	c.VideoFunction = VideoFunction(GetFloat64(obj, "videoFunction"))
	c.PageUrl = GetString(obj, "pageUrl")
	c.ObjectEncoding = ObjectEncoding(GetFloat64(obj, "objectEncoding"))
	return nil
}

func (c *Connect) ToMessageCommand() (message.Command, error) {
	cmd := &message.Amf0CommandMessage{
		Command:       c.CommandName(),
		TransactionId: float64(c.Transaction),
		Object: amf0.Object{
			"app":            c.App,
			"flashver":       c.FlashVer,
			"swfUrl":         c.SwfUrl,
			"tcUrl":          c.TcUrl,
			"fpad":           c.Fpad,
			"audioCodecs":    float64(c.AudioCodecs),
			"videoCodecs":    float64(c.VideoCodecs),
			"videoFunction":  float64(c.VideoFunction),
			"pageUrl":        c.PageUrl,
			"objectEncoding": float64(c.ObjectEncoding),
		},
	}
	return cmd, nil
}

func (c *Connect) MakeResponse(status Status, amfLevel ObjectEncoding) message.Command {
	p0 := status.ToObject()
	p0["objectEncoding"] = amfLevel

	command := "_result"
	if status.Level == LevelError {
		command = "_error"
	}

	cmd := &message.Amf0CommandMessage{
		Command:       command,
		TransactionId: float64(c.Transaction),
		Object:        nil,
		Parameters:    []any{p0},
	}
	return cmd
}
