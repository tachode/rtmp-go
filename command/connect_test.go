package command

import (
	"reflect"
	"testing"

	"github.com/tachode/rtmp-go/amf0"
	"github.com/tachode/rtmp-go/message"
)

func TestConnect_ERTMPRoundTrip(t *testing.T) {
	original := &Connect{
		Transaction:    1,
		App:            "live",
		FlashVer:       "FMLE/3.0",
		TcUrl:          "rtmp://localhost/live",
		ObjectEncoding: ObjectEncodingAMF0,
		FourCcList:     []string{"av01", "vp09", "hvc1", "avc1", "Opus", "mp4a"},
		VideoFourCcInfoMap: FourCcInfoMap{
			"av01": FourCcInfoCanDecode | FourCcInfoCanEncode,
			"vp09": FourCcInfoCanDecode,
			"*":    FourCcInfoCanForward,
		},
		AudioFourCcInfoMap: FourCcInfoMap{
			"Opus": FourCcInfoCanDecode | FourCcInfoCanEncode,
			"mp4a": FourCcInfoCanDecode,
		},
		CapsEx: CapsExReconnect | CapsExMultitrack | CapsExModEx,
	}

	cmd, err := original.ToMessageCommand()
	if err != nil {
		t.Fatal("ToMessageCommand:", err)
	}

	parsed := &Connect{}
	err = parsed.FromMessageCommand(cmd)
	if err != nil {
		t.Fatal("FromMessageCommand:", err)
	}

	if parsed.App != original.App {
		t.Errorf("App: got %q, want %q", parsed.App, original.App)
	}
	if parsed.TcUrl != original.TcUrl {
		t.Errorf("TcUrl: got %q, want %q", parsed.TcUrl, original.TcUrl)
	}
	if parsed.CapsEx != original.CapsEx {
		t.Errorf("CapsEx: got %d, want %d", parsed.CapsEx, original.CapsEx)
	}
	if !reflect.DeepEqual(parsed.FourCcList, original.FourCcList) {
		t.Errorf("FourCcList: got %v, want %v", parsed.FourCcList, original.FourCcList)
	}
	if !reflect.DeepEqual(parsed.VideoFourCcInfoMap, original.VideoFourCcInfoMap) {
		t.Errorf("VideoFourCcInfoMap: got %v, want %v", parsed.VideoFourCcInfoMap, original.VideoFourCcInfoMap)
	}
	if !reflect.DeepEqual(parsed.AudioFourCcInfoMap, original.AudioFourCcInfoMap) {
		t.Errorf("AudioFourCcInfoMap: got %v, want %v", parsed.AudioFourCcInfoMap, original.AudioFourCcInfoMap)
	}
}

func TestConnect_LegacyWithoutERTMP(t *testing.T) {
	original := &Connect{
		Transaction:    1,
		App:            "live",
		FlashVer:       "FMLE/3.0",
		TcUrl:          "rtmp://localhost/live",
		AudioCodecs:    SupportSndAAC,
		VideoCodecs:    SupportVidH264,
		ObjectEncoding: ObjectEncodingAMF0,
	}

	cmd, err := original.ToMessageCommand()
	if err != nil {
		t.Fatal("ToMessageCommand:", err)
	}

	// Verify E-RTMP properties are not present on the wire
	amfCmd := cmd.(*message.Amf0CommandMessage)
	for _, key := range []string{"fourCcList", "videoFourCcInfoMap", "audioFourCcInfoMap", "capsEx"} {
		if _, found := amfCmd.Object[key]; found {
			t.Errorf("legacy connect should not include %q", key)
		}
	}

	parsed := &Connect{}
	err = parsed.FromMessageCommand(cmd)
	if err != nil {
		t.Fatal("FromMessageCommand:", err)
	}

	if parsed.CapsEx != 0 {
		t.Errorf("CapsEx: got %d, want 0", parsed.CapsEx)
	}
	if parsed.FourCcList != nil {
		t.Errorf("FourCcList: got %v, want nil", parsed.FourCcList)
	}
	if parsed.VideoFourCcInfoMap != nil {
		t.Errorf("VideoFourCcInfoMap: got %v, want nil", parsed.VideoFourCcInfoMap)
	}
	if parsed.AudioFourCcInfoMap != nil {
		t.Errorf("AudioFourCcInfoMap: got %v, want nil", parsed.AudioFourCcInfoMap)
	}
}

