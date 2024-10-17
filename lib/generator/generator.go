package generator

import (
	"context"
	"math/big"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/0glabs/evmchainbench/lib/account"
	"github.com/0glabs/evmchainbench/lib/contract_meta_data/erc20"
	"github.com/0glabs/evmchainbench/lib/store"
	"github.com/0glabs/evmchainbench/lib/util"
)

type Generator struct {
	FaucetAccount *account.Account
	Senders       []*account.Account
	Recipients    []string

	RpcUrl   string
	ChainID  *big.Int
	GasPrice *big.Int

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

func (g *Generator) GenerateSimple() (map[int]types.Transactions, error) {
	txsMap := make(map[int]types.Transactions)

	if g.ShouldPersist {
		defer g.Store.PersistPrepareTxs()
	}

	err := g.prepareSenders()
	if err != nil {
		return txsMap, err
	}

	value := big.NewInt(10000000000000) // 1/100,000 ETH

	var mutex sync.Mutex
	ch := make(chan error)

	for index, sender := range g.Senders {
		go func(index int, sender *account.Account) {
			txs := types.Transactions{}
			for _, recipient := range g.Recipients {
				tx, err := GenerateSimpleTransferTx(sender.PrivateKey, recipient, sender.GetNonce(), g.ChainID, g.GasPrice, value)
				if err != nil {
					ch <- err
					return
				}
				txs = append(txs, tx)
			}

			mutex.Lock()
			txsMap[index] = txs
			mutex.Unlock()
			ch <- nil
		}(index, sender)
	}

	for i := 0; i < len(g.Senders); i++ {
		msg := <-ch
		if msg != nil {
			return txsMap, msg
		}
	}

	if g.ShouldPersist {
		err := g.Store.PersistTxsMap(txsMap)
		if err != nil {
			return txsMap, err
		}
	}

	return txsMap, nil
}

func (g *Generator) GenerateERC20() (map[int]types.Transactions, error) {
	txsMap := make(map[int]types.Transactions)

	contractAddress, err := g.prepareContractERC20()
	if err != nil {
		return txsMap, err
	}
	contractAddressStr := contractAddress.Hex()

	err = g.prepareSenders()
	if err != nil {
		return txsMap, err
	}

	amount := big.NewInt(1000) // a random small amount

	var mutex sync.Mutex
	ch := make(chan error)

	for index, sender := range g.Senders {
		go func(index int, sender *account.Account) {
			txs := types.Transactions{}
			for _, recipient := range g.Recipients {
				tx, err := GenerateContractCallingTx(
					sender.PrivateKey,
					contractAddressStr,
					sender.GetNonce(),
					g.ChainID,
					g.GasPrice,
					erc20.MyTokenABI,
					"transfer",
					common.HexToAddress(recipient),
					amount,
				)
				if err != nil {
					ch <- err
					return
				}
				txs = append(txs, tx)
			}

			mutex.Lock()
			txsMap[index] = txs
			mutex.Unlock()
			ch <- nil
		}(index, sender)
	}

	for i := 0; i < len(g.Senders); i++ {
		msg := <-ch
		if msg != nil {
			return txsMap, msg
		}
	}

	if g.ShouldPersist {
		err := g.Store.PersistTxsMap(txsMap)
		if err != nil {
			return txsMap, err
		}
	}

	return txsMap, nil
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

	err = util.WaitForReceiptsOfTxs(client, txs, 5*time.Second)
	if err != nil {
		return err
	}

	return nil
}

func (g *Generator) prepareContractERC20() (common.Address, error) {
	client, err := ethclient.Dial(g.RpcUrl)
	if err != nil {
		return common.Address{}, err
	}
	defer client.Close()

	tx, err := GenerateContractCreationTx(
		g.FaucetAccount.PrivateKey,
		g.FaucetAccount.GetNonce(),
		g.ChainID,
		g.GasPrice,
		erc20.MyTokenBin,
		erc20.MyTokenABI,
		"My Token",
		"MYTOKEN",
	)
	if err != nil {
		return common.Address{}, err
	}

	err = client.SendTransaction(context.Background(), tx)
	if err != nil {
		return common.Address{}, err
	}

	ercContractAddress, err := bind.WaitDeployed(context.Background(), client, tx)
	if err != nil {
		return common.Address{}, err
	}

	if g.ShouldPersist {
		g.Store.AddPrepareTx(tx)
		if err != nil {
			return common.Address{}, err
		}
	}

	return ercContractAddress, nil
}
