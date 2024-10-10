package run

import (
	"context"
	"crypto/ecdsa"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/0glabs/evmchainbench/lib/account"
)

type Generator struct {
	FaucetAccount *account.Account
	Senders       []*account.Account
	Recipients    []string

	RpcUrl   string
	ChainID  *big.Int
	GasPrice *big.Int
}

func NewGenerator(rpcUrl, faucetPrivateKey string, senderCount, txCount int) (Generator, error) {
	client, err := ethclient.Dial(rpcUrl)
	if err != nil {
		return Generator{}, err
	}

	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		return Generator{}, err
	}

	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		return Generator{}, err
	}

	faucetAccount, err := account.CreateFaucetAccount(client, faucetPrivateKey)
	if err != nil {
		return Generator{}, err
	}

	senders := make([]*account.Account, senderCount)
	for i := 0; i < senderCount; i++ {
		s, err := account.NewAccount(client)
		if err != nil {
			return Generator{}, err
		}
		senders[i] = s
	}

	recipients := make([]string, txCount)
	for i := 0; i < txCount; i++ {
		r, err := account.GenerateRandomAddress()
		if err != nil {
			return Generator{}, err
		}
		recipients[i] = r
	}

	client.Close()

	return Generator{
		FaucetAccount: faucetAccount,
		Senders:       senders,
		Recipients:    recipients,
		RpcUrl:        rpcUrl,
		ChainID:       chainID,
		GasPrice:      gasPrice,
	}, nil
}

func (g *Generator) GenerateSimple() (map[int]types.Transactions, error) {
	txsMap := make(map[int]types.Transactions)

	err := g.prepareSenders()
	if err != nil {
		return txsMap, err
	}

	value := big.NewInt(10000000000000) // 1/100,000 ETH

	for index, sender := range g.Senders {
		txs := types.Transactions{}
		for _, recipient := range g.Recipients {
			tx, err := generateSimpleTx(sender.PrivateKey, sender.Address, recipient, sender.GetNonce(), g.ChainID, g.GasPrice, value)
			if err != nil {
				return txsMap, err
			}
			txs = append(txs, tx)
		}
		txsMap[index] = txs
	}

	return txsMap, nil
}

func (g *Generator) prepareSenders() error {
	client, err := ethclient.Dial(g.RpcUrl)
	if err != nil {
		return err
	}

	value := new(big.Int)
	value.Mul(big.NewInt(1000000000000000000), big.NewInt(100)) // 100 Eth

	for _, recipient := range g.Senders {
		signedTx, err := generateSimpleTx(g.FaucetAccount.PrivateKey, g.FaucetAccount.Address, recipient.Address.Hex(), g.FaucetAccount.GetNonce(), g.ChainID, g.GasPrice, value)
		if err != nil {
			return err
		}

		err = client.SendTransaction(context.Background(), signedTx)
		if err != nil {
			return err
		}
	}

	//TODO: to use the transaction receipt
	time.Sleep(5 * time.Second)
	return nil
}

func generateSimpleTx(privateKey *ecdsa.PrivateKey, fromAddress common.Address, recipient string, nonce uint64, chainID, gasPrice, value *big.Int) (*types.Transaction, error) {
	gasLimit := uint64(21000) // Gas limit for ETH transfer

	toAddress := common.HexToAddress(recipient)
	tx := types.NewTransaction(nonce, toAddress, value, gasLimit, gasPrice, nil)

	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		return &types.Transaction{}, err
	}

	return signedTx, nil
}
