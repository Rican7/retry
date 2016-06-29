package strategy

import (
	"math"
	"testing"
	"time"
)

func TestBackoff(t *testing.T) {
	const backoffDuration = time.Duration(10 * time.Millisecond)
	const modifierDurationBase = time.Millisecond

	modifier := func(duration time.Duration, attempt uint) time.Duration {
		return duration - (modifierDurationBase * time.Duration(attempt))
	}

	strategy := Backoff(backoffDuration, modifier)

	if now := time.Now(); !strategy(0) || 0 != time.Since(now) {
		t.Error("strategy expected to return true in 0 time")
	}

	for i := uint(1); i < 10; i++ {
		expectedResult := backoffDuration - (modifierDurationBase * time.Duration(i))

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

	modifier := Incremental(increment)

	duration := time.Millisecond

	for i := uint(0); i < 10; i++ {
		result := modifier(duration, i)
		expected := duration + (increment * time.Duration(i))

		if result != expected {
			t.Errorf("modifier expected to return a %s duration, but received %s instead", expected, result)
		}
	}
}

func TestLinear(t *testing.T) {
	modifier := Linear()

	duration := time.Millisecond

	for i := uint(0); i < 10; i++ {
		result := modifier(duration, i)
		expected := duration * time.Duration(i)

		if result != expected {
			t.Errorf("modifier expected to return a %s duration, but received %s instead", expected, result)
		}
	}
}

func TestExponential(t *testing.T) {
	const base = 2

	modifier := Exponential(base)

	duration := time.Second

	for i := uint(0); i < 10; i++ {
		result := modifier(duration, i)
		expected := duration * time.Duration(math.Pow(base, float64(i)))

		if result != expected {
			t.Errorf("modifier expected to return a %s duration, but received %s instead", expected, result)
		}
	}
}
