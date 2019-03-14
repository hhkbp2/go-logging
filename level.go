package logging

import (
	"fmt"
	"sync"
)

// Type definition for log level
type LogLevelType uint8

// Default levels and level names, these can be replaced with any positive set
// of values having corresponding names. There is a pseudo-level, NOTSET, which
// is only really there as a lower limit for user-defined levels. Handlers and
// loggers are initialized with NOTSET so that they will log all messages, even
// at user-defined levels.
const (
	LevelCritical LogLevelType = 50
	LevelFatal    LogLevelType = LevelCritical
	LevelError    LogLevelType = 40
	LevelWarning  LogLevelType = 30
	LevelWarn     LogLevelType = LevelWarning
	LevelInfo     LogLevelType = 20
	LevelDebug    LogLevelType = 10
	LevelTrace    LogLevelType = 5
	LevelNotset   LogLevelType = 0
)

var (
	levelToNames = map[LogLevelType]string{
		LevelFatal:  "FATAL",
		LevelError:  "ERROR",
		LevelWarn:   "WARN",
		LevelInfo:   "INFO",
		LevelDebug:  "DEBUG",
		LevelTrace:  "TRACE",
		LevelNotset: "NOTSET",
	}
	nameToLevels = map[string]LogLevelType{
		"FATAL":  LevelFatal,
		"ERROR":  LevelError,
		"WARN":   LevelWarn,
		"INFO":   LevelInfo,
		"DEBUG":  LevelDebug,
		"TRACE":  LevelTrace,
		"NOTSET": LevelNotset,
	}
	levelLock sync.RWMutex
)

// Print the name of corresponding log level.
func (level LogLevelType) String() string {
	return GetLevelName(level)
}

// Return the textual representation of specified logging level.
// If the level is not registered by calling AddLevel(), ok would be false.
func getLevelName(level LogLevelType) (name string, ok bool) {
	levelLock.RLock()
	defer levelLock.RUnlock()
	levelName, ok := levelToNames[level]
	return levelName, ok
}

// GetLevelName returns the textual representation of specified logging level.
// If the level is one of the predefined levels (LevelFatal, LevelError,
// LevelWarn, LevelInfo, LevelDebug) then you get the corresponding string.
// If you have registered level with name using AddLevel() then the name you
// associated with level is returned.
// Otherwise, the string "Level %d"(%d is level value) is returned.
func GetLevelName(level LogLevelType) (name string) {
	name, ok := getLevelName(level)
	if !ok {
		return fmt.Sprintf("Level %d", uint8(level))
	}
	return name
}

// Associate levelName with level.
// This is used when converting levels to test during message formatting.
func AddLevel(level LogLevelType, levelName string) {
	levelLock.Lock()
	defer levelLock.Unlock()
	levelToNames[level] = levelName
	nameToLevels[levelName] = level
}