func TestConnect_MakeResponseWithCaps(t *testing.T) {
	c := &Connect{Transaction: 1}
	caps := ConnectResponseCaps{
		VideoFourCcInfoMap: FourCcInfoMap{
			"av01": FourCcInfoCanDecode,
			"*":    FourCcInfoCanForward,
		},
		AudioFourCcInfoMap: FourCcInfoMap{
			"Opus": FourCcInfoCanDecode | FourCcInfoCanEncode,
		},
		CapsEx: CapsExReconnect | CapsExMultitrack,
	}

	resp := c.MakeResponse(NewStatus(NetConnectionConnectSuccess), ObjectEncodingAMF0, caps)
	amfCmd := resp.(*message.Amf0CommandMessage)

	if amfCmd.Command != "_result" {
		t.Errorf("Command: got %q, want %q", amfCmd.Command, "_result")
	}

	params := amfCmd.Parameters
	if len(params) < 1 {
		t.Fatal("expected at least 1 parameter")
	}
	p0, ok := params[0].(amf0.Object)
	if !ok {
		t.Fatal("parameter 0 is not amf0.Object")
	}

	// Check capsEx
	capsExVal, found := p0["capsEx"]
	if !found {
		t.Fatal("capsEx not found in response")
	}
	if capsExVal != float64(CapsExReconnect|CapsExMultitrack) {
		t.Errorf("capsEx: got %v, want %v", capsExVal, float64(CapsExReconnect|CapsExMultitrack))
	}

	// Check videoFourCcInfoMap
	vidMap, found := p0["videoFourCcInfoMap"]
	if !found {
		t.Fatal("videoFourCcInfoMap not found in response")
	}
	vidObj, ok := vidMap.(amf0.Object)
	if !ok {
		t.Fatal("videoFourCcInfoMap is not amf0.Object")
	}
	av01Val, found := vidObj["av01"]
	if !found || av01Val != float64(FourCcInfoCanDecode) {
		t.Errorf("videoFourCcInfoMap[av01]: got %v, want %v", av01Val, float64(FourCcInfoCanDecode))
	}

	// Check audioFourCcInfoMap
	audMap, found := p0["audioFourCcInfoMap"]
	if !found {
		t.Fatal("audioFourCcInfoMap not found in response")
	}
	audObj, ok := audMap.(amf0.Object)
	if !ok {
		t.Fatal("audioFourCcInfoMap is not amf0.Object")
	}
	opusVal, found := audObj["Opus"]
	if !found || opusVal != float64(FourCcInfoCanDecode|FourCcInfoCanEncode) {
		t.Errorf("audioFourCcInfoMap[Opus]: got %v, want %v", opusVal, float64(FourCcInfoCanDecode|FourCcInfoCanEncode))
	}

	// Check fourCcList — wildcard in maps means fourCcList should be ["*"]
	fourCcListVal, found := p0["fourCcList"]
	if !found {
		t.Fatal("fourCcList not found in response")
	}
	fourCcArr, ok := fourCcListVal.(amf0.StrictArray)
	if !ok {
		t.Fatal("fourCcList is not amf0.StrictArray")
	}
	if len(fourCcArr) != 1 {
		t.Fatalf("fourCcList: got %d entries, want 1", len(fourCcArr))
	}
	if fourCcArr[0] != "*" {
		t.Errorf("fourCcList[0]: got %v, want \"*\"", fourCcArr[0])
	}
}

