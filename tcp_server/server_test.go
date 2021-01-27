package tcp_server

import (
	"fmt"
	"testing"
)

func TestServer(t *testing.T) {
	server := NewTcpServer("127.0.0.1:9900")

	server.OnConnection(connFunc)

	server.OnDisConnection(disConnFunc)

	server.OnReceive(receiveMsg)

	err := server.StartServer()
	fmt.Println(err)
}

func connFunc(fd int, addr string) {
	fmt.Println(fd, addr)
}

func disConnFunc(fd int, addr string) {
	fmt.Println(fd, addr)
}

func receiveMsg(fd int, data []byte) {
	fmt.Println(fd, string(data))
}
