package tcp_server

import "time"

type ServerConfig struct {
	Url       string        //server 地址
	Size      int           //可维护的连接数量
	HeartBeat time.Duration //至少1秒
}
