package tcp_server

import (
	"bufio"
	"errors"
	"io"
	"log"
	"net"
	"sync"
	"time"
)

const (
	defaultBufSize   int           = 4096             //默认读取buf
	defaultCycleSize int           = 5000             //默认可维护的连接数量
	defaultHeartBeat time.Duration = 30 * time.Second //默认连接心跳.s
)

var (
	OverMaxConn = errors.New("over max connect amount")
	FdExist     = errors.New("fd not exist")
	FdInvalid   = errors.New("fd is invalid")
)

type tcpServer struct {
	addr           string
	listener       *net.TCPListener
	clients        []*client
	cliSize        int
	heartBeatLimit time.Duration
	cycle          int
	cursor         int
	bufSize        int
	lock           sync.Mutex
	newFd          chan *client
	closeFd        chan *client
	receiver       chan *receiver
	onConn         bool
	onDisConn      bool
	onMsg          bool
}

type client struct {
	fd      int
	conn    *net.TCPConn
	stat    bool   //连钱连接是否可用
	heartAt int64  //上次心跳时间。ms
	addr    string //客户端地址
	lock    sync.Mutex
}

func NewTcpServer(addr string) (*tcpServer, error) {
	conf := ServerConfig{
		addr,
		defaultCycleSize,
		defaultHeartBeat,
	}
	return NewTcpServerWithConfig(conf)
}

func NewTcpServerWithConfig(conf ServerConfig) (*tcpServer, error) {
	if conf.HeartBeat < time.Second {
		return nil, errors.New("heart beat must over one second")
	}
	if conf.Size == 0 {
		conf.Size = defaultCycleSize
	}
	return &tcpServer{
		addr:           conf.Url,
		clients:        make([]*client, conf.Size),
		cliSize:        conf.Size,
		heartBeatLimit: conf.HeartBeat,
		newFd:          make(chan *client),
		closeFd:        make(chan *client),
		receiver:       make(chan *receiver),
	}, nil
}

func (s *tcpServer) StartServer() (err error) {
	var tcpAddr *net.TCPAddr
	tcpAddr, err = net.ResolveTCPAddr("tcp", s.addr)
	if err != nil {
		return
	}
	var tcpListener *net.TCPListener
	tcpListener, err = net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		return
	}
	s.listener = tcpListener
	defer s.listener.Close()

	go s.checkHeartBeat()

	log.Println("Server is already ...")
	for {
		tcpConn, err := s.listener.AcceptTCP()
		if err != nil {
			log.Println("new client conn err", err)
			continue
		}
		c, err := s.addClient(tcpConn)
		if err != nil {
			log.Println("add new client conn err", err)
			tcpConn.Close()
			continue
		}
		if s.onConn {
			s.newFd <- c
		}
		go s.readConn(c)
	}
}

func (s *tcpServer) Send(fd int, msg string) (err error) {
	if fd > s.cliSize || fd < 1 {
		err = FdExist
		return
	}
	var cli *client
	cli, err = s.searchFd(fd)
	if err != nil {
		return
	}
	_, err = cli.send(msg)
	if err != nil {
		return
	}
	return
}

func (s *tcpServer) searchFd(fd int) (cli *client, err error) {
	if fd > s.cliSize || fd < 1 {
		err = FdExist
		return
	}
	cli = s.clients[fd-1]
	if cli == nil {
		err = FdExist
		return
	}
	if !cli.stat {
		err = FdInvalid
		return
	}
	return
}

func (s *tcpServer) addClient(conn *net.TCPConn) (c *client, err error) {
	c = &client{
		conn:    conn,
		stat:    true,
		heartAt: time.Now().Unix(),
		addr:    conn.RemoteAddr().String(),
	}
	s.lock.Lock()
	defer s.lock.Unlock()
	nowCycle := s.cycle
	length := len(s.clients)
	for {
		if ok := s.clients[s.cursor]; ok != nil {
			if ok.stat {
				if s.cursor == length-1 {
					s.cycle++
					s.cursor = 0
					//循环2轮还没有的话就直接抛弃
					if s.cycle >= nowCycle+2 {
						err = OverMaxConn
						return
					}
				}
			} else {
				break
			}
		} else {
			break
		}
	}
	c.fd = s.cursor + 1
	s.clients[s.cursor] = c
	s.cursor++
	if s.cursor >= length {
		s.cursor = 0
	}
	return
}

func (s *tcpServer) readConn(c *client) {
	defer func() {
		c.disable()
		_ = c.conn.Close()
		if s.onDisConn {
			s.closeFd <- c
		}
	}()
	//获取一个连接的reader读取流
	reader := bufio.NewReaderSize(c.conn, defaultBufSize)
	//接收并返回消息
	for {
		message, err := buffReader(reader)
		if err != nil || err == io.EOF {
			return
		}
		if s.onMsg {
			s.receiver <- &receiver{c.fd, message}
		}
		c.beatHeart()
	}
}

func (cli *client) beatHeart() {
	cli.lock.Lock()
	defer cli.lock.Unlock()
	cli.heartAt = time.Now().Unix()
}

func (cli *client) enable() {
	cli.lock.Lock()
	defer cli.lock.Unlock()
	cli.stat = true
}

func (cli *client) disable() {
	cli.lock.Lock()
	defer cli.lock.Unlock()
	cli.stat = false
}

func (cli *client) send(msg string) (len int, err error) {
	len, err = cli.conn.Write([]byte(msg))
	return
}

func (s *tcpServer) checkHeartBeat() {
	ticker := time.NewTicker(time.Second)
	for {
		<-ticker.C
		for _, cli := range s.clients {
			func() {
				if cli == nil {
					return
				}
				cli.lock.Lock()
				defer cli.lock.Unlock()
				now := time.Now().Unix()
				if cli.stat == false {
					return
				}
				//心跳超时。主动关闭连接
				if (now - cli.heartAt) > int64(s.heartBeatLimit)/1e9 {
					cli.stat = false
					_ = cli.conn.Close()
				}
			}()
		}
	}
}
