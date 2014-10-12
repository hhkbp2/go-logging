package handlers

import (
    "github.com/hhkbp2/go-logging"
)

// This handler does nothing. It's intended to be used to avoid the
// "No handlers could be found for logger XXX" one-off warning. This is
// important for library code, which may contain code to log events.
// If a user of the library does not configure logging, the one-off warning
// might be produced; to avoid this, the library developer simply needs to
// instantiate a NullHandler and add it to the top-level logger of the library
// module or package.
type NullHandler struct {
    *logging.BaseHandler
}

// Initialize a NullHandler.
func NewNullHandler() *NullHandler {
    object := &NullHandler{
        BaseHandler: logging.NewBaseHandler("", logging.LevelNotset),
    }
    logging.Closer.AddHandler(object)
    return object
}

func (self *NullHandler) Emit(_ *logging.LogRecord) error {
    // Do nothing
    return nil
}

func (self *NullHandler) Handle(_ *logging.LogRecord) int {
    // Do nothing
    return 0
}
