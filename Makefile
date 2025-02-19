# All tasks are expected be documented with a comment of the form `##
# targetName: usage information`. When in doubt, consult
# https://clarkgrubb.com/makefile-style-guide

MAKEFLAGS += --warn-undefined-variables
SHELL := bash
.SHELLFLAGS := -eu -o pipefail -c
.DEFAULT_GOAL := help
.DELETE_ON_ERROR:
.SUFFIXES:

include .bingo/Variables.mk

.PHONY: help
help:
	@grep '^## ' $(MAKEFILE_LIST) | sed 's/^.*:## //' | sort | column -t -s:

# directories
REPORTDIR = build/reports

# commands
GO         ?= $(shell which go)

# flags
COVERFLAGS  = -covermode $(COVERMODE) -coverprofile $(COVEROUT)
COVERMODE   = atomic
COVEROUT    = $(REPORTDIR)/coverage.out
COVERXML    = $(REPORTDIR)/coverage.xml
COVERHTML   = $(REPORTDIR)/coverage.html
LINTFLAGS   =
TESTFLAGS   =

# recipes

## all: lint and test code
.PHONY: all
all: lint test

## clean: removes generated files
.PHONY: clean
clean:
	@$(RM) -r $(REPORTDIR)/

## format: format source code
.PHONY: format
format: LINTFLAGS += --fix
format: lint

## lint: runs linters on source code
.PHONY: lint
lint: | $(GOLANGCI_LINT)
	@$(GOLANGCI_LINT) run $(LINTFLAGS)

## test-with-coverage: run tests and generate coverage reports
.PHONY: test-with-coverage
test-with-coverage: TESTFLAGS := $(COVERFLAGS) $(TESTFLAGS)
test-with-coverage: test
	$(GO) tool cover -html $(COVEROUT) -o $(COVERHTML)

## test: run tests
.PHONY: test
test: | $(REPORTDIR)
	$(GO) test $(TESTFLAGS) ./...

## version: print the version currently being worked on
.PHONY: version
version: | $(GOTAGGER)
	@$(GOTAGGER)

## tools: install support tools
.PHONY: tools
tools:
	@bingo get --verbose

$(REPORTDIR):
	@mkdir -p $@

