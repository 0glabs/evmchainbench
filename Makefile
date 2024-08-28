BINARY_NAME=evmchainbench

build:
	go build -o bin/${BINARY_NAME} main.go

all: build
