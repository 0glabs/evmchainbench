package lib

import (
	//"github.com/ethanz0g/block-benchmark/config"
	"github.com/ethanz0g/block-benchmark/lib/incrementer_contract"
	"github.com/spf13/viper"
)

func DeployContract() (string, error) {
	client, auth, err := PrepareDeployTransaction()
	if err != nil {
		return "", err
	}

	address, _, _, err := incrementer_contract.DeployIncrementerContract(auth, client)
	if err != nil {
		return "", err
	}

	return address.Hex(), nil
}

