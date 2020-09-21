package dns

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net"
	"sync"
	"time"
)

type DNS interface {
	Listen() error
	AddWeightIpInfo(string, string, int) error
}

//服务模式
const (
	//客户端选择模式。此模式下返回ip列表，有客户端选择
	ClientMode = "client"
	//权重模式。此模式下根据权重仅返回一个ip
	WeightMode = "weight"
)

type dnsServer struct {
	port int64
	ttl  int64
	mode string
	conn *net.UDPConn
	lock sync.RWMutex
}

/**
** port：端口号
** ttl：缓存超时时间。秒级
 */
func NewServer(port, ttl int64, mode string) (DNS, error) {
	switch mode {
	case ClientMode:
	case WeightMode:
		break
	default:
		return nil, errors.New("mode is not available")
	}
	s := &dnsServer{
		port: port,
		ttl:  ttl,
		mode: mode,
	}
	table = map[string]*domainInfo{}
	return s, nil
}

func (s *dnsServer) Listen() error {
	var err error
	//udp监听
	s.conn, err = net.ListenUDP("udp", &net.UDPAddr{Port: int(s.port)})
	if err != nil {
		log.Printf(err.Error())
		return err
	}
	defer s.conn.Close()

	for {
		buf := make([]byte, dnsPacketLen)

		//读取接收到的dns协议包
		_, addr, err := s.conn.ReadFromUDP(buf)
		fmt.Println("conn.ReadFromUDP", addr, err, buf)
		if err != nil {
			log.Printf(err.Error())
			continue
		}

		//解析dns协议包得到domain

		//nowTime := time.Now().Unix()
		//先从缓存找到ip信息，并且看是否超时，超时就重新从dns流程获取

		//缓存没有的话就走正常的dns。本地dns服务器-----(.com)---->顶级dns服务器----(.baidu.com)--->权威dns服务器

		//根据ip获取地理信息

		//存入缓存，更新超时时间

		//拼装结构，返回结果
	}
}

func (s *dnsServer) AddWeightIpInfo(domain, ip string, weight int) error {
	if s.mode != WeightMode {
		return errors.New("server is not in weight mode")
	}
	tableLock.Lock()
	defer tableLock.Unlock()

	//存在的话直接替换权重
	var ipExist bool
	if table[domain] != nil {
		for index, item := range table[domain].ipList {
			if item.ipAddr == ip {
				ipExist = true
				table[domain].ipList[index].weight = weight
			}
		}
	}
	//不存在的话，查ip地理信息后加入进去
	if !ipExist {
		info, err := RequestIpAddr(ip)
		if err != nil {
			return err
		}
		if table[domain] == nil {
			table[domain] = &domainInfo{
				lastTime: time.Now().Unix(),
			}
		}
		table[domain].ipList = append(table[domain].ipList, ipInfo{
			weight:    weight,
			ipAddr:    ip,
			country:   "中国",
			province:  info.Content.AddressDetail.Province,
			city:      info.Content.AddressDetail.City,
			longitude: info.Content.Point.Y,
			latitude:  info.Content.Point.X,
		})
	}
	refreshWeight(domain)
	return nil
}

func (s *dnsServer) Search(domain string) ([]byte, error) {
	switch s.mode {
	case WeightMode:
		info, err := searchWeightMode(domain)
		if err != nil {
			return nil, err
		}
		return []byte(info), nil
	case ClientMode:
		info, err := searchClientMode(domain)
		if err != nil {
			return nil, err
		}
		res, _ := json.Marshal(info)
		return res, nil
	default:
		return nil, errors.New("unknown mode")
	}
}
