package concurrency

import (
	"context"
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
	"time"
)

func TestAdvancedPool(t *testing.T) {
	t.Run("Single Task", func(t *testing.T) {
		var waitGroup sync.WaitGroup
		testCtx := context.Background()
		pool, err := NewAdvancedPoolContext(testCtx, 1, 1)
		assert.Nil(t, err)
		assert.NotNil(t, pool)

		waitGroup.Add(1)
		err = pool.Submit(testCtx, func(_ context.Context) {
			waitGroup.Done()
		})
		assert.Nil(t, err)

		err = waitForWaitGroup(&waitGroup, time.Minute)
		assert.Nil(t, err)

		err = pool.Close(testCtx)
		assert.Nil(t, err)
	})
	t.Run("Test Close Waiting For Task", func(t *testing.T) {
		var waitGroup sync.WaitGroup
		testCtx := context.Background()
		pool, err := NewAdvancedPoolContext(testCtx, 1, 1)
		assert.Nil(t, err)
		assert.NotNil(t, pool)

		waitGroup.Add(1)
		err = pool.Submit(testCtx, func(_ context.Context) {
			time.Sleep(5 * time.Second)
			waitGroup.Done()
		})
		assert.Nil(t, err)

		ctx, cancel := context.WithTimeout(testCtx, 6*time.Second)
		err = pool.Close(ctx)
		cancel()
		assert.Nil(t, err)

		err = waitForWaitGroup(&waitGroup, time.Minute)
		assert.Nil(t, err)
	})
	t.Run("Test Double Close", func(t *testing.T) {
		testCtx := context.Background()
		pool, err := NewAdvancedPoolContext(testCtx, 1, 1)
		assert.Nil(t, err)
		assert.NotNil(t, pool)

		err = pool.Close(testCtx)
		assert.Nil(t, err)
		err = pool.Close(testCtx)
		assert.ErrorIs(t, err, ErrPoolClosed)
	})
	t.Run("Test Submit After Close", func(t *testing.T) {
		testCtx := context.Background()
		pool, err := NewAdvancedPoolContext(testCtx, 1, 1)
		assert.Nil(t, err)
		assert.NotNil(t, pool)

		err = pool.Close(testCtx)
		assert.Nil(t, err)

		err = pool.Submit(testCtx, func(_ context.Context) {
			assert.Fail(t, "Should not be hit!")
		})
		assert.ErrorIs(t, err, ErrPoolClosed)
	})
	t.Run("Test Submit Timeout", func(t *testing.T) {
		var waitGroup sync.WaitGroup
		testCtx := context.Background()
		pool, err := NewAdvancedPoolContext(testCtx, 1, 1)
		assert.Nil(t, err)
		assert.NotNil(t, pool)

		waitGroup.Add(2)
		// Start one task
		err = pool.Submit(testCtx, func(_ context.Context) {
			time.Sleep(10 * time.Second)
			waitGroup.Done()
		})
		assert.Nil(t, err)
		// Put another into the buffered channel
		err = pool.Submit(testCtx, func(_ context.Context) {
			waitGroup.Done()
		})
		assert.Nil(t, err)

		// Try adding another task, but the buffer is full and the context will timeout
		ctx, cancel := context.WithTimeout(testCtx, time.Second)
		err = pool.Submit(ctx, func(_ context.Context) {
			assert.Fail(t, "Should not be hit!")
		})
		cancel()
		assert.ErrorIs(t, err, context.DeadlineExceeded)

		err = pool.Close(testCtx)
		assert.Nil(t, err)

		err = waitForWaitGroup(&waitGroup, time.Minute)
		assert.Nil(t, err)
	})
	t.Run("Test Cancel timeout before chanel buffer drained", func(t *testing.T) {
		var waitGroup sync.WaitGroup
		testCtx := context.Background()
		pool, err := NewAdvancedPoolContext(testCtx, 1, 1)
		assert.Nil(t, err)
		assert.NotNil(t, pool)

		waitGroup.Add(1)
		// Start one task
		err = pool.Submit(testCtx, func(_ context.Context) {
			time.Sleep(5 * time.Second)
			waitGroup.Done()
		})
		assert.Nil(t, err)
		// Put another into the buffered channel that will be lost due to the Close context deadline being reached
		err = pool.Submit(testCtx, func(_ context.Context) {
			assert.Fail(t, "Should not be hit!")
		})
		assert.Nil(t, err)

		// Cancel, but with a timeout that would be hit before submitted tasks are finished
		ctx, cancel := context.WithTimeout(testCtx, time.Microsecond)
		err = pool.Close(ctx)
		cancel()
		assert.ErrorIs(t, err, context.DeadlineExceeded)

		err = waitForWaitGroup(&waitGroup, time.Minute)
		assert.Nil(t, err)
	})
	t.Run("Test Cancel while Submit is blocking", func(t *testing.T) {
		var waitGroup sync.WaitGroup
		testCtx := context.Background()
		pool, err := NewAdvancedPoolContext(testCtx, 1, 1)
		assert.Nil(t, err)
		assert.NotNil(t, pool)

		// Start one task
		err = pool.Submit(testCtx, func(_ context.Context) {
			time.Sleep(10 * time.Second)
		})
		assert.Nil(t, err)
		// Put another into the buffered channel that will be lost due to the Close context being reached
		err = pool.Submit(testCtx, func(_ context.Context) {
			time.Sleep(10 * time.Second)
		})
		assert.Nil(t, err)

		// Submit a new task in a new coroutine while we close the pool in this one
		waitGroup.Add(1)
		go func() {
			err := pool.Submit(context.Background(), func(_ context.Context) {
				assert.Fail(t, "Should not be hit!")
			})
			assert.ErrorIs(t, err, ErrPoolClosed)
			waitGroup.Done()
		}()
		time.Sleep(time.Millisecond)

		// Cancel, but with a timeout that would be hit before submitted tasks are finished
		ctx, cancel := context.WithTimeout(testCtx, time.Second)
		err = pool.Close(ctx)
		cancel()
		assert.ErrorIs(t, err, context.DeadlineExceeded)

		err = waitForWaitGroup(&waitGroup, time.Minute)
		assert.Nil(t, err)
	})
	t.Run("Invalid Task", func(t *testing.T) {
		testCtx := context.Background()
		pool, err := NewAdvancedPoolContext(testCtx, 1, 1)
		assert.Nil(t, err)
		assert.NotNil(t, pool)

		err = pool.Submit(testCtx, nil)
		assert.NotNil(t, err)

		err = pool.Close(testCtx)
		assert.Nil(t, err)
	})
	t.Run("Invalid maxSlots", func(t *testing.T) {
		testCtx := context.Background()
		pool, err := NewAdvancedPoolContext(testCtx, 0, 1)
		assert.NotNil(t, err)
		assert.Nil(t, pool)
	})
	t.Run("Invalid maxConcurrent", func(t *testing.T) {
		testCtx := context.Background()
		pool, err := NewAdvancedPoolContext(testCtx, 1, 0)
		assert.NotNil(t, err)
		assert.Nil(t, pool)
	})
}
