package redigo_sub

import (
	"context"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"io"
	"testing"
	"time"
)

func TestSub(t *testing.T) {
	var sub *Subscribe
	go func() {
		for i := 0; i < 5; i++ {
			conn, err := redis.Dial("tcp", "127.0.0.1:6379")
			if err == nil {
				i = 0
				sub, err = Subbbb(context.Background(), conn)
				if err == nil || (err != nil && err != io.EOF) {
					return
				}
			}
			fmt.Println(i, "retry")
			time.Sleep(5 * time.Second)
		}
	}()
	time.Sleep(3 * time.Second)
	if sub != nil {
		fmt.Println("ShutDown", sub.ShutDown(context.Background()))
	}

	time.Sleep(2 * time.Second)

}

func Subbbb(ctx context.Context, conn redis.Conn) (*Subscribe, error) {
	sub := NewSub(ctx, conn)
	sub.RegisterChannel("aaa", func(ctx context.Context, data []byte) {
		fmt.Println(1, string(data))
	})
	sub.RegisterChannel("bbb", func(ctx context.Context, data []byte) {
		fmt.Println(2, string(data))
	})
	sub.RegisterChannel("ccc", func(ctx context.Context, data []byte) {
		fmt.Println(3, string(data))
	})
	go func() {
		err := sub.ListenAndServe()
		if err != nil {
			return
		}
	}()
	return sub, nil
}
