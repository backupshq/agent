.PHONY: test packages
test:
	go test ./...

packages:
	go get github.com/urfave/cli github.com/BurntSushi/toml

build: backupshq

backupshq:
	go build -o backupshq