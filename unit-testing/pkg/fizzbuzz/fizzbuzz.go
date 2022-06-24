package fizzbuzz

import (
	"strconv"
)

const (
	Fizz = "Fizz"
	Buzz = "Buzz"
)

// FizzBuzz performs a FizzBuzz operation over a range of integers
//
// Given a range of integers:
// - Return "Fizz" if the integer is divisible by the `fizzAt` value.
// - Return "Buzz" if the integer is divisible by the `buzzAt` value.
// - Return "FizzBuzz" if the integer is divisible by both the `fizzAt` and
//   `buzzAt` values.
// - Return the original number if it is not divisible by either the `fizzAt` or
//   the `buzzAt` values.
func FizzBuzz(total, fizzAt, buzzAt int64) []string {
	result := make([]string, total)

	for i := int64(1); i <= total; i++ {
		isFizz := numberDivisibleBy(i, fizzAt)
		isBuzz := numberDivisibleBy(i, buzzAt)

		if !isFizz && !isBuzz {
			result[i-1] = strconv.FormatInt(i, 10)
			continue
		}

		if isFizz {
			result[i-1] = Fizz
		}

		if isBuzz {
			result[i-1] += Buzz
		}
	}

	return result
}

// numberDivisibleBy returns true if the number i is divisible by the provided divisor
// If divisor is 0 this function returns false
func numberDivisibleBy(i, divisor int64) bool {
	if divisor == 0 {
		return false
	}
	return i%divisor == 0
}
