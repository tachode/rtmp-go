package amf0_test

import (
	"bytes"
	"encoding/binary"
	"testing"
	"time"

	"github.com/tachode/rtmp-go/amf0"
)

func TestDate_Read(t *testing.T) {
	date := amf0.Date{}
	buf := new(bytes.Buffer)

	// Write a sample date to the buffer
	expectedTime := time.Date(2021, time.January, 1, 0, 0, 0, 0, time.UTC)
	err := binary.Write(buf, binary.BigEndian, float64(expectedTime.UnixMilli()))
	if err != nil {
		t.Fatalf("Failed to write to buffer: %v", err)
	}
	err = binary.Write(buf, binary.BigEndian, int16(0))
	if err != nil {
		t.Fatalf("Failed to write timezone to buffer: %v", err)
	}

	// Read the date from the buffer
	err = date.Read(buf)
	if err != nil {
		t.Fatalf("Failed to read date: %v", err)
	}

	// Check if the date matches the expected value
	if !time.Time(date).Equal(expectedTime) {
		t.Errorf("Expected %v, got %v", expectedTime, time.Time(date))
	}
}

func TestDate_Write(t *testing.T) {
	expectedTime := time.Date(2021, time.January, 1, 0, 0, 0, 0, time.UTC)
	date := amf0.Date(expectedTime)
	buf := new(bytes.Buffer)

	// Write the date to the buffer
	err := date.Write(buf)
	if err != nil {
		t.Fatalf("Failed to write date: %v", err)
	}

	// Read the date from the buffer
	var value float64
	err = binary.Read(buf, binary.BigEndian, &value)
	if err != nil {
		t.Fatalf("Failed to read from buffer: %v", err)
	}
	var timezone int16
	err = binary.Read(buf, binary.BigEndian, &timezone)
	if err != nil {
		t.Fatalf("Failed to read timezone from buffer: %v", err)
	}

	// Check if the date matches the expected value
	if int64(value) != expectedTime.UnixMilli() {
		t.Errorf("Expected %v, got %v", expectedTime.UnixMilli(), int64(value))
	}
	if timezone != 0 {
		t.Errorf("Expected timezone 0, got %v", timezone)
	}
}
