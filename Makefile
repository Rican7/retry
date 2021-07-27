# Set the mode for code-coverage
GO_TEST_COVERAGE_MODE ?= count
GO_TEST_COVERAGE_FILE_NAME ?= coverage.out

# Set flags for `gofmt`
GOFMT_FLAGS ?= -s

# Set a default `min_confidence` value for `golint`
GOLINT_MIN_CONFIDENCE ?= 0.3


all: build

clean:
	go clean -i -x ./...

build:
	go build -v ./...

install-deps:
	go mod download

install-deps-dev: install-deps
	go get golang.org/x/lint/golint
	go get golang.org/x/tools/cmd/goimports

update-deps:
	go get -d -t -u ./...

update-deps-dev: update-deps
	go get -u golang.org/x/lint/golint
	go get -u golang.org/x/tools/cmd/goimports

test:
	go test -v ./...

test-with-coverage:
	go test -cover ./...

test-with-coverage-formatted:
	go test -cover ./... | column -t | sort -r

test-with-coverage-profile:
	echo "mode: ${GO_TEST_COVERAGE_MODE}" > ${GO_TEST_COVERAGE_FILE_NAME}
	for package in $$(go list ./...); do \
		go test -covermode ${GO_TEST_COVERAGE_MODE} -coverprofile "coverage_$${package##*/}.out" "$${package}"; \
		sed '1d' "coverage_$${package##*/}.out" >> ${GO_TEST_COVERAGE_FILE_NAME}; \
	done

format-lint:
	errors=$$(gofmt -l ${GOFMT_FLAGS} .); if [ "$${errors}" != "" ]; then echo "$${errors}"; exit 1; fi

import-lint:
	errors=$$(goimports -l .); if [ "$${errors}" != "" ]; then echo "$${errors}"; exit 1; fi

style-lint:
	errors=$$(golint -min_confidence=${GOLINT_MIN_CONFIDENCE} ./...); if [ "$${errors}" != "" ]; then echo "$${errors}"; exit 1; fi

lint: install-deps-dev format-lint import-lint style-lint

format-fix:
	gofmt -w ${GOFMT_FLAGS} .

import-fix:
	goimports -w .

vet:
	go vet ./...


.PHONY: all clean build install-deps install-deps-dev update-deps update-deps-dev test test-with-coverage test-with-coverage-formatted test-with-coverage-profile format-lint import-lint style-lint lint format-fix import-fix vet
