package fizzbuzz

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFizzBuzz_EdgeCases(t *testing.T) {
	t.Run("Empty", func(t *testing.T) {
		result := FizzBuzz(0, 3, 5)
		assert.Len(t, result, 0)
	})
	t.Run("Zero Fizz", func(t *testing.T) {
		result := FizzBuzz(20, 0, 5)
		assert.Len(t, result, 20)
	})
	t.Run("Zero Buzz", func(t *testing.T) {
		result := FizzBuzz(20, 3, 0)
		assert.Len(t, result, 20)
	})
}
