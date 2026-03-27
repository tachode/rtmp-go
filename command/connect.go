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
	SupportVidClientSeek              VideoFunction = 0x0001
	SupportVidClientHDR               VideoFunction = 0x0002
	SupportVidClientVideoPacketTypeMD VideoFunction = 0x0004
	SupportVidClientLargeScaleTile    VideoFunction = 0x0008
)

// FourCcInfoMask defines capability flags for a FourCC codec.
type FourCcInfoMask uint16

const (
	FourCcInfoCanDecode  FourCcInfoMask = 0x01
	FourCcInfoCanEncode  FourCcInfoMask = 0x02
	FourCcInfoCanForward FourCcInfoMask = 0x04
)

// FourCcInfoMap maps FourCC codec strings to FourCcInfoMask capability flags.
// A key of "*" acts as a wildcard for any codec.
type FourCcInfoMap map[string]FourCcInfoMask

// CapsExMask defines extended capability flags for E-RTMP.
type CapsExMask uint16

const (
	CapsExReconnect           CapsExMask = 0x01
	CapsExMultitrack          CapsExMask = 0x02
	CapsExModEx               CapsExMask = 0x04
	CapsExTimestampNanoOffset CapsExMask = 0x08
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

	// E-RTMP capability negotiation
	FourCcList         []string      // List of FourCC codec strings the client supports.
	VideoFourCcInfoMap FourCcInfoMap // Per-codec capability flags for video codecs.
	AudioFourCcInfoMap FourCcInfoMap // Per-codec capability flags for audio codecs.
	CapsEx             CapsExMask    // Extended capabilities bitmask.
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

	// E-RTMP capability negotiation
	c.FourCcList = GetStringSlice(obj, "fourCcList")
	c.VideoFourCcInfoMap = GetFourCcInfoMap(obj, "videoFourCcInfoMap")
	c.AudioFourCcInfoMap = GetFourCcInfoMap(obj, "audioFourCcInfoMap")
	c.CapsEx = CapsExMask(GetFloat64(obj, "capsEx"))
	return nil
}

func (c *Connect) ToMessageCommand() (message.Command, error) {
	obj := amf0.Object{
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
	}

	if len(c.FourCcList) > 0 {
		arr := make(amf0.StrictArray, len(c.FourCcList))
		for i, v := range c.FourCcList {
			arr[i] = v
		}
		obj["fourCcList"] = arr
	}
	if len(c.VideoFourCcInfoMap) > 0 {
		obj["videoFourCcInfoMap"] = fourCcInfoMapToAMF(c.VideoFourCcInfoMap)
	}
	if len(c.AudioFourCcInfoMap) > 0 {
		obj["audioFourCcInfoMap"] = fourCcInfoMapToAMF(c.AudioFourCcInfoMap)
	}
	if c.CapsEx != 0 {
		obj["capsEx"] = float64(c.CapsEx)
	}

	cmd := &message.Amf0CommandMessage{
		Command:       c.CommandName(),
		TransactionId: float64(c.Transaction),
		Object:        obj,
	}
	return cmd, nil
}

func (c *Connect) MakeResponse(status Status, amfLevel ObjectEncoding, serverCaps ...ConnectResponseCaps) message.Command {
	p0 := status.ToObject()
	p0["objectEncoding"] = amfLevel

	if len(serverCaps) > 0 {
		caps := serverCaps[0]
		if len(caps.VideoFourCcInfoMap) > 0 {
			p0["videoFourCcInfoMap"] = fourCcInfoMapToAMF(caps.VideoFourCcInfoMap)
		}
		if len(caps.AudioFourCcInfoMap) > 0 {
			p0["audioFourCcInfoMap"] = fourCcInfoMapToAMF(caps.AudioFourCcInfoMap)
		}
		if caps.CapsEx != 0 {
			p0["capsEx"] = float64(caps.CapsEx)
		}
		if fourCcList := caps.fourCcList(); len(fourCcList) > 0 {
			arr := make(amf0.StrictArray, len(fourCcList))
			for i, v := range fourCcList {
				arr[i] = v
			}
			p0["fourCcList"] = arr
		}
	}

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

// ConnectResponseCaps holds E-RTMP capabilities the server wishes to advertise
// in its connect response.
type ConnectResponseCaps struct {
	VideoFourCcInfoMap FourCcInfoMap
	AudioFourCcInfoMap FourCcInfoMap
	CapsEx             CapsExMask
}

// fourCcList returns the deduplicated union of all codec keys from
// VideoFourCcInfoMap and AudioFourCcInfoMap. If an "*" is encountered,
// the returned list will be "*"
func (c ConnectResponseCaps) fourCcList() []string {
	seen := make(map[string]struct{})
	var list []string
	for _, m := range []FourCcInfoMap{c.VideoFourCcInfoMap, c.AudioFourCcInfoMap} {
		for k := range m {
			if k == "*" {
				return []string{"*"}
			}
			if _, dup := seen[k]; !dup {
				seen[k] = struct{}{}
				list = append(list, k)
			}
		}
	}
	return list
}

func fourCcInfoMapToAMF(m FourCcInfoMap) amf0.Object {
	o := make(amf0.Object, len(m))
	for k, v := range m {
		o[k] = float64(v)
	}
	return o
}
