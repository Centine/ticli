.PHONY: build run clean

# Define the name of the binary output
BIN_NAME=ticli
BUILD=`date +%FT%T%z`
VERSION=`git describe --tags --always --dirty`

# Define the go command
GO_CMD=go

build:
	$(GO_CMD) build -o $(BIN_NAME) ./cmd/ticli

run: build
	./$(BIN_NAME)

clean:
	rm -f $(BIN_NAME)

winbuild:
	GOOS=windows GOARCH=amd64 $(GO_CMD) build -o $(BIN_NAME).exe ./cmd/ticli
