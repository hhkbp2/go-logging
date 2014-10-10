package logging

import (
    "sync"
)

type LogLevelType uint8

const (
    LevelCritical LogLevelType = 50
    LevelFatal                 = LevelCritical
    LevelError                 = 40
    LevelWarning               = 30
    LevelWarn                  = LevelWarning
    LevelInfo                  = 20
    LevelDebug                 = 10
    LevelNotset                = 0
)

var (
    levelToNames = map[LogLevelType]string{
        LevelFatal:  "FATAL",
        LevelError:  "ERROR",
        LevelWarn:   "WARN",
        LevelInfo:   "INFO",
        LevelDebug:  "DEBUG",
        LevelNotset: "NOTSET",
    }
    nameToLevels = map[string]LogLevelType{
        "FATAL":  LevelFatal,
        "ERROR":  LevelError,
        "WARN":   LevelWarn,
        "INFO":   LevelInfo,
        "DEBUG":  LevelDebug,
        "NOTSET": LevelNotset,
    }
    levelLock sync.RWMutex
)

func (level LogLevelType) String() string {
    return GetLevelName(level)
}

func getLevelName(level LogLevelType) (name string, ok bool) {
    levelLock.RLock()
    defer levelLock.RUnlock()
    levelName, ok := levelToNames[level]
    return levelName, ok
}

func GetLevelName(level LogLevelType) (name string) {
    name, ok := getLevelName(level)
    if !ok {
        return "UNKNOWN"
    }
    return name
}

func AddLevel(level LogLevelType, levelName string) {
    levelLock.Lock()
    defer levelLock.Unlock()
    levelToNames[level] = levelName
    nameToLevels[levelName] = level
}
