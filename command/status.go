package command

import (
	"fmt"

	"github.com/tachode/rtmp-go/amf0"
)

type Level string

const (
	LevelStatus  Level = "status"
	LevelError   Level = "error"
	LevelWarning Level = "warning"
)

type StatusCode string

const (
	NetConnectionCallFailed              StatusCode = "NetConnection.Call.Failed"
	NetConnectionConnectAppShutdown      StatusCode = "NetConnection.Connect.AppShutdown"
	NetConnectionConnectClosed           StatusCode = "NetConnection.Connect.Closed"
	NetConnectionConnectFailed           StatusCode = "NetConnection.Connect.Failed"
	NetConnectionConnectReconnectRequest StatusCode = "NetConnection.Connect.ReconnectRequest"
	NetConnectionConnectRejected         StatusCode = "NetConnection.Connect.Rejected"
	NetConnectionConnectSuccess          StatusCode = "NetConnection.Connect.Success"
	NetConnectionProxyNotResponding      StatusCode = "NetConnection.Proxy.NotResponding"
	NetStreamClearFailed                 StatusCode = "NetStream.Clear.Failed"
	NetStreamClearSuccess                StatusCode = "NetStream.Clear.Success"
	NetStreamConnectFailed               StatusCode = "NetStream.Connect.Failed"
	NetStreamConnectSuccess              StatusCode = "NetStream.Connect.Success"
	NetStreamFailed                      StatusCode = "NetStream.Failed"
	NetStreamMulticastStreamReset        StatusCode = "NetStream.MulticastStream.Reset"
	NetStreamPauseNotify                 StatusCode = "NetStream.Pause.Notify"
	NetStreamPlayFailed                  StatusCode = "NetStream.Play.Failed"
	NetStreamPlayInsufficientBW          StatusCode = "NetStream.Play.InsufficientBW"
	NetStreamPlayPublishNotify           StatusCode = "NetStream.Play.PublishNotify"
	NetStreamPlayReset                   StatusCode = "NetStream.Play.Reset"
	NetStreamPlayStart                   StatusCode = "NetStream.Play.Start"
	NetStreamPlayStop                    StatusCode = "NetStream.Play.Stop"
	NetStreamPlayStreamNotFound          StatusCode = "NetStream.Play.StreamNotFound"
	NetStreamPlayUnpublishNotify         StatusCode = "NetStream.Play.UnpublishNotify"
	NetStreamPublishBadName              StatusCode = "NetStream.Publish.BadName"
	NetStreamPublishFailed               StatusCode = "NetStream.Publish.Failed"
	NetStreamPublishStart                StatusCode = "NetStream.Publish.Start"
	NetStreamRecordDiskQuotaExceeded     StatusCode = "NetStream.Record.DiskQuotaExceeded"
	NetStreamRecordFailed                StatusCode = "NetStream.Record.Failed"
	NetStreamRecordNoAccess              StatusCode = "NetStream.Record.NoAccess"
	NetStreamRecordStart                 StatusCode = "NetStream.Record.Start"
	NetStreamRecordStop                  StatusCode = "NetStream.Record.Stop"
	NetStreamSeekNotify                  StatusCode = "NetStream.Seek.Notify"
	NetStreamUnpauseNotify               StatusCode = "NetStream.Unpause.Notify"
	NetStreamUnpublishSuccess            StatusCode = "NetStream.Unpublish.Success"
)

type Status struct {
	Level       Level
	Code        StatusCode
	Description string
}

func (s Status) Error() string {
	return fmt.Sprintf("%s: %s", s.Code, s.Description)
}

