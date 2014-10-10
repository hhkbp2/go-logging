package handlers

import (
    "fmt"
    "github.com/hhkbp2/go-logging"
)

type Stream interface {
    Write(s string) error
    Flush() error
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

func (self *StreamHandler) Flush() error {
    self.Lock()
    defer self.Unlock()
    return self.stream.Flush()
}

func (self *StreamHandler) Emit(record *logging.LogRecord) error {
    message := self.Format(record)
    stream := self.stream
    err := stream.Write(fmt.Sprintf("%s\n", message))
    if err != nil {
        self.HandleError(record, err)
    }
    if err = self.Flush(); err != nil {
        self.HandleError(record, err)
    }
    return err
}
