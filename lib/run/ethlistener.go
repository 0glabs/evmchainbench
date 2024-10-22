package run

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/websocket"
)

type BlockInfo struct {
	Time     int64
	TxCount  int64
	GasUsed  int64
	GasLimit int64
}

type EthereumListener struct {
	wsURL            string
	conn             *websocket.Conn
	limiter          *RateLimiter
	blockStat        []BlockInfo
	quit             chan struct{}
	bestTPS          int64
	gasUsedAtBestTPS float64
}

func NewEthereumListener(wsURL string, limiter *RateLimiter) *EthereumListener {
	return &EthereumListener{
		wsURL:   wsURL,
		limiter: limiter,
		quit:    make(chan struct{}),
	}
}

func (el *EthereumListener) Connect() error {
	conn, _, err := websocket.DefaultDialer.Dial(el.wsURL, http.Header{})
	if err != nil {
		return fmt.Errorf("dial error: %v", err)
	}
	el.conn = conn
	return nil
}

func (el *EthereumListener) SubscribeNewHeads() error {
	subscribeMsg := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      1,
		"method":  "eth_subscribe",
		"params":  []interface{}{"newHeads"},
	}
	err := el.conn.WriteJSON(subscribeMsg)
	if err != nil {
		return fmt.Errorf("subscribe error: %v", err)
	}

	go el.listenForMessages()

	return nil
}

func (el *EthereumListener) listenForMessages() {
	for {
		_, message, err := el.conn.ReadMessage()
		if err != nil {
			return
		}

		var response map[string]interface{}
		err = json.Unmarshal(message, &response)
		if err != nil {
			log.Println("unmarshal error:", err)
			continue
		}

		if method, ok := response["method"]; ok && method == "eth_subscription" {
			el.handleNewHead(response)
		} else {
			el.handleBlockResponse(response)
		}
	}
}

func (el *EthereumListener) handleNewHead(response map[string]interface{}) {
	params := response["params"].(map[string]interface{})
	result := params["result"].(map[string]interface{})

	blockNo := result["number"].(string)

	request := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      1,
		"method":  "eth_getBlockByNumber",
		"params":  []interface{}{blockNo, false},
	}
	err := el.conn.WriteJSON(request)
	if err != nil {
		log.Println("Failed to send block request:", err)
	}
}

func (el *EthereumListener) handleBlockResponse(response map[string]interface{}) {
	if result, ok := response["result"].(map[string]interface{}); ok {
		if txns, ok := result["transactions"].([]interface{}); ok {
			el.limiter.IncreaseLimit(len(txns))
			ts, _ := strconv.ParseInt(result["timestamp"].(string)[2:], 16, 64)
			gasUsed, _ := strconv.ParseInt(result["gasUsed"].(string)[2:], 16, 64)
			gasLimit, _ := strconv.ParseInt(result["gasLimit"].(string)[2:], 16, 64)
			el.blockStat = append(el.blockStat, BlockInfo{
				Time:     ts,
				TxCount:  int64(len(txns)),
				GasUsed:  gasUsed,
				GasLimit: gasLimit,
			})
			// keep only the last 60 seconds of blocks
			for {
				if len(el.blockStat) == 1 {
					break
				}
				if el.blockStat[len(el.blockStat)-1].Time-el.blockStat[0].Time > 60 {
					el.blockStat = el.blockStat[1:]
				} else {
					break
				}
			}
			timeSpan := el.blockStat[len(el.blockStat)-1].Time - el.blockStat[0].Time
			// calculate TPS and gas used percentage
			if timeSpan > 50 {
				totalTxCount := int64(0)
				totalGasLimit := int64(0)
				totalGasUsed := int64(0)
				for _, block := range el.blockStat {
					totalTxCount += block.TxCount
					totalGasLimit += block.GasLimit
					totalGasUsed += block.GasUsed
				}
				tps := totalTxCount / timeSpan
				gasUsedPercent := float64(totalGasUsed) / float64(totalGasLimit)
				if tps > el.bestTPS {
					el.bestTPS = tps
					el.gasUsedAtBestTPS = gasUsedPercent
				}
				fmt.Println("TPS:", tps, "GasUsed%:", gasUsedPercent*100)
				if totalTxCount < 100 {
					// exit if total tx count is less than 100
					fmt.Println("Best TPS:", el.bestTPS, "GasUsed%:", el.gasUsedAtBestTPS*100)
					el.Close()
				}
			}
		}
	}
}

func (el *EthereumListener) Close() {
	if el.conn != nil {
		el.conn.Close()
	}
	close(el.quit)
}
