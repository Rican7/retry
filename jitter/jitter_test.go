package jitter

import (
	"math/rand"
	"testing"
	"time"
)

func TestFullRandom(t *testing.T) {
	const seed = 0
	const duration = time.Millisecond

	generator := rand.New(rand.NewSource(seed))

	transformation := FullRandom(generator)

	// Based on constant seed
	expectedDurations := []time.Duration{165505, 393152, 995827, 197794, 376202}

	for _, expected := range expectedDurations {
		result := transformation(duration)

		if result != expected {
			t.Errorf("transformation expected to return a %s duration, but received %s instead", expected, result)
		}
	}
}

func TestFallbackNewRandom(t *testing.T) {
	generator := rand.New(rand.NewSource(0))

	if result := fallbackNewRandom(generator); generator != result {
		t.Errorf("result expected to match parameter, received %+v instead", result)
	}

	if result := fallbackNewRandom(nil); nil == result {
		t.Error("recieved unexpected nil result")
	}
}
