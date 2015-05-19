package logging

import (
	"testing"
)

func TestStdoutHandler(_ *testing.T) {
	logger := GetLogger("stdout")
	handler := NewStdoutHandler()
	logger.AddHandler(handler)
	logger.Warnf("test message")
}
