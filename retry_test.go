package retry

import (
	"errors"
	"testing"
)

func TestRetry(t *testing.T) {
	action := func(attempt uint) error {
		return nil
	}

	err := Retry(action)

	if err != nil {
		t.Error("expected a nil error")
	}
}

func TestRetryAttemptNumberIsAccurate(t *testing.T) {
	var strategyAttemptNumber uint
	var actionAttemptNumber uint

	strategy := func(attempt uint) bool {
		strategyAttemptNumber = attempt

		return true
	}

	action := func(attempt uint) error {
		actionAttemptNumber = attempt

		return nil
	}

	err := Retry(action, strategy)

	if err != nil {
		t.Error("expected a nil error")
	}

	if strategyAttemptNumber != 0 {
		t.Errorf(
			"expected strategy to receive 0, received %v instead",
			strategyAttemptNumber,
		)
	}

	if actionAttemptNumber != 1 {
		t.Errorf(
			"expected action to receive 1, received %v instead",
			actionAttemptNumber,
		)
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

	if err != nil {
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

func TestShouldAttempt(t *testing.T) {
	shouldAttempt := shouldAttempt(1)

	if !shouldAttempt {
		t.Error("expected to return true")
	}
}

func TestShouldAttemptWithStrategy(t *testing.T) {
	const attemptNumberShouldReturnFalse = 7

	strategy := func(attempt uint) bool {
		return (attemptNumberShouldReturnFalse != attempt)
	}

	should := shouldAttempt(1, strategy)

	if !should {
		t.Error("expected to return true")
	}

	should = shouldAttempt(1+attemptNumberShouldReturnFalse, strategy)

	if !should {
		t.Error("expected to return true")
	}

	should = shouldAttempt(attemptNumberShouldReturnFalse, strategy)

	if should {
		t.Error("expected to return false")
	}
}

func TestShouldAttemptWithMultipleStrategies(t *testing.T) {
	trueStrategy := func(attempt uint) bool {
		return true
	}

	falseStrategy := func(attempt uint) bool {
		return false
	}

	should := shouldAttempt(1, trueStrategy)

	if !should {
		t.Error("expected to return true")
	}

	should = shouldAttempt(1, falseStrategy)

	if should {
		t.Error("expected to return false")
	}

	should = shouldAttempt(1, trueStrategy, trueStrategy, trueStrategy)

	if !should {
		t.Error("expected to return true")
	}

	should = shouldAttempt(1, falseStrategy, falseStrategy, falseStrategy)

	if should {
		t.Error("expected to return false")
	}

	should = shouldAttempt(1, trueStrategy, trueStrategy, falseStrategy)

	if should {
		t.Error("expected to return false")
	}
}
