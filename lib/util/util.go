package util

import (
	"context"
	"errors"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

func WaitForReceiptsOfTxs(client *ethclient.Client, txs types.Transactions, timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	for _, tx := range txs {
		txHash := tx.Hash()
		for {
			_, err := client.TransactionReceipt(context.Background(), txHash)
			if err == nil {
				break
			}

			if err != ethereum.NotFound {
				return err
			}

			select {
			case <-ctx.Done():
				return errors.New("Timeout before finding all receipts of txs")
			default:
				time.Sleep(500 * time.Millisecond)
			}
		}
	}

	return nil
}
