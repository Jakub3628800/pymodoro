# some makefile for go

lint:
	golangci-lint run

build:
	go build -o bin/pomodoro main.go

test:
	go test -v ./...


all: lint build
