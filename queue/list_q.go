package queue

import (
	"errors"
	"sync"
)

type listQueue struct {
	l        int //当前队列长度
	size     int
	headNode *listNode
	tailNode *listNode
	lc       sync.RWMutex
}

type listNode struct {
	cont []byte //[]byte：为了使用的时候方便点，但占用内存比string多一倍，可能原因是slice结构体元素多了一个cap
	next *listNode
}

func NewListQueue(size int) (Queue, error) {
	if size < -1 {
		return nil, errors.New("invalid size")
	}
	return &listQueue{size: size}, nil
}

func (q *listQueue) Pull() (cont []byte, err error) {
	if q.headNode == nil {
		err = errors.New("empty queue")
		return
	}
	cont, q.headNode = q.headNode.pull()
	if q.headNode == nil {
		q.tailNode = nil
	}
	q.lc.Lock()
	q.l--
	q.lc.Unlock()
	return
}

func (q *listQueue) Push(cont []byte) (err error) {
	if q.l == q.size {
		err = errors.New("over max size")
		return
	}
	if q.headNode == nil { //队列的第一条
		fir := &listNode{cont: cont}
		q.headNode = fir
		q.tailNode = fir
	} else {
		q.tailNode = q.tailNode.push(cont)
	}
	q.lc.Lock()
	q.l++
	q.lc.Unlock()
	return
}

func (q *listQueue) Len() int {
	q.lc.RLock()
	defer q.lc.RUnlock()
	return q.l
}

func (n *listNode) push(cont []byte) *listNode {
	if n.next == nil {
		n.next = &listNode{cont: cont}
		return n.next
	}
	return n.next.push(cont)
}

func (n *listNode) pull() ([]byte, *listNode) {
	return n.cont, n.next
}
