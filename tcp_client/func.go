package tcp_client

type (
	OnReceiveFunc      func(data []byte)
	OnDisConnectedFunc func()
)

func (cli *client) OnMsg(f OnReceiveFunc) {
	if cli.onMsg {
		return
	}
	cli.onMsg = true
	go func() {
		for d := range cli.msgCh {
			if f != nil {
				go f(d) //考虑加协程池
			}
		}
	}()
}

func (cli *client) OnDisconnected(f OnDisConnectedFunc) {
	if cli.onDisCon {
		return
	}
	cli.onDisCon = true
	go func() {
		for _ = range cli.disConCh {
			if f != nil {
				go f() //考虑加协程池
			}
		}
	}()
}
