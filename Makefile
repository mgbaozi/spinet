all: build

env =

linuxenv:
	$(eval env := GOOS=linux GOARCH=amd64)

linux: linuxenv build

mod:
	go build ./...

.PHONY: build
build: cli

# .PHONY: generate
# generate:
# 	go generate ./ent

.PHONY: run
run: bin/spictl
	bin/spictl

.PHONY: cli
cli: bin/spictl

bin/spictl: cmd/cli/*.go pkg/**/*.go
	$(env) go build -o bin/spictl ./cmd/cli


.PHONY: clean
clean:
	rm -rf bin
