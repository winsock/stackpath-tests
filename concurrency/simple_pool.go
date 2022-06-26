// Package concurrency implements worker pool interfaces, one simple and one a
// bit more complex.
package concurrency

import (
	"errors"
	"log"
)

type Runnable func()

type Pool interface {
	// Submit a task to be executed asynchronously. This function will return as
	// soon as the task is submitted. If the pool does not have an available slot
	// for the task, this blocks until it can submit.
	Submit(runnable Runnable)
}

// SimplePool is a simple worker pool that does not support cancellation or
// closing. All functions are safe to call from multiple goroutines.
type SimplePool struct {
	taskChan chan Runnable
}

// NewSimplePool creates a new SimplePool that only allows the given maximum
// concurrent tasks to run at any one time. maxConcurrent must be greater than
// zero.
func NewSimplePool(maxConcurrent int) (*SimplePool, error) {
	if maxConcurrent <= 0 {
		return nil, errors.New("maxConcurrent must be at least 1")
	}

	pool := &SimplePool{
		taskChan: make(chan Runnable),
	}

	// Start the number of runners
	for i := 0; i < maxConcurrent; i++ {
		log.Printf("Starting pool runner %d\n", i)
		go pool.run()
	}

	return pool, nil
}

func (p *SimplePool) Submit(task Runnable) {
	if task == nil {
		// Would normally return as an error, but the interface spec does not allow for this
		log.Println("Invalid nil task submitted")
		return
	}

	p.taskChan <- task
}

func (p *SimplePool) Close() {
	close(p.taskChan)
}

func (p *SimplePool) run() {
	for task := range p.taskChan {
		if task != nil {
			task()
		}
	}
}
