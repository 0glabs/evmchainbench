package run

import (
	"log"

	"github.com/0glabs/evmchainbench/lib/generator"
	limiterpkg "github.com/0glabs/evmchainbench/lib/limiter"
	"github.com/ethereum/go-ethereum/core/types"
)

func Run(httpRpc, wsRpc, faucetPrivateKey string, senderCount, txCount int, txType string, mempool int) {
	generator, err := generator.NewGenerator(httpRpc, faucetPrivateKey, senderCount, txCount, false, "")
	if err != nil {
		log.Fatalf("Failed to create generator: %v", err)
	}

	var txsMap map[int]types.Transactions

	switch txType {
	case "simple":
		txsMap, err = generator.GenerateSimple()
	case "erc20":
		txsMap, err = generator.GenerateERC20()
	case "uniswap":
		txsMap, err = generator.GenerateUniswap()
	default:
		log.Fatalf("Transaction type \"%v\" is not valid", txType)
	}
	if err != nil {
		log.Fatalf("Failed to generate transactions: %v", err)
	}

	limiter := limiterpkg.NewRateLimiter(mempool)

	ethListener := NewEthereumListener(wsRpc, limiter)
	err = ethListener.Connect()
	if err != nil {
		log.Fatalf("Failed to connect to WebSocket: %v", err)
	}

	// Subscribe new heads
	err = ethListener.SubscribeNewHeads()
	if err != nil {
		log.Fatalf("Failed to subscribe to new heads: %v", err)
	}

	transmitter, err := NewTransmitter(httpRpc, limiter)
	if err != nil {
		log.Fatalf("Failed to create transmitter: %v", err)
	}

	err = transmitter.Broadcast(txsMap)
	if err != nil {
		log.Fatalf("Failed to broadcast transactions: %v", err)
	}

	<-ethListener.quit
}
