package run

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"strconv"
	"time"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
)

type record struct {
	Height    *big.Int
	BlockTime uint64
	TxCount   uint64
	GasLimit  uint64
	GasUsed   uint64
	PendingTxCount uint64
}

type txPoolStatus struct {
        Pending string `json:"pending"` // Number of pending transactions
        Queued  string `json:"queued"`  // Number of queued transactions
}

func MeasureTPS(rpcUrl string) {
	client, err := ethclient.Dial(rpcUrl)
	if err != nil {
		log.Printf("Failed to connect to the Ethereum client: %v", err)
		return
	}

	rpcClient, err := rpc.Dial(rpcUrl)
	if err != nil {
		log.Printf("Failed to connect to the Ethereum client: %v", err)
		return
	}

        one := big.NewInt(1)

        currentBlock, err := client.BlockByNumber(context.Background(), nil)
        if err != nil {
                log.Printf("Failed to get the start block: %v", err)
		return
        }

        currentBlockHeight := currentBlock.Number()
	records := []record{}

        for ; ;  {
                currentBlock, err = client.BlockByNumber(context.Background(), currentBlockHeight)
                if err != nil {
                        time.Sleep(500 * time.Millisecond)
                        continue
                }

		r := record{}
		r.Height    = currentBlockHeight
		r.TxCount   = uint64(len(currentBlock.Transactions()))
		r.BlockTime = currentBlock.Time()
		r.GasLimit  = currentBlock.GasLimit()
		r.GasUsed   = currentBlock.GasUsed()

		pendingTxCount, err := getPendingTxCount(rpcClient)
		if err != nil {
			log.Printf("Failed to get pending txs: %v", err)
			return
		}

		if r.TxCount == 0 && pendingTxCount == 0 {
			return
		}

		r.PendingTxCount = pendingTxCount
		records = append(records, r)

		calculateAndOutput(records)

                currentBlockHeight.Add(currentBlockHeight, one)
                time.Sleep(200 * time.Millisecond)
        }
}

func getPendingTxCount(rpcClient *rpc.Client) (uint64, error) {
	status := txPoolStatus{}
	err := rpcClient.CallContext(context.Background(), &status, "txpool_status")
	if err != nil {
		return 0, err
	}
	pendingTxCount, err := strconv.ParseUint(status.Pending[2:], 16, 64)
	if err != nil {
		return 0, err
	}

	return pendingTxCount, nil
}

func calculateAndOutput(records []record) {
	length := len(records)
	if length == 0 {
		return
	}
	
	r := records[length-1]
	fmt.Printf("\rHeight: %v  TxCount: %v  PendingTx: %v  BlockTime: %v  GasLimit: %v  GasUsed: %v",
		r.Height,
		r.TxCount,
		r.PendingTxCount,
		r.BlockTime,
		r.GasLimit,
		r.GasUsed,
	)
}
