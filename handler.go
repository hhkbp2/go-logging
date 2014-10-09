package logging

import (
    "sync"
)

var (
    handlers      = make(map[string]Handler)
    handlersMutex sync.RWMutex
)

func handlersLock() {
    handlersMutex.Lock()
}

func handlersUnlock() {
    handlersMutex.Unlock()
}

type Handler interface {
    GetName() string
    SetName(name string)
    GetLevel() LogLevelType
    SetLevel(level LogLevelType) error
    SetFormatter(formater Formatter)
    Lock()
    Unlock()
    Format(record *LogRecord) string
    Emit(record *LogRecord) error
    Handle(record *LogRecord) int
    HandleError(record *LogRecord, err error)
    Flush() error
    Close()
}

type BaseHandler struct {
    // TODO add impl
    *Filterer
    name      string
    level     LogLevelType
    formatter Formatter

    lock sync.RWMutex
}

func NewBaseHandler(name string, level LogLevelType) *BaseHandler {
    return &BaseHandler{
        Filterer:  NewFilterer(),
        name:      name,
        level:     level,
        formatter: nil,
    }
}

func (self *BaseHandler) GetName() string {
    return self.name
}

func (self *BaseHandler) SetName(name string, handler Handler) {
    handlersLock()
    defer handlersUnlock()
    if len(self.name) > 0 {
        if _, ok := handlers[self.name]; ok {
            delete(handlers, self.name)
        }
    }
    handler.SetName(name)
    if len(name) > 0 {
        handlers[name] = handler
    }
}

func (self *BaseHandler) GetLevel() LogLevelType {
    return self.level
}

func (self *BaseHandler) SetLevel(level LogLevelType) error {
    _, ok := getLevelName(level)
    if !ok {
        return ErrorNoSuchLevel
    }
    self.level = level
    return nil
}

func (self *BaseHandler) SetFormatter(formatter Formatter) {
    self.formatter = formatter
}

func (self *BaseHandler) Lock() {
    self.lock.Lock()
}

func (self *BaseHandler) Unlock() {
    self.lock.Unlock()
}

func (self *BaseHandler) Format(record *LogRecord) string {
    var formatter Formatter
    if self.formatter != nil {
        formatter = self.formatter
    } else {
        formatter = defaultFormatter
    }
    return formatter.Format(record)
}

func (self *BaseHandler) Handle(handler Handler, record *LogRecord) int {
    rv := self.Filter(record)
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

func (self *BaseHandler) Close() {
    handlersLock()
    defer handlersUnlock()
    if len(self.name) > 0 {
        _, ok := handlers[self.name]
        if ok {
            delete(handlers, self.name)
        }
    }
}
