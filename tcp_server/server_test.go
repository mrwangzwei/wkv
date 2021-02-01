package tcp_server

import (
	"fmt"
	"testing"
)

var svr *tcpServer

func TestServer(t *testing.T) {
	conf := ServerConfig{Url: "127.0.0.1:9900"}
	svr, err := NewTcpServerWithConfig(conf)
	if err != nil {
		fmt.Println(err)
		return
	}

	//注册监听方法，只有第一个生效
	svr.OnConnection(connFunc)
	svr.OnConnection(connFunc2)

	svr.OnDisConnection(disConnFunc)

	svr.OnReceive(receiveMsg)

	err = svr.StartServer()
	fmt.Println(err)
}

func connFunc(fd int, addr string) {
	fmt.Println("connected", fd, addr)
	err := svr.Send(fd, "welcome")
	fmt.Println("send", fd, err)
}

func connFunc2(fd int, addr string) {
	fmt.Println("connected 222222", fd, addr)
	err := svr.Send(fd, "welcome 222222")
	fmt.Println("send 222222", fd, err)
}

func disConnFunc(fd int, addr string) {
	fmt.Println("disconnected", fd, addr)
}

func receiveMsg(fd int, data []byte) {
	fmt.Println("new msg", fd, string(data))
	err := svr.Send(fd, "welcome")
	fmt.Println("answer", fd, err)
}
