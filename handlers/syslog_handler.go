package handlers

import (
	"github.com/hhkbp2/go-logging"
	"log/syslog"
)

// A handler class which sends formatted logging records to a syslog server.
type SyslogHandler struct {
	*logging.BaseHandler
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
		BaseHandler: logging.NewBaseHandler("", logging.LevelNotset),
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
		BaseHandler: logging.NewBaseHandler("", logging.LevelNotset),
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
func (self *SyslogHandler) Emit(record *logging.LogRecord) error {
	message := self.BaseHandler.Format(record)
	var err error
	switch record.Level {
	case logging.LevelFatal:
		err = self.writer.Crit(message)
	case logging.LevelError:
		err = self.writer.Err(message)
	case logging.LevelWarn:
		err = self.writer.Warning(message)
	case logging.LevelInfo:
		err = self.writer.Info(message)
	case logging.LevelDebug:
		err = self.writer.Debug(message)
	default:
		_, err = self.writer.Write([]byte(message))
	}
	return err
}

func (self *SyslogHandler) Handle(record *logging.LogRecord) int {
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
