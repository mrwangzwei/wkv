package main

import (
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"time"
	"wkv/tcp_server"
)

var svr *tcp_server.TcpServer

func main() {
	go func() {
		_ = http.ListenAndServe("127.0.0.1:6060", nil)
	}()
	var err error
	conf := tcp_server.ServerConfig{Url: "127.0.0.1:9900", HeartBeat: 5 * time.Second, Size: 100}
	svr, err = tcp_server.NewTcpServerWithConfig(conf)
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
}

func connFunc2(fd int, addr string) {
	fmt.Println("connected 222222", fd, addr)
}

func disConnFunc(fd int, addr string, err error) {
	fmt.Println("disconnected", fd, addr, err)
}

func receiveMsg(fd int, data []byte) {
	//fmt.Println("new msg", fd, string(data))
}
