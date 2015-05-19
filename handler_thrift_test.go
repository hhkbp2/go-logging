package logging

import (
	"container/list"
	"fmt"
	gen "gen-go/logging"
	"git.apache.org/thrift.git/lib/go/thrift"
	"github.com/hhkbp2/testify/require"
	"testing"
)

type TestThriftServerHandler struct {
	received *list.List
	ch       chan int
}

func NewTestThriftServerHandler(
	received *list.List, ch chan int) *TestThriftServerHandler {

	return &TestThriftServerHandler{
		received: received,
		ch:       ch,
	}
}

func (self *TestThriftServerHandler) Report(record *gen.ThriftLogRecord) error {
	self.received.PushBack(record.Message)
	self.ch <- 1
	return nil
}

func _TestSetupThriftServer(
	t *testing.T, host string, port uint16, received *list.List, ch chan int) thrift.TServer {

	handler := NewTestThriftServerHandler(received, ch)
	processor := gen.NewThriftLoggingServiceProcessor(handler)

	address := fmt.Sprintf("%s:%d", host, port)
	serverTransport, err := thrift.NewTServerSocket(address)
	require.Nil(t, err)
	transportFactory := thrift.NewTFramedTransportFactory(thrift.NewTTransportFactory())
	protocolFactory := thrift.NewTBinaryProtocolFactoryDefault()
	server := thrift.NewTSimpleServer4(
		processor,
		serverTransport,
		transportFactory,
		protocolFactory)
	go func() {
		server.Serve()
	}()
	return server
}

func TestThriftHandler(t *testing.T) {
	host := "127.0.0.1"
	port := uint16(8082)
	serverReceived := list.New()
	ch := make(chan int, 1)
	server := _TestSetupThriftServer(t, host, port, serverReceived, ch)
	require.Equal(t, 0, serverReceived.Len())
	handler := NewThriftHandler(host, port)
	logger := GetLogger("thrift")
	logger.AddHandler(handler)
	message := "test"
	logger.Errorf(message)
	handler.Close()
	<-ch
	server.Stop()
	require.Equal(t, 1, serverReceived.Len())
	m, ok := serverReceived.Front().Value.(string)
	require.True(t, ok)
	require.Equal(t, message, m)
}
