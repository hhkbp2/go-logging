package logging

var (
    root    = NewRootLogger(LevelWarn)
    manager = NewManager(root)
)

func GetLogger(name string) Logger {
    if len(name) > 0 {
        return manager.GetLogger(name)
    } else {
        return root
    }
}

func Fatal(format string, args ...interface{}) {
    root.Fatal(format, args...)
}

func Error(format string, args ...interface{}) {
    root.Error(format, args...)
}

func Warn(format string, args ...interface{}) {
    root.Warn(format, args...)
}

func Info(format string, args ...interface{}) {
    root.Info(format, args...)
}

func Debug(format string, args ...interface{}) {
    root.Debug(format, args...)
}

func Log(level LogLevelType, format string, args ...interface{}) {
    root.Log(level, format, args...)
}

func Shutdown() {
    handlersLock()
    defer handlersUnlock()
    ShutdownHandlers(handlers)
}
func ShutdownHandlers(handlers map[string]Handler) {
    for _, handler := range handlers {
        handler.Lock()
        defer handler.Unlock()
        handler.Flush()
        handler.Close()
    }
}
