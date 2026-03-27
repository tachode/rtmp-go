package command

import (
	"testing"
)

func TestNewReconnectStatus(t *testing.T) {
	s := NewReconnectStatus("rtmp://newserver:1935/app")
	if s.Level != LevelStatus {
		t.Errorf("Level = %q, want %q", s.Level, LevelStatus)
	}
	if s.Code != NetConnectionConnectReconnectRequest {
		t.Errorf("Code = %q, want %q", s.Code, NetConnectionConnectReconnectRequest)
	}
	if s.TcUrl != "rtmp://newserver:1935/app" {
		t.Errorf("TcUrl = %q, want %q", s.TcUrl, "rtmp://newserver:1935/app")
	}
	if s.Description != "" {
		t.Errorf("Description = %q, want empty", s.Description)
	}
}

func TestStatus_ToObject_OmitsEmptyFields(t *testing.T) {
	s := NewStatus(NetConnectionConnectSuccess)
	obj := s.ToObject()

	if _, ok := obj["level"]; !ok {
		t.Error("expected 'level' in object")
	}
	if _, ok := obj["code"]; !ok {
		t.Error("expected 'code' in object")
	}
	if _, ok := obj["description"]; !ok {
		t.Error("expected 'description' in object when non-empty")
	}
	if _, ok := obj["tcUrl"]; ok {
		t.Error("expected 'tcUrl' to be absent when empty")
	}
}

func TestStatus_ToObject_IncludesTcUrl(t *testing.T) {
	s := NewReconnectStatus("rtmp://newhost/app")
	obj := s.ToObject()

	tcUrl, ok := obj["tcUrl"]
	if !ok {
		t.Fatal("expected 'tcUrl' in object")
	}
	if tcUrl != "rtmp://newhost/app" {
		t.Errorf("tcUrl = %q, want %q", tcUrl, "rtmp://newhost/app")
	}
	if _, ok := obj["description"]; ok {
		t.Error("expected 'description' to be absent when empty")
	}
}
