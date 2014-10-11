package handlers

import (
    "fmt"
    "github.com/hhkbp2/go-logging"
    "os"
)

func FileExists(filename string) bool {
    _, err := os.Stat(filename)
    if err != nil {
        if os.IsNotExist(err) {
            return false
        }
    }
    return true
}

type RotatingHandler interface {
    logging.Handler
    ShouldRollover(record *logging.LogRecord) bool
    DoRollover() error
}

type BaseRotatingHandler struct {
    *FileHandler
}

func NewBaseRotatingHandler(
    filepath string, mode int) (*BaseRotatingHandler, error) {

    fileHandler, err := NewFileHandler(filepath, mode)
    if err != nil {
        return nil, err
    }
    object := &BaseRotatingHandler{
        FileHandler: fileHandler,
    }
    return object, nil
}

func (self *BaseRotatingHandler) RolloverEmit(
    handler RotatingHandler, record *logging.LogRecord) error {

    if handler.ShouldRollover(record) {
        if err := handler.DoRollover(); err != nil {
            return err
        }
    }
    if err := self.Emit(record); err != nil {
        return err
    }
    return nil
}

type RotatingFileHandler struct {
    *BaseRotatingHandler
    maxByte     uint64
    backupCount uint32
}

func NewRotatingFileHandler(
    filepath string,
    mode int,
    maxByte uint64,
    backupCount uint32) (*RotatingFileHandler, error) {

    base, err := NewBaseRotatingHandler(filepath, mode)
    if err != nil {
        return nil, err
    }
    object := &RotatingFileHandler{
        BaseRotatingHandler: base,
        maxByte:             maxByte,
        backupCount:         backupCount,
    }
    return object, nil
}

func (self *RotatingFileHandler) ShouldRollover(
    record *logging.LogRecord) bool {

    if self.maxByte > 0 {
        message := fmt.Sprintf("%s\n", self.Format(record))
        offset, err := self.GetStream().Tell()
        if err != nil {
            // don't trigger rollover action if we lose offset info
            return false
        }
        if (uint64(offset) + uint64(len(message))) > self.maxByte {
            return true
        }
    }
    return false
}

func (self *RotatingFileHandler) RotateFile(sourceFile, destFile string) error {
    if FileExists(sourceFile) {
        if FileExists(destFile) {
            if err := os.Remove(destFile); err != nil {
                return err
            }
        }
        if err := os.Rename(sourceFile, destFile); err != nil {
            return err
        }
    }
    return nil
}

func (self *RotatingFileHandler) DoRollover() (err error) {
    self.Close()
    defer func() {
        if e := self.Open(); e != nil {
            if err == nil {
                err = e
            }
        }
    }()
    if self.backupCount > 0 {
        filepath := self.GetFilePath()
        for i := self.backupCount - 1; i > 0; i-- {
            sourceFile := fmt.Sprintf("%s.%d", filepath, i)
            destFile := fmt.Sprintf("%s.%d", filepath, i+1)
            if err := self.RotateFile(sourceFile, destFile); err != nil {
                return err
            }
        }
        destFile := fmt.Sprintf("%s.%d", filepath, 1)
        if err := self.RotateFile(filepath, destFile); err != nil {
            return err
        }
    }
    return nil
}

func (self *RotatingFileHandler) Emit(record *logging.LogRecord) error {
    return self.RolloverEmit(self, record)
}

func (self *RotatingFileHandler) Handle(record *logging.LogRecord) int {
    return self.Handle2(self, record)
}
