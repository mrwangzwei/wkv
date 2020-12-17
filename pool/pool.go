package pool

import (
	"errors"
	"sync"
)

type Pool struct {
	poolSize int
	taskChan chan *task
	//waitRunChan chan *task		//考虑加个待运行缓冲，taskChan就要有大小限制，在pool close的时候，还没有到taskChan的就直接丢弃了，已经在的就先执行完再close
	running int
	lc      *sync.Mutex
	wg      *sync.WaitGroup
}

func NewPool(size int) (*Pool, error) {
	if size <= 0 {
		return nil, errors.New("invalid pool size")
	}
	return &Pool{
		poolSize: size,
		taskChan: make(chan *task),
		lc:       &sync.Mutex{},
		wg:       &sync.WaitGroup{},
	}, nil
}

//任务写入
func (p *Pool) AddTask(t *task) {
	//协程数未满时再开新的协程
	if p.running < p.poolSize {
		p.incr()
		go p.worker()
	}

	p.taskChan <- t
}

func (p *Pool) incr() {
	defer p.lc.Unlock()
	p.lc.Lock()
	p.running++
}

func (p *Pool) decr() {
	defer p.lc.Unlock()
	p.lc.Lock()
	p.running--
}

func (p *Pool) worker() {
	defer p.decr()
	for {
		select {
		case w, ok := <-p.taskChan:
			if ok != true {
				return
			}
			p.wg.Add(1)
			w.execute(p)
		}
	}
}

func (p *Pool) NewTask(f handleF, params interface{}, ph pHandleF) *task {
	return &task{
		params: params,
		h:      f,
		ph:     ph,
	}
}

func (p *Pool) Wait() {
	p.wg.Wait()
}
