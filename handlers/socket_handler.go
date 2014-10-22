package handlers

import (
	"github.com/hhkbp2/go-logging"
	"net"
)

type SocketHandler struct {
	*logging.BaseHandler
	host         string
	port         uint16
	closeOnError bool
	retryTime    int
	conn         net.Conn
}

func NewSocketHandler(host string, port uint16) *SocketHandler {
	object := &SocketHandler{
		BaseHandler:  logging.NewBaseHandler("", logging.LevelNotset),
		host:         host,
		port:         port,
		closeOnError: false,
	}
	logging.Closer.AddHandler(object)
	return object
}

func (self *SocketHandler) Serialize(record *logging.LogRecord) []byte {
	// TODO serialize record using gob
	return make([]byte, 0)
}

func (self *SocketHandler) Send(bin []byte) error {
	// TODO send the bin through socket
	return nil
}

func (self *SocketHandler) Emit(record *logging.LogRecord) error {
	bin := self.Serialize(record)
	return self.Send(bin)
}

func (self *SocketHandler) Handle(record *logging.LogRecord) int {
	return self.Handle2(self, record)
}

func (self *SocketHandler) HandleError(record *logging.LogRecord, err error) {

}

func (self *SocketHandler) Close() {
	// TODO add impl
}
