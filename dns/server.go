package dns

import (
	"encoding/json"
	"errors"
	"log"
	"net"
	"sync"
	"time"
)

type DNS interface {
	Listen() error
	AddWeightIpInfo(domain string, ip string, weight int) error
}

//服务模式
const (
	//客户端选择模式。此模式下返回ip列表，有客户端选择
	ClientMode = "client"
	//权重模式。此模式下根据权重仅返回一个ip
	WeightMode = "weight"
)

type dnsServer struct {
	Port int64
	TTL  int64
	Mode string
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
		Port: port,
		TTL:  ttl,
		Mode: mode,
	}
	table = map[string]*domainInfo{}
	return s, nil
}

func (s *dnsServer) Listen() error {
	var err error
	//udp监听
	s.conn, err = net.ListenUDP("udp", &net.UDPAddr{Port: int(s.Port)})
	if err != nil {
		log.Printf(err.Error())
		return err
	}
	defer s.conn.Close()

	for {
		buf := make([]byte, dnsPacketLen)

		//读取接收到的dns协议包
		length, addr, err := s.conn.ReadFromUDP(buf)
		log.Println("conn.ReadFromUDP", addr, err, string(buf))
		if err != nil {
			log.Printf(err.Error())
			continue
		}

		//todo 解析dns协议包得到domain,暂时解不来。。。。先json代替一下实现过程
		mmp := make(map[string]string)
		err = json.Unmarshal(buf[:length], &mmp)
		var domain string
		if val, ok := mmp["domain"]; ok {
			domain = val
		}
		go searchAndRespone(s, domain, addr)
	}
}

func (s *dnsServer) AddWeightIpInfo(domain, ip string, weight int) error {
	if s.Mode != WeightMode {
		return errors.New("server is not in weight mode")
	}

	tableLock.Lock()
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
		table[domain].domainLock.Lock()
		table[domain].ipList = append(table[domain].ipList, ipInfo{
			weight:    weight,
			ipAddr:    ip,
			country:   "中国",
			province:  info.Content.AddressDetail.Province,
			city:      info.Content.AddressDetail.City,
			longitude: info.Content.Point.Y,
			latitude:  info.Content.Point.X,
		})
		table[domain].domainLock.Unlock()
	}
	tableLock.Unlock()
	//刷权重
	refreshWeight(domain)
	return nil
}

func (s *dnsServer) Search(domain string) ([]byte, error) {
	switch s.Mode {
	case WeightMode:
		info, err := searchWeightMode(domain)
		if err != nil {
			return nil, err
		}
		return []byte(info), nil
	case ClientMode:
		//先从缓存找到ip信息，并且看是否超时，超时就重新从dns流程获取
		tableLock.RLock()
		if table[domain] != nil && table[domain].lastTime+s.TTL < time.Now().Unix() {
			tableLock.RUnlock()
			err := updateDomain(domain)
			if err != nil {
				log.Println(err)
			}
		} else {
			tableLock.RUnlock()
		}
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

func searchAndRespone(s *dnsServer, domain string, addr *net.UDPAddr) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
		}
	}()
	res, err := s.Search(domain)
	if err != nil {
		log.Println(err)
		return
	}
	go writeUdpCtrl(s.conn, res, addr)
}

const BufNextExist = 1 << 7 //10000000	,还没传完
const BufNextNone = 0       //00000000	,传完了
const DataBufSize = 511
const MaxBufCount = (1 << 7) - 1

//控制每次传512字节
//第一个字节做标志位,第一位（是否还有下一个包），后7位表当前第几个包（即最大127个包）
//发包超过最大重试次数就直接跳过
func writeUdpCtrl(conn *net.UDPConn, buf []byte, addr *net.UDPAddr) {
	nowBag := 0
	for {
		leftBufLen := len(buf)
		if (leftBufLen / DataBufSize) > MaxBufCount {
			log.Println("over the max buf count")
			return
		}
		if leftBufLen < 1 {
			return
		}

		sendBuf := make([]byte, 1)
		if leftBufLen <= DataBufSize {
			sendBuf[0] = byte(BufNextNone + nowBag)
			sendBuf = append(sendBuf, buf[:]...)
			buf = append(buf[0:0])
		} else {
			sendBuf[0] = byte(BufNextExist + nowBag)
			sendBuf = append(sendBuf, buf[:DataBufSize]...)
			buf = append(buf[0:0], buf[DataBufSize:]...)
		}

		times := 1
		for _, err := conn.WriteToUDP(sendBuf, addr); err != nil; {
			log.Println(err)
			if times > 10 {
				log.Println("resend over times")
				log.Println("nowBag", nowBag)
				log.Println("leftBufLen", leftBufLen)
				break
			}
			times++
		}
		nowBag++
	}

}
