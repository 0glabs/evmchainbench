package gentx

import (
	"log"

	generatorpkg "github.com/0glabs/evmchainbench/lib/generator"
	limiterpkg "github.com/0glabs/evmchainbench/lib/limiter"
)

func GenTx(rpcUrl, faucetPrivateKey string, senderCount, txCount int, txStoreDir string, mempool int) {
	limiter := limiterpkg.NewRateLimiter(mempool)

	generator, err := generatorpkg.NewGenerator(rpcUrl, faucetPrivateKey, senderCount, txCount, true, txStoreDir, limiter)
	if err != nil {
		log.Fatalf("Failed to create generator: %v", err)
	}

	_, err = generator.GenerateSimple()
	if err != nil {
		log.Fatalf("Failed to generate transactions: %v", err)
	}
}
