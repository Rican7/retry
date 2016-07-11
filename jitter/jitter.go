// Package jitter provides methods of transforming durations.
//
// Copyright © 2016 Trevor N. Suarez (Rican7)
package jitter

import (
	"math/rand"
	"time"
)

// Transformation defines a function that calculates a time.Duration based on
// the given duration.
type Transformation func(duration time.Duration) time.Duration

// FullRandom creates a Transformation that transforms a duration into a
// result duration in [0, n), where n is the given duration.
//
// The given generator is what is used to determine the random transformation.
// If a nil generator is passed, a default one will be provided.
//
// Inspired by https://www.awsarchitectureblog.com/2015/03/backoff.html
func FullRandom(generator *rand.Rand) Transformation {
	random := fallbackNewRandom(generator)

	return func(duration time.Duration) time.Duration {
		return time.Duration(random.Int63n(int64(duration)))
	}
}

// fallbackNewRandom returns the passed in random instance if it's not nil,
// and otherwise returns a new random instance seeded with the current time.
func fallbackNewRandom(random *rand.Rand) *rand.Rand {
	// Return the passed in value if it's already not null
	if nil != random {
		return random
	}

	seed := time.Now().UnixNano()

	return rand.New(rand.NewSource(seed))
}
