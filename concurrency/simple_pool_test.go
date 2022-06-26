package concurrency

import (
	"context"
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
	"time"
)

func TestSimplePool(t *testing.T) {
	t.Run("Single Task", func(t *testing.T) {
		var waitGroup sync.WaitGroup

		pool, err := NewSimplePool(1)
		assert.Nil(t, err)
		assert.NotNil(t, pool)

		waitGroup.Add(1)
		pool.Submit(func() {
			waitGroup.Done()
		})

		err = waitForWaitGroup(&waitGroup, time.Minute)
		assert.Nil(t, err)
	})
	t.Run("Waiting Task", func(t *testing.T) {
		var waitGroup sync.WaitGroup

		pool, err := NewSimplePool(1)
		assert.Nil(t, err)
		assert.NotNil(t, pool)

		waitGroup.Add(2)
		pool.Submit(func() {
			time.Sleep(5 * time.Second)
			waitGroup.Done()
		})
		pool.Submit(func() {
			time.Sleep(5 * time.Second)
			waitGroup.Done()
		})

		err = waitForWaitGroup(&waitGroup, 11*time.Second)
		assert.Nil(t, err)
	})
	t.Run("Concurrent Tasks", func(t *testing.T) {
		var waitGroup sync.WaitGroup

		pool, err := NewSimplePool(2)
		assert.Nil(t, err)
		assert.NotNil(t, pool)

		waitGroup.Add(2)
		pool.Submit(func() {
			time.Sleep(5 * time.Second)
			waitGroup.Done()
		})
		pool.Submit(func() {
			time.Sleep(5 * time.Second)
			waitGroup.Done()
		})

		err = waitForWaitGroup(&waitGroup, 6*time.Second)
		assert.Nil(t, err)
	})
	t.Run("Invalid Concurrency", func(t *testing.T) {
		pool, err := NewSimplePool(0)
		assert.NotNil(t, err)
		assert.Nil(t, pool)
	})
}

func waitForWaitGroup(wg *sync.WaitGroup, timeout time.Duration) error {
	c := make(chan struct{})
	go func() {
		defer close(c)
		wg.Wait()
	}()
	select {
	case <-c:
		return nil
	case <-time.After(timeout):
		return context.DeadlineExceeded
	}
}
