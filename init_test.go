package logging

import (
	"testing"
	"time"

	"github.com/hhkbp2/testify/require"
)

func TestShutdown(t *testing.T) {
	defer Shutdown()
	handler := NewMockHandler(t)
	logger := GetLogger("a")
	logger.AddHandler(handler)
	message := "abcd"
	logger.Errorf(message)
	record, err := handler.GetEmitOnTimeout(time.Second * 1)
	require.Nil(t, err)
	require.Equal(t, message, record.GetMessage())
}
