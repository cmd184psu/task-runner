.PHONY: all build run test

all: build

build:
	go build -o ssh-task-runner .

run:
	./ssh-task-runner

test:
	go test -v ./...

stop:
	kill -9 $(shell cat example.pid)