package redigo_sub

import (
	"context"
	"errors"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"sync"
	"time"
)

type Subscribe struct {
	ctx  context.Context
	conn redis.Conn
	// channel对应的用户自定义函数
	chFunc map[interface{}]func(ctx context.Context, data []byte)
	chList []interface{}
	// chFunc读写时避免并发
	chFuncLc *sync.RWMutex
	// 订阅实例
	pubSubConn redis.PubSubConn
	// 消息分发通道
	dataCh chan struct {
		channel interface{}
		data    []byte
		err     error
	}
	closeCtx    context.Context
	closeCancel context.CancelFunc
}

func NewSub(ctx context.Context, conn redis.Conn) *Subscribe {
	closeCtx, cancel := context.WithCancel(context.Background())
	return &Subscribe{
		ctx:      ctx,
		conn:     conn,
		chFunc:   make(map[interface{}]func(ctx context.Context, data []byte)),
		chFuncLc: &sync.RWMutex{},
		dataCh: make(chan struct {
			channel interface{}
			data    []byte
			err     error
		}),
		closeCtx:    closeCtx,
		closeCancel: cancel,
	}
}

func (sub *Subscribe) RegisterChannel(channel interface{}, f func(ctx context.Context, data []byte)) {
	sub.chFuncLc.Lock()
	defer sub.chFuncLc.Unlock()
	if _, ok := sub.chFunc[channel]; !ok {
		sub.chList = append(sub.chList, channel)
	}
	sub.chFunc[channel] = f
}

// ListenAndServe 启动订阅流程
func (sub *Subscribe) ListenAndServe() (err error) {
	defer func() {
		_ = sub.ShutDown(context.Background())
	}()

	// ping
	_, err = redis.String(sub.conn.Do("PING"))
	if err != nil {
		return
	}

	err = sub.newSub()
	if err != nil {
		return
	}

	// 阻塞监听
	for {
		select {
		case <-sub.ctx.Done():
			return
		case v, ok := <-sub.dataCh:
			if ok == false {
				continue
			}
			if v.err != nil {
				err = v.err
				return
			}
			sub.chFuncLc.RLock()
			f, ok := sub.chFunc[v.channel]
			sub.chFuncLc.RUnlock()
			if ok == false {
				continue
			}
			go func() {
				defer func() {
					if r := recover(); r != nil {
						fmt.Println(r)
					}
				}()
				f(sub.ctx, v.data)
			}()
		}
	}
}

func (sub *Subscribe) newSub() error {
	sub.pubSubConn = redis.PubSubConn{
		Conn: sub.conn,
	}
	if len(sub.chList) == 0 {
		return errors.New("empty subscription channel")
	}

	sub.chFuncLc.RLock()
	err := sub.pubSubConn.Subscribe(sub.chList...)
	sub.chFuncLc.RUnlock()
	if err != nil {
		return err
	}

	go func() {
		for {
			select {
			case <-sub.closeCtx.Done():
				return
			default:
				switch v := sub.pubSubConn.Receive().(type) {
				case redis.Message:
					sub.dataCh <- struct {
						channel interface{}
						data    []byte
						err     error
					}{channel: v.Channel, data: v.Data, err: nil}
				case redis.Subscription:
				case error:
					sub.dataCh <- struct {
						channel interface{}
						data    []byte
						err     error
					}{channel: nil, data: nil, err: v}
					return
				default:
					time.Sleep(time.Millisecond)
				}
			}
		}
	}()

	return nil
}

func (sub *Subscribe) ShutDown(ctx context.Context) error {
	timer := time.NewTimer(500 * time.Millisecond)
	defer timer.Stop()
	sub.closeCancel()
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-timer.C:
			return nil
		}
	}
}
