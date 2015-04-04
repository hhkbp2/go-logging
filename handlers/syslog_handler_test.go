package handlers

import (
	"github.com/hhkbp2/go-logging"
	"github.com/hhkbp2/testify/require"
	"log/syslog"
	"testing"
)

func TestSyslogHandler(t *testing.T) {
	defer logging.Shutdown()
	handler, err := NewSyslogHandler(
		syslog.LOG_USER|syslog.LOG_DEBUG,
		"atag")
	require.Nil(t, err)
	logger := logging.GetLogger("test")
	logger.SetLevel(logging.LevelDebug)
	logger.AddHandler(handler)
	prefix := "test syslog handler "
	logger.Debugf(prefix + "Debug() a message")
	logger.Errorf(prefix + "Error() a message")
}
