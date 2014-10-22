package handlers

import (
	"github.com/hhkbp2/go-logging"
	"os"
)

type StdoutStream struct {
	stdout *os.File
}

func NewStdoutStream() *StdoutStream {
	return &StdoutStream{
		stdout: os.Stdout,
	}
}

func (self *StdoutStream) Tell() (int64, error) {
	// Stdout.Stat().Size() always return 0
	return 0, nil
}

func (self *StdoutStream) Write(s string) error {
	_, err := self.stdout.WriteString(s)
	return err
}

func (self *StdoutStream) Flush() error {
	// Empty body
	return nil
}

func (self *StdoutStream) Close() error {
	// Don't close stdout
	return nil
}

type TerminalHandler struct {
	*StreamHandler
}

func NewTerminalHandler() *TerminalHandler {
	stream := NewStdoutStream()
	handler := NewStreamHandler("stdout", logging.LevelNotset, stream)
	object := &TerminalHandler{
		StreamHandler: handler,
	}
	logging.Closer.RemoveHandler(object.StreamHandler)
	logging.Closer.AddHandler(object)
	return object
}

func (self *TerminalHandler) Emit(record *logging.LogRecord) error {
	return self.StreamHandler.Emit2(self, record)
}

func (self *TerminalHandler) Handle(record *logging.LogRecord) int {
	return self.Handle2(self, record)
}

func (self *TerminalHandler) Close() {
	self.StreamHandler.Close2()
}
