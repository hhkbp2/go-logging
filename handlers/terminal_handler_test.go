package handlers

import (
    "github.com/hhkbp2/go-logging"
    "testing"
)

func TestTerminalHandler(_ *testing.T) {
    logger := logging.GetLogger("a.b")
    handler := NewTerminalHandler()
    logger.AddHandler(handler)
    logger.Warn("test message")
}
