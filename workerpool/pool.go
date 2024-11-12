package workerpool

import (
	"sync"
	"time"
)

// Response contain response object as well as error if any
type Response[T any] struct {
	Res T
	Err error
}

// unit of Work to be provided
type WorkFunc[T any] func() (T, error)

type WorkerPool[T any] struct {
	ch    chan WorkFunc[T]
	ResCh chan *Response[T]
	stats map[int]*Stats
	wg    sync.WaitGroup
	count int
	mu    sync.Mutex
}

type Stats struct {
	TaskCount int // Total tasks completed
	UpTime    time.Duration
}

func New[T any](count int) *WorkerPool[T] {
	return &WorkerPool[T]{
		count: count,
		ch:    make(chan WorkFunc[T], 1),
		ResCh: make(chan *Response[T], 1),
	}
}

// Run is main method on WorkerPool type which will spin up workerpool
// and wait for all the task to complete please make sure this is blocking method
func (wp *WorkerPool[T]) Run() {
	for i := 1; i <= wp.count; i++ {
		wp.wg.Add(1)
		go wp.worker(i)
	}
	wp.wg.Wait()
}

// Close closes the workerpool
func (wp *WorkerPool[T]) Close() {
	close(wp.ch)
	close(wp.ResCh)
}

func (wp *WorkerPool[T]) worker(id int) {
	t := time.Now()
	defer wp.wg.Done()
	for fn := range wp.ch {
		r, err := fn()
		wp.ResCh <- &Response[T]{r, err}
		wp.addStat(id, time.Since(t))
	}
}

func (wp *WorkerPool[T]) addStat(id int, uptime time.Duration) {
	wp.mu.Lock()
	if s, ok := wp.stats[id]; ok {
		s.UpTime = uptime
		s.TaskCount++
	} else {
		wp.stats[id] = &Stats{TaskCount: 1, UpTime: uptime}
	}
	wp.mu.Unlock()
}

func (wp *WorkerPool[T]) Stats() map[int]Stats {
	m := make(map[int]Stats, wp.count)
	wp.mu.Lock()
	for k, v := range wp.stats {
		m[k] = *v
	}
	wp.mu.Unlock()
	return m
}
