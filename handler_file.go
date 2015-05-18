package logging

import (
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
	// NOTICE: For performance optimization, don't sync inner buffer to
	// disk everytime it writes something.
	// return self.File.Sync()
	return nil
}

func (self *FileStream) Close() error {
	self.File.Sync()
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
	handler := NewStreamHandler(filepath, LevelNotset, nil)
	object := &FileHandler{
		StreamHandler: handler,
		filepath:      filepath,
		mode:          mode,
	}
	if err = object.Open(); err != nil {
		return nil, err
	}
	Closer.RemoveHandler(object.StreamHandler)
	Closer.AddHandler(object)
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
	var file *os.File
	var err error
	for {
		file, err = os.OpenFile(
			self.filepath, os.O_WRONLY|os.O_CREATE|self.mode, 0666)
		if err == nil {
			break
		}
		// try to create all the parent directories for specified log file
		// if it doesn't exist
		if os.IsNotExist(err) {
			err2 := os.MkdirAll(filepath.Dir(self.filepath), 0755)
			if err2 != nil {
				return err
			}
			continue
		}
		return err
	}
	stream := &FileStream{
		File: file,
	}
	self.StreamHandler.SetStream(stream)
	return nil
}

// Emit a record.
func (self *FileHandler) Emit(record *LogRecord) error {
	return self.StreamHandler.Emit2(self, record)
}

func (self *FileHandler) Handle(record *LogRecord) int {
	return self.Handle2(self, record)
}

// Close this file handler.
func (self *FileHandler) Close() {
	self.StreamHandler.Close2()
}
