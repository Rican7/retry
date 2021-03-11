package strategy

import (
	"testing"
	"time"
)

func TestLimit(t *testing.T) {
	const attemptLimit = 3

	strategy := Limit(attemptLimit)

	if !strategy(1, time.Sleep) {
		t.Error("strategy expected to return true")
	}

	if !strategy(2, time.Sleep) {
		t.Error("strategy expected to return true")
	}

	if !strategy(3, time.Sleep) {
		t.Error("strategy expected to return true")
	}

	if strategy(4, time.Sleep) {
		t.Error("strategy expected to return false")
	}
}

func TestDelay(t *testing.T) {
	const delayDuration = time.Duration(10)

	strategy := Delay(delayDuration)

	if spy, actual := sleepSpy(); !strategy(0, spy) || delayDuration != *actual {
		t.Errorf("strategy expected to return true in %s", delayDuration)
	}

	if spy, actual := sleepSpy(); !strategy(5, spy) || 0 != *actual {
		t.Error("strategy expected to return true in ~0 time")
	}
}

func TestWait(t *testing.T) {
	strategy := Wait()

	if spy, actual := sleepSpy(); !strategy(0, spy) || 0 != *actual {
		t.Error("strategy expected to return true in ~0 time")
	}

	if spy, actual := sleepSpy(); !strategy(999, spy) || 0 != *actual {
		t.Error("strategy expected to return true in ~0 time")
	}
}

func TestWaitWithDuration(t *testing.T) {
	const waitDuration = time.Duration(10)

	strategy := Wait(waitDuration)

	if spy, actual := sleepSpy(); !strategy(0, spy) || 0 != *actual {
		t.Error("strategy expected to return true in ~0 time")
	}

	if spy, actual := sleepSpy(); !strategy(1, spy) || waitDuration != *actual {
		t.Errorf("strategy expected to return true in %s", waitDuration)
	}
}

func TestWaitWithMultipleDurations(t *testing.T) {
	waitDurations := []time.Duration{
		time.Duration(10),
		time.Duration(20),
		time.Duration(30),
		time.Duration(40),
	}

	strategy := Wait(waitDurations...)

	if spy, actual := sleepSpy(); !strategy(0, spy) || 0 != *actual {
		t.Error("strategy expected to return true in ~0 time")
	}

	if spy, actual := sleepSpy(); !strategy(1, spy) || waitDurations[0] != *actual {
		t.Errorf("strategy expected to return true in %s", waitDurations[0])
	}

	if spy, actual := sleepSpy(); !strategy(3, spy) || waitDurations[2] != *actual {
		t.Errorf("strategy expected to return true in %s", waitDurations[2])
	}

	if spy, actual := sleepSpy(); !strategy(999, spy) || waitDurations[len(waitDurations)-1] != *actual {
		t.Errorf("strategy expected to return true in %s", waitDurations[len(waitDurations)-1])
	}
}

func TestBackoff(t *testing.T) {
	const backoffDuration = time.Duration(10)
	const algorithmDurationBase = time.Duration(1)

	algorithm := func(attempt uint) time.Duration {
		return backoffDuration - (algorithmDurationBase * time.Duration(attempt))
	}

	strategy := Backoff(algorithm)

	if spy, actual := sleepSpy(); !strategy(0, spy) || 0 != *actual {
		t.Error("strategy expected to return true in ~0 time")
	}

	for i := uint(1); i < 10; i++ {
		expectedResult := algorithm(i)

		if spy, actual := sleepSpy(); !strategy(i, spy) || expectedResult != *actual {
			t.Errorf("strategy expected to return true in %s", expectedResult)
		}
	}
}

func TestBackoffWithJitter(t *testing.T) {
	const backoffDuration = time.Duration(20)
	const algorithmDurationBase = time.Duration(1)

	algorithm := func(attempt uint) time.Duration {
		return backoffDuration - (algorithmDurationBase * time.Duration(attempt))
	}

	transformation := func(duration time.Duration) time.Duration {
		return duration - time.Duration(10)
	}

	strategy := BackoffWithJitter(algorithm, transformation)

	if spy, actual := sleepSpy(); !strategy(0, spy) || 0 != *actual {
		t.Error("strategy expected to return true in ~0 time")
	}

	for i := uint(1); i < 10; i++ {
		expectedResult := transformation(algorithm(i))

		if spy, actual := sleepSpy(); !strategy(i, spy) || expectedResult != *actual {
			t.Errorf("strategy expected to return true in %s", expectedResult)
		}
	}
}

func TestNoJitter(t *testing.T) {
	transformation := noJitter()

	for i := uint(0); i < 10; i++ {
		duration := time.Duration(i)
		result := transformation(duration)
		expected := duration

		if result != expected {
			t.Errorf("transformation expected to return a %s duration, but received %s instead", expected, result)
		}
	}
}

// sleepSpy returns a spy for the time.Sleep function that sums the
// durations passed to it.
func sleepSpy() (func(time.Duration), *time.Duration) {
	var actual time.Duration

	return func(d time.Duration) { actual += d }, &actual
}
