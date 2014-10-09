package handlers

import (
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
    file, err := os.OpenFile(filepath, O_WRONLY|O_CREATE|mode, 0666)
    if err != nil {
        return nil, err
    }
    stream := &FileStream{
        File: file,
    }
    return &FileHandler{
        StreamHandler: NewStreamHandler(stream),
        filepath:      filepath,
        mode:          mode,
    }
}

func (self *FileHandler) Emit(record *LogRecord) {
    self.StreamHandler.Emit(record)
}

func (self *FileHandler) Close() {
    self.Acquire()
    defer self.Release()
    self.Flush()
    self.StreamHandler.Close()
}
