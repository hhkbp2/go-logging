package logging

import (
	"container/list"
	"encoding/gob"
	"fmt"
	"github.com/hhkbp2/testify/require"
	"net"
	"testing"
)

func SetupTestSocketServer(
	t *testing.T, host string, port uint16, received *list.List, ch chan int) {
	address := fmt.Sprintf("%s:%d", host, port)
	listener, err := net.Listen("tcp", address)
	require.Nil(t, err)
	go func() {
		conn, err := listener.Accept()
		require.Nil(t, err)
		go func(conn net.Conn) {
			defer conn.Close()
			decoder := gob.NewDecoder(conn)
			var record SocketLogRecord
			err := decoder.Decode(&record)
			require.Nil(t, err)
			received.PushBack(*record.Message)
			ch <- 1
		}(conn)
		listener.Close()
	}()
}

func TestSocketHandler(t *testing.T) {
	host := "127.0.0.1"
	port := uint16(8082)
	serverReceived := list.New()
	ch := make(chan int)
	SetupTestSocketServer(t, host, port, serverReceived, ch)
	handler := NewSocketHandler(host, port)
	logger := GetLogger("a")
	logger.AddHandler(handler)
	message := "test"
	logger.Errorf(message)
	handler.Close()
	<-ch
	require.Equal(t, 1, serverReceived.Len())
	m, ok := serverReceived.Front().Value.(string)
	require.True(t, ok)
	require.Equal(t, message, m)
}
