package tcp_server

import (
	"errors"
	"time"
)

const (
	defaultBufSize   int           = 4096             //默认读取buf
	defaultCycleSize int           = 10000            //默认可维护的连接数量
	defaultHeartBeat time.Duration = 30 * time.Second //默认连接心跳.s
)

var (
	OverMaxConn      = errors.New("over max connect amount")
	FdNotExist       = errors.New("fd not exist")
	FdInvalid        = errors.New("fd is invalid")
	cliDisconnected  = errors.New("client disconnected")
	cliHeartOverTime = errors.New("client heartbeat over time")
	cliClosed        = errors.New("client already closed")
	cliNotExist      = errors.New("client not exist")
)

type ServerConfig struct {
	Url       string        //server 地址
	Size      int           //可维护的连接数量
	HeartBeat time.Duration //至少1秒
}
