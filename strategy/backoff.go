// Package strategy provides a way to change the way that retry is performed.
//
// Copyright © 2016 Trevor N. Suarez (Rican7)
package strategy

import (
	"math"
	"time"
)

// BackoffAlgorithm defines a function that calculates a time.Duration based on
// the given retry attempt number.
type BackoffAlgorithm func(attempt uint) time.Duration

// Backoff creates a Strategy that waits before each attempt, with a duration as
// defined by the given BackoffAlgorithm.
func Backoff(algorithm BackoffAlgorithm) Strategy {
	return func(attempt uint) bool {
		if 0 < attempt {
			time.Sleep(algorithm(attempt))
		}

		return true
	}
}

// Incremental creates a BackoffAlgorithm that increments the initial duration
// by the given increment for each attempt.
func Incremental(initial, increment time.Duration) BackoffAlgorithm {
	return func(attempt uint) time.Duration {
		return initial + (increment * time.Duration(attempt))
	}
}

// Linear creates a BackoffAlgorithm that linearly multiplies the factor
// duration by the attempt number for each attempt.
func Linear(factor time.Duration) BackoffAlgorithm {
	return func(attempt uint) time.Duration {
		return (factor * time.Duration(attempt))
	}
}

// Exponential creates a BackoffAlgorithm that multiplies the factor duration by
// an exponentially increasing factor for each attempt, where the factor is
// calculated as the given base raised to the attempt number.
func Exponential(factor time.Duration, base float64) BackoffAlgorithm {
	return func(attempt uint) time.Duration {
		return (factor * time.Duration(math.Pow(base, float64(attempt))))
	}
}

// BinaryExponential creates a BackoffAlgorithm that multiplies the factor
// duration by an exponentially increasing factor for each attempt, where the
// factor is calculated as `2` raised to the attempt number (2^attempt).
func BinaryExponential(factor time.Duration) BackoffAlgorithm {
	return Exponential(factor, 2)
}

// Fibonacci creates a BackoffAlgorithm that multiplies the factor duration by
// an increasing factor for each attempt, where the factor is the Nth number in
// the Fibonacci sequence.
func Fibonacci(factor time.Duration) BackoffAlgorithm {
	return func(attempt uint) time.Duration {
		return (factor * time.Duration(fibonacciNumber(attempt)))
	}
}

// fibonacciNumber calculates the Fibonacci sequence number for the given
// sequence position.
func fibonacciNumber(n uint) uint {
	if 0 == n {
		return 0
	} else if 1 == n {
		return 1
	} else {
		return fibonacciNumber(n-1) + fibonacciNumber(n-2)
	}
}
