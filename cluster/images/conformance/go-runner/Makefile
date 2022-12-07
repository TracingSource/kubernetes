HOST_GOOS ?= $(shell go env GOOS)
HOST_GOARCH ?= $(shell go env GOARCH)
GO_BUILD ?= go build

.PHONY: all build clean

all: build

build:
	$(GO_BUILD)

clean:
	rm done e2e.log e2e.tar.gz go-runner
