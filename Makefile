SRC_ROOT := $(dir $(realpath $(lastword $(MAKEFILE_LIST))))
BUILD_DIR := $(SRC_ROOT)/build

.DEFAULT_GOAL := all

clean:
	go clean

test:
	go test ./...

race:
	go test -race ./...

server:
	mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/jqsrv

all: test server
