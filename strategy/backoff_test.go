package strategy

import (
	"math"
	"testing"
	"time"
)

func TestBackoff(t *testing.T) {
	const backoffDuration = time.Duration(10 * time.Millisecond)

	strategy := Backoff(backoffDuration, Linear())

	if now := time.Now(); !strategy(0) || 0 != time.Since(now) {
		t.Error("strategy expected to return true in 0 time")
	}

	if now := time.Now(); !strategy(1) || backoffDuration > time.Since(now) {
		t.Errorf(
			"strategy expected to return true in %s",
			time.Duration(backoffDuration),
		)
	}

	if now := time.Now(); !strategy(5) || (5*backoffDuration) > time.Since(now) {
		t.Errorf(
			"strategy expected to return true in %s",
			time.Duration(5*backoffDuration),
		)
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
