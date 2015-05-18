package logging

// This handler does nothing. It's intended to be used to avoid the
// "No handlers could be found for logger XXX" one-off warning. This is
// important for library code, which may contain code to log events.
// If a user of the library does not configure logging, the one-off warning
// might be produced; to avoid this, the library developer simply needs to
// instantiate a NullHandler and add it to the top-level logger of the library
// module or package.
type NullHandler struct {
	*BaseHandler
}

// Initialize a NullHandler.
func NewNullHandler() *NullHandler {
	object := &NullHandler{
		BaseHandler: NewBaseHandler("", LevelNotset),
	}
	Closer.AddHandler(object)
	return object
}

func (self *NullHandler) Emit(_ *LogRecord) error {
	// Do nothing
	return nil
}

func (self *NullHandler) Handle(_ *LogRecord) int {
	// Do nothing
	return 0
}
