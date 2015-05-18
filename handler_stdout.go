package logging

import (
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

type StdoutHandler struct {
	*StreamHandler
}

func NewStdoutHandler() *StdoutHandler {
	stream := NewStdoutStream()
	handler := NewStreamHandler("stdout", LevelNotset, stream)
	object := &StdoutHandler{
		StreamHandler: handler,
	}
	Closer.RemoveHandler(object.StreamHandler)
	Closer.AddHandler(object)
	return object
}

func (self *StdoutHandler) Emit(record *LogRecord) error {
	return self.StreamHandler.Emit2(self, record)
}

func (self *StdoutHandler) Handle(record *LogRecord) int {
	return self.Handle2(self, record)
}

func (self *StdoutHandler) Close() {
	self.StreamHandler.Close2()
}
