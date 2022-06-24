package main

import (
	"flag"
	"fmt"
	"github.com/stackpath/backend-developer-tests/unit-testing/pkg/fizzbuzz"
)

var (
	// The default argument values and will contain the real values after calling `flag.Parse`
	startValue = int64(1)
	endValue   = int64(20)
	fizzAt     = int64(3)
	buzzAt     = int64(5)
)

func main() {
	fmt.Println("SP// Backend Developer Test - FizzBuzz")
	fmt.Println()

	// Read the provided flags
	flag.Parse()

	// FizzBuzz the input
	fmt.Println(fmt.Sprintf("FizzBuzzing from %d to %d, fizzing at %d and buzzing at %d:", startValue, endValue, fizzAt, buzzAt))
	results, err := fizzbuzz.FizzBuzz(startValue, endValue, fizzAt, buzzAt)
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		return
	}

	// Print the results
	for _, result := range results {
		fmt.Println(result)
	}

	fmt.Println()
	fmt.Println("Done")
	fmt.Println()
}

func init() {
	// Configure the flags for this FizzBuzz application
	flag.Int64Var(&startValue, "start", startValue, "What number to start FizzBuzzing at, must be less than the end")
	flag.Int64Var(&endValue, "end", endValue, "What number to end FizzBuzzing at, must be greater than the start")
	flag.Int64Var(&fizzAt, "fizz", fizzAt, "The divisor to determine when to Fizz")
	flag.Int64Var(&buzzAt, "buzz", buzzAt, "The divisor to determine when to Buzz")
}
