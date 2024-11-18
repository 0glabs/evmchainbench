package generator

import (
	"fmt"
	"math/big"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/0glabs/evmchainbench/lib/account"
	"github.com/0glabs/evmchainbench/lib/contract_meta_data/erc20"
)

func (g *Generator) GenerateERC20() (map[int]types.Transactions, error) {
	txsMap := make(map[int]types.Transactions)

	if g.ShouldPersist {
		defer g.Store.PersistPrepareTxs()
	}

	contractAddress, err := g.prepareContractERC20()
	if err != nil {
		return txsMap, err
	}
	contractAddressStr := contractAddress.Hex()

	g.prepareSenders()

	g.prepareERC20(contractAddressStr)

	amount := big.NewInt(1000) // a random small amount

	var mutex sync.Mutex
	ch := make(chan error)

	sender := g.Senders[0]
	tx := GenerateContractCallingTx(
		sender.PrivateKey,
		contractAddressStr,
		1,
		g.ChainID,
		g.GasPrice,
		erc20TransferGasLimit,
		erc20.MyTokenABI,
		"transfer",
		common.HexToAddress(g.Recipients[0]),
		amount,
	)
	ethCallTx := ConvertLegacyTxToCallMsg(tx, sender.Address)
	estimateGas := g.estimateGas(ethCallTx)

	fmt.Println("Estimated gas:", estimateGas)

	for index, sender := range g.Senders {
		go func(index int, sender *account.Account) {
			txs := types.Transactions{}
			for _, recipient := range g.Recipients {
				tx := GenerateContractCallingTx(
					sender.PrivateKey,
					contractAddressStr,
					sender.GetNonce(),
					g.ChainID,
					g.GasPrice,
					estimateGas,
					erc20.MyTokenABI,
					"transfer",
					common.HexToAddress(recipient),
					amount,
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

func (g *Generator) prepareContractERC20() (common.Address, error) {
	return g.deployContract(erc20ContractGasLimit, erc20.MyTokenBin, erc20.MyTokenABI, "My Token", "MYTOKEN")
}
