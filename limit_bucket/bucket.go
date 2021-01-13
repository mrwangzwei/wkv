package limit_bucket

import (
	"errors"
	"sync"
	"time"
)

type bucket struct {
	size     int
	duration int   //纳秒级别
	c        []int //每次访问的时间
	lc       *sync.RWMutex
}

func NewBucket(limitSize, durationNano int) *bucket {
	return &bucket{
		size:     limitSize,
		duration: durationNano,
		c:        make([]int, limitSize),
		lc:       new(sync.RWMutex),
	}
}

func (b *bucket) Inject() error {
	nowTime := int(time.Now().UnixNano())
	enableDura := nowTime - b.duration

	b.lc.RLock()
	tail := b.c[b.size-1]
	b.lc.RUnlock()

	if tail != 0 && tail >= enableDura {
		return errors.New("over limit")
	}

	//没超的话,移位置
	b.lc.Lock()
	b.c = append([]int{nowTime}, b.c[:b.size-1]...)
	b.lc.Unlock()
	return nil
}
