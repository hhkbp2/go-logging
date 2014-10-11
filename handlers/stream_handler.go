package handlers

import (
    "fmt"
    "github.com/hhkbp2/go-logging"
)

type Stream interface {
    Tell() (offset int64, err error)
    Write(s string) error
    Flush() error
    Close() error
}

type StreamHandler struct {
    *logging.BaseHandler
    stream Stream
}

func NewStreamHandler(
    name string, level logging.LogLevelType, stream Stream) *StreamHandler {

    return &StreamHandler{
        BaseHandler: logging.NewBaseHandler(name, level),
        stream:      stream,
    }
}

func (self *StreamHandler) GetStream() Stream {
    return self.stream
}

func (self *StreamHandler) SetStream(s Stream) {
    self.stream = s
}

func (self *StreamHandler) Flush() error {
    return self.stream.Flush()
}

func (self *StreamHandler) Emit(record *logging.LogRecord) error {
    return self.Emit2(self, record)
}

func (self *StreamHandler) Emit2(
    handler logging.Handler, record *logging.LogRecord) error {

    message := handler.Format(record)
    if err := self.stream.Write(fmt.Sprintf("%s\n", message)); err != nil {
        return err
    }
    if err := handler.Flush(); err != nil {
        return err
    }
    return nil
}

func (self *StreamHandler) Handle(record *logging.LogRecord) int {
    return self.BaseHandler.Handle2(self, record)
}
