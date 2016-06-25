package retry

// AttemptLimiter creates a Strategy that limits the number of attempts that
// Retry will make
func AttemptLimiter(attemptLimit uint) Strategy {
	return func(attempt uint) bool {
		return (attempt < attemptLimit)
	}
}
