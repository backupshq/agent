.PHONY: test packages
test:
	go test ./...

build: backupshq

.PHONY: backupshq
backupshq:
	go build -o backupshq
