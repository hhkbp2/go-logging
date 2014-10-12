package handlers

import (
    "github.com/hhkbp2/go-logging"
    "os"
    "path/filepath"
)

// A class wraps os.File to the stream interface.
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
    // TODO to be optimized for performance
    return self.File.Sync()
}

func (self *FileStream) Close() error {
    return self.File.Close()
}

// A handler class which writes formatted logging records to disk files.
type FileHandler struct {
    *StreamHandler

    filepath string
    mode     int
}

// Open the specified file and use it as the stream for logging.
func NewFileHandler(filename string, mode int) (*FileHandler, error) {
    // keep the absolute path, otherwise derived classes which use this
    // may come a cropper when the current directory changes.
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
    logging.Closer.RemoveHandler(object.StreamHandler)
    logging.Closer.AddHandler(object)
    return object, nil
}

// Return the absolute path of logging file.
func (self *FileHandler) GetFilePath() string {
    return self.filepath
}

// Open the current base file with the (original) mode and encoding,
// and set it to the underlying stream handler.
// Return non-nil error if error happens.
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

// Emit a record.
func (self *FileHandler) Emit(record *logging.LogRecord) error {
    return self.StreamHandler.Emit2(self, record)
}

func (self *FileHandler) Handle(record *logging.LogRecord) int {
    return self.Handle2(self, record)
}

// Close this file handler.
func (self *FileHandler) Close() {
    self.StreamHandler.Close2()
}
