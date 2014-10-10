package handlers

import (
    "github.com/hhkbp2/go-logging"
)

type NullHandler struct {
    *logging.BaseHandler
}

func NewNullHandler() *NullHandler {
    return &NullHandler{
        BaseHandler: logging.NewBaseHandler("", logging.LevelNotset),
    }
}

func (self *NullHandler) Emit(_ *logging.LogRecord) error {
    // Do nothing
    return nil
}

func (self *NullHandler) Handle(_ *logging.LogRecord) int {
    // Do nothing
    return 0
}
