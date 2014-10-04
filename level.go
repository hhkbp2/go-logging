package logging

import (
    "sync"
)

type LogLevelType uint8

const (
    Critical LogLevelType = 50
    Fatal                 = Critical
    Error                 = 40
    Warning               = 30
    Warn                  = Warning
    Info                  = 20
    Debug                 = 10
    Notset                = 0
)

var (
    levelToNames = map[LogLevelType]string{
        Fatal:  "FATAL",
        Error:  "ERROR",
        Warn:   "WARN",
        Info:   "INFO",
        Debug:  "DEBUG",
        Notset: "NOTSET",
    }
    nameToLevels = map[string]LogLevelType{
        "FATAL":  Fatal,
        "ERROR":  Error,
        "WARN":   Warn,
        "INFO":   Info,
        "DEBUG":  Debug,
        "NOTSET": Notset,
    }
    levelLock sync.RWMutex
)

func getLevelName(level LogLevelType) (name string, ok bool) {
    levelLock.RLock()
    defer levelLock.RUnlock()
    return levelToNames[level]
}

func addLevel(level LogLevelType, levelName string) error {
    levelLock.Lock()
    defer levelLock.Unlock()
    levelToNames[level] = levelName
    nameToLevels[levelName] = level
}
