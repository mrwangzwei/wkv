package tcp_server

import (
	"log"
)

type (
	OnConnectionFunc    func(fd int, clientAddr string)
	OnDisConnectionFunc func(fd int, clientAddr string, err error)
	OnReceiveFunc       func(fd int, data []byte)
)

type newConn struct {
	fd   int
	addr string
}

type disConn struct {
	fd   int
	addr string
	err  error
}

type receiver struct {
	fd   int
	data []byte
}

func (s *TcpServer) OnConnection(f OnConnectionFunc) {
	if s.onConn {
		return
	}
	s.onConn = true
	go func() {
		log.Println("OnConnection is already")
		for {
			select {
			case client := <-s.newFd:
				if f != nil {
					go f(client.fd, client.addr) //考虑加协程池
				}
			}
		}
	}()
}

func (s *TcpServer) OnDisConnection(f OnDisConnectionFunc) {
	if s.onDisConn {
		return
	}
	s.onDisConn = true
	go func() {
		log.Println("OnDisConnection is already")
		for {
			select {
			case dis := <-s.closeFd:
				if f != nil {
					go f(dis.fd, dis.addr, dis.err)
				}
			}
		}
	}()
}

func (s *TcpServer) OnReceive(f OnReceiveFunc) {
	if s.onMsg {
		return
	}
	s.onMsg = true
	go func() {
		log.Println("OnReceive is already")
		for {
			select {
			case r := <-s.receiver:
				if f != nil {
					go f(r.fd, r.data)
				}
			}
		}
	}()
}
