package strategy

import (
	"testing"
	"time"
)

// timeMarginOfError represents the acceptable amount of time that may pass for
// a time-based (sleep) unit before considering invalid.
const timeMarginOfError = time.Millisecond

func TestLimit(t *testing.T) {
	// Strategy attempts are 0-based.
	// Treat this functionally as n+1.
	const attemptLimit = 3

	strategy := Limit(attemptLimit)

	if !strategy(0) {
		t.Error("strategy expected to return true")
	}

	if !strategy(1) {
		t.Error("strategy expected to return true")
	}

	if !strategy(2) {
		t.Error("strategy expected to return true")
	}

	if strategy(3) {
		t.Error("strategy expected to return false")
	}
}

func TestDelay(t *testing.T) {
	const delayDuration = 10 * timeMarginOfError

	strategy := Delay(delayDuration)

	if now := time.Now(); !strategy(0) || delayDuration > time.Since(now) {
		t.Errorf(
			"strategy expected to return true in %s",
			delayDuration,
		)
	}

	if now := time.Now(); !strategy(5) || (delayDuration/10) < time.Since(now) {
		t.Error("strategy expected to return true in ~0 time")
	}
}

func TestWait(t *testing.T) {
	strategy := Wait()

	if now := time.Now(); !strategy(0) || timeMarginOfError < time.Since(now) {
		t.Error("strategy expected to return true in ~0 time")
	}

	if now := time.Now(); !strategy(999) || timeMarginOfError < time.Since(now) {
		t.Error("strategy expected to return true in ~0 time")
	}
}

func TestWaitWithDuration(t *testing.T) {
	const waitDuration = 10 * timeMarginOfError

	strategy := Wait(waitDuration)

	if now := time.Now(); !strategy(0) || timeMarginOfError < time.Since(now) {
		t.Error("strategy expected to return true in ~0 time")
	}

	if now := time.Now(); !strategy(1) || waitDuration > time.Since(now) {
		t.Errorf(
			"strategy expected to return true in %s",
			waitDuration,
		)
	}
}

func TestWaitWithMultipleDurations(t *testing.T) {
	waitDurations := []time.Duration{
		10 * timeMarginOfError,
		20 * timeMarginOfError,
		30 * timeMarginOfError,
		40 * timeMarginOfError,
	}

	strategy := Wait(waitDurations...)

	if now := time.Now(); !strategy(0) || timeMarginOfError < time.Since(now) {
		t.Error("strategy expected to return true in ~0 time")
	}

	if now := time.Now(); !strategy(1) || waitDurations[0] > time.Since(now) {
		t.Errorf(
			"strategy expected to return true in %s",
			waitDurations[0],
		)
	}

	if now := time.Now(); !strategy(3) || waitDurations[2] > time.Since(now) {
		t.Errorf(
			"strategy expected to return true in %s",
			waitDurations[2],
		)
	}

	if now := time.Now(); !strategy(999) || waitDurations[len(waitDurations)-1] > time.Since(now) {
		t.Errorf(
			"strategy expected to return true in %s",
			waitDurations[len(waitDurations)-1],
		)
	}
}

func TestBackoff(t *testing.T) {
	const testCycles = 10
	const backoffDuration = testCycles * timeMarginOfError
	const algorithmDurationBase = timeMarginOfError

	algorithm := func(attempt uint) time.Duration {
		return backoffDuration - (time.Duration(attempt) * algorithmDurationBase)
	}

	strategy := Backoff(algorithm)

	if now := time.Now(); !strategy(0) || timeMarginOfError < time.Since(now) {
		t.Error("strategy expected to return true in ~0 time")
	}

	for i := uint(1); i < testCycles; i++ {
		expectedResult := algorithm(i)

		if expectedResult < 0 {
			t.Errorf(
				"algorithm returned a negative duration %s",
				expectedResult,
			)
		}

		if now := time.Now(); !strategy(i) || expectedResult > time.Since(now) {
			t.Errorf(
				"strategy expected to return true in %s",
				expectedResult,
			)
		}
	}
}

func TestBackoffWithJitter(t *testing.T) {
	const testCycles = 10
	const backoffDuration = 2 * testCycles * timeMarginOfError
	const algorithmDurationBase = timeMarginOfError

	algorithm := func(attempt uint) time.Duration {
		return backoffDuration - (time.Duration(attempt) * algorithmDurationBase)
	}

	transformation := func(duration time.Duration) time.Duration {
		return duration - (backoffDuration / 2)
	}

	strategy := BackoffWithJitter(algorithm, transformation)

	if now := time.Now(); !strategy(0) || timeMarginOfError < time.Since(now) {
		t.Error("strategy expected to return true in ~0 time")
	}

	for i := uint(1); i < testCycles; i++ {
		expectedResult := transformation(algorithm(i))

		if expectedResult < 0 {
			t.Errorf(
				"transformation returned a negative duration %s",
				expectedResult,
			)
		}

		if now := time.Now(); !strategy(i) || expectedResult > time.Since(now) {
			t.Errorf(
				"strategy expected to return true in %s",
				expectedResult,
			)
		}
	}
}

func TestNoJitter(t *testing.T) {
	transformation := noJitter()

	for i := uint(0); i < 10; i++ {
		duration := time.Duration(i) * timeMarginOfError
		result := transformation(duration)
		expected := duration

		if result != expected {
			t.Errorf("transformation expected to return a %s duration, but received %s instead", expected, result)
		}
	}
}
