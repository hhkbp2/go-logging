package logging

// A handler class which writes logging records, in gob format, to
// a datagram socket.
type DatagramHandler struct {
	*SocketHandler
}

func NewDatagramHandler(host string, port uint16) *DatagramHandler {
	object := &DatagramHandler{
		SocketHandler: NewSocketHandler(host, port),
	}
	object.makeConnFunc = object.makeUDPSocket
	object.sendFunc = object.sendUDP
	object.closeOnError = false
	return object
}
