package run

import (
	"context"
	"time"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"

	limiterpkg "github.com/0glabs/evmchainbench/lib/limiter"
)

type Transmitter struct {
	RpcUrl  string
	limiter *limiterpkg.RateLimiter
}

func NewTransmitter(rpcUrl string, limiter *limiterpkg.RateLimiter) (*Transmitter, error) {
	return &Transmitter{
		RpcUrl:  rpcUrl,
		limiter: limiter,
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
				for {
					if t.limiter == nil || t.limiter.AllowRequest() {
						err := broadcast(client, tx)
						if err != nil {
							ch <- err
							return
						}
						break
					} else {
						time.Sleep(10*time.Millisecond)
					}

				}
			}
			ch <- nil
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
