package logging

import (
	"fmt"
	"git.apache.org/thrift.git/lib/go/thrift"
	gen "github.com/hhkbp2/go-logging/gen-go/logging"
	"time"
)

const (
	ThriftDefaultTimeout = 3 * time.Second
)

// A handler class which write logging records, in thrift.
// the thrift connection is kept open across logging calss. If there is
// any problem with it, an attempt is made to reconnect on the next call.
type ThriftHandler struct {
	*BaseHandler
	host      string
	port      uint16
	timeout   time.Duration
	transport thrift.TTransport
	client    *gen.ThriftLoggingServiceClient
}

// Initializes the handler with a specific host address and port.
func NewThriftHandler(host string, port uint16) *ThriftHandler {
	object := &ThriftHandler{
		BaseHandler: NewBaseHandler("", LevelNotset),
		host:        host,
		port:        port,
		timeout:     ThriftDefaultTimeout,
	}
	Closer.AddHandler(object)
	return object
}

// Try to establish the thrift connection to specific host and port.
func (self *ThriftHandler) connect() error {
	address := fmt.Sprintf("%s:%d", self.host, self.port)
	socket, err := thrift.NewTSocket(address)
	if err != nil {
		return err
	}
	transport := thrift.NewTFramedTransport(socket)
	factory := thrift.NewTBinaryProtocolFactoryDefault()
	client := gen.NewThriftLoggingServiceClientFactory(transport, factory)
	if err := transport.Open(); err != nil {
		return err
	}
	self.transport = transport
	self.client = client
	return nil
}

// Emit a record.
// Report the logging record to server(establish the connectino if necessary).
// If there is an error with connection, silently drop the packet.
func (self *ThriftHandler) Emit(record *LogRecord) error {
	self.Format(record)
	r := &gen.ThriftLogRecord{
		Name:     record.Name,
		Level:    int32(record.Level),
		PathName: record.PathName,
		FileName: record.FileName,
		LineNo:   int32(record.LineNo),
		FuncName: record.FuncName,
		Message:  record.Message,
	}
	if self.client == nil {
		if err := self.connect(); err != nil {
			return err
		}
	}
	return self.client.Report(r)
}

func (self *ThriftHandler) Handle(record *LogRecord) int {
	return self.Handle2(self, record)
}

// Handles an error during logging.
// An error has occurred during logging. Most likely cause connection lost.
// Close the thrift connection so that we can retry on the next event.
func (self *ThriftHandler) HandleError(record *LogRecord, err error) {
	if self.client != nil {
		self.transport.Close()
		self.transport = nil
		self.client = nil
	}
	self.BaseHandler.HandleError(record, err)
}

// Close the thrift client.
func (self *ThriftHandler) Close() {
	self.Lock()
	defer self.Unlock()
	if self.client != nil {
		self.transport.Close()
		self.transport = nil
		self.client = nil
	}
	self.BaseHandler.Close()
}
