GOBIN=$(shell pwd)/bin
GOFILES=main.go
GONAME=$(shell basename "$(PWD)")
PID=/tmp/go-$(GONAME).pid

build:
	@echo "Building $(GOFILES) to ./bin"
	@GOPATH=$(GOPATH) GOBIN=$(GOBIN) go get -d
	@GOPATH=$(GOPATH) GOBIN=$(GOBIN) go build -o bin/$(GONAME) $(GOFILES)

all:
	echo $$GOPATH
	go get -d
	go run *.go