package generator

import (
	"math/big"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/0glabs/evmchainbench/lib/account"
	"github.com/0glabs/evmchainbench/lib/contract_meta_data/erc20"
	"github.com/0glabs/evmchainbench/lib/contract_meta_data/uniswap"
)

func (g *Generator) GenerateUniswap() (map[int]types.Transactions, error) {
	txsMap := make(map[int]types.Transactions)

	pairContract, err := g.prepareContractUniswap()
	if err != nil {
		return txsMap, err
	}
	pairContractStr := pairContract.Hex()

	err = g.prepareSenders()
	if err != nil {
		return txsMap, err
	}

	var mutex sync.Mutex
	ch := make(chan error)

	for index, sender := range g.Senders {
		go func(index int, sender *account.Account) {
			txs := types.Transactions{}
			amount0out := &big.Int{}
			amount1out := &big.Int{}
			for ind, _ := range g.Recipients {
				if ind % 2 == 1 {
					amount0out = big.NewInt(1000)
					amount1out = big.NewInt(0)
				} else {
					amount0out = big.NewInt(0)
					amount1out = big.NewInt(1000)
				}
				tx, err := GenerateContractCallingTx(
					sender.PrivateKey,
					pairContractStr,
					sender.GetNonce(),
					g.ChainID,
					g.GasPrice,
					uniswap.UniswapV2PairABI,
					"swap",
					amount0out,
					amount1out,
					sender.Address,
					[]byte{},
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

func (g *Generator) prepareContractUniswap() (common.Address, error) {
	tokenA, err := g.deployContract(erc20.MyTokenBin, erc20.MyTokenABI, "Token A", "TOKENA")
	if err != nil {
		return common.Address{}, err
	}

	tokenB, err := g.deployContract(erc20.MyTokenBin, erc20.MyTokenABI, "Token B", "TOKENB")
	if err != nil {
		return common.Address{}, err
	}


	uniswapContract, err :=g.deployContract(uniswap.UniswapV2FactoryBin, uniswap.UniswapV2FactoryABI, g.FaucetAccount.Address)
	if err != nil {
		return common.Address{}, err
	}

	err = g.executeContractFunction(uniswapContract, uniswap.UniswapV2FactoryABI, "createPair", tokenA, tokenB)
	if err != nil {
		return common.Address{}, err
	}

	data, err := g.callContractView(uniswapContract, uniswap.UniswapV2FactoryABI, "getPair", tokenA, tokenB)
	if err != nil {
		return common.Address{}, err
	}

	// we know the first returned value is the address of the pair
	pairContract := data[0].(common.Address)

	amount := big.NewInt(1000000000000000000)
	err = g.executeContractFunction(tokenA, erc20.MyTokenABI, "transfer", pairContract, amount)
	if err != nil {
		return common.Address{}, err
	}

	err = g.executeContractFunction(tokenB, erc20.MyTokenABI, "transfer", pairContract, amount)
	if err != nil {
		return common.Address{}, err
	}

	err = g.executeContractFunction(pairContract, uniswap.UniswapV2PairABI, "mint", g.FaucetAccount.Address)
	if err != nil {
		return common.Address{}, err
	}

	return pairContract, nil

}
