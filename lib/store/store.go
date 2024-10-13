package store

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
)

type Store struct {
	TxStoreDir string
}

func NewStore(txStoreDir string) *Store {
	return &Store{
		TxStoreDir: txStoreDir,
	}
}

func (s *Store) PersistPrepareTxs(txs types.Transactions) error {
	return persistTxs(s.prepareFilePath(), txs)
}

func (s *Store) PersistTxsMap(txsMap map[int]types.Transactions) error {
	for index, txs := range txsMap {
		err := persistTxs(s.txsFilePath(index), txs)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *Store) prepareFilePath() string {
	return filepath.Join(s.TxStoreDir, "prepare.rlp")
}

func (s *Store) txsFilePath(index int) string {
	return filepath.Join(s.TxStoreDir, fmt.Sprintf("transactions-%d.rlp", index))
}

func persistTxs(path string, txs types.Transactions) error {
	dir := filepath.Dir(path)
	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return err
	}

	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	err = rlp.Encode(file, txs)
	if err != nil {
		return err
	}

	return nil
}
