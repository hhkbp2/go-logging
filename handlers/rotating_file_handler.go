package handlers

import (
    "fmt"
    "github.com/hhkbp2/go-logging"
    "os"
)

// Check whether the specified directory/file exists or not.
func FileExists(filename string) bool {
    _, err := os.Stat(filename)
    if err != nil {
        if os.IsNotExist(err) {
            return false
        }
    }
    return true
}

// An interface for rotating handler abstraction.
type RotatingHandler interface {
    logging.Handler
    // Determine if rollover should occur.
    ShouldRollover(record *logging.LogRecord) bool
    // Do a rollover.
    DoRollover() error
}

// Base class for handlers that rotate log files at certain point.
// Not meant to be instantiated directly. Insteed, use RotatingFileHandler
// or TimedRotatingFileHandler.
type BaseRotatingHandler struct {
    *FileHandler
}

// Initialize base rotating handler with specified filename for stream logging.
func NewBaseRotatingHandler(
    filepath string, mode int) (*BaseRotatingHandler, error) {

    fileHandler, err := NewFileHandler(filepath, mode)
    if err != nil {
        return nil, err
    }
    object := &BaseRotatingHandler{
        FileHandler: fileHandler,
    }
    logging.Closer.RemoveHandler(object.FileHandler)
    logging.Closer.AddHandler(object)
    return object, nil
}

// A helper function for subclass to emit record.
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

// Handler for logging to a set of files, which switches from one file to
// the next when the current file reaches a certain size.
type RotatingFileHandler struct {
    *BaseRotatingHandler
    maxByte     uint64
    backupCount uint32
}

// Open the specified file and use it as the stream for logging.
//
// By default, the file grows indefinitely. You can specify particular values
// of maxByte and backupCount to allow the file to rollover at a predetermined
// size.
//
// Rollover occurs whenever the current log file is nearly maxByte in length.
// If backupCount is >= 1, the system will successively create new files with
// the same pathname as the base file, but with extensions ".1", ".2" etc.
// append to it. For example, with a backupCount of 5 and a base file name of
// "app.log", you would get "app.log", "app.log.1", "app.log.2", ...
// through to "app.log.5". The file being written to is always "app.log" - when
// it gets filled up, it is closed and renamed to "app.log.1", and if files
// "app.log.1", "app.log.2" etc. exist, then they are renamed to "app.log.2",
// "app.log.3" etc. respectively.
//
// If maxByte is zero, rollover never occurs.
func NewRotatingFileHandler(
    filepath string,
    mode int,
    maxByte uint64,
    backupCount uint32) (*RotatingFileHandler, error) {

    // If rotation/rollover is wanted, it doesn't make sense to use another
    // mode. If for example 'w' were specified, then if there were multiple
    // runs of the calling application, the logs from previous runs would be
    // lost if the "os.O_TRUNC" is respected, because the log file would be
    // truncated on each run.
    if maxByte > 0 {
        mode = os.O_APPEND
    }
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

func MustNewRotatingFileHandler(
    filepath string,
    mode int,
    maxByte uint64,
    backupCount uint32) *RotatingFileHandler {

    handler, err := NewRotatingFileHandler(
        filepath, mode, maxByte, backupCount)
    if err != nil {
        panic("NewRotatingFileHandler(), error: " + err.Error())
    }
    return handler
}

// Determine if rollover should occur.
// Basically, see if the supplied record would cause the file to exceed the
// size limit we have.
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

// Rotate source file to destination file if source file exists.
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

// Do a rollover, as described above.
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
