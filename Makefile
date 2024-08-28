BINARY_NAME=evmchainbench

build:
	go build -o bin/${BINARY_NAME} main.go

contract:
	solc --abi --bin -o contracts/build-incrementer contracts/incrementer.sol

abigen:
	mkdir -p lib/incrementer_contract
	abigen --abi=./contracts/build-incrementer/Incrementer.abi \
	       --bin=./contracts/build-incrementer/Incrementer.bin \
	       --pkg=incrementer_contract \
	       --out=./lib/incrementer_contract/contract.go

all: contract abigen build

clean:
	rm -rf contracts/build-incrementer
