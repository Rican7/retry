package strategy

import (
	"fmt"
	"math"
	"testing"
	"time"
)

func TestBackoff(t *testing.T) {
	const backoffDuration = time.Duration(10 * time.Millisecond)
	const algorithmDurationBase = time.Millisecond

	algorithm := func(duration time.Duration, attempt uint) time.Duration {
		return duration - (algorithmDurationBase * time.Duration(attempt))
	}

	strategy := Backoff(backoffDuration, algorithm)

	if now := time.Now(); !strategy(0) || 0 != time.Since(now) {
		t.Error("strategy expected to return true in 0 time")
	}

	for i := uint(1); i < 10; i++ {
		expectedResult := backoffDuration - (algorithmDurationBase * time.Duration(i))

		if now := time.Now(); !strategy(i) || expectedResult > time.Since(now) {
			t.Errorf(
				"strategy expected to return true in %s",
				expectedResult,
			)
		}
	}
}

func TestIncremental(t *testing.T) {
	const increment = time.Nanosecond

	algorithm := Incremental(increment)

	duration := time.Millisecond

	for i := uint(0); i < 10; i++ {
		result := algorithm(duration, i)
		expected := duration + (increment * time.Duration(i))

		if result != expected {
			t.Errorf("algorithm expected to return a %s duration, but received %s instead", expected, result)
		}
	}
}

func TestLinear(t *testing.T) {
	algorithm := Linear()

	duration := time.Millisecond

	for i := uint(0); i < 10; i++ {
		result := algorithm(duration, i)
		expected := duration * time.Duration(i)

		if result != expected {
			t.Errorf("algorithm expected to return a %s duration, but received %s instead", expected, result)
		}
	}
}

func TestExponential(t *testing.T) {
	const base = 2

	algorithm := Exponential(base)

	duration := time.Second

	for i := uint(0); i < 10; i++ {
		result := algorithm(duration, i)
		expected := duration * time.Duration(math.Pow(base, float64(i)))

		if result != expected {
			t.Errorf("algorithm expected to return a %s duration, but received %s instead", expected, result)
		}
	}
}

func TestFibonacci(t *testing.T) {
	algorithm := Fibonacci()

	duration := time.Millisecond

	for i := uint(0); i < 10; i++ {
		result := algorithm(duration, i)
		expected := duration * time.Duration(fibonacciNumber(i))

		if result != expected {
			t.Errorf("algorithm expected to return a %s duration, but received %s instead", expected, result)
		}
	}
}

func TestFibonacciNumber(t *testing.T) {
	// Fibonacci sequence
	expectedSequence := []uint{0, 1, 1, 2, 3, 5, 8, 13, 21, 34, 55, 89, 144, 233}

	for i, expected := range expectedSequence {
		result := fibonacciNumber(uint(i))

		if result != expected {
			t.Errorf("fibonacci %d number expected %d, but got %d", i, expected, result)
		}
	}
}

func ExampleIncremental() {
	algorithm := Incremental(10 * time.Millisecond)

	for i := uint(1); i <= 5; i++ {
		duration := algorithm(15*time.Millisecond, i)

		fmt.Printf("#%d attempt: %s\n", i, duration)
		// Output:
		// #1 attempt: 25ms
		// #2 attempt: 35ms
		// #3 attempt: 45ms
		// #4 attempt: 55ms
		// #5 attempt: 65ms
	}
}

func ExampleLinear() {
	algorithm := Linear()

	for i := uint(1); i <= 5; i++ {
		duration := algorithm(15*time.Millisecond, i)

		fmt.Printf("#%d attempt: %s\n", i, duration)
		// Output:
		// #1 attempt: 15ms
		// #2 attempt: 30ms
		// #3 attempt: 45ms
		// #4 attempt: 60ms
		// #5 attempt: 75ms
	}
}

func ExampleExponential() {
	algorithm := Exponential(2)

	for i := uint(1); i <= 5; i++ {
		duration := algorithm(15*time.Millisecond, i)

		fmt.Printf("#%d attempt: %s\n", i, duration)
		// Output:
		// #1 attempt: 30ms
		// #2 attempt: 60ms
		// #3 attempt: 120ms
		// #4 attempt: 240ms
		// #5 attempt: 480ms
	}
}

func ExampleFibonacci() {
	algorithm := Fibonacci()

	for i := uint(1); i <= 5; i++ {
		duration := algorithm(15*time.Millisecond, i)

		fmt.Printf("#%d attempt: %s\n", i, duration)
		// Output:
		// #1 attempt: 15ms
		// #2 attempt: 15ms
		// #3 attempt: 30ms
		// #4 attempt: 45ms
		// #5 attempt: 75ms
	}
}
