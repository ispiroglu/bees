package bees

import (
	"log"
	"sync"
)

type IPool interface {
	Start()
	Stop()
	AddWork(ITask)
}

type Pool struct {
	workerCount int
	tasks       chan ITask

	start sync.Once
	stop  sync.Once

	close chan struct{}
}

func (p *Pool) AddWork(t ITask) {
	select {
	case p.tasks <- t:
	case <-p.close:
	}
}

func (p *Pool) AddWorkAsync(t ITask) {
	go p.AddWork(t)
}

func (p *Pool) Start() {
	p.start.Do(func() {
		p.startWorkers()
	})
}

func (p *Pool) Stop() {
	p.stop.Do(func() {
		close(p.close)
	})
}

func NewPool(wc int, cs int) (IPool, error) {
	if wc <= 0 {
		return nil, ERROR_CREATING_POOL_NO_WORKER
	}

	if cs < 0 {
		return nil, ERROR_CREATING_POOL_NEGATIVE_CHANNEL
	}

	tc := make(chan ITask, cs)

	return &Pool{
		workerCount: wc,
		tasks:       tc,
		start:       sync.Once{},
		stop:        sync.Once{},
		close:       make(chan struct{}),
	}, nil
}

func (p *Pool) startWorkers() {
	for i := 0; i < p.workerCount; i++ {
		go p.workerStarter(i)
	}
}

func (p *Pool) workerStarter(wn int) {
	for {
		select {
		case <-p.close:
			log.Printf("stopping worker %d with quit channel\n", wn)
			return
		case t, ok := <-p.tasks:
			if !ok {
				log.Printf("stopping worker %d with closed tasks channel\n", wn)
				return
			}
			if err := t.Execute(); err != nil {
				t.OnFailure(err)
			}
		}
	}
}
