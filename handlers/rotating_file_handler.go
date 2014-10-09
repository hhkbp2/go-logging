package handlers

import (
    "github.com/hhkbp2/go-logging"
)

type RotatingHandler interface {
    ShouldRollover(record *LogRecord) bool
    DoRollover()
}

type BaseRotatingHandler struct {
    *FileHandler
}

func NewBaseRotatingHandler(filepath string, mode int) *BaseRotatingHandler {
    return &BaseRotatingHandler{}
}

func (self *BaseRotatingHandler) Emit(record *LogRecord) {
    if self.ShouldRollover(record) {
        self.DoRollover()
    }
    if err := self.FileHandler.Emit(record); err != nil {
        self.HandleError(record, err)
    }
}

type RotatingFileHandler struct {
    *BaseRotatingHandler
}

func NewRotatingFileHandler(
    filepath string,
    mode int,
    maxByte uint64,
    backupCount uint32) *RotatingFileHandler {

    // TODO add impl
    return &RotatingFileHandler{}
}

func (self *RotatingFileHandler) ShouldRollover(record *LogRecord) bool {
    // TODO add impl
}

func (self *RotatingFileHandler) DoRollover() {
    // TODO add impl
}
