package handlers

import (
	"github.com/hhkbp2/go-logging"
	"log/syslog"
)

type SyslogHandler struct {
	*logging.BaseHandler
	network  string
	raddr    string
	priority syslog.Priority
	tag      string
	writer   *syslog.Writer
}

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
