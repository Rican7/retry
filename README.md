# retry

[![Build Status](https://github.com/Rican7/retry/actions/workflows/main.yml/badge.svg?branch=master)](https://github.com/Rican7/retry/actions/workflows/main.yml)
[![Coverage Status](https://coveralls.io/repos/github/Rican7/retry/badge.svg)](https://coveralls.io/github/Rican7/retry)
[![Go Report Card](https://goreportcard.com/badge/Rican7/retry)](http://goreportcard.com/report/Rican7/retry)
[![Go Reference](https://pkg.go.dev/badge/github.com/Rican7/retry.svg)](https://pkg.go.dev/github.com/Rican7/retry)
[![Latest Stable Version](https://img.shields.io/github/release/Rican7/retry.svg?style=flat)](https://github.com/Rican7/retry/releases)

A simple, stateless, functional mechanism to perform actions repetitively until successful.


## Project Status

This project is currently in "pre-release". While the code is heavily tested, the API may change.
Use a tagged version or vendor this dependency if you plan on using it.

That said, this code has been used in production without issue for years, and has been used by some relatively
[high-profile projects/codebases](https://pkg.go.dev/github.com/Rican7/retry?tab=importedby).


## Examples

### Basic

```go
retry.Retry(func(attempt uint) error {
	return nil // Do something that may or may not cause an error
})
```

### File Open

```go
const logFilePath = "/var/log/myapp.log"

var logFile *os.File

err := retry.Retry(func(attempt uint) error {
	var err error

	logFile, err = os.Open(logFilePath)

	return err
})

if err != nil {
	log.Fatalf("Unable to open file %q with error %q", logFilePath, err)
}

logFile.Chdir() // Do something with the file
```

### HTTP request with strategies and backoff

```go
var response *http.Response

action := func(attempt uint) error {
	var err error

	response, err = http.Get("https://api.github.com/repos/Rican7/retry")

	if err == nil && response != nil && response.StatusCode > 200 {
		err = fmt.Errorf("failed to fetch (attempt #%d) with status code: %d", attempt, response.StatusCode)
	}

	return err
}

err := retry.Retry(
	action,
	strategy.Limit(5),
	strategy.Backoff(backoff.Fibonacci(10*time.Millisecond)),
)

if err != nil {
	log.Fatalf("Failed to fetch repository with error %q", err)
}
```

### Retry with backoff jitter

```go
action := func(attempt uint) error {
	return errors.New("something happened")
}

seed := time.Now().UnixNano()
random := rand.New(rand.NewSource(seed))

retry.Retry(
	action,
	strategy.Limit(5),
	strategy.BackoffWithJitter(
		backoff.BinaryExponential(10*time.Millisecond),
		jitter.Deviation(random, 0.5),
	),
)
```
