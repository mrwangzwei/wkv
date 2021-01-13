package limit_bucket

import (
	"errors"
	"sync"
	"time"
)

const IncrLimit = 10000000000 //每次token注入的最大数量

type TokenBucket struct {
	freq        int //每秒频限
	nSize       int //当前还有的数量
	preInjectAt int //上次注入的时间。ms
	lc          *sync.RWMutex
}

func NewTokenBucket(freq int) *TokenBucket {
	return &TokenBucket{
		freq:  freq,
		nSize: 0,
		lc:    new(sync.RWMutex),
	}
}

func (b *TokenBucket) GetToken() (token int, err error) {
	b.inject()

	b.lc.Lock()
	defer b.lc.Unlock()
	if b.nSize < 1 {
		err = errors.New("empty token")
		return
	}
	token = b.nSize
	b.nSize--
	return
}

func (b *TokenBucket) inject() {
	b.lc.RLock()
	gap := b.freq - b.nSize
	pre := b.preInjectAt
	freq := b.freq
	b.lc.RUnlock()
	if gap < 1 {
		return
	}
	amount := countAmount(pre, freq)
	if amount < 1 {
		return
	}

	b.doInject(amount)
}

func countAmount(pre, freq int) (amount int) {
	if pre == 0 {
		amount = freq
		return
	}
	now := int(time.Now().UnixNano() / 1e6)
	timeGap := now - pre
	if freq >= 1000 {
		//按毫秒级频限计算
		freq = int(freq / 1000)
	} else {
		timeGap = int((now - pre) / 1000)
	}

	if timeGap < 1 {
		return
	}
	if timeGap > int(IncrLimit/freq) {
		amount = IncrLimit
	} else {
		amount = timeGap * freq
	}
	return
}

func (b *TokenBucket) doInject(amount int) {
	b.lc.Lock()
	defer b.lc.Unlock()
	gap := b.freq - b.nSize
	if gap <= 0 {
		return
	}
	if amount < gap {
		gap = amount
	}
	b.nSize += gap
	b.preInjectAt = int(time.Now().UnixNano() / 1e6)
	return
}
