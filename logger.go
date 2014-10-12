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

// Node type definition for the placeholder/logger tree in manager.
type NodeType uint8

const (
    NodeUnknown NodeType = 0 + iota
    NodeLogger
    NodePlaceHolder
)

// Node interface for the placeHolder/logger tree in manager.
type Node interface {
    // Return the node type.
    Type() NodeType
}

// Placeholder instances are used in the manager logger hierarchy to take
// the place of nodes for which no loggers have been defined. This class is
// intended for internal use only and not as part of the public API.
type PlaceHolder struct {
    Loggers mapset.Set
}

// Initialize a PlaceHolder with the specified logger being a child of
// this PlaceHolder.
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

// Add the specified logger as a child of this PlaceHolder.
func (self *PlaceHolder) Append(logger Logger) {
    if !self.Loggers.Contains(logger) {
        self.Loggers.Add(logger)
    }
}

// An interface represents a single logging channel.
// A "logging channel" indicates an area of an application. Exactly how an
// "area" is defined is up to the application developer.
// Since an application can have any number of areas, logging channels are
// identified by a unique string. Application areas can be nested (e.g. an area
// of "input processing" might include sub-areas "read CSV files",
// "read XLS files" and "read Gnumberic files"). To cater for this natural
// nesting, channel names are organized into a namespace hierarchy where levels
// are separated by periods, much like the Java or Python package namespace.
// So in the instance given above, channel names might be "input" for the upper
// level, an "input.csv", "input.xls" and "input.gnu" for the sub-levels.
// There is no arbitrary limit to the depth of nesting.
type Logger interface {
    // A Logger is a node in Manager tree.
    Node

    // Return the name of Logger.
    GetName() string
    // Return the propagate of Logger.
    GetPropagate() bool
    // Set the propagate.
    SetPropagate(v bool)
    // Return the logging level attached to this Logger.
    GetLevel() LogLevelType
    // Set the logging level attached to this Logger.
    SetLevel(level LogLevelType) error
    // Query whether this Logger is enabled for specified logging level.
    IsEnabledFor(level LogLevelType) bool
    // Get the effective level for this Logger.
    // An effective level is the first level value of Logger and its all parent
    // in the Logger hierarchy, which is not equal to LevelNotset.
    GetEffectiveLevel() LogLevelType

    // Log a message with severity "LevelFatal".
    Fatal(format string, args ...interface{})
    // Log a message with severity "LevelError".
    Error(format string, args ...interface{})
    // Log a message with severity "LevelWarn".
    Warn(format string, args ...interface{})
    // Log a message with severity "LevelInfo".
    Info(format string, args ...interface{})
    // Log a message with severity "LevelDebug".
    Debug(format string, args ...interface{})
    // Log a message with specified severity level.
    Log(level LogLevelType, format string, args ...interface{})

    // Add the specified handler to this Logger.
    AddHandler(handler Handler)
    // Remove the specified handler from this Logger.
    RemoveHandler(handler Handler)
    // Return all handler of this Logger.
    GetHandlers() []Handler

    // Filterer
    Filterer

    // Return the Manager of this Logger.
    GetManager() *Manager
    // Set the Manager of this Logger.
    SetManager(manager *Manager)
    // Return the parent Logger of this Logger.
    GetParent() Logger
    // Set the parent Logger of this Logger.
    SetParent(parent Logger)
}

// The standard logger implementation class.
type StandardLogger struct {
    *StandardFilterer
    name      string
    level     LogLevelType
    parent    Logger
    propagate bool
    handlers  mapset.Set
    manager   *Manager
    lock      sync.RWMutex
}

// Initialize a standard logger instance with name and logging level.
func NewStandardLogger(name string, level LogLevelType) *StandardLogger {
    return &StandardLogger{
        StandardFilterer: NewStandardFilterer(),
        parent:           nil,
        name:             name,
        level:            level,
        propagate:        true,
        handlers:         mapset.NewSet(),
        manager:          nil,
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

// Get the effective level for this logger.
// Loop through this logger and its parents in the logger hierarchy,
// looking for a non-zero logging level. Return the first one found.
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

// The informations of caller of this module.
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

// Find the stack frame of the caller so that we can note the source file name,
// line number and function name.
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

// Pass a record to all relevant handlers.
// Loop through all handlers for this logger and its parents in the logger
// hierarchy. Stop searching up the hierarchy whenever a logger with the
// "propagate" attribute set to zero is found - that will be the last
// logger whose handlers are called.
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

// Get a logger which is descendant to this one.
// This is a convenience method, such that
//     logging.GetLogger("abc").GetChild("def.ghi")
// is the same as
//     logging.GetLogger("abc.def.ghi")
// It's useful, for example, when the parent logger is named using
// some string unknown or random.
func (self *StandardLogger) GetChild(suffix string) Logger {
    self.lock.RLock()
    defer self.lock.Unlock()
    fullname := strings.Join([]string{self.name, suffix}, ".")
    return self.manager.GetLogger(fullname)
}

// A root logger is not that different to any other logger, except that
// it must have a logging level and there is only one instance of it in
// the hierarchy.
type RootLogger struct {
    *StandardLogger
}

// Initialize the root logger with the name "root".
func NewRootLogger(level LogLevelType) *RootLogger {
    logger := NewStandardLogger("root", level)
    logger.SetPropagate(false)
    return &RootLogger{
        StandardLogger: logger,
    }
}

// The logger maker function type.
type LoggerMaker func(name string) Logger

// The default logger maker for this module.
func defaultLoggerMaker(name string) Logger {
    return NewStandardLogger(name, LevelNotset)
}

// This is [under normal circumstances] just one manager instance, which
// holds the hierarchy of loggers.
type Manager struct {
    root        Logger
    loggers     map[string]Node
    loggerMaker LoggerMaker
    lock        sync.Mutex
}

// Initialize the manager with the root node of the logger hierarchy.
func NewManager(logger Logger) *Manager {
    return &Manager{
        root:        logger,
        loggers:     make(map[string]Node),
        loggerMaker: defaultLoggerMaker,
    }
}

// Set the logger maker to be used when instantiating
// a logger with this manager.
func (self *Manager) SetLoggerMaker(maker LoggerMaker) {
    self.lock.Lock()
    defer self.lock.Unlock()
    self.loggerMaker = maker
}

// Get a logger with the specified name (channel name), creating it
// if it doesn't yet exists. This name is a dot-separated hierarchical
// name, such as "a", "a.b", "a.b.c" or similar.
//
// If a placeholder existed for the specified name [i.e. the logger didn't
// exist but a child of it did], replace it with the created logger and fix up
// the parent/child references which pointed to the placeholder to now point
// to the logger.
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

// Ensure that there are either loggers or placeholders all the way from
// the specified logger to the root of the logger hierarchy.
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

// Ensure that children of the PlaceHolder placeHolder are connected to the
// specified logger.
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
