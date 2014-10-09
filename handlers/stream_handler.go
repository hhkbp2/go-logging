package handlers

import (
    "github.com/hhkbp2/go-logging"
)

type Stream interface {
    Write(s string) error
    Flush() error
}

type StreamHandler struct {
    *DefaultHandler
    stream Stream
}

func NewStreamHandler(stream Stream) *StreamHandler {
    return &StreamHandler{
        DefaultHandler: NewDefaultHandler(),
        stream:         stream,
    }
}

func (self *StreamHandler) Flush() {
    self.Acquire()
    defer self.Release()
    self.stream.Flush()
}

func (self *StreamHandler) Emit(record *LogRecord) {
    message := self.Format(record)
    stream := self.stream
    err := stream.Write(fmt.Sprintf("%s\n", message))
    if err != nil {
        self.HandleError(record, err)
    }
    if err = self.Flush(); err != nil {
        self.HandleError(record, err)
    }
}
