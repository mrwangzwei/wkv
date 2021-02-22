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

type TcpServer struct {
	addr           string //地址
	listener       *net.TCPListener
	clients        []*client     //cli集合
	cliSize        int           //cli最大数量
	heartBeatLimit time.Duration //心跳限制
	cycle          int           //当前fd循环次数
	cursor         int           //下一个client的游标
	bufSize        int           //read申请的buf大小
	lock           sync.Mutex
	newFd          chan *newConn  //新连接通知
	closeFd        chan *disConn  //连接关闭通知
	receiver       chan *receiver //新消息通知
	cliCloseCh     chan *cliClose //统一关闭连接入口
	onConn         bool           //是否注册连接通知
	onDisConn      bool           //是否注册关闭连接通知
	onMsg          bool           //是否注册新消息通知
}

type cliClose struct {
	cli *client
	err error
}

type client struct {
	fd      int
	conn    *net.TCPConn
	stat    bool   //连钱连接是否可用
	heartAt int64  //上次心跳时间。ms
	addr    string //客户端地址
	lock    sync.Mutex
}

func NewTcpServer(addr string) (*TcpServer, error) {
	conf := ServerConfig{
		addr,
		defaultCycleSize,
		defaultHeartBeat,
	}
	return NewTcpServerWithConfig(conf)
}

func NewTcpServerWithConfig(conf ServerConfig) (*TcpServer, error) {
	if conf.HeartBeat < time.Second {
		return nil, errors.New("heart beat must over one second")
	}
	if conf.Size == 0 {
		conf.Size = defaultCycleSize
	}
	return &TcpServer{
		addr:           conf.Url,
		clients:        make([]*client, conf.Size),
		cliSize:        conf.Size,
		heartBeatLimit: conf.HeartBeat,
		newFd:          make(chan *newConn),
		closeFd:        make(chan *disConn),
		receiver:       make(chan *receiver),
		cliCloseCh:     make(chan *cliClose),
	}, nil
}

func (s *TcpServer) StartServer() (err error) {
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

	go s.listenCloseCli()

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
			log.Println("add new client conn err:", err, tcpConn.RemoteAddr().String())
			tcpConn.Close()
			continue
		}
		if s.onConn {
			s.newFd <- &newConn{fd: c.fd, addr: c.addr}
		}
		go s.readConn(c)
	}
}

func (s *TcpServer) Send(fd int, msg string) (err error) {
	if fd > s.cliSize || fd < 1 {
		err = FdNotExist
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

func (s *TcpServer) searchFd(fd int) (cli *client, err error) {
	if fd > s.cliSize || fd < 1 {
		err = FdNotExist
		return
	}
	cli = s.clients[fd-1]
	if cli == nil {
		err = FdNotExist
		return
	}
	if !cli.stat {
		err = FdInvalid
		return
	}
	return
}

func (s *TcpServer) addClient(conn *net.TCPConn) (c *client, err error) {
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
				s.cursor++
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

func (s *TcpServer) readConn(cli *client) {
	defer s.closeCli(cli, cliDisconnected)
	//获取一个连接的reader读取流
	reader := bufio.NewReaderSize(cli.conn, defaultBufSize)
	//接收并返回消息
	for {
		message, err := buffReader(reader)
		if err != nil || err == io.EOF {
			return
		}
		if string(message) != "B" && s.onMsg {
			s.receiver <- &receiver{cli.fd, message}
		}
		cli.beatHeart()
	}
}

func (s *TcpServer) closeCli(cli *client, err error) {
	s.cliCloseCh <- &cliClose{cli: cli, err: err}
}

func (s *TcpServer) listenCloseCli() {
	for c := range s.cliCloseCh {
		res := c.cli.disable()
		if res {
			_ = c.cli.conn.Close()
			if s.onDisConn {
				s.closeFd <- &disConn{fd: c.cli.fd, addr: c.cli.addr, err: c.err}
			}
		}
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

func (cli *client) disable() bool {
	cli.lock.Lock()
	defer cli.lock.Unlock()
	if cli.stat == true {
		cli.stat = false
		return true
	}
	return false
}

func (cli *client) send(msg string) (l int, err error) {
	l, err = cli.conn.Write([]byte(msg + "\n"))
	return
}

func (s *TcpServer) checkHeartBeat() {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
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
						s.closeCli(cli, cliHeartOverTime)
					}
				}()
			}
		}

	}
}
