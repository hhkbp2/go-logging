package handlers

import (
    "github.com/hhkbp2/go-logging"
    "os"
    "path/filepath"
)

type FileStream struct {
    File *os.File
}

func (self *FileStream) Tell() (int64, error) {
    fileInfo, err := self.File.Stat()
    if err != nil {
        return 0, err
    }
    return fileInfo.Size(), nil
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

func (self *FileStream) Close() error {
    return self.File.Close()
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
    handler := NewStreamHandler(filepath, logging.LevelNotset, nil)
    object := &FileHandler{
        StreamHandler: handler,
        filepath:      filepath,
        mode:          mode,
    }
    if err = object.Open(); err != nil {
        return nil, err
    }
    return object, nil
}

func (self *FileHandler) GetFilePath() string {
    return self.filepath
}

func (self *FileHandler) Open() error {
    file, err := os.OpenFile(
        self.filepath, os.O_WRONLY|os.O_CREATE|self.mode, 0666)
    if err != nil {
        return err
    }
    stream := &FileStream{
        File: file,
    }
    self.StreamHandler.SetStream(stream)
    return nil
}

func (self *FileHandler) Emit(record *logging.LogRecord) error {
    return self.StreamHandler.Emit2(self, record)
}

func (self *FileHandler) Handle(record *logging.LogRecord) int {
    return self.Handle2(self, record)
}

func (self *FileHandler) Close() {
    self.StreamHandler.Close()
}
