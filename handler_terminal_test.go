package logging

import (
	"testing"
)

func TestTerminalHandler(_ *testing.T) {
	logger := GetLogger("a.b")
	handler := NewTerminalHandler()
	logger.AddHandler(handler)
	logger.Warnf("test message")
}
