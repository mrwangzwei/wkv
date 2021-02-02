package tcp_client

import (
	"fmt"
	"strconv"
	"testing"
	"time"
)

func TestNewCli(t *testing.T) {

	for i := 0; i < 3; i++ {
		go func() {
			conf := Config{Addr: "127.0.0.1:9900", HeartBeat: 5 * time.Second}
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
			//_, err = cli.Send("aaaaaaa")
			fmt.Println(err)
		}()
	}

	time.Sleep(10 * time.Second)

}

func Receive(data []byte) {
	fmt.Println(string(data))
}

func Disconnected() {
	fmt.Println("disconnected")
}

//压测
func TestClient(t *testing.T) {
	timer := time.NewTicker(time.Second)
	i := 0
	for {
		i++
		select {
		case <-timer.C:
			go func() {
				conf := Config{Addr: "127.0.0.1:9900", HeartBeat: 15 * time.Second}
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
				for {
					_, err = cli.Send(strconv.Itoa(i))
					time.Sleep(10 * time.Millisecond)
				}
			}()
		}
		fmt.Println("当前client数量", i)
	}
}
