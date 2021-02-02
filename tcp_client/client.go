package tcp_client

import (
	"bufio"
	"errors"
	"io"
	"net"
	"sync"
	"time"
)

type client struct {
	conn      *net.TCPConn
	svrAddr   string
	heartBeat time.Duration
	msgCh     chan []byte
	disConCh  chan bool
	onMsg     bool
	onDisCon  bool
	lock      sync.Mutex
}

func NewCli(host string) (*client, error) {
	conf := Config{Addr: host, HeartBeat: defaultHeartBeat}
	return NewCliWithConfig(conf)
}

func NewCliWithConfig(conf Config) (*client, error) {
	if conf.HeartBeat < time.Second {
		return nil, errors.New("heart beat must over one second")
	}
	return &client{
		svrAddr:   conf.Addr,
		heartBeat: conf.HeartBeat,
		msgCh:     make(chan []byte),
		disConCh:  make(chan bool),
	}, nil
}

func (cli *client) StartClient() (err error) {
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
	go cli.beatHeart()

	go cli.readMsg()
	return nil
}

func (cli *client) Send(msg string) (l int, err error) {
	if len(msg) < 1 {
		err = EmptyMsg
		return
	}
	l, err = cli.conn.Write([]byte(msg + "\n"))
	return
}

func (cli *client) beatHeart() {
	ticker := time.NewTicker(cli.heartBeat)
	for {
		<-ticker.C
		if cli == nil {
			return
		}
		func() {
			cli.lock.Lock()
			defer cli.lock.Unlock()
			if cli.conn == nil {
				return
			}
			_, _ = cli.conn.Write([]byte("B" + "\n"))
		}()
	}
}

func (cli *client) Close() {
	_ = cli.conn.Close()
}

func (cli *client) readMsg() {
	defer func() {
		_ = cli.conn.Close()
		if cli.onDisCon {
			cli.disConCh <- true
		}
	}()
	//获取一个连接的reader读取流
	reader := bufio.NewReaderSize(cli.conn, defaultBufSize)
	//接收并返回消息
	for {
		message, err := buffReader(reader)

		if err != nil || err == io.EOF {
			return
		}
		if cli.onMsg {
			cli.msgCh <- message
		}
	}
}
