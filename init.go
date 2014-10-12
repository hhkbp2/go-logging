package logging

var (
    root    = NewRootLogger(LevelWarn)
    manager = NewManager(root)
)

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

func Shutdown() {
    // TODO shutdown every logger in manager
}
