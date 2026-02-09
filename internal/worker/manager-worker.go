package worker

import (
	"log"
	"sync"
)

type Task[T any] interface {
	Process(workerId int, data T)
}

type PoolManager[T any] struct {
	mu           sync.Mutex
	currentCount int
	jobChan      chan T
	quitChan     chan struct{}
	processor    Task[T]
	poolName     string
}

func NewPoolManager[T any](name string, ch chan T, proc Task[T]) *PoolManager[T] {
	return &PoolManager[T]{
		poolName:  name,
		jobChan:   ch,
		quitChan:  make(chan struct{}),
		processor: proc,
	}
}

func (p *PoolManager[T]) Adjust(target int) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if target > p.currentCount {
		for i := 0; i < target-p.currentCount; i++ {
			go p.startWorker(p.currentCount + i + 1)
		}
	} else if target < p.currentCount {
		for i := 0; i < p.currentCount-target; i++ {
			p.quitChan <- struct{}{}
		}
	}
	p.currentCount = target
	log.Printf("[%s] Pool adjusted to %d workers", p.poolName, target)
}

func (p *PoolManager[T]) startWorker(id int) {
	for {
		select {
		case data, ok := <-p.jobChan:
			if !ok {
				return
			}
			p.processor.Process(id, data)
		case <-p.quitChan:
			return
		}
	}
}

func (p *PoolManager[T]) GetWorkerCount() int {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.currentCount
}
