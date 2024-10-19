package generator

import (
	"context"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/0glabs/evmchainbench/lib/account"
	"github.com/0glabs/evmchainbench/lib/store"
	"github.com/0glabs/evmchainbench/lib/util"
)

type Generator struct {
	FaucetAccount *account.Account
	Senders       []*account.Account
	Recipients    []string
	RpcUrl        string
	ChainID       *big.Int
	GasPrice      *big.Int
	ShouldPersist bool
	Store         *store.Store
}

func NewGenerator(rpcUrl, faucetPrivateKey string, senderCount, txCount int, shouldPersist bool, txStoreDir string) (*Generator, error) {
	client, err := ethclient.Dial(rpcUrl)
	if err != nil {
		return &Generator{}, err
	}

	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		return &Generator{}, err
	}

	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		return &Generator{}, err
	}

	faucetAccount, err := account.CreateFaucetAccount(client, faucetPrivateKey)
	if err != nil {
		return &Generator{}, err
	}

	senders := make([]*account.Account, senderCount)
	for i := 0; i < senderCount; i++ {
		s, err := account.NewAccount(client)
		if err != nil {
			return &Generator{}, err
		}
		senders[i] = s
	}

	recipients := make([]string, txCount)
	for i := 0; i < txCount; i++ {
		r, err := account.GenerateRandomAddress()
		if err != nil {
			return &Generator{}, err
		}
		recipients[i] = r
	}

	client.Close()

	return &Generator{
		FaucetAccount: faucetAccount,
		Senders:       senders,
		Recipients:    recipients,
		RpcUrl:        rpcUrl,
		ChainID:       chainID,
		GasPrice:      gasPrice,
		ShouldPersist: shouldPersist,
		Store:         store.NewStore(txStoreDir),
	}, nil
}

func (g *Generator) prepareSenders() error {
	client, err := ethclient.Dial(g.RpcUrl)
	if err != nil {
		return err
	}
	defer client.Close()

	value := new(big.Int)
	value.Mul(big.NewInt(1000000000000000000), big.NewInt(100)) // 100 Eth

	txs := types.Transactions{}

	for _, recipient := range g.Senders {
		signedTx, err := GenerateSimpleTransferTx(g.FaucetAccount.PrivateKey, recipient.Address.Hex(), g.FaucetAccount.GetNonce(), g.ChainID, g.GasPrice, value)
		if err != nil {
			return err
		}

		err = client.SendTransaction(context.Background(), signedTx)
		if err != nil {
			return err
		}

		if g.ShouldPersist {
			g.Store.AddPrepareTx(signedTx)
			if err != nil {
				return err
			}
		}

		txs = append(txs, signedTx)
	}

	err = util.WaitForReceiptsOfTxs(client, txs, 10*time.Second)
	if err != nil {
		return err
	}

	return nil
}
