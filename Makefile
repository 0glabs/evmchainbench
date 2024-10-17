BINARY_NAME=evmchainbench

build:
	go build -o bin/${BINARY_NAME} main.go

contract:
	solc --optimize --overwrite --abi --bin -o contracts/erc20 contracts/erc20.sol

metadata:
	@./generate_contract_meta_data.sh

all: clean contract metadata build

clean:
	rm -rf contracts/erc20
