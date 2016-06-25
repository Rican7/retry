package retry

import (
	"time"
)

// Limit creates a Strategy that limits the number of attempts that Retry will
// make
func Limit(attemptLimit uint) Strategy {
	return func(attempt uint) bool {
		return (attempt < attemptLimit)
	}
}

// Delay creates a Strategy that waits the given duration on the first attempt
func Delay(duration time.Duration) Strategy {
	return func(attempt uint) bool {
		if 0 == attempt {
			time.Sleep(duration)
		}

		return true
	}
}

// Wait creates a Strategy that waits the given duration after the first attempt
func Wait(duration time.Duration) Strategy {
	return func(attempt uint) bool {
		if 0 < attempt {
			time.Sleep(duration)
		}

		return true
	}
}

// Backoff creates a Strategy that waits an increasing duration
func Backoff(duration time.Duration) Strategy {
	return func(attempt uint) bool {
		time.Sleep(duration * time.Duration(attempt))

		return true
	}
}
