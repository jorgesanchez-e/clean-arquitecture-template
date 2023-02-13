APPNAME ?= diary

# used by `test` target
export REPORTS_DIR=./reports
# used by lint target
export GOLANGCILINT_VERSION=v1.34.1

build: clean
	mkdir -p build
	GOOS=$(GOOS) GOARCH=$(GOARCH) APPNAME=$(APPNAME) ./scripts/build

run: build
	./build/${APPNAME}

test:
	./scripts/unit-test

test-report:
	./scripts/show-tests

lint:
	./scripts/lint

clean:
	APPNAME=$(APPNAME) ./scripts/clean

.PHONY: build run test test-report lint clean
