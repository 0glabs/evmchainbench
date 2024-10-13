package run

import (
	"log"
)

func Run(rpcUrl, faucetPrivateKey string, senderCount, txCount int) {
	generator, err := NewGenerator(rpcUrl, faucetPrivateKey, senderCount, txCount, false, "")
	if err != nil {
		log.Fatalf("Failed to create generator: %v", err)
	}

	txsMap, err := generator.GenerateSimple()
	if err != nil {
		log.Fatalf("Failed to generate transactions: %v", err)
	}

	transmitter, err := NewTransmitter(rpcUrl)
	if err != nil {
		log.Fatalf("Failed to create transmitter: %v", err)
	}

	err = transmitter.Broadcast(txsMap)
	if err != nil {
		log.Fatalf("Failed to broadcast transactions: %v", err)
	}
}
