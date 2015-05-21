package logging

import (
	"sync"
)

// An interface for dispatching logging events to specific destinations.
// Handler can optionally use formatter instances to format records as
// desired.
type Handler interface {
	// Return the name of Handler.
	GetName() string
	// Set the name of Handler.
	SetName(name string)
	// Return the log level of Handler.
	GetLevel() LogLevelType
	// Set the log level of Handler.
	SetLevel(level LogLevelType) error

	// For Formatter.
	// Format the specified record.
	Formatter
	// Set the formatter for this Handler.
	SetFormatter(formatter Formatter)

	// For Filter managing.
	Filterer

	// Do whatever it takes to actually log the specified logging record.
	Emit(record *LogRecord) error
	// Conditionally emit the specified logging record.
	Handle(record *LogRecord) int
	// Handle errors which occur during an Emit() call.
	HandleError(record *LogRecord, err error)
	// Ensure all logging output has been flushed.
	Flush() error
	// Tidy up any resources used by the handler.
	Close()
}

// The base handler class. Acts as a base parent of any concrete handler class.
// By default, no formatter is specified, in this case, the "raw" message as
// determined by record.Message is logged.
type BaseHandler struct {
	*StandardFilterer
	name          string
	nameLock      sync.RWMutex
	level         LogLevelType
	levelLock     sync.RWMutex
	formatter     Formatter
	formatterLock sync.RWMutex

	lock sync.Mutex
}

// Initialize the instance - basically setting the formatter to nil and the
// filterer without filter.
func NewBaseHandler(name string, level LogLevelType) *BaseHandler {
	return &BaseHandler{
		StandardFilterer: NewStandardFilterer(),
		name:             name,
		level:            level,
		formatter:        nil,
	}
}

func (self *BaseHandler) GetName() string {
	self.nameLock.RLock()
	defer self.nameLock.RUnlock()
	return self.name
}

func (self *BaseHandler) SetName(name string) {
	self.nameLock.Lock()
	defer self.nameLock.Unlock()
	self.name = name
}

func (self *BaseHandler) GetLevel() LogLevelType {
	self.levelLock.RLock()
	defer self.levelLock.RUnlock()
	return self.level
}

func (self *BaseHandler) SetLevel(level LogLevelType) error {
	self.levelLock.Lock()
	defer self.levelLock.Unlock()
	_, ok := getLevelName(level)
	if !ok {
		return ErrorNoSuchLevel
	}
	self.level = level
	return nil
}

func (self *BaseHandler) SetFormatter(formatter Formatter) {
	self.formatterLock.Lock()
	defer self.formatterLock.Unlock()
	self.formatter = formatter
}

// Acquire a lock for serializing access to the underlying I/O.
func (self *BaseHandler) Lock() {
	self.lock.Lock()
}

// Release the I/O lock.
func (self *BaseHandler) Unlock() {
	self.lock.Unlock()
}

// Format the specified record.
// If a formatter is set, use it. Otherwise, use the default formatter
// for the module.
func (self *BaseHandler) Format(record *LogRecord) string {
	self.formatterLock.RLock()
	self.formatterLock.RUnlock()
	var formatter Formatter
	if self.formatter != nil {
		formatter = self.formatter
	} else {
		formatter = defaultFormatter
	}
	return formatter.Format(record)
}

// A helper function for any subclass to define its Handle() method.
// Logging event emission depends on filters which may have heen added to
// the handler. Wrap the actual emission of the record and error handling
// with Lock()/Unlock() of the I/O lock. Returns non-zero if the filter passed
// the record for emission, else zero.
func (self *BaseHandler) Handle2(handler Handler, record *LogRecord) int {
	rv := handler.Filter(record)
	if rv > 0 {
		self.Lock()
		defer self.Unlock()
		err := handler.Emit(record)
		if err != nil {
			handler.HandleError(record, err)
		}
	}
	return rv
}

// A doing-nothing implementation as a stub for any subclass.
func (self *BaseHandler) HandleError(_ *LogRecord, _ error) {
	// Empty body
}

// A doing-nothing implementation as a stub for any subclass.
func (self *BaseHandler) Flush() error {
	// Empty body
	return nil
}

// A doing-nothing implementation as a stub for any subclass.
func (self *BaseHandler) Close() {
	// Empty body
}
