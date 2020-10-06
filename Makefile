.PHONY: help
help:
	@echo 'Usage: make [target]'
	@echo 'Available targets:'
	@echo
	@grep -Eo '^[-a-z/]+' Makefile | sort

.PHONY: test
test:
	go test ./...

.PHONY: build
build:
	go build -o backupshq

.PHONY: lint
lint:
	@gofmt -d -s ./.. | diff -u /dev/null -

.PHONY: fix-style
fix-style:
	@gofmt -s -w ./..
