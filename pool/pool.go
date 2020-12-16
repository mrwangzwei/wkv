package pool

import (
	"errors"
	"sync"
)

type pool struct {
	poolSize int
	taskChan chan *task
	running  int
	lc       sync.Mutex
}

var pw sync.WaitGroup

func NewPool(size int) (*pool, error) {
	if size <= 0 {
		return nil, errors.New("invalid pool size")
	}
	return &pool{
		poolSize: size,
		taskChan: make(chan *task),
		lc:       sync.Mutex{},
	}, nil
}

//任务写入
func (p *pool) AddTask(t *task) error {
	//协程数未满时再开新的协程
	if p.running < p.poolSize {
		p.incr()
		go p.worker()
	}

	p.taskChan <- t
	return nil
}

func (p *pool) incr() {
	defer p.lc.Unlock()
	p.lc.Lock()
	p.running++
}

func (p *pool) decr() {
	defer p.lc.Unlock()
	p.lc.Lock()
	p.running--
}

func (p *pool) worker() {
	defer p.decr()
	for {
		select {
		case w, ok := <-p.taskChan:
			if ok != true {
				return
			}
			pw.Add(1)
			w.execute()
		}
	}
}

func (p *pool) NewTask(f HandleF, params interface{}, ph PHandleF) *task {
	return &task{
		params:  params,
		taskId:  newTaskId(),
		handle:  f,
		pHandle: ph,
	}
}

func (p *pool) Wait() {
	pw.Wait()
}
