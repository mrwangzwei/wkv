package dns

import (
	"errors"
	"sync"
	"time"
)

var table map[string]*domainInfo

var tableLock sync.RWMutex

type domainInfo struct {
	domainLock sync.RWMutex
	lastTime   int64
	pointScore int
	times      int
	ipList     []ipInfo
}

type ipInfo struct {
	weight     int
	pointStart int
	pointEnd   int
	ipAddr     string
	country    string
	province   string
	city       string
	longitude  string
	latitude   string
}

//缓存没有的话就走正常的dns。本地dns服务器-----(.com)---->顶级dns服务器----(.baidu.com)--->权威dns服务器
func updateDomain(domain string) error {
	list, err := simpleSend(domain)
	if err != nil {
		return err
	}
	if len(list) < 1 {
		return errors.New("dns not fount ip")
	}

	//重复的更新信息
	repeat := updateRepeat(domain, list)

	//删不存在的
	delNotIN(domain, list)

	//新增的ip
	insertNew(domain, list, repeat)

	tableLock.Lock()
	table[domain].lastTime = time.Now().Unix()
	tableLock.Unlock()

	return nil
}

func updateRepeat(domain string, newIpList []string) []string {
	tableLock.Lock()
	defer tableLock.Unlock()

	if table[domain] == nil {
		emp := make([]string, 0, 0)
		return emp
	}

	table[domain].domainLock.Lock()
	defer table[domain].domainLock.Unlock()

	repeat := make([]string, 0, len(table[domain].ipList))
	for _, newIP := range newIpList {
		for index, old := range table[domain].ipList {
			if newIP == old.ipAddr {
				repeat = append(repeat, newIP)
				//更新旧的
				addr, err := RequestIpAddr(newIP)
				if err != nil {
					continue
				}
				table[domain].ipList[index].city = addr.Content.AddressDetail.City
				table[domain].ipList[index].country = "中国"
				table[domain].ipList[index].province = addr.Content.AddressDetail.Province
				table[domain].ipList[index].longitude = addr.Content.Point.Y
				table[domain].ipList[index].latitude = addr.Content.Point.X
			}
		}
	}
	return repeat
}

func delNotIN(domain string, newIpList []string) {
	tableLock.Lock()
	defer tableLock.Unlock()

	var oldExist bool
	var delIndex []int
	if table[domain] == nil {
		return
	}

	table[domain].domainLock.Lock()
	defer table[domain].domainLock.Unlock()

	for index, old := range table[domain].ipList {
		oldExist = false
		for _, newIP := range newIpList {
			if newIP == old.ipAddr {
				oldExist = true
			}
		}
		if !oldExist {
			delIndex = append(delIndex, index)
		}
	}
	for _, i := range delIndex {
		table[domain].ipList = append(table[domain].ipList[0:i], table[domain].ipList[i+1:]...)
	}
}

func insertNew(domain string, newIpList, repeat []string) {
	tableLock.Lock()
	defer tableLock.Unlock()

	var isNew bool
	for _, newIP := range newIpList {
		isNew = true
		for _, i := range repeat {
			if i == newIP {
				isNew = false
			}
		}
		if isNew {
			if table[domain] == nil {
				table[domain] = &domainInfo{
					lastTime: time.Now().Unix(),
				}
			}

			//根据ip获取地理信息
			addr, _ := RequestIpAddr(newIP)
			table[domain].domainLock.Lock()
			table[domain].ipList = append(table[domain].ipList, ipInfo{
				weight:    0,
				ipAddr:    newIP,
				country:   "中国",
				province:  addr.Content.AddressDetail.Province,
				city:      addr.Content.AddressDetail.City,
				longitude: addr.Content.Point.Y,
				latitude:  addr.Content.Point.X,
			})
			table[domain].domainLock.Unlock()
		}
	}
}

func refreshWeight(domain string) {
	tableLock.Lock()
	defer tableLock.Unlock()

	if table[domain] == nil || len(table[domain].ipList) < 1 {
		return
	}

	table[domain].domainLock.Lock()
	defer table[domain].domainLock.Unlock()

	pointScore := 0
	for index, item := range table[domain].ipList {
		table[domain].ipList[index].pointStart = pointScore
		table[domain].ipList[index].pointEnd = pointScore + item.weight
		pointScore += item.weight
	}
	table[domain].pointScore = pointScore
	table[domain].times = 0
}

func searchWeightMode(domain string) (string, error) {
	tableLock.RLock()
	defer tableLock.RUnlock()

	//不存在的话走dns流程
	if table[domain] == nil || len(table[domain].ipList) < 1 {
		list, err := simpleSend(domain)
		if err != nil {
			return "", err
		}
		if len(list) < 1 {
			return "", nil
		}
		return list[0], nil
	}

	table[domain].domainLock.Lock()
	defer table[domain].domainLock.Unlock()

	point := table[domain].times % table[domain].pointScore
	table[domain].times++
	for _, item := range table[domain].ipList {
		if item.pointStart <= point && point < item.pointEnd {
			return item.ipAddr, nil
		}
	}
	return "", nil
}

type clientIPInfo struct {
	IP       string
	Country  string
	Province string
	City     string
}

func searchClientMode(domain string) ([]clientIPInfo, error) {
	tableLock.RLock()
	//不存在的话走dns流程

	if table[domain] == nil || len(table[domain].ipList) < 1 {
		tableLock.RUnlock()
		list, err := simpleSend(domain)
		if err != nil {
			return nil, err
		}
		insertNew(domain, list, nil)
	} else {
		tableLock.RUnlock()
	}

	if table[domain] == nil {
		return nil, nil
	}

	table[domain].domainLock.RLock()
	defer table[domain].domainLock.RUnlock()

	clientIpList := make([]clientIPInfo, len(table[domain].ipList))
	for index, ip := range table[domain].ipList {
		clientIpList[index] = clientIPInfo{
			IP:       ip.ipAddr,
			Country:  ip.country,
			Province: ip.province,
			City:     ip.city,
		}
	}
	return clientIpList, nil
}
