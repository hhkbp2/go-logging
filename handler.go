package logging

import (
    "github.com/deckarep/golang-set"
)

var (
    handlers = mapset.Set
)

type Handler interface {
    GetName() string
    SetName(name string)
    CreateLock()
    Acquire()
    Release()
    Format(record *LogRecord) string
    Emit(record *LogRecord)
    Handle(record *LogRecord)
    SetFormatter(formater Formatter)
    Flush()
    Close()
    HandleError(record *LogRecord)
}
