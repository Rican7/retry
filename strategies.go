package retry

import (
	"math"
	"time"
)

// BackoffModifier (TODO: name?!) defines a function that calculates a
// time.Duration based on a given retry attempt number.
type BackoffModifier func(duration time.Duration, attempt uint) time.Duration

// Limit creates a Strategy that limits the number of attempts that Retry will
// make.
func Limit(attemptLimit uint) Strategy {
	return func(attempt uint) bool {
		return (attempt < attemptLimit)
	}
}

// Delay creates a Strategy that waits the given duration before the first
// attempt is made.
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

			if len(durations) <= durationIndex {
				durationIndex = len(durations) - 1
			}

			time.Sleep(durations[durationIndex])
		}

		return true
	}
}

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
