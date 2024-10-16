package generator

import (
	"math/big"
	"crypto/ecdsa"
	"encoding/hex"
	"strings"

	abipkg "github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

const GasLimit = uint64(2100000) // a random big enough gasLimit

func GenerateSimpleTransferTx(privateKey *ecdsa.PrivateKey, recipient string, nonce uint64, chainID, gasPrice, value *big.Int) (*types.Transaction, error) {
	return generateGenericTx(
		privateKey,
		recipient,
		nonce,
		chainID,
		gasPrice,
		value,
		[]byte{},
	)
}

func GenerateContractCreationTx(privateKey *ecdsa.PrivateKey, nonce uint64, chainID, gasPrice *big.Int, contractBin, contractABI string, args ...interface{}) (*types.Transaction, error) {
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
	return generateGenericTx(
		privateKey,
		"",
		nonce,
		chainID,
		gasPrice,
		big.NewInt(0),
		bytecode,
	)
}

func GenerateContractCallingTx(privateKey *ecdsa.PrivateKey, contractAddress string, nonce uint64, chainID, gasPrice *big.Int, contractABI, method string, args ...interface{}) (*types.Transaction, error) {
	abi, err := abipkg.JSON(strings.NewReader(contractABI))
	if err != nil {
		return &types.Transaction{}, err
	}

	data, err := abi.Pack(method, args...)
	if err != nil {
		return &types.Transaction{}, err
	}

	return generateGenericTx(
		privateKey,
		contractAddress,
		nonce,
		chainID,
		gasPrice,
		big.NewInt(0),
		data,
	)
}

func generateGenericTx(privateKey *ecdsa.PrivateKey, recipient string, nonce uint64, chainID, gasPrice, value *big.Int, data []byte) (*types.Transaction, error) {
	var tx *types.Transaction
	if recipient == "" {
		tx = types.NewContractCreation(
			nonce,
			big.NewInt(0),
			GasLimit,
			gasPrice,
			data,
		)
	} else {
		toAddress := common.HexToAddress(recipient)
		tx = types.NewTransaction(
			nonce,
			toAddress,
			value,
			GasLimit,
			gasPrice,
			nil,
		)
	}

	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		return &types.Transaction{}, err
	}

	return signedTx, nil
}
