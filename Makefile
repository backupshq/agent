.PHONY: test packages
test:
	go test ./...

packages:
	go get github.com/urfave/cli github.com/BurntSushi/toml github.com/robfig/cron

build: backupshq

.PHONY: backupshq
backupshq:
	go build -o backupshq
