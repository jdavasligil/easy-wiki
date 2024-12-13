ifeq ($(PREFIX),)
	PREFIX := /usr/local
endif

.PHONY: all fmt tidy build clean install
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

install:
	install -D ./bin/ewik $(DESTDIR)$(PREFIX)/bin/

uninstall:
	rm /usr/local/bin/ewik
