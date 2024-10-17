package run

import (
	"log"
)

func Run(httprpc, wsrpc, faucetPrivateKey string, senderCount, txCount int, mempool int) {
	limiter := NewRateLimiter(mempool)

	ethListener := NewEthereumListener(wsrpc, limiter)
	err := ethListener.Connect()
	if err != nil {
		log.Fatalf("Failed to connect to WebSocket: %v", err)
	}

	// 订阅新块事件
	err = ethListener.SubscribeNewHeads()
	if err != nil {
		log.Fatalf("Failed to subscribe to new heads: %v", err)
	}

	generator, err := NewGenerator(httprpc, faucetPrivateKey, senderCount, txCount, false, "", limiter)
	if err != nil {
		log.Fatalf("Failed to create generator: %v", err)
	}

	txsMap, err := generator.GenerateSimple()
	if err != nil {
		log.Fatalf("Failed to generate transactions: %v", err)
	}

	transmitter, err := NewTransmitter(httprpc, limiter)
	if err != nil {
		log.Fatalf("Failed to create transmitter: %v", err)
	}

	err = transmitter.Broadcast(txsMap)
	if err != nil {
		log.Fatalf("Failed to broadcast transactions: %v", err)
	}

	<-ethListener.quit
}
