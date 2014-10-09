package logging

import (
    "github.com/deckarep/golang-set"
    "sync"
)

var (
    ErrorNoSuchLevel = errors.New("no such level")
)

type NodeType uint8

const (
    NodeUnknown NodeType = 0 + iota
    NodeLogger
    NodePlaceHolder
)

type Node interface {
    Type() NodeType
    GetParent() Node
    SetParent(node Node)
}

type BaseNode struct {
    parent Node
    lock   sync.RWMutex
}

func NewBaseNode() *BaseNode {
    return &BaseNode{
        parent: nil,
    }
}

func (self *BaseNode) GetParent() Node {
    self.lock.Rlock()
    defer self.lock.RUnlock()
    return self.parent
}

func (self *BaseNode) SetParent(node Node) {
    self.lock.Lock()
    defer self.lock.Unlock()
    self.parent = node
}

type PlaceHolder struct {
    *BaseNode
    Loggers mapset.Set
}

func NewPlaceHolder(logger Logger) *PlaceHolder {
    object := &PlaceHolder{
        BaseNode: NewBaseNode(),
        Loggers:  mapset.NewThreadUnsafeSet(),
    }
    object.Append(logger)
    return object
}

func (self *PlaceHolder) Type() NodeType {
    return NodePlaceHolder
}

func (self *PlaceHolder) Append(logger Logger) {
    if !self.Loggers.Contains(logger) {
        self.Loggers.Add(logger)
    }
}

type Logger interface {
    Node

    GetName() string
    GetPropagate() bool
    SetPropagate(v bool)
    GetLevel() LogLevelType
    SetLevel(level LogLevelType)
    IsEnableFor(level LogLevelType) bool
    GetEffectiveLevel() LogLevelType

    Fatal(format string, args ...interface{})
    Error(format string, args ...interface{})
    Warn(format string, args ...interface{})
    Info(format string, args ...interface{})
    Debug(format string, args ...interface{})
    Log(level LogLevelType, format string, args ...interface{})

    AddHandler(handler Handler)
    RemoveHandler(handler Handler)
    AddFilter(filter Filter)
    RemoveFilter(filter Filter)

    GetManager() *Manager
    SetManager(manager *Manager)
}

type StandardLogger struct {
    *BaseNode
    *Filterer
    name      string
    level     LogLevelType
    parent    Node
    propagate bool
    handlers  mapset.Set
    manager   *Manager
    lock      sync.RWMutex
}

func NewStandardLogger(name string, level LogLevelType) *StandardLogger {
    return &StandardLogger{
        BaseNode:  NewBaseNode(),
        Filterer:  NewDefaultFilterer(),
        name:      name,
        level:     level,
        parent:    nil,
        propagate: true,
        handlers:  mapset.NewThreadSafeSet(),
        manager:   nil,
    }
}

func (self *StandardLogger) Type() NodeType {
    return NodeLogger
}

func (self *StandardLogger) GetName() string {
    self.lock.RLock()
    defer self.lock.RUnlock()
    return self.name
}

func (self *StandardLogger) GetPropagate() bool {
    self.lock.RLock()
    defer self.lock.RUnlock()
    return self.propagate
}

func (self *StandardLogger) SetPropagate(v bool) {
    self.lock.Lock()
    defer self.lock.Unlock()
    self.propagate = v
}

func (self *StandardLogger) GetLevel() LogLevelType {
    self.lock.RLock()
    defer self.lock.RUnlock()
    return self.level
}

func (self *StandardLogger) SetLevel(level LogLevelType) error {
    levelName, ok = getLevelName(level)
    if !ok {
        return ErrorNoSuchLevel
    }
    self.lock.Lock()
    defer self.lock.Unlock()
    self.level = level
    return nil
}

func (self *StandardLogger) IsEnabledFor(level LogLevelType) bool {
    return level >= self.GetEffectiveLevel()
}

func (self *StandardLogger) GetEffectiveLevel() LogLevelType {
    self.lock.RLock()
    defer self.lock.RUnlock()
    logger := self
    for logger != nil {
        if logger.level != Notset {
            return logger.level
        }
        logger = logger.parent
    }
    return Notset
}

func (self *StandardLogger) Fatal(format string, args ...interface{}) {
    if self.IsEnabledFor(Fatal) {
        self.log(Fatal, format, args...)
    }
}

func (self *StandardLogger) Error(format string, args ...interface{}) {
    if self.IsEnabledFor(Error) {
        self.log(Error, format, args...)
    }
}

func (self *StandardLogger) Warn(format string, args ...interface{}) {
    if self.IsEnabledFor(Warn) {
        self.log(Warn, format, args...)
    }
}

func (self *StandardLogger) Info(format string, args ...interface{}) {
    if self.IsEnabledFor(Info) {
        self.log(Info, format, args...)
    }
}

