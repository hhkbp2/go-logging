package handlers

import (
    "container/list"
    "github.com/hhkbp2/go-logging"
    "reflect"
)

func IsNil(i interface{}) bool {
    return (i == nil || reflect.ValueOf(i).IsNil())
}

func IsNotNil(i interface{}) bool {
    return !IsNil(i)
}

type BufferingHandler interface {
    logging.Handler
    ShouldFlush(record *logging.LogRecord) bool
}

type BaseBufferingHandler struct {
    *logging.BaseHandler
    capacity uint64
    buffer   *list.List
}

func NewBaseBufferingHandler(capacity uint64) *BaseBufferingHandler {
    return &BaseBufferingHandler{
        BaseHandler: logging.NewBaseHandler("", logging.LevelNotset),
        capacity:    capacity,
        buffer:      list.New(),
    }
}

func (self *BaseBufferingHandler) GetBuffer() *list.List {
    return self.buffer
}

func (self *BaseBufferingHandler) ShouldFlush(_ *logging.LogRecord) bool {
    return uint64(self.buffer.Len()) >= self.capacity
}

func (self *BaseBufferingHandler) Emit2(
    handler BufferingHandler, record *logging.LogRecord) error {

    self.buffer.PushBack(record)
    if handler.ShouldFlush(record) {
        return handler.Flush()
    }
    return nil
}

func (self *BaseBufferingHandler) Flush() error {
    self.buffer.Init()
    return nil
}

func (self *BaseBufferingHandler) Close() {
    self.Flush()
}

type MemoryHandler struct {
    *BaseBufferingHandler
    flushLevel logging.LogLevelType
    target     logging.Handler
}

func NewMemoryHandler(
    capacity uint64,
    flushLevel logging.LogLevelType,
    target logging.Handler) *MemoryHandler {

    object := &MemoryHandler{
        BaseBufferingHandler: NewBaseBufferingHandler(capacity),
        flushLevel:           flushLevel,
        target:               target,
    }
    logging.Closer.AddHandler(object)
    return object
}

func (self *MemoryHandler) ShouldFlush(record *logging.LogRecord) bool {
    return ((self.BaseBufferingHandler.ShouldFlush(record)) ||
        (record.Level >= self.flushLevel))
}

func (self *MemoryHandler) SetTarget(target logging.Handler) {
    self.target = target
}

func (self *MemoryHandler) Emit(record *logging.LogRecord) error {
    return self.BaseBufferingHandler.Emit2(self, record)
}

func (self *MemoryHandler) Handle(record *logging.LogRecord) int {
    return self.Handle2(self, record)
}

func (self *MemoryHandler) Flush() error {
    if IsNotNil(self.target) {
        buffer := self.BaseBufferingHandler.GetBuffer()
        for e := buffer.Front(); e != nil; e = e.Next() {
            record, _ := e.Value.(*logging.LogRecord)
            self.target.Handle(record)
        }
        return self.BaseBufferingHandler.Flush()
    }
    return nil
}

func (self *MemoryHandler) Close() {
    self.Flush()
    self.BaseBufferingHandler.Close()
}
