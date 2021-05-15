SHELL := /bin/bash
PWD := $(shell pwd)

BINDIR  := $(CURDIR)/bin
BINNAME ?= contagion-updater

default: build

all:

deps:
	go mod tidy
	go mod vendor

.PHONY: format
format: ## Fix format code style
	go fmt ./...

.PHONY: build
build: deps
	GOOS=linux go build $(GOFLAGS) -ldflags '$(LDFLAGS)' -o $(BINDIR)/$(BINNAME) .

.PHONY: clean
clean: ## Clean workspace
	rm -rf $(BINDIR)
