package fizzbuzz

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFizzBuzz_EdgeCases(t *testing.T) {
	t.Run("Empty", func(t *testing.T) {
		result, err := FizzBuzz(0, 0, 3, 5)
		assert.Nil(t, err)
		assert.Len(t, result, 1)
	})
	t.Run("Negative", func(t *testing.T) {
		result, err := FizzBuzz(-20, 0, 3, 5)
		assert.Nil(t, err)
		assert.Len(t, result, 20)
	})
	t.Run("Zero Fizz", func(t *testing.T) {
		result, err := FizzBuzz(0, 20, 0, 5)
		assert.Nil(t, err)
		assert.Len(t, result, 20)
	})
	t.Run("Zero Buzz", func(t *testing.T) {
		result, err := FizzBuzz(0, 20, 3, 0)
		assert.Nil(t, err)
		assert.Len(t, result, 20)
	})
}
