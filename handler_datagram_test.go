package logging

import (
	"bytes"
	"container/list"
	"encoding/gob"
	"fmt"
	"github.com/hhkbp2/testify/require"
	"net"
	"testing"
)

func _testSetupDatagramServer(
	t *testing.T, host string, port uint16, received *list.List, ch chan int) {

	address, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", host, port))
	require.Nil(t, err)
	conn, err := net.ListenUDP("udp", address)
	require.Nil(t, err)
	go func() {
		bin := make([]byte, 1024)
		n, _, err := conn.ReadFromUDP(bin)
		require.Nil(t, err)
		defer conn.Close()
		buf := bytes.NewBuffer(bin[:n])
		decoder := gob.NewDecoder(buf)
		var record SocketLogRecord
		err = decoder.Decode(&record)
		require.Nil(t, err)
		received.PushBack(*record.Message)
		ch <- 1
	}()
}

func TestDatagramHandler(t *testing.T) {
	host := "127.0.0.1"
	port := uint16(8082)
	serverReceived := list.New()
	ch := make(chan int)
	_testSetupDatagramServer(t, host, port, serverReceived, ch)
	require.Equal(t, 0, serverReceived.Len())
	handler := NewDatagramHandler(host, port)
	logger := GetLogger("datagram")
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
