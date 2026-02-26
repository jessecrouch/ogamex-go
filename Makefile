.PHONY: build run dev swagger clean

export PATH := $(PATH):$(shell go env GOPATH)/bin

build:
	swag init -g cmd/server/main.go -o docs
	go build -o ogamex cmd/server/main.go

run:
	LD_LIBRARY_PATH=./storage/rust-libs ./ogamex

dev: build run

swagger:
	swag init -g cmd/server/main.go -o docs

clean:
	rm -f ogamex
