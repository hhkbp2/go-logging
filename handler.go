package logging

import ()

var (
    handlers     = make(map[string]Handler)
    handlersLock sync.RWMutex
)

func handlersLock() {
    handlersLock.Lock()
}

func handlersUnlock() {
    handlersLock.Unlock()
}

type Handler interface {
    GetName() string
    SetName(name string)
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
    Name      string
    Level     LogLevelType
    Formatter Formatter

    lock sync.RWMutex
}

func NewBaseHandler(name string, level LogLevelType) *BaseHandler {
    return &BaseHandler{
        Filterer:  NewFilterer(),
        Name:      name,
        Level:     level,
        Formatter: nil,
    }
}

func (self *BaseHandler) GetName() {
    return self.Name
}

func (self *BaseHandler) SetName(name string) {
    handlersLock()
    defer handlersUnlock()
    if len(self.name) > 0 {
        if handler, ok := handlers[self.name]; ok {
            delete(handlers, self.name)
        }
    }
    self.name = name
    if len(name) > 0 {
        handlers[name] = self
    }
}

func (self *BaseHandler) SetFormatter(formatter Formatter) {
    self.Formatter = formatter
}

func (self *BaseHandler) Lock() {
    self.lock.Lock()
}

func (self *BaseHandler) Unlock() {
    self.lock.Unlock()
}

func (self *BaseHandler) Format(record *LogRecord) string {
    var formatter Formatter
    if self.Formatter {
        formatter = self.Formatter
    } else {
        formatter = defaultFormatter
    }
    return formatter.Format(record)
}

func (self *BaseHandler) Handle(handler Handler, record *LogRecord) int {
    rv := self.Filter(record)
    if rv {
        self.Lock()
        defer self.Unlock()
        handler.Emit(record)
    }
    return rv
}

func (self *BaseHandler) HandleError(record *LogRecord, err error) {
    // TODO add impl
}

func (self *BaseHandler) Close() {
    handlersLock()
    defer handlersUnlock()
    if len(self.name) > 0 {
        handler, ok := handlers[self.name]
        if ok {
            delete(handlers, self.name)
        }
    }
}
