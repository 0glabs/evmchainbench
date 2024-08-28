package run

import (
	"context"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

type Transmitter struct {
	RpcUrl string
}

func NewTransmitter(rpcUrl string) (*Transmitter, error) {
	return &Transmitter{
		RpcUrl: rpcUrl,
	}, nil
}

func (t *Transmitter) Broadcast(txsMap map[int]types.Transactions) error {
	ch := make(chan error)

	for _, txs := range txsMap {
		go func(txs []*types.Transaction) {
			client, err := ethclient.Dial(t.RpcUrl)
			if err != nil {
				ch <- err
				return
			}

			for _, tx := range txs {
				err := broadcast(client, tx)
				if err != nil {
					ch <- err
					return
				}
			}

			ch <- nil
			return
		}(txs)

	}

	senderCount := len(txsMap)
	for i := 0; i < senderCount; i++ {
		err := <-ch
		if err != nil {
			return err
		}
	}

	return nil
}

func broadcast(client *ethclient.Client, tx *types.Transaction) error {
	err := client.SendTransaction(context.Background(), tx)
	if err != nil {
		return err
	}

	// Check tx hash
	// the hash can be abtained: tx.Hash().Hex()
	return nil
}
