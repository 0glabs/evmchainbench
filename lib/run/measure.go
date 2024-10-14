package run

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"golang.org/x/term"
)

type record struct {
	Height         *big.Int
	BlockTime      uint64
	TxCount        uint64
	GasLimit       uint64
	GasUsed        uint64
	PendingTxCount uint64
}

type txPoolStatus struct {
	Pending string `json:"pending"` // Number of pending transactions
	Queued  string `json:"queued"`  // Number of queued transactions
}

var terminalWith int

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

	for {
		currentBlock, err = client.BlockByNumber(context.Background(), currentBlockHeight)
		if err != nil {
			time.Sleep(500 * time.Millisecond)
			continue
		}

		r := record{}
		r.Height = currentBlockHeight
		r.TxCount = uint64(len(currentBlock.Transactions()))
		r.BlockTime = currentBlock.Time()
		r.GasLimit = currentBlock.GasLimit()
		r.GasUsed = currentBlock.GasUsed()

		pendingTxCount, err := getPendingTxCount(rpcClient)
		if err != nil {
			log.Printf("Failed to get pending txs: %v", err)
			return
		}

		if r.TxCount == 0 && pendingTxCount == 0 {
			break
		}

		r.PendingTxCount = pendingTxCount
		records = append(records, r)

		calculateAndOutput(records)

		currentBlockHeight.Add(currentBlockHeight, one)
		time.Sleep(200 * time.Millisecond)
	}
	fmt.Printf("\n")
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

	r0 := records[0]
	r := records[length-1]
	output1 := fmt.Sprintf("Height: %v  Tx: %v  PendingTx: %v  BlockTime: %v  GasLimit: %v  GasUsed: %v",
		r.Height,
		r.TxCount,
		r.PendingTxCount,
		r.BlockTime-r0.BlockTime,
		r.GasLimit,
		r.GasUsed,
	)

	output3 := fmt.Sprintf("TPS: %d", calculateTPS(records))

	spaceLength := terminalWith - len(output1) - len(output3) - 1
	if spaceLength < 0 {
		spaceLength = 0
	}
	output2 := strings.Repeat(" ", spaceLength)
	fmt.Printf("\r%s%s%s", output1, output2, output3)
}

func calculateTPS(records []record) int {
	length := len(records)
	if length <= 1 {
		return 0
	}

	count := 0
	for i := 1; i < length; i++ {
		count += int(records[i].TxCount)
	}

	timeSpan := int(records[length-1].BlockTime - records[0].BlockTime)
	if timeSpan <= 0 {
		return 0
	}

	return count / timeSpan
}

func init() {
	width, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		fmt.Printf("Error getting terminal size:", err)
		os.Exit(1)
	}
	terminalWith = width
}
