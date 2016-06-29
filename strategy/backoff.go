// Package strategy provides a way to change the way that retry is performed.
//
// Copyright © 2016 Trevor N. Suarez (Rican7)
package strategy

import (
	"math"
	"time"
)

// BackoffModifier (TODO: name?!) defines a function that calculates a
// time.Duration based on a given retry attempt number.
type BackoffModifier func(duration time.Duration, attempt uint) time.Duration

// Backoff creates a Strategy that waits before each attempt, with an increasing
// duration.
func Backoff(initial time.Duration, modifier BackoffModifier) Strategy {
	return func(attempt uint) bool {
		time.Sleep(modifier(initial, attempt))

		return true
	}
}

// Incremental TODO
func Incremental(increment time.Duration) BackoffModifier {
	return func(duration time.Duration, attempt uint) time.Duration {
		return duration + (increment * time.Duration(attempt))
	}
}

// Linear TODO
func Linear() BackoffModifier {
	return func(duration time.Duration, attempt uint) time.Duration {
		return (duration * time.Duration(attempt))
	}
}

// Exponential TODO
func Exponential() BackoffModifier {
	return func(duration time.Duration, attempt uint) time.Duration {
		return time.Duration(math.Pow(float64(duration), float64(attempt)))
	}
}
