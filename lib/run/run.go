package run

import (
	"fmt"
	"log"
)

func Run(rpcUrl, faucetPrivateKey string, senderCount, txCount int) {
	generator, err := NewGenerator(rpcUrl, faucetPrivateKey, senderCount, txCount)
	if err != nil {
		log.Fatalf("Failed to generate generator: %v", err)
	}

	txsMap, err := generator.GenerateSimple()
	fmt.Printf("%v", txsMap)
}
