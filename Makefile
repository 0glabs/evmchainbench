BINARY_NAME=evmchainbench

build:
	go build -o bin/${BINARY_NAME} main.go

contract-erc20:
	docker run \
		--rm \
		-v $$(pwd)/contracts:/src ethereum/solc:0.8.19 \
		--optimize --bin --abi --overwrite \
		-o /src/build/erc20 \
		/src/erc20.sol

contract-uniswap:
	docker run \
		--rm \
		-v $$(pwd)/contracts:/src ethereum/solc:0.5.16 \
		--optimize --bin --abi --overwrite \
		-o /src/build/uniswap \
		/src/uniswap_source/UniswapV2ERC20.sol \
		/src/uniswap_source/UniswapV2Factory.sol \
		/src/uniswap_source/UniswapV2Pair.sol

metadata:
	@./generate_contract_meta_data.sh

contract: contract-erc20 contract-uniswap

all: clean contract metadata build

clean:
	rm -rf contracts/erc20
