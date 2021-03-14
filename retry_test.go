package retry

import (
	"context"
	"errors"
	"testing"
	"time"
)

// timeMarginOfError represents the acceptable amount of time that may pass for
// a time-based (sleep) unit before considering invalid.
const timeMarginOfError = time.Millisecond

func TestRetry(t *testing.T) {
	action := func(attempt uint) error {
		return nil
	}

	err := Retry(action)

	if nil != err {
		t.Error("expected a nil error")
	}
}

func TestRetryRetriesUntilNoErrorReturned(t *testing.T) {
	const errorUntilAttemptNumber = 5

	var attemptsMade uint

	action := func(attempt uint) error {
		attemptsMade = attempt

		if errorUntilAttemptNumber == attempt {
			return nil
		}

		return errors.New("erroring")
	}

	err := Retry(action)

	if nil != err {
		t.Error("expected a nil error")
	}

	if errorUntilAttemptNumber != attemptsMade {
		t.Errorf(
			"expected %d attempts to be made, but %d were made instead",
			errorUntilAttemptNumber,
			attemptsMade,
		)
	}
}

func TestRetryWithContextChecksContextAfterLastAttempt(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	strategy := func(attempt uint, sleep func(time.Duration)) bool {
		if attempt == 0 {
			return true
		}

		cancel()
		return false
	}

	action := func(ctx context.Context, attempt uint) error {
		return errors.New("erroring")
	}

	err := RetryWithContext(ctx, action, strategy)

	if context.Canceled != err {
		t.Error("expected a context error")
	}
}

func TestRetryWithContextCancelStopsAttempts(t *testing.T) {
	var numCalls int

	ctx, cancel := context.WithCancel(context.Background())

	action := func(ctx context.Context, attempt uint) error {
		numCalls++

		if numCalls == 1 {
			cancel()
			return ctx.Err()
		}

		return nil
	}

	err := RetryWithContext(ctx, action)

	if 1 != numCalls {
		t.Errorf("expected the action to be tried once, not %d times", numCalls)
	}

	if context.Canceled != err {
		t.Error("expected a context error")
	}
}

func TestRetryWithContextSleepIsInterrupted(t *testing.T) {
	const sleepDuration = 100 * timeMarginOfError
	fullySleptBy := time.Now().Add(sleepDuration)

	strategy := func(attempt uint, sleep func(time.Duration)) bool {
		if attempt > 0 {
			sleep(sleepDuration)
		}
		return attempt <= 1
	}

	var numCalls int

	action := func(ctx context.Context, attempt uint) error {
		numCalls++
		return errors.New("erroring")
	}

	stopAfter := 10 * timeMarginOfError
	deadline := time.Now().Add(stopAfter)
	ctx, _ := context.WithDeadline(context.Background(), deadline)

	err := RetryWithContext(ctx, action, strategy)

	if time.Now().Before(deadline) {
		t.Errorf("expected to stop after %s", stopAfter)
	}

	if time.Now().After(fullySleptBy) {
		t.Errorf("expected to stop before %s", sleepDuration)
	}

	if 1 != numCalls {
		t.Errorf("expected the action to be tried once, not %d times", numCalls)
	}

	if context.DeadlineExceeded != err {
		t.Error("expected a context error")
	}
}

func TestShouldAttempt(t *testing.T) {
	shouldAttempt := shouldAttempt(1, time.Sleep)

	if !shouldAttempt {
		t.Error("expected to return true")
	}
}

func TestShouldAttemptWithStrategy(t *testing.T) {
	const attemptNumberShouldReturnFalse = 7

	strategy := func(attempt uint, sleep func(time.Duration)) bool {
		return (attemptNumberShouldReturnFalse != attempt)
	}

	should := shouldAttempt(1, time.Sleep, strategy)

	if !should {
		t.Error("expected to return true")
	}

	should = shouldAttempt(1+attemptNumberShouldReturnFalse, time.Sleep, strategy)

	if !should {
		t.Error("expected to return true")
	}

	should = shouldAttempt(attemptNumberShouldReturnFalse, time.Sleep, strategy)

	if should {
		t.Error("expected to return false")
	}
}

func TestShouldAttemptWithMultipleStrategies(t *testing.T) {
	trueStrategy := func(attempt uint, sleep func(time.Duration)) bool {
		return true
	}

	falseStrategy := func(attempt uint, sleep func(time.Duration)) bool {
		return false
	}

	should := shouldAttempt(1, time.Sleep, trueStrategy)

	if !should {
		t.Error("expected to return true")
	}

	should = shouldAttempt(1, time.Sleep, falseStrategy)

	if should {
		t.Error("expected to return false")
	}

	should = shouldAttempt(1, time.Sleep, trueStrategy, trueStrategy, trueStrategy)

	if !should {
		t.Error("expected to return true")
	}

	should = shouldAttempt(1, time.Sleep, falseStrategy, falseStrategy, falseStrategy)

	if should {
		t.Error("expected to return false")
	}

	should = shouldAttempt(1, time.Sleep, trueStrategy, trueStrategy, falseStrategy)

	if should {
		t.Error("expected to return false")
	}
}
