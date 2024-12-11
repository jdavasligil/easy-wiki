.PHONY: all fmt tidy build clean
all: fmt tidy build

fmt:
	@go fmt ./...

tidy:
	@go mod tidy

test:
	@go test -v ./...

build:
	@go build -o ./bin/ewik ./cmd/ewik

run:
	@./bin/ewik

clean:
	@go clean && rm -rf ./bin/*
