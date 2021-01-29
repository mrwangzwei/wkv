package tcp_client

import (
	"fmt"
	"net"
	"time"
)

const (
	defaultBufSize   int           = 4096             //默认读取buf
	defaultCycleSize int           = 5000             //默认可维护的连接数量
	defaultHeartBeat time.Duration = 15 * time.Second //默认连接心跳.s
)

type Config struct {
	addr string
}

type client struct {
	conn      *net.TCPConn
	addr      string
	heartBeat time.Duration
}

func NewClient(host string) *client {
	return &client{addr: host, heartBeat: defaultHeartBeat}
}

func (cli *client) StartClient(str string) (err error) {
	var tcpAddr *net.TCPAddr
	tcpAddr, err = net.ResolveTCPAddr("tcp", cli.addr)
	if err != nil {
		return
	}
	var conn *net.TCPConn
	conn, err = net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		return
	}
	cli.conn = conn
	fmt.Println(cli.conn.LocalAddr().String() + " : Client connected!")
	_, err = cli.conn.Write([]byte(str + "\n"))
	return nil
}

func (cli *client) Close() {
	cli.conn.Close()
}

func (cli *client) OnReceive() {
	defer func() {
		cli.conn.Close()
	}()

	//reader := bufio.NewReaderSize(cli.conn, defaultBufSize)
	//for {
	//
	//}
}
