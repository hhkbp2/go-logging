package logging

import (
    "errors"
    "github.com/deckarep/golang-set"
    "runtime"
    "strings"
    "sync"
)

var (
    ErrorNoSuchLevel = errors.New("no such level")
)

const (
    thisPackageName = "go-logging"
    thisFileName    = "logger.go"
)

type NodeType uint8

const (
    NodeUnknown NodeType = 0 + iota
    NodeLogger
    NodePlaceHolder
)

type Node interface {
    Type() NodeType
}

type PlaceHolder struct {
    Loggers mapset.Set
}

func NewPlaceHolder(logger Logger) *PlaceHolder {
    object := &PlaceHolder{
        Loggers: mapset.NewThreadUnsafeSet(),
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
    SetLevel(level LogLevelType) error
    IsEnabledFor(level LogLevelType) bool
    GetEffectiveLevel() LogLevelType

    Fatal(format string, args ...interface{})
    Error(format string, args ...interface{})
    Warn(format string, args ...interface{})
    Info(format string, args ...interface{})
    Debug(format string, args ...interface{})
    Log(level LogLevelType, format string, args ...interface{})

    AddHandler(handler Handler)
    RemoveHandler(handler Handler)
    GetHandlers() []Handler
    AddFilter(filter Filter)
    RemoveFilter(filter Filter)
    GetFilters() []Filter

    GetManager() *Manager
    SetManager(manager *Manager)
    GetParent() Logger
    SetParent(parent Logger)
}

type StandardLogger struct {
    *Filterer
    name      string
    level     LogLevelType
    parent    Logger
    propagate bool
    handlers  mapset.Set
    manager   *Manager
    lock      sync.RWMutex
}

func NewStandardLogger(name string, level LogLevelType) *StandardLogger {
    return &StandardLogger{
        Filterer:  NewFilterer(),
        parent:    nil,
        name:      name,
        level:     level,
        propagate: true,
        handlers:  mapset.NewSet(),
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
    _, ok := getLevelName(level)
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
    var logger Logger = self
    for logger != nil {
        level := logger.GetLevel()
        if level != LevelNotset {
            return level
        }
        logger = logger.GetParent()
    }
    return LevelNotset
}

func (self *StandardLogger) Fatal(format string, args ...interface{}) {
    if self.IsEnabledFor(LevelFatal) {
        self.log(LevelFatal, format, args...)
    }
}

func (self *StandardLogger) Error(format string, args ...interface{}) {
    if self.IsEnabledFor(LevelError) {
        self.log(LevelError, format, args...)
    }
}

func (self *StandardLogger) Warn(format string, args ...interface{}) {
    if self.IsEnabledFor(LevelWarn) {
        self.log(LevelWarn, format, args...)
    }
}

func (self *StandardLogger) Info(format string, args ...interface{}) {
    if self.IsEnabledFor(LevelInfo) {
        self.log(LevelInfo, format, args...)
    }
}

func (self *StandardLogger) Debug(format string, args ...interface{}) {
    if self.IsEnabledFor(LevelDebug) {
        self.log(LevelDebug, format, args...)
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

    callerInfo := self.findCaller()
    record := NewLogRecord(
        self.name,
        level,
        callerInfo.PathName,
        callerInfo.FileName,
        callerInfo.LineNo,
        callerInfo.FuncName,
        format,
        args)
    self.Handle(record)
}

type CallerInfo struct {
    PathName string
    FileName string
    LineNo   uint32
    FuncName string
}

var (
    UnknownCallerInfo = &CallerInfo{
        PathName: "(unknown path)",
        FileName: "(unknown file)",
        LineNo:   0,
        FuncName: "(unknown function)",
    }
)

func (self *StandardLogger) findCaller() *CallerInfo {
    for i := 1; ; i++ {
        pc, filepath, line, ok := runtime.Caller(i)
        if !ok {
            return UnknownCallerInfo
        }
        parts := strings.Split(filepath, "/")
        dir := parts[len(parts)-2]
        file := parts[len(parts)-1]
        if (dir != thisPackageName) || (file != thisFileName) {
            funcName := runtime.FuncForPC(pc).Name()
            return &CallerInfo{
                PathName: filepath,
                FileName: file,
                LineNo:   uint32(line),
                FuncName: funcName,
            }
        }
    }
}

func (self *StandardLogger) Handle(record *LogRecord) {
    if self.Filter(record) > 0 {
        self.callHandlers(record)
    }
}

func (self *StandardLogger) callHandlers(record *LogRecord) {
    var call Logger = self
    for call != nil {
        for _, handler := range call.GetHandlers() {
            if record.Level >= handler.GetLevel() {
                handler.Handle(record)
            }
        }
        if !call.GetPropagate() {
            call = nil
        } else {
            call = call.GetParent()
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

func (self *StandardLogger) GetHandlers() []Handler {
    self.lock.RLock()
    defer self.lock.RUnlock()
    result := make([]Handler, 0, self.handlers.Cardinality())
    for i := range self.handlers.Iter() {
        handler, _ := i.(Handler)
        result = append(result, handler)
    }
    return result
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

func (self *StandardLogger) GetParent() Logger {
    self.lock.RLock()
    defer self.lock.RUnlock()
    return self.parent
}

func (self *StandardLogger) SetParent(parent Logger) {
    self.lock.Lock()
    defer self.lock.Unlock()
    self.parent = parent
}

type RootLogger struct {
    *StandardLogger
}

func NewRootLogger(level LogLevelType) *RootLogger {
    logger := NewStandardLogger("root", level)
    logger.SetPropagate(false)
    return &RootLogger{
        StandardLogger: logger,
    }
}

type LoggerMaker func(name string) Logger

func defaultLoggerMaker(name string) Logger {
    return NewStandardLogger(name, LevelNotset)
}

type Manager struct {
    root        Logger
    loggers     map[string]Node
    loggerMaker LoggerMaker
    lock        sync.Mutex
}

func NewManager(logger Logger) *Manager {
    return &Manager{
        root:        logger,
        loggers:     make(map[string]Node),
        loggerMaker: defaultLoggerMaker,
    }
}

func (self *Manager) SetLoggerMaker(maker LoggerMaker) {
    self.lock.Lock()
    defer self.lock.Unlock()
    self.loggerMaker = maker
}

func (self *Manager) GetLogger(name string) Logger {
    self.lock.Lock()
    defer self.lock.Unlock()
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

func (self *Manager) fixupParents(logger Logger) {
    name := logger.GetName()
    index := strings.LastIndex(name, ".")
    var parent Logger
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

func (self *Manager) fixupChildren(placeHolder *PlaceHolder, logger Logger) {
    name := logger.GetName()
    for i := range placeHolder.Loggers.Iter() {
        l, _ := i.(Logger)
        parent := l.GetParent()
        if !strings.HasPrefix(parent.GetName(), name) {
            logger.SetParent(parent)
            l.SetParent(logger)
        }
    }
}
