package logging

import (
	"errors"
	"github.com/hhkbp2/testify/require"
	"testing"
	"time"
)

var (
	ErrorTimeout = errors.New("timeout")
)

type MockHandler struct {
	*BaseHandler
	emitChan chan *LogRecord
	t        *testing.T
}

func NewMockHandler(t *testing.T) *MockHandler {
	object := &MockHandler{
		BaseHandler: NewBaseHandler("", LevelDebug),
		emitChan:    make(chan *LogRecord, 100),
	}
	Closer.AddHandler(object)
	return object
}

func (self *MockHandler) Emit(record *LogRecord) error {
	self.emitChan <- record
	return nil
}

func (self *MockHandler) Handle(record *LogRecord) int {
	return self.BaseHandler.Handle2(self, record)
}

func (self *MockHandler) HandleError(record *LogRecord, err error) {
	require.True(self.t, false, "should not be any error")
}

func (self *MockHandler) GetEmitOnTimeout(
	timeout time.Duration) (record *LogRecord, err error) {

	select {
	case record = <-self.emitChan:
		return record, nil
	case <-time.After(timeout):
		return nil, ErrorTimeout
	}
}

func TestLoggerLogToHandler(t *testing.T) {
	defer Shutdown()
	logger := GetLogger("b")
	logger.SetLevel(LevelDebug)
	require.Equal(t, 0, len(logger.GetHandlers()))
	handler := NewMockHandler(t)
	logger.AddHandler(handler)
	require.Equal(t, 1, len(logger.GetHandlers()))
	message := "abcd"
	logger.Debug(message)
	record, err := handler.GetEmitOnTimeout(time.Second * 0)
	require.Nil(t, err)
	require.Equal(t, message, record.GetMessage())
}
