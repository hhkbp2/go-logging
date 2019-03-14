package logging

// An interface to stream abstraction.
type Stream interface {
	// Report the current offset in the stream.
	Tell() (offset int64, err error)
	// Write a string into the stream.
	Write(s string) error
	// Flush the stream.
	Flush() error
	// Close the stream.
	Close() error
}

// A handler class with writes logging records, appropriately formatted,
// to a stream. Note that this class doesn't close the stream, as os.Stdin or
// os.Stdout my be used. However a Close2() method is there for subclass.
type StreamHandler struct {
	*BaseHandler
	stream Stream
}

// Initialize a stream handler with name, logging level and underlying stream.
func NewStreamHandler(
	name string, level LogLevelType, stream Stream) *StreamHandler {

	object := &StreamHandler{
		BaseHandler: NewBaseHandler(name, level),
		stream:      stream,
	}
	Closer.AddHandler(object)
	return object
}

// GetStream returns the underlying stream.
func (self *StreamHandler) GetStream() Stream {
	return self.stream
}

// Set the underlying stream.
func (self *StreamHandler) SetStream(s Stream) {
	self.stream = s
}

// Emit a record.
func (self *StreamHandler) Emit(record *LogRecord) error {
	return self.Emit2(self, record)
}

// A helper function to emit a record.
// If a formatter is specified, it is used to format the record.
// The record is then written to the stream with a trailing newline.
func (self *StreamHandler) Emit2(
	handler Handler, record *LogRecord) error {

	message := handler.Format(record)
	if err := self.stream.Write(message); err != nil {
		return err
	}
	return nil
}

// Handle() function is for the usage of stream handler on its own.
func (self *StreamHandler) Handle(record *LogRecord) int {
	return self.BaseHandler.Handle2(self, record)
}

// Flush the stream.
func (self *StreamHandler) Flush() error {
	return self.stream.Flush()
}

// A helper function for subclass implementation to close stream.
func (self *StreamHandler) Close2() {
	self.Flush()
	self.stream.Close()
}
