package tcp_server

import (
	"fmt"
	"testing"
	"time"
)

var svr *TcpServer

func TestServer(t *testing.T) {
	var err error
	conf := ServerConfig{Url: "127.0.0.1:9900", HeartBeat: 5 * time.Second}
	svr, err = NewTcpServerWithConfig(conf)
	if err != nil {
		fmt.Println(err)
		return
	}

	//注册监听方法，只有第一个生效
	svr.OnConnection(connFunc)
	svr.OnConnection(connFunc2)

	svr.OnDisConnection(disConnFunc)

	svr.OnReceive(receiveMsg)

	if err = svr.StartServer(); err != nil {
		fmt.Println(err)
		return
	}
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

func disConnFunc(fd int, addr string, err error) {
	fmt.Println("disconnected", fd, addr, err)
	cerr := svr.Close(fd)
	fmt.Println(cerr)
}

func receiveMsg(fd int, data []byte) {
	fmt.Println("new msg", fd, string(data))
	err := svr.Send(fd, "receiveMsg")
	fmt.Println("answer", fd, err)
}

func TestChannel(t *testing.T) {
	cc := make(chan int)
	go func() {
		for {
			select {
			case a := <-cc:
				fmt.Println(a)
			default:
			}
		}

	}()
	cc <- 2
	close(cc)
	time.Sleep(2 * time.Second)

}
