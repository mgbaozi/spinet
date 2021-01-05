all: build

env =

linuxenv:
	$(eval env := GOOS=linux GOARCH=amd64)

linux: linuxenv build

mod:
	go build ./...

.PHONY: build
build: spinet spictl

.PHONY: run
run: bin/spinet
	bin/spinet

.PHONY: spinet
spinet: bin/spinet

.PHONY: spictl
spictl: bin/spictl

bin/spinet: cmd/spinet/*.go pkg/**/*.go
	$(env) go build -o bin/spinet ./cmd/spinet

bin/spictl: cmd/spictl/*.go pkg/**/*.go
	$(env) go build -o bin/spictl ./cmd/spictl

.PHONY: clean
clean:
	rm -rf bin
