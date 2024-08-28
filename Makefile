BINARY_NAME=block-benchmark

build:
	go build -o bin/${BINARY_NAME} main.go

all: build
