// +build !windows

package logging

import (
	"log/syslog"
)

// A handler class which sends formatted logging records to a syslog server.
type SyslogHandler struct {
	*BaseHandler
	network  string
	raddr    string
	priority syslog.Priority
	tag      string
	writer   *syslog.Writer
}

// Initialize a syslog handler.
// The arguements are the same as New() in package log/syslog.
func NewSyslogHandler(
	priority syslog.Priority,
	tag string) (*SyslogHandler, error) {

	writer, err := syslog.New(priority, tag)
	if err != nil {
		return nil, err
	}
	object := &SyslogHandler{
		BaseHandler: NewBaseHandler("", LevelNotset),
		network:     "",
		raddr:       "",
		priority:    priority,
		tag:         tag,
		writer:      writer,
	}
	return object, nil
}

// Initialize a syslog handler with connection to a specified syslog server.
// The arguements are the same as Dial() in package log/syslog.
func NewSyslogHandlerToAddr(
	network, raddr string,
	priority syslog.Priority,
	tag string) (*SyslogHandler, error) {

	writer, err := syslog.Dial(network, raddr, priority, tag)
	if err != nil {
		return nil, err
	}
	object := &SyslogHandler{
		BaseHandler: NewBaseHandler("", LevelNotset),
		network:     network,
		raddr:       raddr,
		priority:    priority,
		tag:         tag,
		writer:      writer,
	}
	return object, nil
}

// Emit a record.
// The record is formatted, and then sent to the syslog server
// in specified log level.
func (self *SyslogHandler) Emit(record *LogRecord) error {
	message := self.BaseHandler.Format(record)
	var err error
	switch record.Level {
	case LevelFatal:
		err = self.writer.Crit(message)
	case LevelError:
		err = self.writer.Err(message)
	case LevelWarn:
		err = self.writer.Warning(message)
	case LevelInfo:
		err = self.writer.Info(message)
	case LevelDebug:
		err = self.writer.Debug(message)
	case LevelTrace:
		err = self.writer.Debug(message)
	default:
		_, err = self.writer.Write([]byte(message))
	}
	return err
}

func (self *SyslogHandler) Handle(record *LogRecord) int {
	return self.BaseHandler.Handle2(self, record)
}

func (self *SyslogHandler) Flush() error {
	// Nothing to do
	return nil
}

func (self *SyslogHandler) Close() {
	// ignore the error return code
	self.writer.Close()
}
