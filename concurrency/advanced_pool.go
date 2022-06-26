package concurrency

import (
	"context"
	"errors"
	"log"
	"sync"
	"sync/atomic"
)

type RunnableContext func(context.Context)

// ErrPoolClosed is returned from AdvancedPool.Submit when the pool is closed
// before submission can be sent.
var ErrPoolClosed = errors.New("pool closed")

type CancelablePool interface {
	// Submit submits the given task to the pool, blocking until a slot becomes
	// available or the context is closed. The given context and its lifetime only
	// affects this function and is not the context passed to the callback. If the
	// context is closed before a slot becomes available, the context error is
	// returned. If the pool is closed before a slot becomes available,
	// ErrPoolClosed is returned. Otherwise the task is submitted to the pool and
	// no error is returned. The context passed to the callback will be closed
	// when the pool is closed.
	Submit(context.Context, RunnableContext) error

	// Close closes the pool and waits until all submitted tasks have completed
	// before returning. If the pool is already closed, ErrPoolClosed is returned.
	// If the given context is closed before all tasks have finished, the context
	// error is returned. Otherwise, no error is returned.
	Close(context.Context) error
}

// AdvancedPool is a more advanced worker pool that supports cancelling the
// submission and closing the pool. All functions are safe to call from multiple
// goroutines.
type AdvancedPool struct {
	ctx      context.Context
	cancel   context.CancelFunc
	taskChan chan RunnableContext
	runners  sync.WaitGroup
	open     atomic.Value
}

// NewAdvancedPool creates a new AdvancedPool. maxSlots is the maximum total
// submitted tasks, running or waiting, that can be submitted before Submit
// blocks waiting for more room. maxConcurrent is the maximum tasks that can be
// running at any one time. An error is returned if maxSlots is less than
// maxConcurrent or if either value is not greater than zero.
func NewAdvancedPool(maxSlots, maxConcurrent int) (*AdvancedPool, error) {
	return NewAdvancedPoolContext(context.Background(), maxSlots, maxConcurrent)
}

func NewAdvancedPoolContext(parentCtx context.Context, maxSlots, maxConcurrent int) (*AdvancedPool, error) {
	if maxConcurrent <= 0 {
		return nil, errors.New("maxConcurrent must be at least 1")
	}
	if maxSlots < maxConcurrent {
		return nil, errors.New("maxSlots must me greater than maxConcurrent")
	}

	ctx, cancel := context.WithCancel(parentCtx)
	pool := &AdvancedPool{
		ctx:      ctx,
		cancel:   cancel,
		taskChan: make(chan RunnableContext, maxSlots),
	}

	pool.runners.Add(maxConcurrent)
	for i := 0; i < maxConcurrent; i++ {
		log.Printf("Starting pool runner %d\n", i)
		go pool.run()
	}
	pool.open.Store(true)

	return pool, nil
}

func (p *AdvancedPool) Submit(ctx context.Context, task RunnableContext) error {
	if task == nil {
		return errors.New("invalid nil task submitted")
	}

	select {
	case p.taskChan <- task:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	case <-p.ctx.Done():
		return ErrPoolClosed
	}
}

func (p *AdvancedPool) IsOpen() bool {
	return p.open.Load().(bool)
}

func (p *AdvancedPool) Close(ctx context.Context) error {
	if !p.open.CompareAndSwap(true, false) {
		return ErrPoolClosed
	}

	// Cancel our context and close our task submission channel
	p.cancel()
	close(p.taskChan)

	c := make(chan struct{})
	go func() {
		defer close(c)
		p.runners.Wait()
	}()
	select {
	case <-c:
		// All runners have stopped and finished running their tasks
		return nil
	case <-ctx.Done():
		// Pool is closed, but some tasks may be outstanding
		return ctx.Err()
	}
}

func (p *AdvancedPool) run() {
	defer p.runners.Done()
	for {
		select {
		case task, ok := <-p.taskChan:
			if !ok {
				// Closed channel
				return
			}
			// We got a task run it if not nil
			if task != nil {
				task(p.ctx)
			}
		case <-p.ctx.Done():
			// Our context has been canceled, stop the runner
			return
		}
	}
}
