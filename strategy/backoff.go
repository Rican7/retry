// Package strategy provides a way to change the way that retry is performed.
//
// Copyright © 2016 Trevor N. Suarez (Rican7)
package strategy

import (
	"math"
	"time"
)

// BackoffAlgorithm defines a function that calculates a time.Duration based on
// a given duration and retry attempt number.
type BackoffAlgorithm func(initial time.Duration, attempt uint) time.Duration

// Backoff creates a Strategy that waits before each attempt, with an increasing
// duration.
func Backoff(initial time.Duration, algorithm BackoffAlgorithm) Strategy {
	return func(attempt uint) bool {
		if 0 < attempt {
			time.Sleep(algorithm(initial, attempt))
		}

		return true
	}
}

// Incremental TODO
func Incremental(increment time.Duration) BackoffAlgorithm {
	return func(initial time.Duration, attempt uint) time.Duration {
		return initial + (increment * time.Duration(attempt))
	}
}

// Linear TODO
func Linear() BackoffAlgorithm {
	return func(initial time.Duration, attempt uint) time.Duration {
		return (initial * time.Duration(attempt))
	}
}

// Exponential TODO
func Exponential(base float64) BackoffAlgorithm {
	return func(initial time.Duration, attempt uint) time.Duration {
		return (initial * time.Duration(math.Pow(base, float64(attempt))))
	}
}
