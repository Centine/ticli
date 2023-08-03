.PHONY: build run clean

# Define the name of the binary output
BINARY_NAME=ticli
BUILD=`date +%FT%T%z`
VERSION=`git describe --tags --always --dirty`

LDFLAGS=-ldflags "-X main.Version=${VERSION} -X main.Build=${BUILD}"

# Define the go command
GO_CMD=go

all: test build
build:
	go build $(LDFLAGS) -o $(BINARY_NAME) -v ./cmd/ticli

test:
	go test -v ./...

clean:
	go clean
	rm -f $(BINARY_NAME)

run: build
	./$(BINARY_NAME)

deps:
	go mod download

cross:
	GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o $(BINARY_NAME)_windows_amd64 -v ./cmd/ticli
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o $(BINARY_NAME)_linux_amd64 -v ./cmd/ticli
	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o $(BINARY_NAME)_darwin_amd64 -v ./cmd/ticli
