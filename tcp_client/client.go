package tcp_client

import (
	"errors"
	"fmt"
	"net"
	"time"
)

const (
	defaultBufSize   int           = 4096             //默认读取buf
	defaultCycleSize int           = 5000             //默认可维护的连接数量
	defaultHeartBeat time.Duration = 15 * time.Second //默认连接心跳.s
)

type client struct {
	conn      *net.TCPConn
	svrAddr   string
	heartBeat time.Duration
	msgCh     chan []byte
	disConCh  chan bool
	onMsg     bool
	onDisCon  bool
}

func NewCli(host string) (*client, error) {
	conf := Config{Addr: host, HeartBeat: defaultHeartBeat}
	return NewCliWithConfig(conf)
}

func NewCliWithConfig(conf Config) (*client, error) {
	if conf.HeartBeat < time.Second {
		return nil, errors.New("heart beat must over one second")
	}
	return &client{svrAddr: conf.Addr, heartBeat: conf.HeartBeat}, nil
}

func (cli *client) StartClient(str string) (err error) {
	var tcpAddr *net.TCPAddr
	tcpAddr, err = net.ResolveTCPAddr("tcp", cli.svrAddr)
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