var defaultDescriptions = map[StatusCode]string{
	NetConnectionCallFailed:              "The NetConnection.call() method was not able to invoke the server-side method or command.",
	NetConnectionConnectAppShutdown:      "The application has been shut down (for example, if the application is out of memory resources and must shut down to prevent the server from crashing) or the server has shut down.",
	NetConnectionConnectClosed:           "The connection was closed successfully.",
	NetConnectionConnectFailed:           "The connection attempt failed.",
	NetConnectionConnectReconnectRequest: "The server is requesting that the client reconnect.",
	NetConnectionConnectRejected:         "The client does not have permission to connect to the application, or the application name specified during the connection attempt was not found on the server.",
	NetConnectionConnectSuccess:          "The connection attempt succeeded.",
	NetConnectionProxyNotResponding:      "The proxy server is not responding.",
	NetStreamClearFailed:                 "A call to application.clearStreams() failed to delete a stream.",
	NetStreamClearSuccess:                "A call to application.clearStreams() successfully deleted a stream.",
	NetStreamConnectFailed:               "Dispatched when NetStream creation or connection fails (for example, if there is an error in the GroupSpecifier).",
	NetStreamConnectSuccess:              "Dispatched when a NetStream is created successfully.",
	NetStreamFailed:                      "An attempt to use a Stream method failed.",
	NetStreamMulticastStreamReset:        "A multicast subscription has changed focus to a different stream published with the same name in the same group.",
	NetStreamPauseNotify:                 "The stream is paused.",
	NetStreamPlayFailed:                  "A NetStream cannot play the stream.",
	NetStreamPlayInsufficientBW:          "Data is playing behind the normal speed.",
	NetStreamPlayPublishNotify:           "The initial publish operation to a stream was successful. This message is sent to all subscribers.",
	NetStreamPlayReset:                   "A playlist was reset.",
	NetStreamPlayStart:                   "Play was started.",
	NetStreamPlayStop:                    "Play was stopped.",
	NetStreamPlayStreamNotFound:          "An attempt was made to play a stream that does not exist.",
	NetStreamPlayUnpublishNotify:         "An unpublish operation from a stream was successful. This message is sent to all subscribers.",
	NetStreamPublishBadName:              "An attempt was made to publish to a stream that is already being published by someone else.",
	NetStreamPublishFailed:               "A call to NetStream.publish() was attempted and failed.",
	NetStreamPublishStart:                "An attempt to publish was successful.",
	NetStreamRecordDiskQuotaExceeded:     "An attempt to record a stream failed because the disk quota was exceeded.",
	NetStreamRecordFailed:                "An attempt to record a stream failed.",
	NetStreamRecordNoAccess:              "An attempt was made to record a read-only stream.",
	NetStreamRecordStart:                 "Recording was started.",
	NetStreamRecordStop:                  "Recording was stopped.",
	NetStreamSeekNotify:                  "The seek operation is complete.",
	NetStreamUnpauseNotify:               "The stream is resumed.",
	NetStreamUnpublishSuccess:            "The stream was unpublished successfully.",
}

func NewStatus(code StatusCode, description ...string) Status {
	desc := ""
	if len(description) > 0 {
		desc = description[0]
	}
	if len(desc) == 0 {
		desc = defaultDescriptions[code]
	}
	level := LevelError
	switch code {
	case NetConnectionConnectClosed,
		NetConnectionConnectReconnectRequest,
		NetConnectionConnectSuccess,
		NetStreamClearSuccess,
		NetStreamConnectSuccess,
		NetStreamMulticastStreamReset,
		NetStreamPauseNotify,
		NetStreamPlayPublishNotify,
		NetStreamPlayReset,
		NetStreamPlayStart,
		NetStreamPlayStop,
		NetStreamPlayUnpublishNotify,
		NetStreamPublishStart,
		NetStreamRecordStart,
		NetStreamRecordStop,
		NetStreamSeekNotify,
		NetStreamUnpauseNotify,
		NetStreamUnpublishSuccess:
		level = LevelStatus
	case NetStreamPlayInsufficientBW:
		level = LevelWarning
	}
	return Status{
		Level:       level,
		Code:        code,
		Description: desc,
	}
}

func (s Status) ToObject() amf0.Object {
	return amf0.Object{
		"level":       s.Level,
		"code":        s.Code,
		"description": s.Description,
	}
}
