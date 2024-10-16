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

var (
	terminalWith            int
	slidingWindowBeginIndex int
	finalTPS                int
)

func init() {
	width, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		fmt.Printf("Error getting terminal size:", err)
		os.Exit(1)
	}
	terminalWith = width
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
	if finalTPS == 0 {
		fmt.Printf("\nNot enough transactions for sampling. You may want to specify more txs with option \"-t\"\n")
	} else {
		fmt.Printf("\nThe best one minute TPS: %d\n", finalTPS)
		fmt.Printf("The GasUsed/GasLimit is: %d%%\n", int(100*getRatioOfGasUsedByGasLimit(records)))
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
	output1 := fmt.Sprintf("Time: %v    Tx: %v    PendingTx: %v    GasLimit: %v    GasUsed: %v",
		r.BlockTime-records[0].BlockTime,
		r.TxCount,
		r.PendingTxCount,
		r.GasLimit,
		r.GasUsed,
	)

	oneMinTPS := calculateOneMinTPS(records)
	output3 := fmt.Sprintf("TPS: %d", oneMinTPS)

	spaceLength := terminalWith - len(output1) - len(output3) - 1
	if spaceLength < 0 {
		spaceLength = 0
	}
	output2 := strings.Repeat(" ", spaceLength)
	fmt.Printf("\r%s%s%s", output1, output2, output3)
}

func calculateOneMinTPS(records []record) int {
	length := len(records)
	if length <= 1 {
		return 0
	}

	i := slidingWindowBeginIndex // begining of the one min sliding window
	j := length - 1              // end of the one min window

	for {
		if i+1 >= j {
			break
		}

		if records[j].BlockTime-records[i+1].BlockTime < 60 { // if possible, we sample data no less than 1 minute
			break
		}

		i += 1
	}

	slidingWindowBeginIndex = i // cache the beginning index of the sliding window

	count := 0
	for k := i + 1; k <= j; k++ {
		count += int(records[k].TxCount)
	}

	timeSpan := int(records[j].BlockTime - records[i].BlockTime)
	if timeSpan < 50 {
		return 0
	}

	tps := count / timeSpan

	if tps > finalTPS {
		finalTPS = tps
	}

	return tps
}

func getRatioOfGasUsedByGasLimit(records []record) float64 {
	length := len(records)
	if length <= 1 {
		return 0.0
	}

	var totalGasLimit uint64
	var totalGasUsed uint64

	for i := 1; i < length; i++ {
		totalGasLimit += records[i].GasLimit
		totalGasUsed += records[i].GasUsed
	}

	return float64(totalGasUsed) / float64(totalGasLimit)
}
