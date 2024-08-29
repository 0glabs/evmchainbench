package lib

import (
	"context"
	"errors"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethanz0g/block-benchmark/lib/incrementer_contract"
	"github.com/spf13/viper"
)

func LoadIncrementerContractForRead() (*payment_contract.PaymentContract, error) {
	_, instance, err := loadIncrementerContractBase()
	if err != nil {
		return nil, err
	}
	return instance, nil
}

func LoadIncrementerContractForWrite() (*ethclient.Client, *payment_contract.PaymentContract, *bind.TransactOpts, error) {
	client, instance, err := loadIncrementerContractBase()
	if err != nil {
		return nil, nil, nil, err
	}

	auth, err := prepareTransactionAuth(client)
	if err != nil {
		return nil, nil, nil, err
	}

	return client, instance, auth, err // let caller to handle err
}

func PrepareDeployTransaction() (*ethclient.Client, *bind.TransactOpts, error) {
	client, err := getClient()
	if err != nil {
		return nil, nil, err
	}

	auth, err := prepareTransactionAuth(client)
	if err != nil {
		return nil, nil, err
	}

	return client, auth, nil
}

func prepareTransactionAuth(client *ethclient.Client) (*bind.TransactOpts, error) {
	privKey, err := LoadPrivKey()
	if err != nil {
		return nil, err
	}

	fromAddress := PrivKeyToAddress(privKey)

	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		return nil, err
	}

	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		return nil, err
	}

	gasPrice.Add(gasPrice, big.NewInt(2000000000))

	chainId := viper.GetInt(config.BiofiChainId)
	auth, err := bind.NewKeyedTransactorWithChainID(privKey, big.NewInt(int64(chainId)))
	if err != nil {
		return nil, err
	}
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)       // in wei
	auth.GasLimit = uint64(15000000) // in units
	auth.GasPrice = gasPrice

	return auth, nil
}

func loadIncrementerContractBase() (*ethclient.Client, *payment_contract.PaymentContract, error) {
	client, err := getClient()
	if err != nil {
		return nil, nil, err
	}

	address, err := loadIncrementerContractAddress()
	if err != nil {
		return nil, nil, err
	}

	instance, err := incrementer_contract.NewIncrementerContract(address, client)
	if err != nil {
		return nil, nil, err
	}
	return client, instance, nil
}

func loadIncrementerContractAddress() (common.Address, error) {
	hexAddress := viper.GetString(config.IncrementerContractAddress)
	if !IsValidHexAddress(hexAddress) {
		return common.Address{}, errors.New("Payment contract address has a wrong format!")
	}
	address := HexToAddress(hexAddress)
	return address, nil
}

func getClient() (*ethclient.Client, error) {
	rpcUrl := viper.GetString(config.BiofiRpcUrl)
	client, err := ethclient.Dial(rpcUrl)
	if err != nil {
		return &ethclient.Client{}, err
	}

	return client, nil
}
