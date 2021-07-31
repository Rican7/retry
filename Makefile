# Define directories
ROOT_DIR ?= ${CURDIR}
TOOLS_DIR ?= ${ROOT_DIR}/tools

# Set a local GOBIN, since the default value can't be trusted
# (See https://github.com/golang/go/issues/23439)
export GOBIN ?= ${TOOLS_DIR}/bin

# Set the mode for code-coverage
GO_TEST_COVERAGE_MODE ?= count
GO_TEST_COVERAGE_FILE_NAME ?= coverage.out

# Set flags for `gofmt`
GOFMT_FLAGS ?= -s

# Set a default `min_confidence` value for `golint`
GOLINT_MIN_CONFIDENCE ?= 0.1


all: build

clean:
	go clean -i -x ./...

build:
	go build -v ./...

install-deps:
	go mod download

tools install-deps-dev:
	cd tools && go install \
		golang.org/x/lint/golint \
		golang.org/x/tools/cmd/goimports \
		honnef.co/go/tools/cmd/staticcheck

update-deps:
	go get ./...

test:
	go test -v ./...

test-with-coverage:
	go test -cover -covermode ${GO_TEST_COVERAGE_MODE} ./...

test-with-coverage-formatted:
	go test -cover -covermode ${GO_TEST_COVERAGE_MODE} ./... | column -t | sort -r

test-with-coverage-profile:
	go test -covermode ${GO_TEST_COVERAGE_MODE} -coverprofile ${GO_TEST_COVERAGE_FILE_NAME} ./...

format-lint:
	errors=$$(gofmt -l ${GOFMT_FLAGS} .); if [ "$${errors}" != "" ]; then echo "$${errors}"; exit 1; fi

import-lint: install-deps-dev
	errors=$$(${GOBIN}/goimports -l .); if [ "$${errors}" != "" ]; then echo "$${errors}"; exit 1; fi

style-lint: install-deps-dev
	${GOBIN}/golint -min_confidence=${GOLINT_MIN_CONFIDENCE} -set_exit_status ./...
	${GOBIN}/staticcheck ./...

lint: install-deps-dev format-lint import-lint style-lint

vet:
	go vet ./...

format-fix:
	gofmt -w ${GOFMT_FLAGS} .

import-fix:
	goimports -w .


.PHONY: all clean build install-deps tools install-deps-dev update-deps test test-with-coverage test-with-coverage-formatted test-with-coverage-profile format-lint import-lint style-lint lint vet format-fix import-fix
