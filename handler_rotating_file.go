package logging

import (
	"fmt"
	"os"
	"sync"
	"time"
)

// FileExists checks whether the specified directory/file exists or not.
func FileExists(filename string) bool {
	if _, err := os.Stat(filename); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

// An interface for rotating handler abstraction.
type RotatingHandler interface {
	Handler
	// Determine if rollover should occur.
	ShouldRollover(record *LogRecord) (doRollover bool, message string)
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
	filepath string, mode, bufferSize int) (*BaseRotatingHandler, error) {

	fileHandler, err := NewFileHandler(filepath, mode, bufferSize)
	if err != nil {
		return nil, err
	}
	object := &BaseRotatingHandler{
		FileHandler: fileHandler,
	}
	Closer.RemoveHandler(object.FileHandler)
	Closer.AddHandler(object)
	return object, nil
}

// A helper function for subclass to emit record.
func (self *BaseRotatingHandler) RolloverEmit(
	handler RotatingHandler, record *LogRecord) error {

	// We don't use the implementation of StreamHandler.Emit2() but directly
	// write to stream here in order to avoid calling self.Format() twice
	// for performance optimization.
	doRollover, message := handler.ShouldRollover(record)
	if doRollover {
		if err := handler.Flush(); err != nil {
			return err
		}
		if err := handler.DoRollover(); err != nil {
			return err
		}
	}
	// Message already has a trailing '\n'.
	err := self.GetStream().Write(message)
	if err != nil {
		return err
	}
	return nil
}

type HandleFunc func(record *LogRecord) int

// Handler for logging to a set of files, which switches from one file to
// the next when the current file reaches a certain size.
type RotatingFileHandler struct {
	*BaseRotatingHandler
	maxBytes        uint64
	backupCount     uint32
	bufferFlushTime time.Duration
	inputChanSize   int
	handleFunc      HandleFunc
	inputChan       chan *LogRecord
	group           *sync.WaitGroup
}

// Open the specified file and use it as the stream for logging.
//
// By default, the file grows indefinitely. You can specify particular values
// of maxBytes and backupCount to allow the file to rollover at a predetermined
// size.
//
// Rollover occurs whenever the current log file is nearly maxBytes in length.
// If backupCount is >= 1, the system will successively create new files with
// the same pathname as the base file, but with extensions ".1", ".2" etc.
// append to it. For example, with a backupCount of 5 and a base file name of
// "app.log", you would get "app.log", "app.log.1", "app.log.2", ...
// through to "app.log.5". The file being written to is always "app.log" - when
// it gets filled up, it is closed and renamed to "app.log.1", and if files
// "app.log.1", "app.log.2" etc. exist, then they are renamed to "app.log.2",
// "app.log.3" etc. respectively.
//
// If maxBytes is zero, rollover never occurs.
//
// bufferSize specifies the size of the internal buffer. If it is positive,
// the internal buffer will be enabled, the logs will be first written into
// the internal buffer, when the internal buffer is full all buffer content
// will be flushed to file.
// bufferFlushTime specifies the time for flushing the internal buffer
// in period, no matter the buffer is full or not.
// inputChanSize specifies the chan size of the handler. If it is positive,
// this handler will be initialized as a standardlone go routine to handle
// log message.
func NewRotatingFileHandler(
	filepath string,
	mode int,
	bufferSize int,
	bufferFlushTime time.Duration,
	inputChanSize int,
	maxBytes uint64,
	backupCount uint32) (*RotatingFileHandler, error) {

	// If rotation/rollover is wanted, it doesn't make sense to use another
	// mode. If for example 'w' were specified, then if there were multiple
	// runs of the calling application, the logs from previous runs would be
	// lost if the "os.O_TRUNC" is respected, because the log file would be
	// truncated on each run.
	if maxBytes > 0 {
		mode = os.O_APPEND
	}
	base, err := NewBaseRotatingHandler(filepath, mode, bufferSize)
	if err != nil {
		return nil, err
	}
	object := &RotatingFileHandler{
		BaseRotatingHandler: base,
		maxBytes:            maxBytes,
		backupCount:         backupCount,
		bufferFlushTime:     bufferFlushTime,
		inputChanSize:       inputChanSize,
	}
	// register object to closer
	Closer.RemoveHandler(object.BaseRotatingHandler)
	Closer.AddHandler(object)
	if inputChanSize > 0 {
		object.handleFunc = object.handleChan
		object.inputChan = make(chan *LogRecord, inputChanSize)
		object.group = &sync.WaitGroup{}
		object.group.Add(1)
		go func() {
			defer object.group.Done()
			object.loop()
		}()
	} else {
		object.handleFunc = object.handleCall
	}
	return object, nil
}

func MustNewRotatingFileHandler(
	filepath string,
	mode int,
	bufferSize int,
	bufferFlushTime time.Duration,
	inputChanSize int,
	maxBytes uint64,
	backupCount uint32) *RotatingFileHandler {

	handler, err := NewRotatingFileHandler(
		filepath,
		mode,
		bufferSize,
		bufferFlushTime,
		inputChanSize,
		maxBytes,
		backupCount)
	if err != nil {
		panic("NewRotatingFileHandler(), error: " + err.Error())
	}
	return handler
}

// Determine if rollover should occur.
// Basically, see if the supplied record would cause the file to exceed the
// size limit we have.
func (self *RotatingFileHandler) ShouldRollover(
	record *LogRecord) (bool, string) {

	message := self.Format(record)
	if self.maxBytes > 0 {
		offset, err := self.GetStream().Tell()
		if err != nil {
			// don't trigger rollover action if we lose offset info
			return false, message
		}
		if (uint64(offset) + uint64(len(message))) > self.maxBytes {
			return true, message
		}
	}
	return false, message
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
	self.FileHandler.Close()
	defer func() {
		if e := self.FileHandler.Open(); e != nil {
			if e == nil {
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

// Emit a record.
func (self *RotatingFileHandler) Emit(record *LogRecord) error {
	return self.RolloverEmit(self, record)
}

func (self *RotatingFileHandler) handleCall(record *LogRecord) int {
	return self.Handle2(self, record)
}

func (self *RotatingFileHandler) handleChan(record *LogRecord) int {
	self.inputChan <- record
	return 0
}

func (self *RotatingFileHandler) loop() {
	ticker := time.NewTicker(self.bufferFlushTime)
	for {
		select {
		case r := <-self.inputChan:
			if r == nil {
				return
			}

			self.Handle2(self, r)
		case <-ticker.C:
			self.Flush()
		}
	}
}

func (self *RotatingFileHandler) Handle(record *LogRecord) int {
	return self.handleFunc(record)
}

func (self *RotatingFileHandler) Close() {
	if self.inputChanSize > 0 {
		// send a nil record as "stop signal" to exit loop.
		self.inputChan <- nil
		self.group.Wait()
	}
	self.BaseRotatingHandler.Close()
}
