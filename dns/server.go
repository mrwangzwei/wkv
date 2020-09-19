package dns

import (
	"fmt"
	"net"
	"wkv/wkv_a"
)

type DNS interface {
	Listen()
}

type dnsServer struct {
	Port  int
	conn  *net.UDPConn
	table *wkv_a.Boot
}

func NewServer(port int) DNS {
	return &dnsServer{
		Port: port,
	}
}

//dns数据包长度
const dnsPacketLen = 512

func (s *dnsServer) Listen() {
	var err error
	s.table = &wkv_a.Boot{}
	s.conn, err = net.ListenUDP("udp", &net.UDPAddr{Port: s.Port})
	fmt.Println("net.ListenUDP", err)
	defer s.conn.Close()

	for {
		buf := make([]byte, dnsPacketLen)
		_, addr, err := s.conn.ReadFromUDP(buf)
		fmt.Println("conn.ReadFromUDP", addr, err, buf)
		if err != nil {
			continue
		}

	}
}

func (s *dnsServer) Search() {

}
