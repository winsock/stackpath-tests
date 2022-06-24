package fizzbuzz

import (
	"github.com/stretchr/testify/assert"
	"math"
	"testing"
)

func TestFizzBuzz(t *testing.T) {
	t.Run("Normal", func(t *testing.T) {
		result, err := FizzBuzz(1, 20, 3, 5)
		assert.Nil(t, err)
		assert.Len(t, result, 20)
		assert.Equal(t, []string{"1", "2", Fizz, "4", Buzz, Fizz, "7", "8", Fizz, Buzz, "11", Fizz, "13", "14", Fizz + Buzz, "16", "17", Fizz, "19", Buzz}, result)
	})
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
		assert.Equal(t, []string{Buzz, "-19", Fizz, "-17", "-16", Fizz + Buzz, "-14", "-13", Fizz, "-11", Buzz, Fizz, "-8", "-7", Fizz, Buzz, "-4", Fizz, "-2", "-1"}, result)
	})
	t.Run("Zero Fizz", func(t *testing.T) {
		result, err := FizzBuzz(1, 20, 0, 5)
		assert.Nil(t, err)
		assert.Len(t, result, 20)
		assert.Equal(t, []string{"1", "2", "3", "4", Buzz, "6", "7", "8", "9", Buzz, "11", "12", "13", "14", Buzz, "16", "17", "18", "19", Buzz}, result)
	})
	t.Run("Zero Buzz", func(t *testing.T) {
		result, err := FizzBuzz(1, 20, 3, 0)
		assert.Nil(t, err)
		assert.Len(t, result, 20)
		assert.Equal(t, []string{"1", "2", Fizz, "4", "5", Fizz, "7", "8", Fizz, "10", "11", Fizz, "13", "14", Fizz, "16", "17", Fizz, "19", "20"}, result)
	})
	t.Run("Negative Fizz", func(t *testing.T) {
		result, err := FizzBuzz(1, 20, -3, 5)
		assert.Nil(t, err)
		assert.Len(t, result, 20)
		assert.Equal(t, []string{"1", "2", Fizz, "4", Buzz, Fizz, "7", "8", Fizz, Buzz, "11", Fizz, "13", "14", Fizz + Buzz, "16", "17", Fizz, "19", Buzz}, result)
	})
	t.Run("Negative Buzz", func(t *testing.T) {
		result, err := FizzBuzz(1, 20, 3, -5)
		assert.Nil(t, err)
		assert.Len(t, result, 20)
		assert.Equal(t, []string{"1", "2", Fizz, "4", Buzz, Fizz, "7", "8", Fizz, Buzz, "11", Fizz, "13", "14", Fizz + Buzz, "16", "17", Fizz, "19", Buzz}, result)
	})
	t.Run("Invalid Interval", func(t *testing.T) {
		result, err := FizzBuzz(20, 0, 3, 5)
		assert.NotNil(t, err)
		assert.Len(t, result, 0)
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
