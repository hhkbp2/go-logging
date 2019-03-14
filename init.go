package logging

import (
	"sync"
)

type HandlerCloser struct {
	handlers *ListSet
	lock     sync.Mutex
}

func NewHandlerCloser() *HandlerCloser {
	return &HandlerCloser{
		handlers: NewListSet(),
	}
}

func (self *HandlerCloser) AddHandler(handler Handler) {
	self.lock.Lock()
	defer self.lock.Unlock()
	if !self.handlers.SetContains(handler) {
		self.handlers.SetAdd(handler)
	}
}

func (self *HandlerCloser) RemoveHandler(handler Handler) {
	self.lock.Lock()
	defer self.lock.Unlock()
	if self.handlers.SetContains(handler) {
		self.handlers.SetRemove(handler)
	}
}

func (self *HandlerCloser) Close() {
	self.lock.Lock()
	defer self.lock.Unlock()
	for e := self.handlers.Front(); e != nil; e = e.Next() {
		handler, _ := e.Value.(Handler)
		handler.Close()
	}
}

var (
	root    Logger
	manager *Manager
	Closer  *HandlerCloser
)

func init() {
	initialize()
}

func initialize() {
	root = NewRootLogger(LevelWarn)
	manager = NewManager(root)
	Closer = NewHandlerCloser()
}

// Shutdown ensures all log messages are flushed before program exits.
func Shutdown() {
	Closer.Close()
	initialize()
}

// Set logger maker for default manager.
func SetLoggerMaker(maker LoggerMaker) {
	manager.SetLoggerMaker(maker)
}

// GetLogger returns a logger with the specified name, creating it if necessary.
// If empty name is specified, return the root logger.
func GetLogger(name string) Logger {
	if len(name) > 0 {
		return manager.GetLogger(name)
	} else {
		return root
	}
}

// Log a message with severity "LevelFatal" on the root logger.
func Fatalf(format string, args ...interface{}) {
	root.Fatalf(format, args...)
}

// Log a message with severity "LevelError" on the root logger.
func Errorf(format string, args ...interface{}) {
	root.Errorf(format, args...)
}

// Log a message with severity "LevelWarn" on the root logger.
func Warnf(format string, args ...interface{}) {
	root.Warnf(format, args...)
}

// Log a message with severity "LevelInfo" on the root logger.
func Infof(format string, args ...interface{}) {
	root.Infof(format, args...)
}

// Log a message with severity "LevelDebug" on the root logger.
func Debugf(format string, args ...interface{}) {
	root.Debugf(format, args...)
}

// Log a message with severity "LevelTrace" on the root logger.
func Tracef(format string, args ...interface{}) {
	root.Tracef(format, args...)
}

// Log a message with specified severity level on the root logger.
func Logf(level LogLevelType, format string, args ...interface{}) {
	root.Logf(level, format, args...)
}
