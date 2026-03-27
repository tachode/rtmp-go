package command

import (
	"testing"

	"github.com/tachode/rtmp-go/amf0"
	"github.com/tachode/rtmp-go/message"
)

func TestOnStatus_RoundTrip(t *testing.T) {
	original := NewOnStatus(0, NewStatus(NetStreamPlayStart))

	cmd, err := original.ToMessageCommand()
	if err != nil {
		t.Fatal("ToMessageCommand:", err)
	}

	parsed := &OnStatus{}
	err = parsed.FromMessageCommand(cmd)
	if err != nil {
		t.Fatal("FromMessageCommand:", err)
	}

	if parsed.Transaction != original.Transaction {
		t.Errorf("Transaction: got %d, want %d", parsed.Transaction, original.Transaction)
	}
	if parsed.Level != original.Level {
		t.Errorf("Level: got %q, want %q", parsed.Level, original.Level)
	}
	if parsed.Code != original.Code {
		t.Errorf("Code: got %q, want %q", parsed.Code, original.Code)
	}
	if parsed.Description != original.Description {
		t.Errorf("Description: got %q, want %q", parsed.Description, original.Description)
	}
}

func TestOnStatus_RoundTripWithTcUrl(t *testing.T) {
	original := NewOnStatus(5, NewReconnectStatus("rtmp://newserver:1935/app"))

	cmd, err := original.ToMessageCommand()
	if err != nil {
		t.Fatal("ToMessageCommand:", err)
	}

	parsed := &OnStatus{}
	err = parsed.FromMessageCommand(cmd)
	if err != nil {
		t.Fatal("FromMessageCommand:", err)
	}

	if parsed.Transaction != 5 {
		t.Errorf("Transaction: got %d, want 5", parsed.Transaction)
	}
	if parsed.TcUrl != "rtmp://newserver:1935/app" {
		t.Errorf("TcUrl: got %q, want %q", parsed.TcUrl, "rtmp://newserver:1935/app")
	}
	if parsed.Description != "" {
		t.Errorf("Description: got %q, want empty", parsed.Description)
	}
}

func TestOnStatus_FromMessageCommand_NoParams(t *testing.T) {
	cmd := &message.Amf0CommandMessage{
		Command:       "onStatus",
		TransactionId: 3,
	}

	parsed := &OnStatus{}
	err := parsed.FromMessageCommand(cmd)
	if err != nil {
		t.Fatal("FromMessageCommand:", err)
	}

	if parsed.Transaction != 3 {
		t.Errorf("Transaction: got %d, want 3", parsed.Transaction)
	}
	if parsed.Level != "" {
		t.Errorf("Level: got %q, want empty", parsed.Level)
	}
	if parsed.Code != "" {
		t.Errorf("Code: got %q, want empty", parsed.Code)
	}
}

func TestOnStatus_FromMessageCommand_RawObject(t *testing.T) {
	cmd := &message.Amf0CommandMessage{
		Command:       "onStatus",
		TransactionId: 1,
		Parameters: []any{
			amf0.Object{
				"level":       "error",
				"code":        "NetStream.Play.StreamNotFound",
				"description": "stream not found",
			},
		},
	}

	parsed := &OnStatus{}
	err := parsed.FromMessageCommand(cmd)
	if err != nil {
		t.Fatal("FromMessageCommand:", err)
	}

	if parsed.Level != LevelError {
		t.Errorf("Level: got %q, want %q", parsed.Level, LevelError)
	}
	if parsed.Code != NetStreamPlayStreamNotFound {
		t.Errorf("Code: got %q, want %q", parsed.Code, NetStreamPlayStreamNotFound)
	}
	if parsed.Description != "stream not found" {
		t.Errorf("Description: got %q, want %q", parsed.Description, "stream not found")
	}
}

func TestOnStatus_CommandName(t *testing.T) {
	o := OnStatus{}
	if o.CommandName() != "onStatus" {
		t.Errorf("CommandName: got %q, want %q", o.CommandName(), "onStatus")
	}
}

func TestOnStatus_OmitsEmptyDescription(t *testing.T) {
	original := NewOnStatus(0, NewReconnectStatus("rtmp://host/app"))

	cmd, err := original.ToMessageCommand()
	if err != nil {
		t.Fatal("ToMessageCommand:", err)
	}

	amfCmd := cmd.(*message.Amf0CommandMessage)
	if len(amfCmd.Parameters) < 1 {
		t.Fatal("expected at least 1 parameter")
	}
	obj, ok := amfCmd.Parameters[0].(amf0.Object)
	if !ok {
		t.Fatalf("parameter 0: got %T, want amf0.Object", amfCmd.Parameters[0])
	}

	if _, found := obj["description"]; found {
		t.Error("expected 'description' to be absent when empty")
	}
	if _, found := obj["tcUrl"]; !found {
		t.Error("expected 'tcUrl' to be present")
	}
}

func TestNewOnStatus(t *testing.T) {
	status := NewStatus(NetStreamPublishStart)
	o := NewOnStatus(42, status)

	if o.Transaction != 42 {
		t.Errorf("Transaction: got %d, want 42", o.Transaction)
	}
	if o.Code != NetStreamPublishStart {
		t.Errorf("Code: got %q, want %q", o.Code, NetStreamPublishStart)
	}
	if o.Level != LevelStatus {
		t.Errorf("Level: got %q, want %q", o.Level, LevelStatus)
	}
}
