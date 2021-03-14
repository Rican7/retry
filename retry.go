// Package retry provides a simple, stateless, functional mechanism to perform
// actions repetitively until successful.
//
// Copyright © 2016 Trevor N. Suarez (Rican7)
package retry

import (
	"context"
	"time"

	"github.com/Rican7/retry/strategy"
)

// Action defines a callable function that package retry can handle.
type Action func(attempt uint) error

// ActionWithContext defines a callable function that package retry can handle.
type ActionWithContext func(ctx context.Context, attempt uint) error

// Retry takes an action and performs it, repetitively, until successful.
//
// Optionally, strategies may be passed that assess whether or not an attempt
// should be made.
func Retry(action Action, strategies ...strategy.Strategy) error {
	return RetryWithContext(context.Background(), func(ctx context.Context, attempt uint) error { return action(attempt) }, strategies...)
}

// RetryWithContext takes an action and performs it, repetitively, until successful
// or the context is done.
//
// Optionally, strategies may be passed that assess whether or not an attempt
// should be made.
//
// Context errors take precedence over action errors so this commonplace test:
//
//     err := retry.RetryWithContext(...)
//     if err != nil { return err }
//
// will pass cancellation errors up the call chain.
func RetryWithContext(ctx context.Context, action ActionWithContext, strategies ...strategy.Strategy) error {
	var err error

	for attempt := uint(0); (0 == attempt || nil != err) && shouldAttempt(attempt, sleepFunc(ctx), strategies...) && nil == ctx.Err(); attempt++ {
		err = action(ctx, attempt)
	}

	if ctx.Err() != nil {
		return ctx.Err()
	}

	return err
}

// shouldAttempt evaluates the provided strategies with the given attempt to
// determine if the Retry loop should make another attempt.
func shouldAttempt(attempt uint, sleep func(time.Duration), strategies ...strategy.Strategy) bool {
	shouldAttempt := true

	for i := 0; shouldAttempt && i < len(strategies); i++ {
		shouldAttempt = shouldAttempt && strategies[i](attempt, sleep)
	}

	return shouldAttempt
}

// sleepFunc returns a function with the same signature as time.Sleep()
// that blocks for the given duration, but will return sooner if the context is
// cancelled or its deadline passes.
func sleepFunc(ctx context.Context) func(time.Duration) {
	return func(d time.Duration) {
		select {
		case <-ctx.Done():
		case <-time.After(d):
		}
	}
}
