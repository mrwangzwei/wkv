package queue

import (
	"errors"
	"sync"
)

type arrQueue struct {
	cont      [][]byte
	headArrow int
	tailArrow int
	freeSize  uint
	useSize   uint
	lc        sync.RWMutex
}

func NewArrQueue(size uint) *arrQueue {
	return &arrQueue{
		cont:     make([][]byte, size),
		freeSize: size,
	}
}

func (q *arrQueue) Push(cont []byte) error {
	q.lc.Lock()
	defer q.lc.Unlock()
	if q.freeSize == 0 {
		return errors.New("queue is full")
	}
	q.freeSize--
	q.useSize++
	q.cont[q.tailArrow] = cont
	q.tailArrow++
	if int(q.tailArrow) >= len(q.cont) {
		q.tailArrow = 0
	}
	return nil
}

//从头拿
func (q *arrQueue) Pull() ([]byte, error) {
	q.lc.Lock()
	defer q.lc.Unlock()
	if q.useSize == 0 {
		return nil, errors.New("empty queue")
	}
	q.freeSize++
	q.useSize--
	cont := q.cont[q.headArrow]
	q.cont[q.headArrow] = nil
	q.headArrow++
	if q.headArrow >= len(q.cont) {
		q.headArrow = 0
	}
	return cont, nil
}

//从尾拿
func (q *arrQueue) RPull() ([]byte, error) {
	q.lc.Lock()
	defer q.lc.Unlock()
	if q.useSize == 0 {
		return nil, errors.New("empty queue")
	}
	q.freeSize++
	q.useSize--
	arrow := q.tailArrow
	if arrow == 0 {
		arrow = len(q.cont) - 1
	} else {
		arrow--
	}
	cont := q.cont[arrow]
	q.cont[arrow] = nil
	q.tailArrow = arrow
	return cont, nil
}

func (q *arrQueue) Len() int {
	q.lc.RLock()
	defer q.lc.RUnlock()
	return int(q.useSize)
}
