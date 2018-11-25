
TARGET := gpu-mon

GOPATH ?= $(shell go env GOPATH)

# Ensure GOPATH is set before running build process.
ifeq "$(GOPATH)" ""
  $(error Please set the environment variable GOPATH before running `make`)
endif
PATH := ${GOPATH}/bin:$(PATH)

DCGM ?= $(shell locate nv-hostengine)
ifeq "$(DCGM)" ""
  $(error Please confirm that DCGM has been installed before running `make`)
endif

GitVersion := $(shell git rev-parse --short HEAD || echo unsupported)
DATE := $(shell date "+%Y-%m-%d %H:%M:%S")
VERSION := $(shell cat VERSION)
GOVERSION := $(shell go version)

GOFMT ?= gofmt -s
GOFILES := $(shell find . -name "*.go" -type f -not -path "./vendor/*")
PACKAGES ?= $(shell go list ./... | grep -v /vendor/)

.PHONY: all build fmt vet

all: fmt vet build

fmt:
	$(GOFMT) -w $(GOFILES)

vet:
	@go vet $(PACKAGES)

build:
	@echo "building gpu-mon ..."
	@go build -ldflags "-X main.GitCommit=$(GitVersion) \
			-X 'main.BuildTime=$(DATE)' -X 'main.Version=$(VERSION)' \
			-X 'main.GoVersion=$(GOVERSION)' " \
			-o $(TARGET)

# Run golang test cases
.PHONY: test
test:
	@echo "Run all test cases ..."
	go test ./common/ ./send/
	@echo "test Success!"

.PHONY: clean
clean:
	@rm -rf $(TARGET)
