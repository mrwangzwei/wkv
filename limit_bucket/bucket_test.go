package limit_bucket

import (
	"fmt"
	"testing"
	"time"
)

func TestTimeBucket(t *testing.T) {
	bu := NewBucket(90, int(time.Second))
	go func() {
		i := 0
		for {
			i++
			fmt.Println(i)
			err := bu.Inject()
			if err != nil {
				fmt.Println(err)
			}
			time.Sleep(time.Second / 100)
		}
	}()
	time.Sleep(30 * time.Second)
}

func TestTokenBucket(t *testing.T) {
	b := NewTokenBucket(10)
	go func() {
		i := 0
		for {
			i++
			fmt.Println(i)
			token, err := b.GetToken()
			if err != nil {
				fmt.Println("======", err)
			} else {
				fmt.Println("token", token)
			}
			time.Sleep(time.Second / 100)
		}
	}()
	time.Sleep(30 * time.Second)
}
