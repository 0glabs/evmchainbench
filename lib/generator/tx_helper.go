package generator

import (
	"crypto/ecdsa"
	"encoding/hex"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum"
	abipkg "github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

func GenerateSimpleTransferTx(privateKey *ecdsa.PrivateKey, recipient string, nonce uint64, chainID, gasPrice, value *big.Int, eip1559 bool) (*types.Transaction, error) {
	toAddress := common.HexToAddress(recipient)

	var signedTx *types.Transaction
	var err error
	if eip1559 {
		tx := types.NewTx(&types.DynamicFeeTx{
			ChainID:   chainID,
			Nonce:     nonce,
			To:        &toAddress,
			Value:     value,
			GasFeeCap: gasPrice,
			GasTipCap: gasPrice,
			Gas:       simpleTransferGasLimit,
			Data:      nil,
		})
		signedTx, err = types.SignTx(tx, types.NewLondonSigner(chainID), privateKey)
	} else {
		tx := types.NewTransaction(nonce, toAddress, value, simpleTransferGasLimit, gasPrice, nil)
		signedTx, err = types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	}

	if err != nil {
		return &types.Transaction{}, err
	}

	return signedTx, nil
}

func GenerateContractCreationTx(privateKey *ecdsa.PrivateKey, nonce uint64, chainID, gasPrice *big.Int, gasLimit uint64, contractBin, contractABI string, args ...interface{}) (*types.Transaction, error) {
	bytecode, err := hex.DecodeString(contractBin)
	if err != nil {
		return &types.Transaction{}, err
	}

	if len(args) > 0 {
		abi, err := abipkg.JSON(strings.NewReader(contractABI))
		if err != nil {
			return &types.Transaction{}, err
		}

		inputData, err := abi.Pack("", args...)
		if err != nil {
			return &types.Transaction{}, err
		}

		bytecode = append(bytecode, inputData...)

	}

	tx := types.NewContractCreation(
		nonce,
		big.NewInt(0),
		gasLimit,
		gasPrice,
		bytecode,
	)

	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		return &types.Transaction{}, err
	}

	return signedTx, nil
}

func ConvertLegacyTxToCallMsg(tx *types.Transaction, from common.Address) ethereum.CallMsg {
	return ethereum.CallMsg{
		From:     from,
		To:       tx.To(),
		Gas:      tx.Gas(),
		GasPrice: tx.GasPrice(),
		Value:    tx.Value(),
		Data:     tx.Data(),
	}
}

func GenerateContractCallingTx(privateKey *ecdsa.PrivateKey, contractAddress string, nonce uint64, chainID, gasPrice *big.Int, gasLimit uint64, contractABI, method string, args ...interface{}) (*types.Transaction, error) {
	abi, err := abipkg.JSON(strings.NewReader(contractABI))
	if err != nil {
		return &types.Transaction{}, err
	}

	data, err := abi.Pack(method, args...)
	if err != nil {
		return &types.Transaction{}, err
	}

	toAddress := common.HexToAddress(contractAddress)
	tx := types.NewTransaction(
		nonce,
		toAddress,
		big.NewInt(0),
		gasLimit,
		gasPrice,
		data,
	)

	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		return &types.Transaction{}, err
	}

	return signedTx, nil
}