func TestConnect_MakeResponseFourCcListUnion(t *testing.T) {
	c := &Connect{Transaction: 1}
	caps := ConnectResponseCaps{
		VideoFourCcInfoMap: FourCcInfoMap{
			"av01": FourCcInfoCanDecode,
			"hvc1": FourCcInfoCanDecode,
		},
		AudioFourCcInfoMap: FourCcInfoMap{
			"Opus": FourCcInfoCanDecode | FourCcInfoCanEncode,
			"mp4a": FourCcInfoCanDecode,
		},
	}

	resp := c.MakeResponse(NewStatus(NetConnectionConnectSuccess), ObjectEncodingAMF0, caps)
	amfCmd := resp.(*message.Amf0CommandMessage)
	p0 := amfCmd.Parameters[0].(amf0.Object)

	fourCcListVal, found := p0["fourCcList"]
	if !found {
		t.Fatal("fourCcList not found in response")
	}
	fourCcArr, ok := fourCcListVal.(amf0.StrictArray)
	if !ok {
		t.Fatal("fourCcList is not amf0.StrictArray")
	}
	gotCodecs := make(map[string]bool)
	for _, item := range fourCcArr {
		if s, ok := item.(string); ok {
			gotCodecs[s] = true
		}
	}
	for _, expected := range []string{"av01", "hvc1", "Opus", "mp4a"} {
		if !gotCodecs[expected] {
			t.Errorf("fourCcList missing %q, got %v", expected, fourCcArr)
		}
	}
	if len(fourCcArr) != 4 {
		t.Errorf("fourCcList: got %d entries, want 4", len(fourCcArr))
	}
}

func TestConnect_MakeResponseWithoutCaps(t *testing.T) {
	c := &Connect{Transaction: 1}
	resp := c.MakeResponse(NewStatus(NetConnectionConnectSuccess), ObjectEncodingAMF0)
	amfCmd := resp.(*message.Amf0CommandMessage)

	params := amfCmd.Parameters
	if len(params) < 1 {
		t.Fatal("expected at least 1 parameter")
	}
	p0, ok := params[0].(amf0.Object)
	if !ok {
		t.Fatal("parameter 0 is not amf0.Object")
	}

	for _, key := range []string{"fourCcList", "videoFourCcInfoMap", "audioFourCcInfoMap", "capsEx"} {
		if _, found := p0[key]; found {
			t.Errorf("response without caps should not include %q", key)
		}
	}
}

func TestCapsExMask_Values(t *testing.T) {
	if CapsExReconnect != 0x01 {
		t.Errorf("CapsExReconnect: got %#x, want 0x01", CapsExReconnect)
	}
	if CapsExMultitrack != 0x02 {
		t.Errorf("CapsExMultitrack: got %#x, want 0x02", CapsExMultitrack)
	}
	if CapsExModEx != 0x04 {
		t.Errorf("CapsExModEx: got %#x, want 0x04", CapsExModEx)
	}
	if CapsExTimestampNanoOffset != 0x08 {
		t.Errorf("CapsExTimestampNanoOffset: got %#x, want 0x08", CapsExTimestampNanoOffset)
	}
}

func TestFourCcInfoMask_Values(t *testing.T) {
	if FourCcInfoCanDecode != 0x01 {
		t.Errorf("FourCcInfoCanDecode: got %#x, want 0x01", FourCcInfoCanDecode)
	}
	if FourCcInfoCanEncode != 0x02 {
		t.Errorf("FourCcInfoCanEncode: got %#x, want 0x02", FourCcInfoCanEncode)
	}
	if FourCcInfoCanForward != 0x04 {
		t.Errorf("FourCcInfoCanForward: got %#x, want 0x04", FourCcInfoCanForward)
	}
}

