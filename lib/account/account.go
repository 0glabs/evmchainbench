package account

import (
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

type Account struct {
	Nonce      uint64
	Address    common.Address
	PrivateKey *ecdsa.PrivateKey
}

func NewAccount(client *ethclient.Client) (*Account, error) {
	pk, _ := crypto.GenerateKey()
	addr := crypto.PubkeyToAddress(pk.PublicKey)

	nonce, err := client.PendingNonceAt(context.Background(), addr)
	if err != nil {
		return nil, err
	}

	return &Account{
		Nonce:      nonce,
		Address:    addr,
		PrivateKey: pk,
	}, nil
}

func CreateFaucetAccount(client *ethclient.Client, privateKey string) (*Account, error) {
	pk, err := convertPrivateKeyFromStringForm(privateKey) 
	if err != nil {
		return &Account{}, err
	}

	addr := crypto.PubkeyToAddress(pk.PublicKey)

	nonce, err := client.PendingNonceAt(context.Background(), addr)
	if err != nil {
		return nil, fmt.Errorf("failed to get pending nonce: %w", err)
	}

	return &Account{
		Nonce:      nonce,
		Address:    addr,
		PrivateKey: pk,
	}, nil
}

func (account *Account) GetNonce() uint64 {
	now := account.Nonce
	account.Nonce += 1
	return now
}

func convertPrivateKeyFromStringForm(privateKey string) (*ecdsa.PrivateKey, error) {
	pks := strings.TrimPrefix(privateKey, "0x")
	pkBytes, _ := hex.DecodeString(pks)
	pk, err := crypto.ToECDSA(pkBytes)
	if err != nil {
		return &ecdsa.PrivateKey{}, err
	}

	return pk, nil
}
