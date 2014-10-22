package logging

import (
	"github.com/deckarep/golang-set"
	"sync"
)

type HandlerCloser struct {
	handlers mapset.Set
	lock     sync.Mutex
}

func NewHandlerCloser() *HandlerCloser {
	return &HandlerCloser{
		handlers: mapset.NewThreadUnsafeSet(),
	}
}

func (self *HandlerCloser) AddHandler(handler Handler) {
	self.lock.Lock()
	defer self.lock.Unlock()
	if !self.handlers.Contains(handler) {
		self.handlers.Add(handler)
	}
}

func (self *HandlerCloser) RemoveHandler(handler Handler) {
	self.lock.Lock()
	defer self.lock.Unlock()
	if self.handlers.Contains(handler) {
		self.handlers.Remove(handler)
	}
}

func (self *HandlerCloser) Close() {
	self.lock.Lock()
	defer self.lock.Unlock()
	for i := range self.handlers.Iter() {
		handler, _ := i.(Handler)
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

func Shutdown() {
	Closer.Close()
	initialize()
}

// Return a logger with the specified name, creating it if necessary.
// If empty name is specified, return the root logger.
func GetLogger(name string) Logger {
	if len(name) > 0 {
		return manager.GetLogger(name)
	} else {
		return root
	}
}

// Log a message with severity "LevelFatal" on the root logger.
func Fatal(format string, args ...interface{}) {
	root.Fatal(format, args...)
}

// Log a message with severity "LevelError" on the root logger.
func Error(format string, args ...interface{}) {
	root.Error(format, args...)
}

// Log a message with severity "LevelWarn" on the root logger.
func Warn(format string, args ...interface{}) {
	root.Warn(format, args...)
}

// Log a message with severity "LevelInfo" on the root logger.
func Info(format string, args ...interface{}) {
	root.Info(format, args...)
}

// Log a message with severity "LevelDebug" on the root logger.
func Debug(format string, args ...interface{}) {
	root.Debug(format, args...)
}

// Log a message with specified severity level on the root logger.
func Log(level LogLevelType, format string, args ...interface{}) {
	root.Log(level, format, args...)
}
