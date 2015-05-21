package logging

import (
	"log/syslog"
)

var (
	SyslogNameToPriorities = map[string]syslog.Priority{
		"LOG_EMERG":    syslog.LOG_EMERG,
		"LOG_ALERT":    syslog.LOG_ALERT,
		"LOG_CRIT":     syslog.LOG_CRIT,
		"LOG_ERR":      syslog.LOG_ERR,
		"LOG_WARNING":  syslog.LOG_WARNING,
		"LOG_NOTICE":   syslog.LOG_NOTICE,
		"LOG_INFO":     syslog.LOG_INFO,
		"LOG_DEBUG":    syslog.LOG_DEBUG,
		"LOG_KERN":     syslog.LOG_KERN,
		"LOG_USER":     syslog.LOG_USER,
		"LOG_MAIL":     syslog.LOG_MAIL,
		"LOG_DAEMON":   syslog.LOG_DAEMON,
		"LOG_AUTH":     syslog.LOG_AUTH,
		"LOG_SYSLOG":   syslog.LOG_SYSLOG,
		"LOG_LPR":      syslog.LOG_LPR,
		"LOG_NEWS":     syslog.LOG_NEWS,
		"LOG_UUCP":     syslog.LOG_UUCP,
		"LOG_CRON":     syslog.LOG_CRON,
		"LOG_AUTHPRIV": syslog.LOG_AUTHPRIV,
		"LOG_FTP":      syslog.LOG_FTP,
		"LOG_LOCAL0":   syslog.LOG_LOCAL0,
		"LOG_LOCAL1":   syslog.LOG_LOCAL1,
		"LOG_LOCAL2":   syslog.LOG_LOCAL2,
		"LOG_LOCAL3":   syslog.LOG_LOCAL3,
		"LOG_LOCAL4":   syslog.LOG_LOCAL4,
		"LOG_LOCAL5":   syslog.LOG_LOCAL5,
		"LOG_LOCAL6":   syslog.LOG_LOCAL6,
		"LOG_LOCAL7":   syslog.LOG_LOCAL7,
	}
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
