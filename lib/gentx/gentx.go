package gentx

import (
	"log"

	"github.com/ethanz0g/block-benchmark/lib/run"
)

func GenTx(rpcUrl, faucetPrivateKey string, senderCount, txCount int, txStoreDir string) {
	generator, err := run.NewGenerator(rpcUrl, faucetPrivateKey, senderCount, txCount, true, txStoreDir)
	if err != nil {
		log.Fatalf("Failed to create generator: %v", err)
	}

	_, err = generator.GenerateSimple()
	if err != nil {
		log.Fatalf("Failed to generate transactions: %v", err)
	}
}
