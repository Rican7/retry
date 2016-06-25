package retry

import (
	"time"
)

// Limit creates a Strategy that limits the number of attempts that Retry will
// make.
func Limit(attemptLimit uint) Strategy {
	return func(attempt uint) bool {
		return (attempt < attemptLimit)
	}
}

// Delay creates a Strategy that waits the given duration on the first attempt.
func Delay(duration time.Duration) Strategy {
	return func(attempt uint) bool {
		if 0 == attempt {
			time.Sleep(duration)
		}

		return true
	}
}

// Wait creates a Strategy that waits the given durations for each attempt after
// the first. If the number of attempts is greater than the number of durations
// provided, then the strategy uses the last duration provided.
func Wait(durations ...time.Duration) Strategy {
	return func(attempt uint) bool {
		if 0 < attempt {
			durationIndex := int(attempt - 1)

			if len(durations) >= durationIndex {
				durationIndex = len(durations) - 1
			}

			time.Sleep(durations[durationIndex])
		}

		return true
	}
}

// Backoff creates a Strategy that waits an increasing duration.
func Backoff(duration time.Duration) Strategy {
	return func(attempt uint) bool {
		time.Sleep(duration * time.Duration(attempt))

		return true
	}
}
