#!/usr/bin/make -f

DOCKER := $(shell which docker)

export GO111MODULE = on

###############################################################################
###                          Tools & Dependencies                           ###
###############################################################################

go-mod-cache: go.summake
	@echo "--> Download go modules to local cache"
	@go mod download

go.sum: go.mod
	@echo "--> Ensure dependencies have not been modified"
	@go mod verify
	@go mod tidy

clean:
	rm -rf \
	$(BUILDDIR)/ \
	artifacts/ \
	tmp-swagger-gen/

.PHONY: clean

###############################################################################
###                                Linting                                  ###
###############################################################################

lint:
	golangci-lint run --out-format=tab --timeout=10m

lint-fix:
	golangci-lint run --fix --out-format=tab --issues-exit-code=0 --timeout=10m
.PHONY: lint lint-fix

format:
	find . -name '*.go' -type f -not -path "./vendor*" -not -path "*.git*" | xargs gofmt -w -s
	find . -name '*.go' -type f -not -path "./vendor*" -not -path "*.git*" | xargs misspell -w
	find . -name '*.go' -type f -not -path "./vendor*" -not -path "*.git*" | xargs goimports -w -local github.com/desmos-labs/themis
.PHONY: format
