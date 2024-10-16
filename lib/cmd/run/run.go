package run

import (
	"log"
)

func Run(httpRpc, wsRpc, faucetPrivateKey string, senderCount, txCount int, mempool int) {
	limiter := NewRateLimiter(mempool)

	ethListener := NewEthereumListener(wsRpc, limiter)
	err := ethListener.Connect()
	if err != nil {
		log.Fatalf("Failed to connect to WebSocket: %v", err)
	}

	// Subscribe new heads
	err = ethListener.SubscribeNewHeads()
	if err != nil {
		log.Fatalf("Failed to subscribe to new heads: %v", err)
	}

	generator, err := NewGenerator(httpRpc, faucetPrivateKey, senderCount, txCount, false, "", limiter)
	if err != nil {
		log.Fatalf("Failed to create generator: %v", err)
	}

	txsMap, err := generator.GenerateSimple()
	if err != nil {
		log.Fatalf("Failed to generate transactions: %v", err)
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
