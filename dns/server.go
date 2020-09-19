package dns

import (
	"fmt"
	"log"
	"net"
	"time"
)

type DNS interface {
	Listen()
}

type dnsServer struct {
	Port  int
	TTL   int
	conn  *net.UDPConn
	table []map[string][]ipInfo
}

type ipInfo struct {
	ipAddr   string
	country  string
	province string
	city     string
}

func NewServer(port, ttl int) DNS {
	return &dnsServer{
		Port: port,
		TTL:  ttl,
	}
}

func (s *dnsServer) Listen() {
	var err error
	//udp监听
	s.conn, err = net.ListenUDP("udp", &net.UDPAddr{Port: s.Port})
	if err != nil {
		log.Printf(err.Error())
		return
	}
	defer s.conn.Close()

	for {
		buf := make([]byte, dnsPacketLen)

		//读取接收到的dns协议包
		_, addr, err := s.conn.ReadFromUDP(buf)
		fmt.Println("conn.ReadFromUDP", addr, err, buf)
		if err != nil {
			continue
		}

		//解析dns协议包得到domain

		nowTime := time.Now().Unix()
		//先从缓存找到ip信息，并且看是否超时，超时就重新从dns流程获取

		//缓存没有的话就走正常的dns。本地dns服务器-----(.com)---->顶级dns服务器----(.baidu.com)--->权威dns服务器

		//存入缓存，更新超时时间

		//拼装结构，返回结果
	}
}

func (s *dnsServer) Search() {

}
