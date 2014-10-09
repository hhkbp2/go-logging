package handlers

import (
    "github.com/hhkbp2/go-logging"
)

type NullHandler struct {
}

func (self *NullHandler) Handle(_ *logging.LogRecord) {
    // Do nothing
}

func (self *NullHandler) Emit(_ *logging.LogRecord) {
    // Do nothing
}
