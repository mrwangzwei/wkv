package tcp_server

import (
	"errors"
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

type ServerConfig struct {
	Url       string        //server 地址
	Size      int           //可维护的连接数量
	HeartBeat time.Duration //至少1秒
}