func TestAudioCodecFlag_String(t *testing.T) {
	tests := []struct {
		val  AudioCodecFlag
		want string
	}{
		{0, "0"},
		{SupportSndMP3, "MP3"},
		{SupportSndAAC, "AAC"},
		{SupportSndMP3 | SupportSndAAC, "MP3|AAC"},
		{SupportSndAll, "PCM|ADPCM|MP3|PCM-LE|Unused|Nelly8k|Nelly|G711A|G711U|Nelly16k|AAC|Speex"},
		{SupportSndAAC | AudioCodecFlag(0x8000), "AAC|0x8000"},
	}
	for _, tt := range tests {
		if got := tt.val.String(); got != tt.want {
			t.Errorf("AudioCodecFlag(%#x).String() = %q, want %q", uint16(tt.val), got, tt.want)
		}
	}
}

func TestVideoCodecFlag_String(t *testing.T) {
	tests := []struct {
		val  VideoCodecFlag
		want string
	}{
		{0, "0"},
		{SupportVidH264, "H264"},
		{SupportVidSorenson | SupportVidH264, "Sorenson|H264"},
		{SupportVidAll, "Unused|JPEG|Sorenson|Screen|VP6|VP6Alpha|ScreenV2|H264"},
		{SupportVidVP6 | VideoCodecFlag(0x4000), "VP6|0x4000"},
	}
	for _, tt := range tests {
		if got := tt.val.String(); got != tt.want {
			t.Errorf("VideoCodecFlag(%#x).String() = %q, want %q", uint16(tt.val), got, tt.want)
		}
	}
}

func TestVideoFunction_String(t *testing.T) {
	tests := []struct {
		val  VideoFunction
		want string
	}{
		{0, "0"},
		{SupportVidClientSeek, "Seek"},
		{SupportVidClientSeek | SupportVidClientHDR, "Seek|HDR"},
		{SupportVidClientSeek | SupportVidClientHDR | SupportVidClientVideoPacketTypeMetadata | SupportVidClientLargeScaleTile, "Seek|HDR|Metadata|LargeScaleTile"},
		{SupportVidClientSeek | VideoFunction(0x0100), "Seek|0x100"},
	}
	for _, tt := range tests {
		if got := tt.val.String(); got != tt.want {
			t.Errorf("VideoFunction(%#x).String() = %q, want %q", uint16(tt.val), got, tt.want)
		}
	}
}

func TestFourCcInfoMask_String(t *testing.T) {
	tests := []struct {
		val  FourCcInfoMask
		want string
	}{
		{0, "0"},
		{FourCcInfoCanDecode, "Decode"},
		{FourCcInfoCanDecode | FourCcInfoCanEncode, "Decode|Encode"},
		{FourCcInfoCanDecode | FourCcInfoCanEncode | FourCcInfoCanForward, "Decode|Encode|Forward"},
		{FourCcInfoCanForward | FourCcInfoMask(0x80), "Forward|0x80"},
	}
	for _, tt := range tests {
		if got := tt.val.String(); got != tt.want {
			t.Errorf("FourCcInfoMask(%#x).String() = %q, want %q", uint16(tt.val), got, tt.want)
		}
	}
}

func TestCapsExMask_String(t *testing.T) {
	tests := []struct {
		val  CapsExMask
		want string
	}{
		{0, "0"},
		{CapsExReconnect, "Reconnect"},
		{CapsExMultitrack | CapsExModEx, "Multitrack|ModEx"},
		{CapsExReconnect | CapsExMultitrack | CapsExModEx | CapsExTimestampNanoOffset, "Reconnect|Multitrack|ModEx|TimestampNanoOffset"},
		{CapsExModEx | CapsExMask(0xF0), "ModEx|0xf0"},
	}
	for _, tt := range tests {
		if got := tt.val.String(); got != tt.want {
			t.Errorf("CapsExMask(%#x).String() = %q, want %q", uint16(tt.val), got, tt.want)
		}
	}
}
