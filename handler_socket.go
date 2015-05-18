package logging

import (
	"net"
)

type SocketHandler struct {
	*BaseHandler
	host         string
	port         uint16
	closeOnError bool
	retryTime    int
	conn         net.Conn
}

func NewSocketHandler(host string, port uint16) *SocketHandler {
	object := &SocketHandler{
		BaseHandler:  NewBaseHandler("", LevelNotset),
		host:         host,
		port:         port,
		closeOnError: false,
	}
	Closer.AddHandler(object)
	return object
}

func (self *SocketHandler) Serialize(record *LogRecord) []byte {
	// TODO serialize record using gob
	return make([]byte, 0)
}

func (self *SocketHandler) Send(bin []byte) error {
	// TODO send the bin through socket
	return nil
}

func (self *SocketHandler) Emit(record *LogRecord) error {
	bin := self.Serialize(record)
	return self.Send(bin)
}

func (self *SocketHandler) Handle(record *LogRecord) int {
	return self.Handle2(self, record)
}

func (self *SocketHandler) HandleError(record *LogRecord, err error) {

}

func (self *SocketHandler) Close() {
	// TODO add impl
}
