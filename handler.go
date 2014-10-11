package logging

import (
    "sync"
)

type Handler interface {
    GetName() string
    SetName(name string)
    GetLevel() LogLevelType
    SetLevel(level LogLevelType) error

    Formatter
    SetFormatter(formater Formatter)

    Filterer

    Emit(record *LogRecord) error
    Handle(record *LogRecord) int
    HandleError(record *LogRecord, err error)
    Flush() error
    Close()
}

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

func (self *BaseHandler) Lock() {
    self.lock.Lock()
}

func (self *BaseHandler) Unlock() {
    self.lock.Unlock()
}

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

func (self *BaseHandler) HandleError(_ *LogRecord, _ error) {
    // Empty body
}

func (self *BaseHandler) Flush() error {
    // Empty body
    return nil
}

func (self *BaseHandler) Close() {
    // Empty body
}
