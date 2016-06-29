package retry

import (
	"testing"
	"time"
)

func TestLimit(t *testing.T) {
	const attemptLimit = 3

	strategy := Limit(attemptLimit)

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
	const delayDuration = time.Duration(10 * time.Millisecond)

	strategy := Delay(delayDuration)

	if now := time.Now(); !strategy(0) || delayDuration > time.Since(now) {
		t.Errorf(
			"strategy expected to return true in %s",
			time.Duration(delayDuration),
		)
	}

	if now := time.Now(); !strategy(5) || 0 != time.Since(now) {
		t.Error("strategy expected to return true in 0 time")
	}
}

func TestWait(t *testing.T) {
	const waitDuration = time.Duration(10 * time.Millisecond)

	strategy := Wait(waitDuration)

	if now := time.Now(); !strategy(0) || 0 != time.Since(now) {
		t.Error("strategy expected to return true in 0 time")
	}

	if now := time.Now(); !strategy(1) || waitDuration > time.Since(now) {
		t.Errorf(
			"strategy expected to return true in %s",
			time.Duration(waitDuration),
		)
	}
}

func TestWaitWithMultipleDurations(t *testing.T) {
	waitDurations := []time.Duration{
		time.Duration(10 * time.Millisecond),
		time.Duration(20 * time.Millisecond),
		time.Duration(30 * time.Millisecond),
		time.Duration(40 * time.Millisecond),
	}

	strategy := Wait(waitDurations...)

	if now := time.Now(); !strategy(0) || 0 != time.Since(now) {
		t.Error("strategy expected to return true in 0 time")
	}

	if now := time.Now(); !strategy(1) || waitDurations[0] > time.Since(now) {
		t.Errorf(
			"strategy expected to return true in %s",
			time.Duration(waitDurations[0]),
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
