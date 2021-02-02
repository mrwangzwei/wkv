package tcp_client

import (
	"errors"
	"time"
)

var (
	EmptyMsg = errors.New("send msg can`t be empty")
)

const (
	defaultBufSize   int           = 4096             //默认读取buf
	defaultHeartBeat time.Duration = 15 * time.Second //默认连接心跳.s
)

type Config struct {
	Addr      string
	HeartBeat time.Duration
}
