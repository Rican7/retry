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
type BackoffAlgorithm func(factor time.Duration, attempt uint) time.Duration

// Backoff creates a Strategy that waits before each attempt, with a duration as
// defined by the given factor duration and BackoffAlgorithm.
func Backoff(factor time.Duration, algorithm BackoffAlgorithm) Strategy {
	return func(attempt uint) bool {
		if 0 < attempt {
			time.Sleep(algorithm(factor, attempt))
		}

		return true
	}
}

// Incremental creates a BackoffAlgorithm that increments the factor duration
// by the given increment for each attempt.
func Incremental(increment time.Duration) BackoffAlgorithm {
	return func(factor time.Duration, attempt uint) time.Duration {
		return factor + (increment * time.Duration(attempt))
	}
}

// Linear creates a BackoffAlgorithm that linearly multiplies the factor
// duration by the attempt number for each attempt.
func Linear() BackoffAlgorithm {
	return func(factor time.Duration, attempt uint) time.Duration {
		return (factor * time.Duration(attempt))
	}
}

// Exponential creates a BackoffAlgorithm that multiplies the factor duration
// by an exponentially increasing factor for each attempt, where the factor is
// calculated as the given base raised to the attempt number.
func Exponential(base float64) BackoffAlgorithm {
	return func(factor time.Duration, attempt uint) time.Duration {
		return (factor * time.Duration(math.Pow(base, float64(attempt))))
	}
}
