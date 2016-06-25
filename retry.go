// Package retry provides a simple, stateless mechanism to perform actions
// repetitively until successful.
//
// Copyright © 2016 Trevor N. Suarez (Rican7)
package retry

// Action defines a callable function that package retry can handle.
type Action func(attempt uint) error

// Strategy defines a function that Retry calls before every successive attempt
// to determine whether it should make the next attempt or not. Returning `true`
// allows for the next attempt to be made. Returning `false` halts the retrying
// process and returns the last error returned by the `Action`.
type Strategy func(attempt uint) bool

// Retry takes an action and performs it until successful, as defined by the
// provided strategy.
func Retry(action Action, strategies ...Strategy) error {
	var err error

	for attempt := uint(0); (0 == attempt || nil != err) && shouldAttempt(attempt, strategies...); attempt++ {
		err = action(attempt)
	}

	return err
}

// shouldAttempt evaluates the provided strategies with the given attempt to
// determine if the Retry loop should make another attempt
func shouldAttempt(attempt uint, strategies ...Strategy) bool {
	shouldAttempt := true

	for i := 0; shouldAttempt && i < len(strategies); i++ {
		shouldAttempt = shouldAttempt && strategies[i](attempt)
	}

	return shouldAttempt
}
