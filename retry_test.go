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

	if nil != err {
		t.Error("expected a nil error")
	}
}

func TestRetryRetriesUntilNoErrorReturned(t *testing.T) {
	const errorUntilAttemptNumber = 5;

	var attemptsMade uint

	action := func(attempt uint) error {
		attemptsMade = attempt

		if (errorUntilAttemptNumber == attempt) {
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
