package fizzbuzz

import (
	"github.com/stretchr/testify/assert"
	"math"
	"testing"
)

func TestFizzBuzz_EdgeCases(t *testing.T) {
	t.Run("Zero to Zero", func(t *testing.T) {
		result, err := FizzBuzz(0, 0, 3, 5)
		assert.Nil(t, err)
		assert.Len(t, result, 1)
		assert.Equal(t, result[0], Fizz+Buzz)
	})
	t.Run("Negative", func(t *testing.T) {
		result, err := FizzBuzz(-20, -1, 3, 5)
		assert.Nil(t, err)
		assert.Len(t, result, 20)
	})
	t.Run("Zero Fizz", func(t *testing.T) {
		result, err := FizzBuzz(1, 20, 0, 5)
		assert.Nil(t, err)
		assert.Len(t, result, 20)
	})
	t.Run("Zero Buzz", func(t *testing.T) {
		result, err := FizzBuzz(1, 20, 3, 0)
		assert.Nil(t, err)
		assert.Len(t, result, 20)
	})
}

func TestNumberDivisibleBy(t *testing.T) {
	t.Run("Divide by Zero", func(t *testing.T) {
		result := numberDivisibleBy(0, 0)
		assert.False(t, result)
	})
	t.Run("Divisible by Itself", func(t *testing.T) {
		for i := int64(-20); i <= int64(20); i++ {
			if i == 0 {
				// Were testing zero separately
				continue
			}
			result := numberDivisibleBy(i, i)
			assert.True(t, result)
		}
	})
	t.Run("MinInt64", func(t *testing.T) {
		result := numberDivisibleBy(math.MinInt64, math.MinInt64)
		assert.True(t, result)
	})
	t.Run("MaxInt64", func(t *testing.T) {
		result := numberDivisibleBy(math.MaxInt64, math.MaxInt64)
		assert.True(t, result)
	})
}
