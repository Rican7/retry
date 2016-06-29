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
	const increment = time.Duration(1 * time.Nanosecond)

	modifier := Incremental(increment)

	duration := time.Duration(1 * time.Millisecond)
	result := modifier(duration, 3)
	expected := (increment * 3) + duration

	if result != expected {
		t.Errorf("modifier expected to return a %s duration, but received %s instead", expected, result)
	}
}

func TestLinear(t *testing.T) {
	modifier := Linear()

	duration := time.Duration(1 * time.Millisecond)
	result := modifier(duration, 3)
	expected := 3 * duration

	if result != expected {
		t.Errorf("modifier expected to return a %s duration, but received %s instead", expected, result)
	}
}

func TestExponential(t *testing.T) {
	modifier := Exponential()

	duration := time.Duration(1 * time.Millisecond)
	result := modifier(duration, 3)
	expected := time.Duration(math.Pow(float64(duration), 3))

	if result != expected {
		t.Errorf("modifier expected to return a %s duration, but received %s instead", expected, result)
	}
}
