.PHONY: build test

build:
	go build -o protog cmd/protog/main.go

test:
	go test -race