func (self *StandardLogger) Debug(format string, args ...interface{}) {
    if self.IsEnabledFor(Debug) {
        self.log(Debug, format, args...)
    }
}

func (self *StandardLogger) Log(
    level LogLevelType, format string, args ...interface{}) {

    if self.IsEnabledFor(level) {
        self.log(level, format, args...)
    }
}

func (self *StandardLogger) log(
    level LogLevelType, format string, args ...interface{}) {

    // TODO get pathName and lineNo
    pathName := ""
    lineNo := 0
    record = NewLogRecord(self.name, level, pathName, lineNo, format, args)
    self.Handle(record)
}

func (self *StandardLogger) Handle(record *LogRecord) {
    if self.Filter(record) > 0 {
        self.callHandlers(record)
    }
}

func (self *StandardLogger) callHandlers(record *LogRecord) {
    self.lock.Lock()
    defer self.lock.Unlock()
    call := self
    found := 0
    for call != nil {
        for handler := range call.handlers {
            found += 1
            if record.Level >= handler.Level {
                handler.Handle(record)
            }
        }
        if !c.propagate {
            c = nil
        } else {
            c = c.parent
        }
    }
}

func (self *StandardLogger) AddHandler(handler Handler) {
    self.lock.Lock()
    defer self.lock.Unlock()
    if !self.handlers.Contains(handler) {
        self.handlers.Add(handler)
    }
}

func (self *StandardLogger) RemoveHandler(handler Handler) {
    self.lock.Lock()
    defer self.lock.Unlock()
    if self.handlers.Contains(handler) {
        self.handlers.Remove(handler)
    }
}

func (self *StandardLogger) GetManager() *Manager {
    self.lock.RLock()
    defer self.lock.RUnlock()
    return self.manager
}

func (self *StandardLogger) SetManager(manager *Manager) {
    self.lock.Lock()
    defer self.lock.Unlock()
    self.manager = manager
}

type RootLogger struct {
    *StandardLogger
}

func NewRootLogger(level LogLevelType) *RootLogger {
    return &RootLogger{
        StandardLogger: NewStandardLogger("root", level),
    }
}

type LoggerMaker func(name string) Logger

func defaultLoggerMaker(name string) Logger {
    return NewStandardLogger(name, Notset)
}

type Manager struct {
    root        Logger
    loggers     map[string]Node
    loggerMaker LoggerMaker
    lock        sync.Mutex
}

func NewManager(logger Logger) *Manager {
    return &Manager{
        root:            logger,
        loggers:         make(map[string]Node),
        loggerGenerator: defaultLoggerGenerator,
    }
}

func (self *Manager) SetLoggerMaker(maker LoggerMaker) {
    self.lock.Lock()
    defer self.lock.Unlock()
    self.loggerMaker = maker
}

func (self *Manager) GetLogger(name string) Logger {
    lock.Lock()
    defer lock.Unlock()
    var logger Logger
    node, ok := self.loggers[name]
    if ok {
        switch node.Type() {
        case NodePlaceHolder:
            placeHolder, _ := node.(*PlaceHolder)
            logger = self.loggerMaker(name)
            logger.SetManager(self)
            self.loggers[name] = logger
            self.fixupChildren(placeHolder, logger)
            self.fixupParents(logger)
        case NodeLogger:
            logger, _ = node.(*StandardLogger)
        default:
            panic("invalid node type")
        }
    } else {
        logger = self.loggerMaker(name)
        logger.SetManager(self)
        self.loggers[name] = logger
        self.fixupParents(logger)
    }
    return logger
}

func (self *StandardLogger) fixupParents(logger Logger) {
    name := logger.GetName()
    index := strings.LastIndex(name, ".")
    var parent Node
    if (index > 0) && (parent == nil) {
        parentStr := name[:index]
        node, ok := self.loggers[name]
        if !ok {
            self.loggers[name] = NewPlaceHolder(logger)
        } else {
            switch node.Type() {
            case NodePlaceHolder:
                placeHolder, _ := node.(*PlaceHolder)
                placeHolder.Append(logger)
            case NodeLogger:
                parent = logger
            default:
                panic("invalid node type")
            }
        }
        index = strings.LastIndex(parentStr, ".")
    }
    if parent == nil {
        parent = root
    }
    logger.SetParent(parent)
}

func (self *StandardLogger) fixupChildren(
    placeHolder *PlaceHolder, logger Logger) {

    name := logger.GetName()
    for l := range placeHolder.Loggers.Iter() {
        parent, _ := l.GetParent().(*Logger)
        if !strings.HasPrefix(parent.GetName(), name) {
            logger.SetParent(parent)
            l.SetParent(logger)
        }
    }
}
