SHELL := /bin/bash
PWD := $(shell pwd)

BINDIR  := $(CURDIR)/bin
BINNAME ?= contagion-updater

HEROKU_APP_NAME:=ct-contagion-updater

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
	GOOS=linux go build $(GOFLAGS) -ldflags '$(LDFLAGS)' -o $(BINDIR)/$(BINNAME) ./cmd

.PHONY: clean
clean: ## Clean workspace
	rm -rf $(BINDIR)

# -- Heroku related commands
# You need to be logged in Heroku CLI before doing this
#   heroku login
#   heroku container:login
.PHONY: heroku-push
heroku-push:
	heroku container:push worker --recursive --app=$(HEROKU_APP_NAME) --verbose

.PHONY: heroku-release
heroku-release:
	heroku container:release worker --app $(HEROKU_APP_NAME) --verbose
