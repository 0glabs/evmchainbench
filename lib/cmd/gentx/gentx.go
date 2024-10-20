package gentx

import (
	"log"

	generatorpkg "github.com/0glabs/evmchainbench/lib/generator"
)

func GenTx(rpcUrl, faucetPrivateKey string, senderCount, txCount int, txType string, txStoreDir string) {
	generator, err := generatorpkg.NewGenerator(rpcUrl, faucetPrivateKey, senderCount, txCount, true, txStoreDir)
	if err != nil {
		log.Fatalf("Failed to create generator: %v", err)
	}

	switch txType {
	case "simple":
		_, err = generator.GenerateSimple()
	case "erc20":
		_, err = generator.GenerateERC20()
	case "uniswap":
		_, err = generator.GenerateUniswap()
	default:
		log.Fatalf("Transaction type \"%v\" is not valid", txType)
	}
	if err != nil {
		log.Fatalf("Failed to generate transactions: %v", err)
	}
}
