package logging

import (
	"container/list"
	"reflect"
)

func IsNil(i interface{}) bool {
	return (i == nil || reflect.ValueOf(i).IsNil())
}

func IsNotNil(i interface{}) bool {
	return !IsNil(i)
}

type BufferingHandler interface {
	Handler
	ShouldFlush(record *LogRecord) bool
}

type BaseBufferingHandler struct {
	*BaseHandler
	capacity uint64
	buffer   *list.List
}

func NewBaseBufferingHandler(capacity uint64) *BaseBufferingHandler {
	return &BaseBufferingHandler{
		BaseHandler: NewBaseHandler("", LevelNotset),
		capacity:    capacity,
		buffer:      list.New(),
	}
}

func (self *BaseBufferingHandler) GetBuffer() *list.List {
	return self.buffer
}

func (self *BaseBufferingHandler) ShouldFlush(_ *LogRecord) bool {
	return uint64(self.buffer.Len()) >= self.capacity
}

func (self *BaseBufferingHandler) Emit2(
	handler BufferingHandler, record *LogRecord) error {

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
	flushLevel LogLevelType
	target     Handler
}

func NewMemoryHandler(
	capacity uint64, flushLevel LogLevelType, target Handler) *MemoryHandler {

	object := &MemoryHandler{
		BaseBufferingHandler: NewBaseBufferingHandler(capacity),
		flushLevel:           flushLevel,
		target:               target,
	}
	Closer.AddHandler(object)
	return object
}

func (self *MemoryHandler) ShouldFlush(record *LogRecord) bool {
	return ((self.BaseBufferingHandler.ShouldFlush(record)) ||
		(record.Level >= self.flushLevel))
}

func (self *MemoryHandler) SetTarget(target Handler) {
	self.target = target
}

func (self *MemoryHandler) Emit(record *LogRecord) error {
	return self.BaseBufferingHandler.Emit2(self, record)
}

func (self *MemoryHandler) Handle(record *LogRecord) int {
	return self.Handle2(self, record)
}

func (self *MemoryHandler) Flush() error {
	if IsNotNil(self.target) {
		buffer := self.BaseBufferingHandler.GetBuffer()
		for e := buffer.Front(); e != nil; e = e.Next() {
			record, _ := e.Value.(*LogRecord)
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
