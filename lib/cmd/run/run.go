package run

import (
	"log"

	"github.com/0glabs/evmchainbench/lib/generator"
	limiterpkg "github.com/0glabs/evmchainbench/lib/limiter"
)

func Run(httpRpc, wsRpc, faucetPrivateKey string, senderCount, txCount int, mempool int) {
	generator, err := generator.NewGenerator(httpRpc, faucetPrivateKey, senderCount, txCount, false, "")
	if err != nil {
		log.Fatalf("Failed to create generator: %v", err)
	}

	txsMap, err := generator.GenerateSimple()
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
