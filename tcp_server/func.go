package tcp_server

import (
	"log"
)

type (
	OnConnectionFunc     func(fd int, clientAddr string)
	OnDisConnnectionFunc func(fd int, clientAddr string)
	OnReceiveFunc        func(fd int, data []byte)
)

type receiver struct {
	fd   int
	data []byte
}

func (s *tcpServer) OnConnection(f OnConnectionFunc) {
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

func (s *tcpServer) OnDisConnection(f OnDisConnnectionFunc) {
	s.onDisConn = true
	go func() {
		log.Println("OnDisConnection is already")
		for {
			select {
			case client := <-s.closeFd:
				if f != nil {
					go f(client.fd, client.addr)
				}
			}
		}
	}()
}

func (s *tcpServer) OnReceive(f OnReceiveFunc) {
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