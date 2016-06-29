package strategy

import (
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
