INSTALL_PATH := $(HOME)/.local/bin/td

all: test build

test:
	go test -v ./...

build: test
	go build -o bin/td main.go


clean:
	rm -f bin/td

run:
	go run main.go

install: build
	cp -f bin/td $(INSTALL_PATH)


.PHONY: all build test clean run deps build-linux
