package retry

import (
	"testing"
)

func TestAttemptLimiter(t *testing.T) {
	const attemptLimit = 3;

	strategy := AttemptLimiter(attemptLimit)

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
