all: build

env =

linuxenv:
	$(eval env := GOOS=linux GOARCH=amd64)

linux: linuxenv build

mod:
	go build ./...

.PHONY: build
build: api cli
#build: api cli

# .PHONY: generate
# generate:
# 	go generate ./ent

.PHONY: run
run: api
	bin/api

.PHONY: api
api: bin/api

.PHONY: cli
cli: bin/cli

bin/api: cmd/api/*.go pkg/**/*.go
	$(env) go build -o bin/api ./cmd/api

bin/cli: cmd/cli/*.go pkg/**/*.go
	$(env) go build -o bin/cli ./cmd/cli


.PHONY: clean
clean:
	rm -rf bin
