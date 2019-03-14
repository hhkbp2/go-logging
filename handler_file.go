package logging

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
)

func ReadN(reader io.Reader, b []byte) (int, error) {
	length := len(b)
	i := 0
	for i < length {
		n, err := reader.Read(b[i:])
		i += n
		if err != nil {
			return i, err
		}
	}
	return i, nil
}

func WriteN(writer io.Writer, b []byte) (int, error) {
	length := len(b)
	i := 0
	for i < length {
		n, err := writer.Write(b[i:])
		i += n
		if err != nil {
			return i, err
		}
	}
	return i, nil
}

// A class wraps os.File to the stream interface.
type FileStream struct {
	File       *os.File
	BufferSize int
	Buffer     *bytes.Buffer
	Offset     int64
}

func NewFileStream(f *os.File, bufferSize int) *FileStream {
	var buf *bytes.Buffer
	if bufferSize > 0 {
		buf = bytes.NewBuffer(make([]byte, 0, bufferSize))
	}
	return &FileStream{
		File:       f,
		BufferSize: bufferSize,
		Buffer:     buf,
	}
}

func (self *FileStream) Tell() (int64, error) {
	if self.Offset == 0 {
		fileInfo, err := self.File.Stat()
		if err != nil {
			return 0, err
		}
		self.Offset = fileInfo.Size()
		if self.BufferSize > 0 {
			self.Offset += int64(self.Buffer.Len())
		}
	}
	return self.Offset, nil
}

func (self *FileStream) Write(s string) error {
	// NOTICE: For performance optimization, don't sync inner buffer to
	// disk everytime it writes something.
	if self.BufferSize > 0 {
		length := len(s)
		if self.Buffer.Len()+length > self.BufferSize {
			self.doFlushBuffer()
			if length > self.BufferSize {
				n, err := WriteN(self.File, []byte(s))
				self.Offset += int64(n)
				return err
			}
		}
		_, err := self.Buffer.Write([]byte(s))
		return err
	} else {
		n, err := WriteN(self.File, []byte(s))
		self.Offset += int64(n)
		return err
	}
}

func (self *FileStream) doFlushBuffer() error {
	if self.Buffer.Len() > 0 {
		n, err := WriteN(self.File, self.Buffer.Bytes())
		self.Offset += int64(n)
		if err != nil {
			return err
		}
		self.Buffer.Reset()
	}
	return nil
}

func (self *FileStream) Flush() error {
	if self.BufferSize > 0 {
		return self.doFlushBuffer()
	}
	return nil
}

func (self *FileStream) Close() error {
	self.Flush()
	return self.File.Close()
}

// A handler class which writes formatted logging records to disk files.
type FileHandler struct {
	*StreamHandler

	filepath   string
	mode       int
	bufferSize int
}

// Open the specified file and use it as the stream for logging.
func NewFileHandler(filename string, mode int, bufferSize int) (*FileHandler, error) {
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
		bufferSize:    bufferSize,
	}
	if err = object.Open(); err != nil {
		return nil, err
	}
	Closer.RemoveHandler(object.StreamHandler)
	Closer.AddHandler(object)
	return object, nil
}

// GetFilePath returns the absolute path of logging file.
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
	stream := NewFileStream(file, self.bufferSize)
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
