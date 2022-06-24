package fizzbuzz

import (
	"errors"
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
func FizzBuzz(start, end, fizzAt, buzzAt int64) ([]string, error) {
	if start > end {
		return make([]string, 0), errors.New("start cannot be after the end")
	}

	total := end - start + 1
	result := make([]string, total)
	for i := int64(0); i < total; i++ {
		number := start + i
		isFizz := numberDivisibleBy(number, fizzAt)
		isBuzz := numberDivisibleBy(number, buzzAt)

		if !isFizz && !isBuzz {
			result[i] = strconv.FormatInt(number, 10)
			continue
		}

		if isFizz {
			result[i] = Fizz
		}

		if isBuzz {
			result[i] += Buzz
		}
	}

	return result, nil
}

// numberDivisibleBy returns true if the number i is divisible by the provided divisor
// If divisor is 0 this function returns false
func numberDivisibleBy(i, divisor int64) bool {
	if divisor == 0 {
		return false
	}
	return i%divisor == 0
}
