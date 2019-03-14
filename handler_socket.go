package logging

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"net"
	"time"
)

const (
	SocketDefaultTimeout     = 1 * time.Second
	SocketDefaultDelay       = 1 * time.Second
	SocketDefaultMaxDeadline = 30 * time.Second
)

// A SocketLogRecord instance contains all LogRecord fields tailored for
// uploading to socket server. We could keep the interested fields and
// remove all others to minimize the network bandwidth usage.
type SocketLogRecord struct {
	CreatedTime time.Time
	AscTime     string
	Name        string
	Level       LogLevelType
	PathName    string
	FileName    string
	LineNo      uint32
	FuncName    string
	Format      string
	UseFormat   bool
	Message     string
}

// A handler class which write logging records, in gob format, to
// a streaming socket. The socket is kept open across logging calls.
// If the peer resets it, an attempt is made to reconnect on the next call.
type SocketHandler struct {
	*BaseHandler
	host         string
	port         uint16
	closeOnError bool
	retry        Retry
	makeConnFunc func() error
	sendFunc     func(bin []byte) error
	conn         net.Conn
}

// Initializes the handler with a specific host address and port.
// The attribute 'closeOnError' is set to true by default, which means that
// if a socket error occurs, the socket is silently closed and then reopen
// on the next loggging call.
func NewSocketHandler(host string, port uint16) *SocketHandler {
	retry := NewErrorRetry().
		Delay(SocketDefaultDelay).
		Deadline(SocketDefaultMaxDeadline)
	object := &SocketHandler{
		BaseHandler:  NewBaseHandler("", LevelNotset),
		host:         host,
		port:         port,
		closeOnError: true,
		retry:        retry,
	}
	object.makeConnFunc = object.makeTCPSocket
	object.sendFunc = object.sendTCP
	Closer.AddHandler(object)
	return object
}

// Marshals the record in gob binary format and returns it ready for
// transmission across socket.
func (self *SocketHandler) Marshal(record *LogRecord) ([]byte, error) {
	r := SocketLogRecord{
		CreatedTime: record.CreatedTime,
		AscTime:     record.AscTime,
		Name:        record.Name,
		Level:       record.Level,
		PathName:    record.PathName,
		FileName:    record.FileName,
		LineNo:      record.LineNo,
		FuncName:    record.FuncName,
		Format:      record.Format,
		UseFormat:   record.UseFormat,
		Message:     record.Message,
	}
	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)
	err := encoder.Encode(r)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// A factory method which allows succlasses to define the precise type of
// socket they want.
func (self *SocketHandler) makeSocket(network string) error {
	address := fmt.Sprintf("%s:%d", self.host, self.port)
	conn, err := net.DialTimeout(network, address, SocketDefaultTimeout)
	if err != nil {
		return err
	}
	self.conn = conn
	return nil
}

func (self *SocketHandler) makeTCPSocket() error {
	return self.makeSocket("tcp")
}

func (self *SocketHandler) makeUDPSocket() error {
	return self.makeSocket("udp")
}

// Try to create a socket, using an exponential backoff with a deadline time.
func (self *SocketHandler) createSocket() error {
	return self.retry.Do(self.makeConnFunc)
}

// Send a marshaled binary to the tcp socket.
func (self *SocketHandler) sendTCP(bin []byte) error {
	if self.conn == nil {
		err := self.createSocket()
		if err != nil {
			return err
		}
	}
	sentSoFar, left := 0, len(bin)
	for left > 0 {
		sent, err := self.conn.Write(bin[sentSoFar:])
		if err != nil {
			return err
		}
		sentSoFar += sent
		left -= sent
	}
	return nil
}

// Send a marshaled binary to the udp socket.
func (self *SocketHandler) sendUDP(bin []byte) error {
	if self.conn == nil {
		err := self.createSocket()
		if err != nil {
			return err
		}
	}
	_, err := self.conn.Write(bin)
	return err
}

// Emit a record.
// Marshals the record and writes it to the socket in binary format.
// If there is an error with the socket, silently drop the packet.
// If there was a problem with the socket, re-establishes the socket.
func (self *SocketHandler) Emit(record *LogRecord) error {
	self.Format(record)
	bin, err := self.Marshal(record)
	if err != nil {
		return err
	}
	return self.sendFunc(bin)
}

func (self *SocketHandler) Handle(record *LogRecord) int {
	return self.Handle2(self, record)
}

// HandleError handles an error during logging.
// An error has occurred during logging. Most likely cause connection lost.
// Close the socket so that we can retry on the next event.
func (self *SocketHandler) HandleError(record *LogRecord, err error) {
	if self.closeOnError && (self.conn != nil) {
		self.conn.Close()
		self.conn = nil
	} else {
		self.BaseHandler.HandleError(record, err)
	}
}

// Close the socket.
func (self *SocketHandler) Close() {
	self.Lock()
	defer self.Unlock()
	if self.conn != nil {
		self.conn.Close()
		self.conn = nil
	}
	self.BaseHandler.Close()
}
