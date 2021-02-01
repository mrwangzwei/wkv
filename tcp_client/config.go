package tcp_client

import "time"

type Config struct {
	Addr      string
	HeartBeat time.Duration
}
