package gentx

import (
	"log"

	"github.com/0glabs/evmchainbench/lib/run"
)

func GenTx(rpcUrl, faucetPrivateKey string, senderCount, txCount int, txStoreDir string) {
	generator, err := run.NewGenerator(rpcUrl, faucetPrivateKey, senderCount, txCount, true, txStoreDir, nil)
	if err != nil {
		log.Fatalf("Failed to create generator: %v", err)
	}

	_, err = generator.GenerateSimple()
	if err != nil {
		log.Fatalf("Failed to generate transactions: %v", err)
	}
}
