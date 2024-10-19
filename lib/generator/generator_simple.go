package generator

import (
	"math/big"
	"sync"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/0glabs/evmchainbench/lib/account"
)

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
				tx, err := GenerateSimpleTransferTx(sender.PrivateKey, recipient, sender.GetNonce(), g.ChainID, g.GasPrice, value, g.EIP1559)
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
