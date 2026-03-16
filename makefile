# Makefile for common Go commands

fmt:
	go fmt .

tidy:
	go mod tidy

vet:
	go vet .

test:
	go test -v tests/order_signing/signing_test.go

build:
	mkdir -p bin
	go build -o bin/example_account_balance examples/account_balance/main.go
	go build -o bin/example_limit_multiple examples/limit_multiple/main.go
	go build -o bin/example_limit_single examples/limit_single/main.go
	go build -o bin/example_socketio examples/socketio/main.go
	go build -o bin/example_websocket examples/websocket/main.go

deps:
	go get -u ./...

all: fmt tidy vet test build
