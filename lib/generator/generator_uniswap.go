package generator

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/big"
	"os"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/0glabs/evmchainbench/lib/account"
	"github.com/0glabs/evmchainbench/lib/contract_meta_data/erc20"
	"github.com/0glabs/evmchainbench/lib/contract_meta_data/uniswap"
)

func (g *Generator) GenerateUniswap() (map[int]types.Transactions, error) {
	txsMap := make(map[int]types.Transactions)

	if g.ShouldPersist {
		defer g.Store.PersistPrepareTxs()
	}

	tokenA, err := g.deployContract(erc20ContractGasLimit, erc20.MyTokenBin, erc20.MyTokenABI, "Token A", "TOKENA")
	if err != nil {
		return txsMap, err
	}

	tokenB, err := g.deployContract(erc20ContractGasLimit, erc20.MyTokenBin, erc20.MyTokenABI, "Token B", "TOKENB")
	if err != nil {
		return txsMap, err
	}

	fmt.Println("Token A:", tokenA.Hex(), "Token B:", tokenB.Hex())

	g.prepareSenders()
	g.prepareERC20(tokenA.Hex())
	g.prepareERC20(tokenB.Hex())

	var data []interface{}

	factory, router := g.prepareContractUniswap()
	fmt.Println("Factory contract:", factory.Hex())
	fmt.Println("Router contract:", router.Hex())

	g.approveERC20(tokenA, router)
	g.approveERC20(tokenB, router)

	data = g.callContractView(tokenA, uniswap.UniswapV2ERC20ABI, "balanceOf", g.FaucetAccount.Address)
	fmt.Println("Token A balance: ", data[0].(*big.Int).String())
	data = g.callContractView(tokenA, uniswap.UniswapV2ERC20ABI, "allowance", g.FaucetAccount.Address, router)
	fmt.Println("Token A allowance: ", data[0].(*big.Int).String())
	data = g.callContractView(tokenB, uniswap.UniswapV2ERC20ABI, "balanceOf", g.FaucetAccount.Address)
	fmt.Println("Token B balance: ", data[0].(*big.Int).String())
	data = g.callContractView(tokenB, uniswap.UniswapV2ERC20ABI, "allowance", g.FaucetAccount.Address, router)
	fmt.Println("Token B allowance: ", data[0].(*big.Int).String())

	g.executeContractFunction(uniswapCreatePairGasLimit, factory, uniswap.UniswapV2FactoryABI, "createPair", tokenA, tokenB)
	data = g.callContractView(factory, uniswap.UniswapV2FactoryABI, "getPair", tokenA, tokenB)
	fmt.Println("Pair address: ", data[0].(common.Address).Hex())

	var tx *types.Transaction
	var ethCallTx ethereum.CallMsg
	var estimateGas uint64

	fmt.Println("Add liquidity")

	g.executeContractFunction(uniswapCreatePairGasLimit, router, uniswap.UniswapV2RouterABI, "addLiquidity",
		tokenA, tokenB, big.NewInt(1000000000), big.NewInt(1000000000), big.NewInt(0), big.NewInt(0), g.FaucetAccount.Address,
		big.NewInt(time.Now().Unix()+15*60))

	var mutex sync.Mutex
	ch := make(chan error)

	sender := g.Senders[0]
	path := []common.Address{
		common.HexToAddress(tokenA.Hex()),
		common.HexToAddress(tokenB.Hex()),
	}
	deadline := big.NewInt(time.Now().Unix() + 15*60)

	tx = GenerateContractCallingTx(
		sender.PrivateKey,
		router.Hex(),
		0,
		g.ChainID,
		g.GasPrice,
		uniswapSwapGasLimit,
		uniswap.UniswapV2RouterABI,
		"swapExactTokensForTokens",
		big.NewInt(1000),
		big.NewInt(0),
		path,
		sender.Address,
		deadline,
	)
	ethCallTx = ConvertLegacyTxToCallMsg(tx, sender.Address)
	estimateGas = g.estimateGas(ethCallTx)
	estimateGas = (uint64)(1.2 * float64(estimateGas))

	fmt.Println("Estimated gas:", estimateGas)

	for index, sender := range g.Senders {
		go func(index int, sender *account.Account) {
			txs := types.Transactions{}
			for range g.Recipients {
				tx := GenerateContractCallingTx(
					sender.PrivateKey,
					router.Hex(),
					sender.GetNonce(),
					g.ChainID,
					g.GasPrice,
					estimateGas,
					uniswap.UniswapV2RouterABI,
					"swapExactTokensForTokens",
					big.NewInt(1000),
					big.NewInt(0),
					path,
					sender.Address,
					deadline,
				)
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

type Contract struct {
	Abi      []interface{} `json:"abi"`
	Bytecode string        `json:"bytecode"`
}

func ReadContract(filePath string) (string, string) {
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("Failed to open file: %v", err)
	}
	defer file.Close()

	fileData, err := io.ReadAll(file)
	if err != nil {
		log.Fatalf("Failed to read file: %v", err)
	}

	var contract Contract
	err = json.Unmarshal(fileData, &contract)
	if err != nil {
		log.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	abiJSON, err := json.Marshal(contract.Abi)
	if err != nil {
		log.Fatalf("Failed to marshal ABI: %v", err)
	}

	return string(abiJSON), contract.Bytecode
}

func (g *Generator) prepareContractUniswap() (common.Address, common.Address) {
	factoryABI, factoryBin := ReadContract("contracts/UniswapV2Factory.json")
	factoryContract, err := g.deployContract(uniswapContractGasLimit, factoryBin, factoryABI, g.FaucetAccount.Address)
	if err != nil {
		panic(err)
	}
	fmt.Println("Uniswap Factory:", factoryContract.Hex())

	routerABI, routerBin := ReadContract("contracts/UniswapV2Router02.json")

	uniswap.UniswapV2RouterABI = routerABI
	routerContract, err := g.deployContract(uniswapContractGasLimit, routerBin, uniswap.UniswapV2RouterABI, factoryContract, factoryContract)
	if err != nil {
		panic(err)
	}
	fmt.Println("Uniswap Router:", routerContract.Hex())

	return factoryContract, routerContract
}
