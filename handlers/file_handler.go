package handlers

import (
    "github.com/hhkbp2/go-logging"
    "os"
    "path/filepath"
)

type FileStream struct {
    File *os.File
}

func (self *FileStream) Write(s string) error {
    length := len(s)
    for i := 0; i < length; {
        n, err := self.File.Write([]byte(s))
        if err != nil {
            return err
        }
        i += n
    }
    return nil
}

func (self *FileStream) Flush() error {
    // TODO temporarily sync to disk
    return self.File.Sync()
}

type FileHandler struct {
    *StreamHandler

    filepath string
    mode     int
}

func NewFileHandler(filename string, mode int) (*FileHandler, error) {
    filepath, err := filepath.Abs(filename)
    if err != nil {
        return nil, err
    }
    file, err := os.OpenFile(filepath, os.O_WRONLY|os.O_CREATE|mode, 0666)
    if err != nil {
        return nil, err
    }
    stream := &FileStream{
        File: file,
    }
    object := &FileHandler{
        StreamHandler: NewStreamHandler(filename, logging.LevelNotset, stream),
        filepath:      filepath,
        mode:          mode,
    }
    return object, nil
}

func (self *FileHandler) Emit(record *logging.LogRecord) error {
    return self.StreamHandler.Emit2(self, record)
}

func (self *FileHandler) Handle(record *logging.LogRecord) int {
    return self.Handle2(self, record)
}

func (self *FileHandler) Close() {
    self.Lock()
    defer self.Unlock()
    self.Flush()
    self.StreamHandler.Close()
}
