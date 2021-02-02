package tcp_client

import (
	"fmt"
	"testing"
	"time"
)

func TestNewCli(t *testing.T) {
	conf := Config{Addr: "127.0.0.1:9900", HeartBeat: 7 * time.Second}
	cli, err := NewCliWithConfig(conf)
	if err != nil {
		fmt.Println(err)
		return
	}

	cli.OnMsg(Receive)
	cli.OnDisconnected(Disconnected)

	if err = cli.StartClient(); err != nil {
		fmt.Println(err)
		return
	}

	//time.Sleep(10 * time.Second)
	_, err = cli.Send("aaaaaaa")
	fmt.Println(err)

	time.Sleep(10 * time.Second)

}

func Receive(data []byte) {
	fmt.Println(string(data))
}

func Disconnected() {
	fmt.Println("disconnected")
}
