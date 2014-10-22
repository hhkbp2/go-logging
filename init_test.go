package logging

import (
	"github.com/hhkbp2/testify/require"
	"testing"
	"time"
)

func TestShutdown(t *testing.T) {
	defer Shutdown()
	handler := NewMockHandler(t)
	logger := GetLogger("a")
	logger.AddHandler(handler)
	message := "abcd"
	logger.Error(message)
	record, err := handler.GetEmitOnTimeout(time.Second * 0)
	require.Nil(t, err)
	require.Equal(t, message, record.GetMessage())
}
