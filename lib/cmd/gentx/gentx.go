package gentx

import (
	"log"

	generatorpkg "github.com/0glabs/evmchainbench/lib/generator"
)

func GenTx(rpcUrl, faucetPrivateKey string, senderCount, txCount int, txStoreDir string) {
	generator, err := generatorpkg.NewGenerator(rpcUrl, faucetPrivateKey, senderCount, txCount, true, txStoreDir)
	if err != nil {
		log.Fatalf("Failed to create generator: %v", err)
	}

	_, err = generator.GenerateSimple()
	if err != nil {
		log.Fatalf("Failed to generate transactions: %v", err)
	}
}
