package account

import (
	"crypto/rand"
	"encoding/hex"
)

// GenerateRandomAddress generates a random Ethereum address without a private key
func GenerateRandomAddress() (string, error) {
	address := make([]byte, 20) // Ethereum address is 20 bytes long
	_, err := rand.Read(address)
	if err != nil {
		return "", err
	}

	return "0x" + hex.EncodeToString(address), nil
}
