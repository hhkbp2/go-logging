package handlers

import (
	"github.com/hhkbp2/go-logging"
	"github.com/hhkbp2/testify/require"
	"testing"
)

func TestNullHandler(t *testing.T) {
	handler := NewNullHandler()
	logger := logging.GetLogger("a")
	logger.AddHandler(handler)
	require.Equal(t, 1, len(logger.GetHandlers()))
	message := "test"
	logger.Debug(message)
	logger.Fatal(message)
}
