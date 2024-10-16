BINARY_NAME=evmchainbench

build:
	go build -o bin/${BINARY_NAME} main.go

contract:
	solc --abi --bin -o contracts/incrementer contracts/incrementer.sol

metadata:
	@./generate_contract_meta_data.sh

all: clean contract metadata build

clean:
	rm -rf contracts/incrementer
